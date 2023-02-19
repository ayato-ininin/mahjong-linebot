package firestore

import (
	"context"
	"errors"
	"log"
	"mahjong-linebot/app/models"
	logger "mahjong-linebot/logs"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

/*
*

	firestoreの設定をroomIdを元に検索
	roomIdをfirestoreから受け取るときは、int64になる。(firestoreの仕様)

*
*/
func GetMatchSetting(ctx context.Context, roomId int) (*models.MatchSetting, error) {
	// contextにセットした値はinterface{}型のため.(string)でassertionが必要
	traceId := ctx.Value("traceId").(string)
	client, err := firebaseInit(ctx)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "firebaseInit失敗", err))
		return nil, err
	}
	// 切断
	defer client.Close()

	// https://cloud.google.com/firestore/docs/query-data/queries?hl=ja
	// ここ、他のfirestoreアクセスできてるのに、permision deniedでてて、結果、IAMの権限不足やった。
	// 最小単位にするのもいいけど、こういう余計なエラーで止まるのは手間。
	iter := client.Collection("matchSettings").Where("roomId", "==", roomId).Limit(1).Documents(ctx)
	var docList []models.MatchSetting
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
		ms := models.MatchSetting{
			RoomId:           m["roomId"].(int64),
			MahjongNumber:    m["mahjongNumber"].(string),
			Name1:            m["name1"].(string),
			Name2:            m["name2"].(string),
			Name3:            m["name3"].(string),
			Name4:            m["name4"].(string),
			Uma:              m["uma"].(string),
			Oka:              m["oka"].(int64),
			IsYakitori:       m["isYakitori"].(bool),
			YakitoriPoint:    m["yakitoriPoint"].(int64),
			IsTobishou:       m["isTobishou"].(bool),
			TobishouPoint:    m["tobishouPoint"].(int64),
			Rate:             m["rate"].(int64),
			IsTip:            m["isTip"].(bool),
			TipInitialNumber: m["tipInitialNumber"].(int64),
			TipRate:          m["tipRate"].(int64),
			CreateTimestamp:  m["createTimestamp"].(time.Time),
			UpdateTimestamp:  m["updateTimestamp"].(time.Time),
		}
		docList = append(docList, ms)
	}
	if len(docList) == 0 {
		log.Printf(logger.ErrorLogEntry(traceId, "Not Found matchSetting data"))
		return nil, errors.New("Not found roomId in database" + strconv.Itoa(roomId))
	}

	// エラーなしは成功
	return &docList[0], err
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
		log.Printf(logger.ErrorLogEntry(traceId, "firebaseInit失敗", err))
		return err
	}
	// 切断
	defer client.Close()

	nextRoomNumber, err := GetNextRoomNumber(ctx, client)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed Get:nextRoomNumber in firestore", err))
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

	err = changeNextRoomNumber(ctx, client, nextRoomNumber, time)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed Change:nextRoomNumber in firestore", err))
	}

	// エラーなしは成功
	return err
}

/*
*

	次のルーム番号を返す

*
*/
func GetNextRoomNumber(ctx context.Context, client *firestore.Client) (int64, error) {
	dsnap, err := client.Collection("roomNumber").Doc("current").Get(ctx)
	if err != nil {
		return 0, err
	}
	m := dsnap.Data()
	nextRoomNumber := m["nextRoomNumber"].(int64)

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
func changeNextRoomNumber(ctx context.Context, client *firestore.Client, i int64, time time.Time) error {
	_, err := client.Collection("roomNumber").Doc("current").Update(ctx, []firestore.Update{
		{
			Path:  "nextRoomNumber",
			Value: i + 1,
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
