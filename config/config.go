package config

import (
	"fmt"
	"log"
	logger "mahjong-linebot/utils"

	"context"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type ConfigList struct {
	ChannelSecret string
	AccessToken   string
}

var Config ConfigList //グローバル変数

// パッケージを読み込むときに、一回だけ読み込まれる。
// main.goからimportされたとき、設定ファイルを読み込むことができる。
// それをグローバル変数に入れてるから、main.goからグローバル変数として呼び出せる仕組み。
// 別途、config.iniファイルの作成が必要。
func init() {
	channel_secret := *GetDataFromSecretManager("LINEBOT_CHANNEL_SECRET")
	access_token := *GetDataFromSecretManager("LINEBOT_ACCESS_TOKEN")

	Config = ConfigList{
		ChannelSecret: string(channel_secret),
		AccessToken:   string(access_token),
	}
}

/*
*
secret managerに保存しているデータをバイト配列で返す
*
*/
func GetDataFromSecretManager(secretName string) *[]byte {
	// Use a service account
	projectID := "1033476136185"

	// クライアント生成
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Secret Manager設定失敗 err=%v", err)))
	}
	defer client.Close()

	// シークレット、test-appの最新バージョンにアクセス
	resourceName := "projects/" + projectID + "/secrets/" + secretName + "/versions/latest"
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: resourceName,
	}

	// シークレット上にアクセスする
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		log.Printf(logger.ErrorLogEntry(fmt.Sprintf("Failed to access secret version err=%v", err)))
	}

	return &result.Payload.Data
}
