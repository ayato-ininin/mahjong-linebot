package firestore

import (
	"context"
	"fmt"
	"log"
	logger "mahjong-linebot/utils"

	"cloud.google.com/go/firestore"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)


func firebaseInit(ctx context.Context) (*firestore.Client, error) {
	decJson := getFirebaseServiceAccountKey()
	sa := option.WithCredentialsJSON(decJson)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Firebase.NewApp失敗 err=%v", err)))
		return nil, err
	}

	faclient, err := app.Firestore(ctx)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("app.Firestore失敗 err=%v", err)))
		return nil, err
	}

	return faclient, nil
}
