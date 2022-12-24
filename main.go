package main

/*
【参考文献】
https://speakerdeck.com/yagieng/go-to-line-bot-ni-matometeru-men-suru
https://github.com/line/line-bot-sdk-go
*/

import (
	"fmt"
	"log"
	"mahjong-linebot/config"
	"mahjong-linebot/firestore"
	"mahjong-linebot/utils"
	"net/http"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
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
				switch (message.Text){
					case "東風戦", "半荘戦":
						log.Print("gamestatus:game register");
						err := firestore.AddGameStatusData(message.Text, "game", time.Now())
						//リプライを返さないと何度も再送される（と思われる）ので返信
						if(err == nil){
							_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("登録")).Do()
							if err != nil {
								log.Print(err)
							}
						}
						break;
					case "三麻", "四麻":
						log.Print("gamestatus:number register");
						err := firestore.AddGameStatusData(message.Text, "number", time.Now())
						//リプライを返さないと何度も再送される（と思われる）ので返信
						if(err == nil){
							_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("登録")).Do()
							if err != nil {
								log.Print(err)
							}
						}
						break;
					case "1", "2", "3", "4":
						err := firestore.AddRankData(message.Text,time.Now())
						//リプライを返さないと何度も再送される（と思われる）ので返信
						if(err == nil){
							_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("登録")).Do()
							if err != nil {
								log.Print(err)
							}
						}
						log.Print("rank register");
						break;
					default:
						//リプライを返さないと何度も再送される（と思われる）ので返信
						_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("登録できません")).Do()
						if err != nil {
							log.Fatal(err)
						}
				}
			}
		}
	}
}

func main() {
	utils.LoggingSettings(config.Config.LogFile)

	http.HandleFunc("/", handler)
	http.HandleFunc("/callback", lineHandler)
	fmt.Println("起動中...")
	http.ListenAndServe(":8080", nil)
}
