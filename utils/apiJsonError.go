package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// jsonでレスポンスしたいならこれ使う。(クライアントで内容詳しく把握するときとかかな。)
// もっとプロパティ増やす必要あり？
type JSONError struct {
	Error string `json:"error"`
	Code  int    `json:"code"` //エラーコード
}

// エラーがあった際に、jsonで返すapiエラー自作
//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//http.Errorはstringになるので、json化
func APIError(w http.ResponseWriter, errMessage string, code int) {
	w.Header().Set("Content-Type", "application/json") //レスポンスヘッダ
	w.WriteHeader(code)                                //エラーコード
	jsonError, err := json.Marshal(JSONError{Error: errMessage, Code: code})
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonError) //jsonをreturn
}

