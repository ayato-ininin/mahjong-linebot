package firestore

import (
	"context"
	"errors"
	"log"
	"mahjong-linebot/app/models"
	logger "mahjong-linebot/logs"
	"mahjong-linebot/utils"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

/*
*

	firestoreに試合結果を保存

*
*/
func AddMatchResult(ctx context.Context, m *models.MatchResult, time time.Time) error {
	traceId, err := utils.GetTraceID(ctx)
	if err != nil {
		return err
	}
	client, err := firebaseInit(ctx)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "firebaseInit失敗", err))
		return err
	}
	// 切断
	defer client.Close()

	err = setMatchResult(ctx, client, m, time)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed Add:matchResult in firestore", err))
		return err
	}
	return nil
}

// firestoreに試合結果を保存(DB接続)
func setMatchResult(ctx context.Context, client *firestore.Client, m *models.MatchResult, time time.Time) error {
	m.CreateTimestamp = time
	m.UpdateTimestamp = time
	_, err := client.Collection("matchResults").Doc(m.DocId).Set(ctx, m)
	return err
}

/*
*

	firestoreの試合結果を更新

*
*/
func UpdateMatchResult(ctx context.Context, m *models.MatchResult, time time.Time) error {
	traceId, err := utils.GetTraceID(ctx)
	if err != nil {
		return err
	}
	client, err := firebaseInit(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	err = updateMatchResultInFirestore(ctx, client, m, time)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "FirestoreのmatchResultsコレクションの更新に失敗しました: %v", err))
		return err
	}

	return nil
}

// firestoreの試合結果を更新(DB接続)
func updateMatchResultInFirestore(ctx context.Context, client *firestore.Client, m *models.MatchResult, time time.Time) error {
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

/*
*

	firestoreの試合結果をroomIdを元に検索

*
*/
func GetMatchResult(ctx context.Context, roomId int) (*[]models.MatchResult, error) {
	traceId, err := utils.GetTraceID(ctx)
	if err != nil {
		return nil, err
	}
	client, err := firebaseInit(ctx)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "firebaseInit失敗", err))
		return nil, err
	}
	// 切断
	defer client.Close()

	matchResult, err := getMatchResultByRoomId(ctx, client, roomId)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "failed to get matchSetting", err))
		return nil, errors.New("failed to get matchResult")
	}

	return matchResult, err
}

// firestoreの試合結果をroomIdを元に検索(DB接続)
func getMatchResultByRoomId(ctx context.Context, client *firestore.Client, roomId int) (*[]models.MatchResult, error) {
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
			RoomId:          m["roomId"].(int64),
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

/*
**

	[]interface{} 型の pointList スライスを []models.PointOfPerson 型に変換するために、pointsOfPerson 関数を作成する
	そのままだと、型判定でmodelの方と合わなくてエラーになる。

**
*/
func pointsOfPerson(ps []interface{}) []models.PointOfPerson {
	result := make([]models.PointOfPerson, 0, len(ps))
	for _, p := range ps {
		m := p.(map[string]interface{})
		result = append(result, models.PointOfPerson{
			NameIndex:  int64(m["nameIndex"].(float64)), //なぜかfloat64でくるので、int64に変換
			Point:      int64(m["point"].(float64)),     //なぜかfloat64でくるので、int64に変換
			IsYakitori: m["isYakitori"].(bool),
			IsTobishou: m["isTobishou"].(bool),
		})
	}
	return result
}
