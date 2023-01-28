package firestore

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	logger "mahjong-linebot/utils"
	"time"

	"mahjong-linebot/config"

	"cloud.google.com/go/firestore"
)

type GameStatus struct {
	Game   string //東風戦、半荘戦
	Number string //三麻、四麻
	Style  string //リアル、ネット
	Place  string //場所
}

type GameResult struct {
	Rank      string    `firestore:"rank"`
	Game      string    `firestore:"game"`
	Number    string    `firestore:"number"`
	Style     string    `firestore:"style"`
	Place     string    `firestore:"place"`
	Timestamp time.Time `firestore:"timestamp"`
}

const (
	//ちなみに、Z09:00をなくすと、Formatにしたら自動でUTCの時間になってしまう。
	RFC3339 = "2006-01-02T15:04:05Z09:00"
)

/*
*

	firestoreのcurrentデータを更新。
	（東風戦、半荘戦）、（三麻、四麻）を試合前に登録しておく。

*
*/
func UpdateGameStatusData(text, param string, time time.Time) error {
	ctx := context.Background()
	client, err := firebaseInit(ctx)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("firebaseInit失敗 err=%v", err)))
		return err
	}
	_, err = client.Collection("gameStatus").Doc("current").Update(ctx, []firestore.Update{
		{
			Path:  param,
			Value: text,
		},
		{
			Path:  "timestamp",
			Value: time,
		},
	})
	if err != nil {
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Failed Change:gameStatus in firestore err=%v", err)))
		return err
	}

	// 切断
	defer client.Close()

	// エラーなしは成功
	return err
}

/*
*

	firestoreのcurrentデータから、試合の種類を取得して順位を登録。

*
*/
func AddRankData(text string, time time.Time) error {
	ctx := context.Background()
	client, err := firebaseInit(ctx)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("firebaseInit失敗 err=%v", err)))
		return err
	}
	// 切断
	defer client.Close()

	dsnap, err := client.Collection("gameStatus").Doc("current").Get(ctx)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Failed Get:gameStatus in firestore err=%v", err)))
		return err
	}
	m := dsnap.Data()
	_, err = client.Collection("ranks").Doc(time.Format(RFC3339)[0:19]).Set(ctx, GameResult{
		Rank:      text,
		Game:      m["game"].(string),
		Number:    m["number"].(string),
		Style:     m["style"].(string),
		Place:     m["place"].(string),
		Timestamp: time,
	})
	if err != nil {
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Failed Add:ranks in firestore err=%v", err)))
		return err
	}

	// エラーなしは成功
	return err
}

/*
*

	現在のステータスを返す

*
*/
func GetCurrentStatus() (*GameStatus, error) {
	ctx := context.Background()
	client, err := firebaseInit(ctx)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("firebaseInit失敗 err=%v", err)))
		return nil, err
	}
	dsnap, err := client.Collection("gameStatus").Doc("current").Get(ctx)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Failed Get:gameStatus in firestore err=%v", err)))
		return nil, err
	}
	m := dsnap.Data()
	status := GameStatus{
		Game:   m["game"].(string),
		Number: m["number"].(string),
		Style:  m["style"].(string),
		Place:  m["place"].(string),
	}

	// 切断
	defer client.Close()

	// エラーなしは成功
	return &status, err
}

/*
*

	ステータス：レスポンス用メッセージ作成

*
*/
func CreateStatusMsg(status *GameStatus) string {
	msg := "【設定】 \n" +
		"・試合　　：" + status.Game + "\n" +
		"・人数　　：" + status.Number + "\n" +
		"・スタイル：" + status.Style + "\n" +
		"・場所　　：" + status.Place

	return msg
}

/*
*
secret managerにbase64でencodeして保存している、
firebaseの認証用jsonをバイト配列にて返す
*
*/
func getFirebaseServiceAccountKey() []byte {
	data := *config.GetDataFromSecretManager("MAHJONG_LINEBOT_FIREBASE_SA")
	dec, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Failed firebaseSA decode: err=%v", err)))
	}

	return dec
}
