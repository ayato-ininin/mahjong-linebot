package firestore

import (
	"context"
	"fmt"
	"log"
	logger "mahjong-linebot/utils"
	"strconv"
	"time"
	"mahjong-linebot/app/models"

	"cloud.google.com/go/firestore"
)

/*
*

	firestoreのに試合の設定を追加
	matchSettingsに追加するのと、nextRoomNumberの更新をトランザクションしないと
	片方ずれたらroomIdかぶるとかになるので、要検討

*
*/
func AddMatchSetting(m *models.MatchSetting, time time.Time) error {
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
	m.RoomId = nextRoomNumber
	m.CreateTimestamp = time
	m.UpdateTimestamp = time
	_, err = client.Collection("matchSettings").Doc(time.Format(RFC3339)[0:19]).Set(ctx, m)
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

// func GetMatchSettingResponse(status int32, message string, matchSetting *models.MatchSetting) *models.RespMatchSettingModel {
// 	res := &models.RespMatchSettingModel{
// 		Status:    status,
// 		Message:   message,
// 		MatchSetting: matchSetting,
// 	}

// 	return res
// }
