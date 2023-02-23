package controllers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"mahjong-linebot/app/firestore"
	"mahjong-linebot/app/models"
	logger "mahjong-linebot/logs"
	"mahjong-linebot/utils"
	"net/http"
	"strconv"
	"time"
)

func GetMatchSettingByRoomId(w http.ResponseWriter, r *http.Request) {
	traceId := logger.GetTraceId(r)
	log.Printf(logger.InfoLogEntry(traceId, "GET:MATCHSETTING START ==========="))

	roomid, err := strconv.Atoi(r.URL.Query().Get("roomid"))
	if err != nil {
		//クエリパラメータが数字でないか空文字
		log.Printf(logger.ErrorLogEntry(traceId, "Not valid query: required number", err))
		utils.APIError(w, "Not valid query: required number", http.StatusBadRequest)
		return
	}

	//ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	//https://www.wakuwakubank.com/posts/867-go-context/→contextの種類
	ctx := context.WithValue(context.Background(), "traceId", traceId)//withvalueはtraceIdを保つために必要
	log.Printf(logger.InfoLogEntry(traceId, "取得部屋番号 : "+ strconv.Itoa(roomid)))
	m, err := firestore.GetMatchSetting(ctx, roomid)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed getMatchSetting", err))
		utils.APIError(w, "Failed getMatchSetting", http.StatusInternalServerError)//ここは500でいいのか？サーバのエラーでもないかも
		return
	}

	// 構造体をJSONに変換する(log出力用)
	jsonData, err := json.Marshal(&m)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed json marshal", err))
		utils.APIError(w, "Failed json marshal", http.StatusInternalServerError)
	}
	log.Printf(logger.InfoLogEntry(traceId, "検索されたデータ : %s", jsonData))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
	log.Printf(logger.InfoLogEntry(traceId, "GET:MATCHSETTING END ==========="))
}

func PostMatchSetting(w http.ResponseWriter, r *http.Request) {
	traceId := logger.GetTraceId(r)
	log.Printf(logger.InfoLogEntry(traceId, "POST:MATCHSETTING START ==========="))
	// contextを作成しtraceIdをセットする(リクエストを渡すのではなく、contextにしてfirestoreに渡す。traceIdにて追跡のため)
	ctx := context.WithValue(context.Background(), "traceId", traceId)
	//JSONから構造体へ
	body, err := io.ReadAll(r.Body)//読み切りが必要なのでio.ReadAllを使う(コネクションの再利用)
	if err != nil {
		utils.APIError(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()//TCPコネクションを閉じて、ファイルディスクリプタの枯渇を防ぐ

	var m models.MatchSetting //構造体
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed json unmarshal", err))
		utils.APIError(w, "Failed json unmarshal", http.StatusInternalServerError)
		return
	}

	jst := time.FixedZone("JST", 9*60*60)
	err = firestore.AddMatchSetting(ctx, &m, time.Now().In(jst))
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed postMatchSetting", err))
		utils.APIError(w, "Failed postMatchSetting", http.StatusInternalServerError)
		return
	}

	// 構造体をJSONに変換する(log出力用)
	jsonData, err := json.Marshal(&m)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed json marshal", err))
		utils.APIError(w, "Failed json marshal", http.StatusInternalServerError)
	}
	log.Printf(logger.InfoLogEntry(traceId, "追加データ : %s", jsonData))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
	log.Printf(logger.InfoLogEntry(traceId, "POST:MATCHSETTING END ==========="))
}

func OptionsMatchSettingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
}
