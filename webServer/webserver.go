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
	"strconv"
)

func StartWebServer() error {
	http.HandleFunc("/v1/api/linebot", lineBotApiHandler)
	http.HandleFunc("/v1/api/matchSetting", func(w http.ResponseWriter, r *http.Request) {
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
	http.HandleFunc("/v1/api/matchResult", matchResultApiHandler)
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

func matchResultApiHandler(w http.ResponseWriter, r *http.Request) {
	traceId := logger.GetTraceId(r)
	log.Printf(logger.InfoLogEntry(traceId, "MATCHRESULT START ==========="))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	switch r.Method {
	case http.MethodGet:
		// ルームIDの取得
		v := r.URL.Query().Get("roomid")
		if v == "" {
			utils.APIError(w, "Don't exist query", http.StatusBadRequest)
			return
		}
		roomid, err := strconv.Atoi(v)
		if err != nil {
			//クエリパラメータが数字でない
			log.Printf(logger.ErrorLogEntry(traceId, "Not valid query: required number", err))
			utils.APIError(w, "Not valid query: required number", http.StatusBadRequest)
			return
		}
		controllers.GetMatchResultByRoomId(w, r, roomid)
		return
	case http.MethodPost:
		controllers.MatchResultPost(w, r)
	case http.MethodPut:
		controllers.MatchResultUpdate(w, r)
	case http.MethodDelete:
		utils.APIError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	case http.MethodOptions:
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Content-Typeヘッダの使用を許可する
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST")    // pre-flightリクエストに対応する
		//これプリフライトして一回目のレスポンス何もないから、クライアント側一回目失敗するかも。
		w.WriteHeader(http.StatusOK)
		return
	default:
		utils.APIError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	log.Printf(logger.InfoLogEntry(traceId, "MATCHRESULT END ==========="))
}
