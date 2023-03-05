package main

/*
【参考文献】
https://speakerdeck.com/yagieng/go-to-line-bot-ni-matometeru-men-suru
https://github.com/line/line-bot-sdk-go
*/

//cmdはアプリケーションのエントリーポイントを持つ
import (
	logger "mahjong-linebot/internal/logs"
	"log"
	"mahjong-linebot/internal/controllers"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	logger.LoggingSettings("mahjong_linebot.log")

	r := chi.NewRouter()

	r.Post("/v1/api/linebot", controllers.LineBotApiPost)

	r.Route("/v1/api/matchSetting", func(r chi.Router){
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://mahjong-linebot.firebaseapp.com"},
			AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders:   []string{"Content-Type"},
		}))
		r.Get("/{roomid}", controllers.GetMatchSettingByRoomId)
		r.Post("/", controllers.PostMatchSetting)
	})

	r.Route("/v1/api/matchResult", func(r chi.Router){
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://mahjong-linebot.firebaseapp.com"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "OPTIONS"},
			AllowedHeaders:   []string{"Content-Type"},
		}))
		r.Get("/{roomid}", controllers.GetMatchResultByRoomId)
		r.Post("/", controllers.PostMatchResult)
		r.Put("/", controllers.UpdateMatchResult)
	})

	log.Printf("コンテナ起動...")
	http.ListenAndServe(":8080", r)
}
