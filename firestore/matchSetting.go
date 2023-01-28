package firestore

import (
	"context"
	"fmt"
	"log"
	logger "mahjong-linebot/utils"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
)

type MatchSetting struct {
	RoomId           string    `firestore:"roomId"`
	MahjongNumber    string    `firestore:"mahjongNumber"`
	Uma              string    `firestore:"uma"`
	Oka              int       `firestore:"oka"`
	IsYakitori       bool      `firestore:"isYakitori"`
	IsTobishou       bool      `firestore:"isTobishou"`
	TobishouPoint    int       `firestore:"tobishouPoint"`
	Rate             int       `firestore:"rate"`
	IsTip            bool      `firestore:"isTip"`
	TipInitialNumber int       `firestore:"tipInitialNumber"`
	TipRate          int       `firestore:"tipRate"`
	CreateTimestamp  time.Time `firestore:"createTimestamp"`
	UpdateTimestamp  time.Time `firestore:"updateTimestamp"`
}

/*
*

	firestoreのに試合の設定を追加
	matchSettingsに追加するのと、nextRoomNumberの更新をトランザクションしないと
	片方ずれたらroomIdかぶるとかになるので、要検討

*
*/
func AddMatchSetting(m *MatchSetting, time time.Time) error {
	ctx := context.Background()
	client, err := firebaseInit(ctx)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("firebaseInit失敗 err=%v", err)))
		return err
	}
	// 切断
	defer client.Close()

	nextRoomNumber, err := GetNextRoomNumber(ctx, client)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Failed Get:nextRoomNumber in firestore err=%v", err)))
		return err
	}
	_, err = client.Collection("matchSettings").Doc(time.Format(RFC3339)[0:19]).Set(ctx, MatchSetting{
		RoomId:           nextRoomNumber,
		MahjongNumber:    m.MahjongNumber,
		Uma:              m.Uma,
		Oka:              m.Oka,
		IsYakitori:       m.IsYakitori,
		IsTobishou:       m.IsTobishou,
		TobishouPoint:    m.TobishouPoint,
		Rate:             m.Rate,
		IsTip:            m.IsTip,
		TipInitialNumber: m.TipInitialNumber,
		TipRate:          m.TipRate,
		CreateTimestamp:  time,
		UpdateTimestamp:  time,
	})
	if err != nil {
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Failed Add:matchSetting in firestore err=%v", err)))
		return err
	}

	err = changeNextRoomNumber(ctx, client, nextRoomNumber, time)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Failed Change:nextRoomNumber in firestore err=%v", err)))
	}

	// エラーなしは成功
	return err
}

/*
*

	次のルーム番号を返す

*
*/
func GetNextRoomNumber(ctx context.Context, client *firestore.Client) (string, error) {
	dsnap, err := client.Collection("roomNumber").Doc("current").Get(ctx)
	if err != nil {
		return "0", err
	}
	m := dsnap.Data()
	nextRoomNumber := m["nextRoomNumber"].(string)

	// エラーなしは成功
	return nextRoomNumber, err
}

/*
*

	nextRoomNumberを+1する処理
	firestoreでnumberで保存するとint64扱いになるから、
	文字列で一旦保存するためにitoaをしている。
	数字で扱えるようにしてもいいかも。

*
*/
func changeNextRoomNumber(ctx context.Context, client *firestore.Client, s string, time time.Time) error {
	i, _ := strconv.Atoi(s)
	_, err := client.Collection("roomNumber").Doc("current").Update(ctx, []firestore.Update{
		{
			Path:  "nextRoomNumber",
			Value: strconv.Itoa(i + 1),
		},
		{
			Path:  "timestamp",
			Value: time,
		},
	})
	if err != nil {
		return err
	}

	// エラーなしは成功
	return err
}
