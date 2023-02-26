package repositories

import (
	"context"
	"errors"
	"mahjong-linebot/internal/models"
	"time"

	"cloud.google.com/go/firestore"
)

const (
	//Z09:00をなくすと、Formatにしたら自動でUTCの時間になってしまう。
	RFC3339 = "2006-01-02T15:04:05Z09:00"
)

// roomIdを元にfirestoreから試合設定を取得(DB接続)
func GetMatchSettingByRoomId(ctx context.Context, client *firestore.Client, roomId int) (*models.MatchSetting, error) {
	iter := client.Collection("matchSettings").Where("roomId", "==", roomId).Limit(1).Documents(ctx)
	docs, err := iter.GetAll() //イテレータを使う必要がなくなり、コードの簡素化
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, errors.New("roomId not found in database")
	}

	var m models.MatchSetting
	doc := docs[0]
	err = doc.DataTo(&m) //各フィールドの初期化が必要なくなる
	if err != nil {
		return nil, err
	}

	return &m, nil
}

// firestoreに試合結果を保存(DB接続)
func SetMatchSetting(ctx context.Context, client *firestore.Client, m *models.MatchSetting, time time.Time, nextNum int64) error {
	m.RoomId = nextNum
	m.CreateTimestamp = time
	m.UpdateTimestamp = time
	_, err := client.Collection("matchSettings").Doc(time.Format(RFC3339)[0:19]).Set(ctx, m)
	return err
}

// 次のルーム番号を返す(DB接続)
func GetNextRoomNumber(ctx context.Context, client *firestore.Client) (int64, error) {
	dsnap, err := client.Collection("roomNumber").Doc("current").Get(ctx)
	if err != nil {
		return 0, err
	}
	return dsnap.Data()["nextRoomNumber"].(int64), nil //ここでエラーになることはない
}

// nextRoomNumberを+1する処理(DB接続)
func ChangeNextRoomNumber(ctx context.Context, client *firestore.Client) error {
	_, err := client.Collection("roomNumber").Doc("current").Update(ctx, []firestore.Update{
		{
			Path:  "nextRoomNumber",
			Value: firestore.Increment(1), //インクリメント
		},
		{
			Path:  "timestamp",
			Value: firestore.ServerTimestamp, //サーバーのタイムスタンプ
		},
	})
	if err != nil {
		return err
	}
	return nil
}
