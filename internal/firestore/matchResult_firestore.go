package firestore

import (
	"context"
	"errors"
	"log"
	logger "mahjong-linebot/internal/logs"
	"mahjong-linebot/internal/models"
	"mahjong-linebot/internal/repositories"
	"mahjong-linebot/internal/utils"
	"time"
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

	err = repositories.SetMatchResult(ctx, client, m, time)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed Add:matchResult in firestore", err))
		return err
	}
	return nil
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

	err = repositories.UpdateMatchResultInFirestore(ctx, client, m, time)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "FirestoreのmatchResultsコレクションの更新に失敗しました: %v", err))
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

	matchResult, err := repositories.GetMatchResultByRoomId(ctx, client, roomId)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "failed to get matchSetting", err))
		return nil, errors.New("failed to get matchResult")
	}

	return matchResult, err
}
