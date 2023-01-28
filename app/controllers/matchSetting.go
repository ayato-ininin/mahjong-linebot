package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mahjong-linebot/firestore"
	"mahjong-linebot/app/models"
	logger "mahjong-linebot/utils"
	"net/http"
	"time"
)

func matchSettingPost(w http.ResponseWriter, r *http.Request) {
	// 	//データ保存処理
	traceId := logger.GetTraceId(r)
	log.Printf(logger.InfoLogEntryTest("[/v1/api/matchSetting: POST] START ===========",traceId))
	//JSONから構造体へ
	body, _ := io.ReadAll(r.Body)
	m := new(models.MatchSetting) //構造体
	err := json.Unmarshal(body, m)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("[/v1/api/matchSetting: POST] Failed json unmarshal: err=%v", err)))
		APIError(w, "Failed json unmarshal:", http.StatusInternalServerError)
		return
	}

	jst := time.FixedZone("JST", 9*60*60)
	err = firestore.AddMatchSetting(m, time.Now().In(jst))
	//リプライを返さないと何度も再送される（と思われる）ので返信
	if err != nil {
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("[/v1/api/matchSetting: POST] Failed addMatchSetting: err=%v", err)))
		APIError(w, "Failed addMatchSetting:", http.StatusInternalServerError)
		return
	}
	res, _ := json.Marshal(m) //json化
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
	log.Printf(logger.InfoLogEntry("[/v1/api/matchSetting: POST] END ==========="))
}
