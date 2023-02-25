package webServer

/*
【参考文献】
https://speakerdeck.com/yagieng/go-to-line-bot-ni-matometeru-men-suru
https://github.com/line/line-bot-sdk-go
*/

import (
	"fmt"
	"log"
	"mahjong-linebot/app/controllers"
	logger "mahjong-linebot/logs"
	"mahjong-linebot/utils"
	"net/http"
)

func StartWebServer() error {
	http.HandleFunc("/v1/api/linebot", lineBotApiHandler)
	http.HandleFunc("/v1/api/matchSetting", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://mahjong-linebot.firebaseapp.com")
		switch r.Method {
		case http.MethodGet:
			controllers.GetMatchSettingByRoomId(w, r)
		case http.MethodPost:
			controllers.PostMatchSetting(w, r)
		case http.MethodOptions:
			controllers.OptionsMatchSettingHandler(w, r)
		default:
			utils.APIError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/v1/api/matchResult", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://mahjong-linebot.firebaseapp.com")
		switch r.Method {
		case http.MethodGet:
			controllers.GetMatchResultByRoomId(w, r)
		case http.MethodPost:
			controllers.PostMatchResult(w, r)
		case http.MethodPut:
			controllers.UpdateMatchResult(w, r)
		case http.MethodOptions:
			controllers.OptionsMatchResultHandler(w, r)
		default:
			utils.APIError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})
	// パスが一致するものがない場合は404を返す
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	log.Printf("コンテナ起動...")
	return http.ListenAndServe(fmt.Sprintf(":%d", 8080), nil)
}

func lineBotApiHandler(w http.ResponseWriter, r *http.Request) {
	traceId := logger.GetTraceId(r)
	log.Printf(logger.InfoLogEntry(traceId, "LINEBOT API START ==========="))
	switch r.Method {
	case http.MethodPost:
		//正しいlinebotからリクエストが送られない場合にerrorを返す
		h := r.Header["User-Agent"][0]
		if h != "LineBotWebhook/2.0" {
			utils.APIError(w, "Error client agent", http.StatusBadRequest)
			log.Printf(logger.ErrorLogEntry(traceId, "Error client agent "+r.Header["User-Agent"][0]))
			return
		} else {
			controllers.LineBotApiPost(w, r)
		}
	default:
		utils.APIError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		log.Printf(logger.ErrorLogEntry(traceId, "Error bad method "+r.Method))
		return
	}
	log.Printf(logger.InfoLogEntry(traceId, "LINEBOT API END ==========="))
}
