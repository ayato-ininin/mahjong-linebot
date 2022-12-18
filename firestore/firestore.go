package firestore

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
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

func AddData() error {
	ctx := context.Background()
	client, err := firebaseInit(ctx)
	if err != nil {
			log.Fatal(err)
	}
	_, _, err = client.Collection("users").Add(ctx, map[string]interface{}{
		"first": "Ada",
		"last":  "Lovelace",
		"born":  1815,
	})
	if err != nil {
			log.Fatalf("Failed adding alovelace: %v", err)
	}

	// 切断
	defer client.Close()

	// エラーなしは成功
	return err
}
