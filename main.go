package main

import (
	"mahjong-linebot/app/controllers"
	"mahjong-linebot/config"
	"mahjong-linebot/utils"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	controllers.StartWebServer()
}
