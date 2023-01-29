package main

import (
	"mahjong-linebot/app/controllers"
	"mahjong-linebot/logs"
)

func main() {
	logs.LoggingSettings("mahjong_linebot.log")
	controllers.StartWebServer()
}
