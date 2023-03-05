package controllers

/*
【参考文献】
https://speakerdeck.com/yagieng/go-to-line-bot-ni-matometeru-men-suru
https://github.com/line/line-bot-sdk-go
*/

import (
	"context"
	"log"
	"mahjong-linebot/internal/config"
	"mahjong-linebot/internal/firestore"
	logger "mahjong-linebot/internal/logs"
	"mahjong-linebot/internal/utils"
	"net/http"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func LineBotApiPost(w http.ResponseWriter, r *http.Request) {
	traceId := logger.GetTraceId(r)
	log.Printf(logger.InfoLogEntry(traceId, "LINEBOT API START ==========="))

	// contextを作成しtraceIdをセットする(リクエストを渡すのではなく、contextにしてfirestoreに渡す。traceIdにて追跡のため)
	ctx := context.WithValue(context.Background(), "traceId", traceId)
	h := r.Header["User-Agent"][0]
	//正しいlinebotからリクエストが送られない場合にerrorを返す
	if h != "LineBotWebhook/2.0" {
		utils.APIError(w, "Error client agent", http.StatusBadRequest)
		log.Printf(logger.ErrorLogEntry(traceId, "Error client agent "+r.Header["User-Agent"][0]))
		return
	}

	jst := time.FixedZone("JST", 9*60*60)
	config, err := config.InitConfig()
	if err != nil {
		utils.APIError(w, "Error init config", http.StatusBadRequest)
		log.Printf(logger.ErrorLogEntry(traceId, "configの初期化に失敗", err))
	}
	bot, err := linebot.New(
		config.ChannelSecret, //channel secret
		config.AccessToken,   //access token
	)
	if err != nil {
		utils.APIError(w, "Error init client", http.StatusBadRequest)
		log.Printf(logger.ErrorLogEntry(traceId, "linebotの認証に失敗", err))
	}

	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			utils.APIError(w, "Invalid signature", http.StatusBadRequest)
			log.Printf(logger.ErrorLogEntry(traceId, "Invalid signature", err))
		} else {
			utils.APIError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Printf(logger.ErrorLogEntry(traceId, "Internal server error", err))
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
					log.Printf(logger.InfoLogEntry(traceId, "[Update gameStatus:game] start"))
					err := firestore.UpdateGameStatusData(ctx, message.Text, "game", time.Now().In(jst))
					//リプライを返さないと何度も再送される（と思われる）ので返信
					if err == nil {
						_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("登録")).Do()
						if err != nil {
							log.Printf(logger.ErrorLogEntry(traceId, "Failed replyMessage", err))
						}
					}
					break
				case "三麻", "四麻":
					log.Printf(logger.InfoLogEntry(traceId, "[Update gameStatus:number] start"))
					err := firestore.UpdateGameStatusData(ctx, message.Text, "number", time.Now().In(jst))
					//リプライを返さないと何度も再送される（と思われる）ので返信
					if err == nil {
						_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("登録")).Do()
						if err != nil {
							log.Printf(logger.ErrorLogEntry(traceId, "Failed replyMessage", err))
						}
					}
					break
				case "リアル", "ネット":
					log.Printf(logger.InfoLogEntry(traceId, "[Update gameStatus:style] start"))
					err := firestore.UpdateGameStatusData(ctx, message.Text, "style", time.Now().In(jst))
					//リプライを返さないと何度も再送される（と思われる）ので返信
					if err == nil {
						_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("登録")).Do()
						if err != nil {
							log.Printf(logger.ErrorLogEntry(traceId, "Failed replyMessage", err))
						}
					}
					break
				case "1", "2", "3", "4":
					log.Printf(logger.InfoLogEntry(traceId, "[Add rankData] start"))
					err := firestore.AddRankData(ctx, message.Text, time.Now().In(jst))
					//リプライを返さないと何度も再送される（と思われる）ので返信
					if err == nil {
						_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("登録")).Do()
						if err != nil {
							log.Printf(logger.ErrorLogEntry(traceId, "Failed replyMessage", err))
						}
					}
					break
				case "設定":
					log.Printf(logger.InfoLogEntry(traceId, "[Get current setting] start"))
					status, err := firestore.GetCurrentStatus(ctx)
					//リプライを返さないと何度も再送される（と思われる）ので返信
					if err != nil {
						_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("設定を取得できませんでした。")).Do()
						if err != nil {
							log.Printf(logger.ErrorLogEntry(traceId, "Failed replyMessage", err))
						}
					}
					if status == nil {
						_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("設定がありません。")).Do()
						if err != nil {
							log.Printf(logger.ErrorLogEntry(traceId, "Failed replyMessage", err))
						}
					} else {
						msg := firestore.CreateStatusMsg(status)
						_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(msg)).Do()
						if err != nil {
							log.Printf(logger.ErrorLogEntry(traceId, "Failed replyMessage", err))
						}
					}
					break
				default:
					// https://qiita.com/seihmd/items/4a878e7fa340d7963fee
					str := string([]rune(message.Text)[:2])
					if str == "場所" {
						log.Printf(logger.InfoLogEntry(traceId, "[Update gameStatus:place] start"))
						err := firestore.UpdateGameStatusData(ctx, string([]rune(message.Text)[2:]), "place", time.Now().In(jst))
						//リプライを返さないと何度も再送される（と思われる）ので返信
						if err == nil {
							_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("登録")).Do()
							if err != nil {
								log.Printf(logger.ErrorLogEntry(traceId, "Failed replyMessage", err))
							}
						}
					} else {
						//リプライを返さないと何度も再送される（と思われる）ので返信
						_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("登録できません")).Do()
						if err != nil {
							log.Printf(logger.ErrorLogEntry(traceId, "Failed replyMessage", err))
						}
					}
				}
			}
			log.Printf(logger.InfoLogEntry(traceId, "LINEBOT API END ==========="))
		}
	}
}
