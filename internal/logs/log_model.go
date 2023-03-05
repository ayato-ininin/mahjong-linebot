package logs

import (
	"encoding/json"
	logpb "google.golang.org/genproto/googleapis/logging/v2"
	"log"
)

// ログレベルのCONSTを定義
// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#logseverity
var (
	INFO  = "INFO"
	WARN  = "WARNING"
	ERROR = "ERROR"
)

// GCPのLogEntryに則った構造化ログモデル
type LogEntryTest struct {
	// GCP上でLogLevelを表す
	Severity string `json:"severity"`
	// ログの内容
	Message        string                        `json:"message"`
	SourceLocation *logpb.LogEntrySourceLocation `json:"sourceLocation"`
	Trace          string                        `json:"logging.googleapis.com/trace,omitempty"`
}

// 構造体をJSON形式の文字列へ変換
// 参考: https://cloud.google.com/run/docs/logging#run_manual_logging-go
func (l LogEntryTest) String() string {
	if l.Severity == "" {
		l.Severity = INFO
	}
	out, err := json.Marshal(l)
	if err != nil {
		log.Printf("json.Marshal: %v", err)
	}
	return string(out)
}
