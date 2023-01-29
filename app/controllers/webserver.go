package controllers

/*
【参考文献】
https://speakerdeck.com/yagieng/go-to-line-bot-ni-matometeru-men-suru
https://github.com/line/line-bot-sdk-go
*/

import (
	"fmt"
	"log"
	logger "mahjong-linebot/logs"
	"mahjong-linebot/utils"
	"net/http"
)

func StartWebServer() error {
	http.HandleFunc("/home", handler)
	http.HandleFunc("/v1/api/linebot", lineBotApiHandler)
	http.HandleFunc("/v1/api/matchSetting", matchSettingApiHandler)
	log.Printf("コンテナ起動...")
	return http.ListenAndServe(fmt.Sprintf(":%d", 8080), nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	//http.Redirect(w, r, "/", http.StatusFound)リダイレクト
	fmt.Fprintf(w, "hello world")
	log.Printf("Hello world")
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
			lineBotApiPost(w, r)
		}
	default:
		utils.APIError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		log.Printf(logger.ErrorLogEntry(traceId, "Error bad method "+r.Method))
		return
	}
	log.Printf(logger.InfoLogEntry(traceId, "LINEBOT API END ==========="))
}

func matchSettingApiHandler(w http.ResponseWriter, r *http.Request) {
	traceId := logger.GetTraceId(r)
	log.Printf(logger.InfoLogEntry(traceId, "MATCHSETTING START ==========="))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	switch r.Method {
	case http.MethodGet:
		utils.APIError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	case http.MethodPost:
		matchSettingPost(w, r)
	case http.MethodDelete:
		utils.APIError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	case http.MethodOptions:
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Content-Typeヘッダの使用を許可する
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")    // pre-flightリクエストに対応する
		//これプリフライトして一回目のレスポンス何もないから、クライアント側一回目失敗するかも。
		w.WriteHeader(http.StatusOK)
		return
	default:
		utils.APIError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	log.Printf(logger.InfoLogEntry(traceId, "MATCHSETTING END ==========="))
}
