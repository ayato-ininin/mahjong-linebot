package firestore

import (
	"context"
	"errors"
	"log"
	logger "mahjong-linebot/pkg/logs"
	"mahjong-linebot/pkg/models"
	"mahjong-linebot/pkg/repositories"
	"mahjong-linebot/pkg/utils"
	"time"
)

/*
*

	firestoreの設定をroomIdを元に検索
	roomIdをfirestoreから受け取るときは、int64になる。(firestoreの仕様)

*
*/
func GetMatchSetting(ctx context.Context, roomId int) (*models.MatchSetting, error) {
	// contextにセットした値はinterface{}型のため.(string)でassertionが必要
	traceId, err := utils.GetTraceID(ctx)
	if err != nil {
		return nil, err
	}
	client, err := firebaseInit(ctx)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "firebaseInit failed", err))
		return nil, errors.New("failed to initialize Firestore client")
	}
	// 切断
	defer client.Close()

	matchSetting, err := repositories.GetMatchSettingByRoomId(ctx, client, roomId)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "failed to get matchSetting", err))
		return nil, errors.New("failed to get matchSetting")
	}

	return matchSetting, nil
}

/*
*

	firestoreのに試合の設定を追加
	matchSettingsに追加するのと、nextRoomNumberの更新をトランザクションしないと
	片方ずれたらroomIdかぶるとかになるので、要検討

*
*/
func AddMatchSetting(ctx context.Context, m *models.MatchSetting, time time.Time) error {
	// contextにセットした値はinterface{}型のため.(string)でassertionが必要
	traceId, err := utils.GetTraceID(ctx)
	if err != nil {
		return err
	}
	client, err := firebaseInit(ctx)
	if err != nil {
		return err
	}
	// 切断
	defer client.Close()

	nextRoomNumber, err := repositories.GetNextRoomNumber(ctx, client)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed Get:nextRoomNumber in firestore", err))
		return err
	}

	err = repositories.SetMatchSetting(ctx, client, m, time, nextRoomNumber)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed Add:matchSetting in firestore", err))
		return err
	}

	err = repositories.ChangeNextRoomNumber(ctx, client)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed Change:nextRoomNumber in firestore", err))
		return err
	}

	return nil
}
