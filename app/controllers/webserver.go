package controllers

/*
【参考文献】
https://speakerdeck.com/yagieng/go-to-line-bot-ni-matometeru-men-suru
https://github.com/line/line-bot-sdk-go
*/

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mahjong-linebot/config"
	"mahjong-linebot/firestore"
	"mahjong-linebot/app/models"
	logger "mahjong-linebot/utils"
	"net/http"
	"regexp"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// jsonでレスポンスしたいならこれ使う。(クライアントで内容詳しく把握するときとかかな。)
// もっとプロパティ増やす必要あり？
// http.Errorはstringになる
type JSONError struct {
	Error string `json:"error"`
	Code  int    `json:"code"` //エラーコード
}

// jsonになにかあったときに、jsonで返すapiエラー自作
func APIError(w http.ResponseWriter, errMessage string, code int) {
	w.Header().Set("Content-Type", "application/json") //レスポンスヘッダ
	w.WriteHeader(code)                                //エラーコード
	jsonError, err := json.Marshal(JSONError{Error: errMessage, Code: code})
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonError) //jsonをreturn
}

func handler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world")
	log.Printf(logger.InfoLogEntry("Hello world"))
}

func lineHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf(logger.InfoLogEntry("[/lineHandler] START ==========="))
	jst := time.FixedZone("JST", 9*60*60)
	bot, err := linebot.New(
		config.Config.ChannelSecret, //channel secret
		config.Config.AccessToken,   //access token
	)
	if err != nil {
		http.Error(w, "Error init client", http.StatusBadRequest)
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("linebotの認証に失敗 err=%v", err)))
	}

	//POSTでない場合にエラー
	if req.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		log.Print("Error bad method " + req.Method)
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Error bad method " + req.Method)))
	}

	//正しいlinebotからリクエストが送られない場合にerrorを返す
	h := req.Header["User-Agent"][0]
	if h != "LineBotWebhook/2.0" {
		http.Error(w, "Error client agent ", http.StatusBadRequest)
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Error client agent " + req.Header["User-Agent"][0])))
	}

	events, err := bot.ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			//w.WriteHeader(400)
			APIError(w, "Invalid signature:", http.StatusBadRequest)
			log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Invalid signature: err=%v", err)))
		} else {
			//w.WriteHeader(500)
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

// pathチェック用
var apiValidPath = regexp.MustCompile("^/api/setting/$")

// ハンドラーのラップ(section13で解説されている。)
func apiMakeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		switch r.Method {
		case http.MethodGet:
			APIError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		case http.MethodPost:
			fn(w, r)
		case http.MethodDelete:
			APIError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")               // Content-Typeヘッダの使用を許可する
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS") // pre-flightリクエストに対応する
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}

func apiPostHandler(w http.ResponseWriter, r *http.Request) {
	// 	//データ保存処理
	log.Printf(logger.InfoLogEntry("[/api/matchSetting] START ==========="))
	//JSONから構造体へ
	body, _ := io.ReadAll(r.Body)
	m := new(models.MatchSetting) //構造体
	err := json.Unmarshal(body, m)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("[/api/matchSetting] Failed json unmarshal: err=%v", err)))
		APIError(w, "Failed json unmarshal:", http.StatusInternalServerError)
		return
	}

	jst := time.FixedZone("JST", 9*60*60)
	err = firestore.AddMatchSetting(m, time.Now().In(jst))
	//リプライを返さないと何度も再送される（と思われる）ので返信
	if err != nil {
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("[/api/matchSetting] Failed addMatchSetting: err=%v", err)))
		APIError(w, "Failed addMatchSetting:", http.StatusInternalServerError)
		return
	}
	res, _ := json.Marshal(m) //json化
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
	log.Printf(logger.InfoLogEntry("[/api/matchSetting] END ==========="))
}

func StartWebServer() error {
	http.HandleFunc("/", handler)
	http.HandleFunc("/lineCallback", lineHandler)
	http.HandleFunc("/api/matchSetting", apiMakeHandler(apiPostHandler))
	log.Printf(logger.InfoLogEntry("コンテナ起動..."))
	return http.ListenAndServe(fmt.Sprintf(":%d", 8080), nil)
}
