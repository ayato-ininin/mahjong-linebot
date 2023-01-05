package config

import (
	"log"
	_ "gopkg.in/ini.v1"
	_ "os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
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

	// cfg, err := ini.Load("config.ini")
	// if err != nil {
	// 	log.Printf("Failed to read file: $v", err)
	// 	os.Exit(1) //configファイルが読めなかったら出る。
	// }
	// channel_secret := cfg.Section("linebot").Key("channel_secret").String()
	// access_token := cfg.Section("linebot").Key("access_token").String()

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
		log.Fatalf("failed to setup client: %v", err)
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
		log.Fatalf("failed to access secret version: %v", err)
	}

	return &result.Payload.Data
}
