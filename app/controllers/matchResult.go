package controllers

import (
	"context"
	"encoding/json"
	"log"
	"mahjong-linebot/app/firestore"
	"mahjong-linebot/app/models"
	logger "mahjong-linebot/logs"
	"mahjong-linebot/utils"
	"net/http"
	"strconv"
	"time"
)

func GetMatchResultByRoomId(w http.ResponseWriter, r *http.Request) {
	traceId := logger.GetTraceId(r)
	log.Printf(logger.InfoLogEntry(traceId, "GET:MATCHRESULT START ==========="))

	roomid, err := strconv.Atoi(r.URL.Query().Get("roomid"))
	if err != nil {
		//クエリパラメータが数字でないか空文字
		log.Printf(logger.ErrorLogEntry(traceId, "Not valid query: required number", err))
		utils.APIError(w, "Not valid query: required number", http.StatusBadRequest)
		return
	}

	// contextを作成しtraceIdをセットする(リクエストを渡すのではなく、contextにしてfirestoreに渡す。traceIdにて追跡のため)
	ctx := context.WithValue(r.Context(), "traceId", traceId) //r.Context()でリクエストのcontextを再利用
	log.Printf(logger.InfoLogEntry(traceId, "取得部屋番号 : "+strconv.Itoa(roomid)))
	m, err := firestore.GetMatchResult(ctx, roomid)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed getMatchSetting", err))
		utils.APIError(w, "Failed getMatchSetting", http.StatusInternalServerError)
		return
	}

	// 構造体をJSONに変換する(log出力用)
	jsonData, err := json.Marshal(&m)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed json marshal", err))
		utils.APIError(w, "Failed json marshal", http.StatusInternalServerError)
		return
	}
	log.Printf(logger.InfoLogEntry(traceId, "検索されたデータ(result) : %s", jsonData))

	w.Header().Set("Content-Type", "application/json")
	//エラー処理
	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed json encode", err))
		utils.APIError(w, "Failed json encode", http.StatusInternalServerError)
		return
	}
	log.Printf(logger.InfoLogEntry(traceId, "GET:MATCHRESULT END ==========="))
}

func PostMatchResult(w http.ResponseWriter, r *http.Request) {
	// 	//データ保存処理
	traceId := logger.GetTraceId(r)
	log.Printf(logger.InfoLogEntry(traceId, "POST:MATCHRESULT START ==========="))
	// contextを作成しtraceIdをセットする(リクエストを渡すのではなく、contextにしてfirestoreに渡す。traceIdにて追跡のため)
	ctx := context.WithValue(r.Context(), "traceId", traceId) //r.Context()でリクエストのcontextを再利用
	//JSONから構造体へ
	var m models.MatchResult                  //構造体
	err := json.NewDecoder(r.Body).Decode(&m) //io.readAllよりも効率的(メモリ使用量が少ない)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed json unmarshal", err))
		utils.APIError(w, "Failed json unmarshal", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close() //TCPコネクションを閉じて、ファイルディスクリプタの枯渇を防ぐ

	jst := time.FixedZone("JST", 9*60*60)
	err = firestore.AddMatchResult(ctx, &m, time.Now().In(jst))
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed addMatchSetting", err))
		utils.APIError(w, "Failed addMatchSetting", http.StatusInternalServerError)
		return
	}

	// 構造体をJSONに変換する(log出力用)
	jsonData, err := json.Marshal(&m)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed json marshal", err))
		utils.APIError(w, "Failed json marshal", http.StatusInternalServerError)
		return
	}
	log.Printf(logger.InfoLogEntry(traceId, "追加データ : %s", jsonData))

	w.Header().Set("Content-Type", "application/json")
	//エラー処理
	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed json encode", err))
		utils.APIError(w, "Failed json encode", http.StatusInternalServerError)
		return
	}
	log.Printf(logger.InfoLogEntry(traceId, "POST:MATCHRESULT END ==========="))
}

