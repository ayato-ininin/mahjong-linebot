package main

import (
	"mahjong-linebot/logs"
	"mahjong-linebot/webServer"
)

func main() {
	logs.LoggingSettings("mahjong_linebot.log")
	webServer.StartWebServer()
}
