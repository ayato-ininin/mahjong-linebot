package main

/*
【参考文献】
https://speakerdeck.com/yagieng/go-to-line-bot-ni-matometeru-men-suru
https://github.com/line/line-bot-sdk-go
*/

import (
	"context"
	"fmt"
	"log"
	"mahjong-linebot/config"
	"mahjong-linebot/utils"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"

	firebase "firebase.google.com/go"
  "google.golang.org/api/option"
)

func handler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func lineHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "d world")
	bot, err := linebot.New(
		config.Config.ChannelSecret, //channel secret
		config.Config.AccessToken,   //access token
	)
	if err != nil {
		http.Error(w, "Error init client", http.StatusBadRequest)
		log.Fatal(err)
	}

	events, err := bot.ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
			log.Fatal(err)
		} else {
			w.WriteHeader(500)
			log.Fatal(err)
		}
		return
	}

	for _, event := range events {
		//イベントがメッセージの受信だった場合
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				replyMessage := message.Text
				_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()
				if err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func main() {
	ctx := context.Background()
	sa := option.WithCredentialsFile("./serviceAccounts/mahjong-linebot-a15af8e60164.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
			log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	_, _, err = client.Collection("users").Add(ctx, map[string]interface{}{
		"first": "Ada",
		"last":  "Lovelace",
		"born":  1815,
	})
	if err != nil {
			log.Fatalf("Failed adding alovelace: %v", err)
	}

	defer client.Close()

	utils.LoggingSettings(config.Config.LogFile)
	http.HandleFunc("/", handler)
	http.HandleFunc("/callback", lineHandler)
	fmt.Println("起動中...")
	http.ListenAndServe(":8080", nil)
}
