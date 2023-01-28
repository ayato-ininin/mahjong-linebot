package controllers

/*
【参考文献】
https://speakerdeck.com/yagieng/go-to-line-bot-ni-matometeru-men-suru
https://github.com/line/line-bot-sdk-go
*/

import (
	"encoding/json"
	"fmt"
	"log"
	logger "mahjong-linebot/utils"
	"net/http"
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

func StartWebServer() error {
	http.HandleFunc("/", handler)
	http.HandleFunc("/v1/api/linebot", lineBotApiHandler)
	http.HandleFunc("/v1/api/matchSetting", matchSettingApiHandler)
	log.Printf(logger.InfoLogEntry("コンテナ起動..."))
	return http.ListenAndServe(fmt.Sprintf(":%d", 8080), nil)
}

func handler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world")
	log.Printf(logger.InfoLogEntry("Hello world"))
}

func lineBotApiHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
			//正しいlinebotからリクエストが送られない場合にerrorを返す
		h := r.Header["User-Agent"][0]
		if h != "LineBotWebhook/2.0" {
			http.Error(w, "Error client agent ", http.StatusBadRequest)
			log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Error client agent " + r.Header["User-Agent"][0])))
			return
		} else {
			lineBotApiPost(w, r)
		}
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Error bad method " + r.Method)))
		return
	}
}

// ハンドラーのラップ(section13で解説されている。)
func matchSettingApiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	switch r.Method {
	case http.MethodGet:
		APIError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	case http.MethodPost:
		matchSettingPost(w, r)
	case http.MethodDelete:
		APIError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	case http.MethodOptions:
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")               // Content-Typeヘッダの使用を許可する
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS") // pre-flightリクエストに対応する
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}
