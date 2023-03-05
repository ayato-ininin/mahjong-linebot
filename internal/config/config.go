package config

import (
	"log"

	"context"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type ConfigList struct {
	ChannelSecret string
	AccessToken   string
}

// initだと他のファイルのテスト時に配下パッケージの中でconfig読んでると,
// 　これが呼び出されてテストエラーになるので、initConfigに変更
func InitConfig() (*ConfigList, error) {
	var channel_secret []byte
	var access_token []byte
	var err error

	channel_secret, err = GetDataFromSecretManager("LINEBOT_CHANNEL_SECRET")
	if err != nil {
		log.Printf("channel_secret取得失敗 err=%v", err)
		return nil, err
	}
	access_token, err = GetDataFromSecretManager("LINEBOT_ACCESS_TOKEN")
	if err != nil {
		log.Printf("access_token取得失敗 err=%v", err)
		return nil, err
	}

	var config ConfigList
	config.ChannelSecret = string(channel_secret)
	config.AccessToken = string(access_token)

	return &config, nil
}

// secret manager(GCP)に保存しているデータをバイト配列で返す
func GetDataFromSecretManager(secretName string) ([]byte, error) {
	// Use a service account
	projectID := "1033476136185"

	// クライアント生成
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Printf("Secret Manager設定失敗 err=%v", err)
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
		log.Printf("Failed to access secret version err=%v", err)
	}

	return result.Payload.Data, err
}
