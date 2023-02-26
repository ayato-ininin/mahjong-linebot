package repositories

import (
	"context"
	"mahjong-linebot/internal/models"
	"mahjong-linebot/internal/utils"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// firestoreに試合結果を保存
func SetMatchResult(ctx context.Context, client *firestore.Client, m *models.MatchResult, time time.Time) error {
	m.CreateTimestamp = time
	m.UpdateTimestamp = time
	_, err := client.Collection("matchResults").Doc(m.DocId).Set(ctx, m)
	return err
}

// firestoreの試合結果を更新
func UpdateMatchResultInFirestore(ctx context.Context, client *firestore.Client, m *models.MatchResult, time time.Time) error {
	_, err := client.Collection("matchResults").Doc(m.DocId).Update(ctx, []firestore.Update{
		{
			Path:  "pointList",
			Value: m.PointList,
		},
		{
			Path:  "updateTimestamp",
			Value: firestore.ServerTimestamp,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// firestoreの試合結果をroomIdを元に検索(DB接続)
func GetMatchResultByRoomId(ctx context.Context, client *firestore.Client, roomId int) (*[]models.MatchResult, error) {
	iter := client.Collection("matchResults").Where("roomId", "==", roomId).Documents(ctx)
	docList := make([]models.MatchResult, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		m := doc.Data()
		ms := models.MatchResult{
			DocId:           m["docId"].(string),
			RoomId:          utils.ConvertInt64(m["roomId"]),
			PointList:       pointsOfPerson(m["pointList"].([]interface{})), //[]interface{}型のスライスでくるので、[]models.PointOfPerson型に変換
			CreateTimestamp: m["createTimestamp"].(time.Time),
			UpdateTimestamp: m["updateTimestamp"].(time.Time),
		}
		docList = append(docList, ms)
	}
	if len(docList) == 0 {
		return nil, nil
	}
	return &docList, nil
}

// []interface{} 型の pointList スライスを []models.PointOfPerson 型に変換するために、pointsOfPerson 関数を作成する
// そのままだと、型判定でmodelの方と合わなくてエラーになる。
func pointsOfPerson(ps []interface{}) []models.PointOfPerson {
	result := make([]models.PointOfPerson, 0, len(ps))
	for _, p := range ps {
		m := p.(map[string]interface{})
		result = append(result, models.PointOfPerson{
			NameIndex:  utils.ConvertInt64(m["nameIndex"]),
			Point:      utils.ConvertInt64(m["point"]),
			IsYakitori: m["isYakitori"].(bool),
			IsTobishou: m["isTobishou"].(bool),
		})
	}
	return result
}