func UpdateMatchResult(w http.ResponseWriter, r *http.Request) {
	// 	//データ保存処理
	traceId := logger.GetTraceId(r)
	log.Printf(logger.InfoLogEntry(traceId, "UPDATE:MATCHRESULT START ==========="))
	// contextを作成しtraceIdをセットする(リクエストを渡すのではなく、contextにしてfirestoreに渡す。traceIdにて追跡のため)
	ctx := context.WithValue(r.Context(), "traceId", traceId) //r.Context()でリクエストのcontextを再利用
	//JSONから構造体へ
	var m models.MatchResult                  //構造体
	err := json.NewDecoder(r.Body).Decode(&m) //io.readAllよりも効率的(メモリ使用量が少ない)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed json unmarshal", err))
		utils.APIError(w, "Failed json unmarshal", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close() //TCPコネクションを閉じて、ファイルディスクリプタの枯渇を防ぐ

	jst := time.FixedZone("JST", 9*60*60)
	err = firestore.UpdateMatchResult(ctx, &m, time.Now().In(jst))
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed updateMatchSetting", err))
		utils.APIError(w, "Failed updateMatchSetting", http.StatusInternalServerError)
		return
	}

	// 構造体をJSONに変換する(log出力用)
	jsonData, err := json.Marshal(&m)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed json marshal", err))
		utils.APIError(w, "Failed json marshal", http.StatusInternalServerError)
		return
	}
	log.Printf(logger.InfoLogEntry(traceId, "追加データ : %s", jsonData))

	w.Header().Set("Content-Type", "application/json")
	//エラー処理
	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed json encode", err))
		utils.APIError(w, "Failed json encode", http.StatusInternalServerError)
		return
	}
	log.Printf(logger.InfoLogEntry(traceId, "UPDATE:MATCHRESULT END ==========="))
}

func OptionsMatchResultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
}

// func PostMatchResult(w http.ResponseWriter, r *http.Request) {
// 	post(w, r, firestore.AddMatchResult)
// }

// func UpdateMatchResult(w http.ResponseWriter, r *http.Request) {
// 	traceID := logger.GetTraceId(r)
// 	// contextを作成しtraceIdをセットする(リクエストを渡すのではなく、contextにしてfirestoreに渡す。traceIdにて追跡のため)
// 	ctx := context.WithValue(r.Context(), "traceId", traceID)//r.Context()でリクエストのcontextを再利用
// 	log.Printf(logger.InfoLogEntry(traceID, "UPDATE:MATCHRESULT START ==========="))
// 	post(ctx, w, r, traceID, firestore.UpdateMatchResult)
// 	log.Printf(logger.InfoLogEntry(traceID, "UPDATE:MATCHRESULT END ==========="))
// }

// /***
// 	firebaseへ追加、更新処理の共通処理
// 	interface{}はswtich等で型を判定しないといけないかも
//  ***/
// func post(ctx context.Context, w http.ResponseWriter, r *http.Request, traceID string, f func(context.Context, interface{}, time.Time) error) {
// 	var m interface{}
// 	err := json.NewDecoder(r.Body).Decode(&m)
// 	if err != nil {
// 			log.Printf(logger.ErrorLogEntry(traceID, "Failed json unmarshal", err))
// 			utils.APIError(w, "Failed json unmarshal", http.StatusInternalServerError)
// 			return
// 	}
// 	defer r.Body.Close()

// 	jst := time.FixedZone("JST", 9*60*60)
// 	err = f(ctx, &m, time.Now().In(jst))
// 	if err != nil {
// 			log.Printf(logger.ErrorLogEntry(traceID, "Failed", err))
// 			utils.APIError(w, "Failed", http.StatusInternalServerError)
// 			return
// 	}

// 	jsonData, err := json.Marshal(&m)
// 	if err != nil {
// 			log.Printf(logger.ErrorLogEntry(traceID, "Failed json marshal", err))
// 			utils.APIError(w, "Failed json marshal", http.StatusInternalServerError)
// 			return
// 	}
// 	log.Printf(logger.InfoLogEntry(traceID, "追加データ : %s", jsonData))

// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(m); err != nil {
// 			log.Printf(logger.ErrorLogEntry(traceID, "Failed json encode", err))
// 			utils.APIError(w, "Failed json encode", http.StatusInternalServerError)
// 			return
// 	}
// }
