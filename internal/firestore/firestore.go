package firestore

import (
	"context"
	"log"
	logger "mahjong-linebot/internal/logs"
	"mahjong-linebot/internal/utils"

	"cloud.google.com/go/firestore"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// firebaseのサービスアカウントキーを取得、クライアントの初期化
func firebaseInit(ctx context.Context) (*firestore.Client, error) {
	traceId, err := utils.GetTraceID(ctx)
	if err != nil {
		return nil, err
	}
	// Firebaseのサービスアカウントキーの取得
	jsonBytes := getFirebaseServiceAccountKey()
	sa := option.WithCredentialsJSON(jsonBytes)

	// Firebaseアプリケーションの初期化
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Firebaseアプリケーションの初期化に失敗しました", err))
		return nil, err
	}

	// Firestoreクライアントの初期化
	fsClient, err := app.Firestore(ctx)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Firestoreクライアントの初期化に失敗しました", err))
		return nil, err
	}

	return fsClient, nil
}
