package firestore

import (
	"context"
	"log"
	"mahjong-linebot/app/models"
	logger "mahjong-linebot/logs"
	"time"
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

	_, err = client.Collection("matchResults").Doc(time.Format(RFC3339)[0:19]).Set(ctx, m)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed Add:matchSetting in firestore", err))
		return err
	}
	// エラーなしは成功
	return err
}
