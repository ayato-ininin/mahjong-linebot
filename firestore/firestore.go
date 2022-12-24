package firestore

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type GameResult struct {
	Rank      string    `firestore:"rank"`
	Game      string    `firestore:"game"`
	Number    string    `firestore:"number"`
	Timestamp time.Time `firestore:"timestamp"`
}

const (
	//ちなみに、Z09:00をなくすと、Formatにしたら自動でUTCの時間になってしまう。
	RFC3339 = "2006-01-02T15:04:05Z09:00"
)

func firebaseInit(ctx context.Context) (*firestore.Client, error) {
	// Use a service account
	sa := option.WithCredentialsFile("./serviceAccounts/mahjong-linebot-a15af8e60164.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return client, nil
}

/*
*

	firestoreのcurrentデータを更新。
	（東風戦、半荘戦）、（三麻、四麻）を試合前に登録しておく。

*
*/
func AddGameStatusData(text, param string, time time.Time) error {
	ctx := context.Background()
	client, err := firebaseInit(ctx)
	if err != nil {
		log.Fatal(err)
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
		log.Fatalf("Failed adding alovelace: %v", err)
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
		log.Fatal(err)
	}
	// 切断
	defer client.Close()

	dsnap, err := client.Collection("gameStatus").Doc("current").Get(ctx)
	if err != nil {
		log.Fatalf("Failed getting currentStatus: %v", err)
		return err
	}
	m := dsnap.Data()
	_, err = client.Collection("ranks").Doc(time.Format(RFC3339)[0:19]).Set(ctx, GameResult{
		Rank:      text,
		Game:      m["game"].(string),
		Number:    m["number"].(string),
		Timestamp: time,
	})
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
		return err
	}

	// エラーなしは成功
	return err
}
