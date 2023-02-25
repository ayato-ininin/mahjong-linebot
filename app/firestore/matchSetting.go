package firestore

import (
	"context"
	"errors"
	"log"
	"mahjong-linebot/app/models"
	logger "mahjong-linebot/logs"
	"time"

	"cloud.google.com/go/firestore"
)

/*
*

	firestoreの設定をroomIdを元に検索
	roomIdをfirestoreから受け取るときは、int64になる。(firestoreの仕様)

*
*/
func GetMatchSetting(ctx context.Context, roomId int) (*models.MatchSetting, error) {
	// contextにセットした値はinterface{}型のため.(string)でassertionが必要
	traceId, ok := ctx.Value("traceId").(string)
	if !ok {
		return nil, errors.New("traceId not found in context")
	}
	client, err := firebaseInit(ctx)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "firebaseInit failed", err))
		return nil, errors.New("failed to initialize Firestore client")
	}
	// 切断
	defer client.Close()

	matchSetting, err := getMatchSettingByRoomId(ctx, client, roomId)
	if err != nil {
			log.Printf(logger.ErrorLogEntry(traceId, "failed to get matchSetting", err))
			return nil, errors.New("failed to get matchSetting")
	}

	return matchSetting, nil
}

func getMatchSettingByRoomId(ctx context.Context, client *firestore.Client, roomId int) (*models.MatchSetting, error) {
	iter := client.Collection("matchSettings").Where("roomId", "==", roomId).Limit(1).Documents(ctx)
	docs, err := iter.GetAll()//イテレータを使う必要がなくなり、コードの簡素化
	if err != nil {
			return nil, err
	}

	if len(docs) == 0 {
			return nil, errors.New("roomId not found in database")
	}

	var m models.MatchSetting
	doc := docs[0]
	err = doc.DataTo(&m)//各フィールドの初期化が必要なくなる
	if err != nil {
			return nil, err
	}

	return &m, nil
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
	traceId := ctx.Value("traceId").(string)
	client, err := firebaseInit(ctx)
	if err != nil {
		return err
	}
	// 切断
	defer client.Close()

	nextRoomNumber, err := GetNextRoomNumber(ctx, client, traceId)
	if err != nil {
		return err
	}
	m.RoomId = nextRoomNumber
	m.CreateTimestamp = time
	m.UpdateTimestamp = time
	_, err = client.Collection("matchSettings").Doc(time.Format(RFC3339)[0:19]).Set(ctx, m)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed Add:matchSetting in firestore", err))
		return err
	}

	err = changeNextRoomNumber(ctx, client, traceId)
	if err != nil {
		return err
	}

	// エラーなしは成功
	return err
}

/*
*

	次のルーム番号を返す

*
*/
func GetNextRoomNumber(ctx context.Context, client *firestore.Client, traceId string) (int64, error) {
	dsnap, err := client.Collection("roomNumber").Doc("current").Get(ctx)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed Get:nextRoomNumber in firestore", err))
		return 0, err
	}
	return dsnap.Data()["nextRoomNumber"].(int64), nil//ここでエラーになることはない
}

/*
*

	nextRoomNumberを+1する処理

*
*/
func changeNextRoomNumber(ctx context.Context, client *firestore.Client, traceId string) error {
	_, err := client.Collection("roomNumber").Doc("current").Update(ctx, []firestore.Update{
		{
			Path:  "nextRoomNumber",
			Value: firestore.Increment(1),//インクリメント
		},
		{
			Path:  "timestamp",
			Value: firestore.ServerTimestamp,//サーバーのタイムスタンプ
		},
	})
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed Change:nextRoomNumber in firestore", err))
		return err
	}
	return nil
}
