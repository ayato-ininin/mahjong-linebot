package logs

import (
	"fmt"
	logpb "google.golang.org/genproto/googleapis/logging/v2"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// log出力設定
func LoggingSettings(logFile string) {
	//O_RDWR : ファイルの読み込みと書き込み両方
	//O_CREATE: ファイルがなければ作成
	//O_APPEND: 上書きではなく追記する
	//0666: ファイルのパーミッション、「rw-rw-rw-」なら「0666」、自分-グループ-他人
	//r:4(読み)、w:2(書き)、x:1(実行)
	logfile, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("file=logFile err=%s", err.Error())
	}
	//Stdout→standard output(標準出力:コンソールに出る！！これがないと、ログファイルにしかかかれない。multiwriterは標準出力とログファイルの両方に書き込むという設定)
	multiLogFile := io.MultiWriter(os.Stdout, logfile)
	//出力時の情報を追加で付加したい場合.
	//定数はビットフラグで定義されているので、| でまとめて設定できます：
	log.SetFlags(0)
	log.SetOutput(multiLogFile) //出力先
}

// INFOレベルのログ出力
func InfoLogEntry(traceId string, message string, args ...interface{}) string {
	pt, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	rFuncName := regexp.MustCompile("^.*/")
	funcName := rFuncName.ReplaceAllString(runtime.FuncForPC(pt).Name(), "")
	msg := fmt.Sprintf("["+path.Base(file)+":"+strconv.Itoa(line)+":"+funcName+"] - "+message, args...)
	entry := &LogEntryTest{
		Severity: INFO,
		Message:  msg,
		Trace:    traceId,
		SourceLocation: &logpb.LogEntrySourceLocation{
			File:     file,
			Line:     int64(line),
			Function: runtime.FuncForPC(pt).Name(),
		},
	}

	return entry.String()
}

// WARNレベルのログ出力
func WarnLogEntry(traceId string, message string, args ...interface{}) string {
	pt, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	rFuncName := regexp.MustCompile("^.*/")
	funcName := rFuncName.ReplaceAllString(runtime.FuncForPC(pt).Name(), "")
	msg := fmt.Sprintf("["+path.Base(file)+":"+strconv.Itoa(line)+":"+funcName+"] - "+message, args...)
	entry := &LogEntryTest{
		Severity: WARN,
		Message:  msg,
		Trace:    traceId,
		SourceLocation: &logpb.LogEntrySourceLocation{
			File:     file,
			Line:     int64(line),
			Function: runtime.FuncForPC(pt).Name(),
		},
	}

	return entry.String()
}

// ERRORレベルのログ出力
func ErrorLogEntry(traceId string, message string, args ...interface{}) string {
	pt, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	rFuncName := regexp.MustCompile("^.*/")
	funcName := rFuncName.ReplaceAllString(runtime.FuncForPC(pt).Name(), "")
	msg := fmt.Sprintf("["+path.Base(file)+":"+strconv.Itoa(line)+":"+funcName+"] - "+message, args...)
	entry := &LogEntryTest{
		Severity: ERROR,
		Message:  msg,
		Trace:    traceId,
		SourceLocation: &logpb.LogEntrySourceLocation{
			File:     file,
			Line:     int64(line),
			Function: runtime.FuncForPC(pt).Name(),
		},
	}

	return entry.String()
}

// http "X-Cloud-Trace-Context" headerからtraceIdを抜き出す
func GetTraceId(r *http.Request) string {
	traceHeader := r.Header.Get("X-Cloud-Trace-Context")
	traceParts := strings.Split(traceHeader, "/")
	traceId := ""
	if len(traceParts) > 0 {
		traceId = traceParts[0]
	}
	return traceId
}
