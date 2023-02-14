package firestore

import (
	"context"
	"log"
	"mahjong-linebot/app/models"
	logger "mahjong-linebot/logs"
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
	// contextにセットした値はinterface{}型のため.(string)でassertionが必要
	traceId := ctx.Value("traceId").(string)
	client, err := firebaseInit(ctx)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "firebaseInit失敗", err))
		return err
	}
	// 切断
	defer client.Close()
	m.CreateTimestamp = time
	m.UpdateTimestamp = time
	_, err = client.Collection("matchResults").Doc(m.DocId).Set(ctx, m)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed Add:matchSetting in firestore", err))
		return err
	}
	// エラーなしは成功
	return err
}

/*
*

	firestoreの試合結果を更新

*
*/
func UpdateMatchResult(ctx context.Context, m *models.MatchResult, time time.Time) error {
	// contextにセットした値はinterface{}型のため.(string)でassertionが必要
	traceId := ctx.Value("traceId").(string)
	client, err := firebaseInit(ctx)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "firebaseInit失敗", err))
		return err
	}
	// 切断
	defer client.Close()
	_, err = client.Collection("matchResults").Doc(m.DocId).Update(ctx, []firestore.Update{
		{
			Path:  "pointList",
			Value: m.PointList,
		},
		{
			Path:  "updateTimestamp",
			Value: time,
		},
	})
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed Add:matchSetting in firestore", err))
		return err
	}
	// エラーなしは成功
	return err
}

/*
*

	firestoreの試合結果をroomIdを元に検索

*
*/
func GetMatchResult(ctx context.Context, roomId int) (*[]models.MatchResult, error) {
	// contextにセットした値はinterface{}型のため.(string)でassertionが必要
	traceId := ctx.Value("traceId").(string)
	client, err := firebaseInit(ctx)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "firebaseInit失敗", err))
		return nil, err
	}
	// 切断
	defer client.Close()

	iter := client.Collection("matchResults").Where("roomId", "==", roomId).Documents(ctx)
	var docList []models.MatchResult
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf(logger.ErrorLogEntry(traceId, "Failed documents iterator", err))
			return nil, err
		}
		m := doc.Data()
		//TODO:PointListの型判定どうするか。
		ms := models.MatchResult{
			DocId:            m["docId"].(string),
			RoomId:           m["roomId"].(int64),
			PointList:        m["pointList"].([]interface{}),//[]models.PointOfPersonにすると型判定でpanicになる
			CreateTimestamp:  m["createTimestamp"].(time.Time),
			UpdateTimestamp:  m["updateTimestamp"].(time.Time),
		}
		docList = append(docList, ms)
	}
	if len(docList) == 0 {
		log.Printf(logger.InfoLogEntry(traceId, "Not Found matchResult data in roomId:", roomId))
		return nil, err
	}

	// エラーなしは成功
	return &docList, err
}
