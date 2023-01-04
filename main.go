package main

import (
	"mahjong-linebot/app/controllers"
	"mahjong-linebot/utils"
)

func main() {
	utils.LoggingSettings("mahjong_linebot.log")
	controllers.StartWebServer()
}
