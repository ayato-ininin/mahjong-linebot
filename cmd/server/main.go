package main

//cmdはアプリケーションのエントリーポイントを持つ
import (
	logger "mahjong-linebot/pkg/logs"
	"mahjong-linebot/pkg/router"
)

func main() {
	logger.LoggingSettings("mahjong_linebot.log")
	router.StartWebServer()
}
