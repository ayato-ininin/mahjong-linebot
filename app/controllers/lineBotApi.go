package controllers

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
	logger "mahjong-linebot/utils"
	"net/http"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)


func lineBotApiPost(w http.ResponseWriter, req *http.Request) {
	log.Printf(logger.InfoLogEntry("[/v1/api/linebot] START ==========="))
	jst := time.FixedZone("JST", 9*60*60)
	bot, err := linebot.New(
		config.Config.ChannelSecret, //channel secret
		config.Config.AccessToken,   //access token
	)
	if err != nil {
		http.Error(w, "Error init client", http.StatusBadRequest)
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("linebotの認証に失敗 err=%v", err)))
	}

	events, err := bot.ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			APIError(w, "Invalid signature:", http.StatusBadRequest)
			log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Invalid signature: err=%v", err)))
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Internal server error: err=%v", err)))
		}
		return
	}

	for _, event := range events {
		//イベントがメッセージの受信だった場合
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				switch message.Text {
				case "東風戦", "半荘戦":
					log.Printf(logger.InfoLogEntry("[Update gameStatus:game] start"))
					err := firestore.UpdateGameStatusData(message.Text, "game", time.Now().In(jst))
					//リプライを返さないと何度も再送される（と思われる）ので返信
					if err != nil {
						log.Printf(logger.InfoLogEntry("[/lineHandler] END ==========="))
					} else {
						_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("登録")).Do()
						if err != nil {
							log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Failed replyMessage: err=%v", err)))
						} else {
							log.Printf(logger.InfoLogEntry("[/lineHandler] END ==========="))
						}
					}
					break
				case "三麻", "四麻":
					log.Printf(logger.InfoLogEntry("[Update gameStatus:number] start"))
					err := firestore.UpdateGameStatusData(message.Text, "number", time.Now().In(jst))
					//リプライを返さないと何度も再送される（と思われる）ので返信
					if err != nil {
						log.Printf(logger.InfoLogEntry("[/lineHandler] END ==========="))
					} else {
						_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("登録")).Do()
						if err != nil {
							log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Failed replyMessage: err=%v", err)))
						} else {
							log.Printf(logger.InfoLogEntry("[/lineHandler] END ==========="))
						}
					}
					break
				case "リアル", "ネット":
					log.Printf(logger.InfoLogEntry("[Update gameStatus:style] start"))
					err := firestore.UpdateGameStatusData(message.Text, "style", time.Now().In(jst))
					//リプライを返さないと何度も再送される（と思われる）ので返信
					if err != nil {
						log.Printf(logger.InfoLogEntry("[/lineHandler] END ==========="))
					} else {
						_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("登録")).Do()
						if err != nil {
							log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Failed replyMessage: err=%v", err)))
						} else {
							log.Printf(logger.InfoLogEntry("[/lineHandler] END ==========="))
						}
					}
					break
				case "1", "2", "3", "4":
					log.Printf(logger.InfoLogEntry("[Add rankData] start"))
					err := firestore.AddRankData(message.Text, time.Now().In(jst))
					//リプライを返さないと何度も再送される（と思われる）ので返信
					if err != nil {
						log.Printf(logger.InfoLogEntry("[/lineHandler] END ==========="))
					} else {
						_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("登録")).Do()
						if err != nil {
							log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Failed replyMessage: err=%v", err)))
						} else {
							log.Printf(logger.InfoLogEntry("[/lineHandler] END ==========="))
						}
					}
					break
				case "設定":
					log.Printf(logger.InfoLogEntry("[Get current setting] start"))
					status, err := firestore.GetCurrentStatus()
					//リプライを返さないと何度も再送される（と思われる）ので返信
					if err != nil {
						_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("設定を取得できませんでした。")).Do()
						if err != nil {
							log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Failed replyMessage: err=%v", err)))
						} else {
							log.Printf(logger.InfoLogEntry("[/lineHandler] END ==========="))
						}
					}
					if status == nil {
						_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("設定がありません。")).Do()
						if err != nil {
							log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Failed replyMessage: err=%v", err)))
						} else {
							log.Printf(logger.InfoLogEntry("[/lineHandler] END ==========="))
						}
					} else {
						msg := firestore.CreateStatusMsg(status)
						_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(msg)).Do()
						if err != nil {
							log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Failed replyMessage: err=%v", err)))
						} else {
							log.Printf(logger.InfoLogEntry("[/lineHandler] END ==========="))
						}
					}
					break
				default:
					// https://qiita.com/seihmd/items/4a878e7fa340d7963fee
					str := string([]rune(message.Text)[:2])
					if str == "場所" {
						log.Printf(logger.InfoLogEntry("[Update gameStatus:place] start"))
						err := firestore.UpdateGameStatusData(string([]rune(message.Text)[2:]), "place", time.Now().In(jst))
						//リプライを返さないと何度も再送される（と思われる）ので返信
						if err != nil {
							log.Printf(logger.InfoLogEntry("[/lineHandler] END ==========="))
						} else {
							_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("登録")).Do()
							if err != nil {
								log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Failed replyMessage: err=%v", err)))
							} else {
								log.Printf(logger.InfoLogEntry("[/lineHandler] END ==========="))
							}
						}
					} else {
						//リプライを返さないと何度も再送される（と思われる）ので返信
						_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("登録できません")).Do()
						if err != nil {
							log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Failed replyMessage: err=%v", err)))
						} else {
							log.Printf(logger.InfoLogEntry("[/lineHandler] END ==========="))
						}
					}
				}
			}
		}
	}
}
