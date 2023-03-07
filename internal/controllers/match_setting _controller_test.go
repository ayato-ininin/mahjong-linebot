package controllers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mahjong-linebot/internal/config"
	"mahjong-linebot/internal/models"
	"mahjong-linebot/internal/repositories"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/go-chi/chi/v5"
	"google.golang.org/api/option"
)

func TestName(t *testing.T) {
	tests := []struct {
		name           string
		roomID         int
		wantStatusCode int
		wantBody       string
		want           models.MatchSetting
	}{
		{
			name:           "first case",
			roomID:         1,
			wantStatusCode: http.StatusOK,
			want: models.MatchSetting{
				RoomId: 1,
				MahjongNumber: "三麻",
				Name1: "name1",
				Name2: "name2",
				Name3: "name3",
				Name4: "name4",
				Uma: "10-20",
				Oka: 30000,
				IsYakitori: false,
				YakitoriPoint: 10,
				IsTobishou: false,
				TobishouPoint: 10,
				Rate: 50,
				IsTip: false,
				TipInitialNumber: 20,
				TipRate: 2,
				CreateTimestamp:  time.Date(2023, 3, 5, 12, 34, 56, 0, time.UTC),
				UpdateTimestamp: time.Date(2023, 3, 5, 12, 34, 56, 0, time.UTC),
			},
		},
	}

	if len(os.Getenv("FIRESTORE_EMULATOR_HOST")) == 0 {
		t.Fatal("firestore emulator is not running")
	}
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		t.Fatalf("firestore emulator is not running: %v", err)
	}
	conn.Close()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			// このへんで data を add する
			m := &models.MatchSetting{
				RoomId: 1,
				MahjongNumber: "三麻",
				Name1: "name1",
				Name2: "name2",
				Name3: "name3",
				Name4: "name4",
				Uma: "10-20",
				Oka: 30000,
				IsYakitori: false,
				YakitoriPoint: 10,
				IsTobishou: false,
				TobishouPoint: 10,
				Rate: 50,
				IsTip: false,
				TipInitialNumber: 20,
				TipRate: 2,
				CreateTimestamp: time.Date(2023, 3, 5, 12, 34, 56, 0, time.UTC),
				UpdateTimestamp: time.Date(2023, 3, 5, 12, 34, 56, 0, time.UTC),
			}
			ctx := context.Background()
			client, err := firebaseInit(ctx)
			if err != nil {
				t.Fatal(err)
			}
			err = repositories.SetMatchSetting(ctx, client, m, time.Date(2023, 3, 5, 12, 34, 56, 0, time.UTC), 1)
			if err != nil {
				t.Fatal(err)
			}

			t.Log("### Set DONE ###")

			reqBody := bytes.NewBufferString("request body")
			url := fmt.Sprintf("http://localhost:8080/%d", tt.roomID)
			req := httptest.NewRequest(http.MethodGet, url, reqBody)

			// パスパラメータを設定
			r := chi.NewRouter()
			r.Get("/{roomid}", GetMatchSettingByRoomId)

			got := httptest.NewRecorder()

			// テスト対象のハンドラー関数を呼び出し
			r.ServeHTTP(got, req)

			// 帰ってきた結果が意図するものと一緒かどうか比較
			if got.Code != http.StatusOK {
				t.Errorf("want %d, but %d", tt.wantStatusCode, got.Code)
			}

			jsonStr, err := json.Marshal(tt.want)
			if err != nil {
				return
			}

			//なぜかgotには末尾に空白が入っているので、取り除く
			if got := got.Body.String(); strings.TrimSpace(got) != string(jsonStr) {
				t.Errorf("want %s, but %s", tt.wantBody, got)
			}
		})
	}
}

func firebaseInit(ctx context.Context) (*firestore.Client, error) {
	// Firebaseのサービスアカウントキーの取得
	jsonBytes := getFirebaseServiceAccountKey()
	sa := option.WithCredentialsJSON(jsonBytes)

	// Firebaseアプリケーションの初期化
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return nil, err
	}

	// Firestoreクライアントの初期化
	fsClient, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	return fsClient, nil
}


func getFirebaseServiceAccountKey() []byte {
	data, err := config.GetDataFromSecretManager("MAHJONG_LINEBOT_FIREBASE_SA")
	if err != nil {
		return nil
	}
	dec, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil
	}

	return dec
}
