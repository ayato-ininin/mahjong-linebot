package main

//cmdはアプリケーションのエントリーポイントを持つ
import (
	logger "mahjong-linebot/internal/logs"
	"mahjong-linebot/internal/router"
)

func main() {
	logger.LoggingSettings("mahjong_linebot.log")
	router.StartWebServer()
}
