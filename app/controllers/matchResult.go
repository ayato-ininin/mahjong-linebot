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
	"time"
)

func MatchResultPost(w http.ResponseWriter, r *http.Request) {
	// 	//データ保存処理
	traceId := logger.GetTraceId(r)
	// contextを作成しtraceIdをセットする(リクエストを渡すのではなく、contextにしてfirestoreに渡す。traceIdにて追跡のため)
	ctx := context.WithValue(context.Background(), "traceId", traceId)
	//JSONから構造体へ
	body, _ := io.ReadAll(r.Body)
	m := new(models.MatchResult) //構造体
	err := json.Unmarshal(body, m)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed json unmarshal", err))
		utils.APIError(w, "Failed json unmarshal", http.StatusInternalServerError)
		return
	}

	jst := time.FixedZone("JST", 9*60*60)
	err = firestore.AddMatchResult(ctx, m, time.Now().In(jst))
	//リプライを返さないと何度も再送される（と思われる）ので返信
	if err != nil {
		log.Printf(logger.ErrorLogEntry(traceId, "Failed addMatchSetting", err))
		utils.APIError(w, "Failed addMatchSetting", http.StatusInternalServerError)
		return
	}
	res, _ := json.Marshal(m) //json化
	log.Printf(logger.InfoLogEntry(traceId, "追加データ : "+string(res)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}
