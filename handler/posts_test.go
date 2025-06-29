package handler_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/yusuke-hoguro/BlogApi/handler"
	"github.com/yusuke-hoguro/BlogApi/testutils"
)

// 初期化処理
func init() {
	// 環境変数の読み込みを実施
	godotenv.Load("../.env")
}

// テスト用のDBを設定
func setupTestDB() (*sql.DB, error) {
	// DB接続設定
	dbHost := "localhost"
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Data Souce Nameの設定
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	return sql.Open("postgres", dsn)
}

// 全投稿を取得するAPIのテスト
func TestGetAllPostsHandler(t *testing.T) {
	//テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	//テスト用のサーバーを作成する
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	// HTTP Getリクエストの送信
	resp, err := http.Get(server.URL + "/posts")
	if err != nil {
		t.Fatalf("HTTPリクエスト失敗: %v", err)
	}
	defer resp.Body.Close()

	// 実行結果の確認
	if resp.StatusCode != http.StatusOK {
		t.Errorf("期待するテストコード %d, 実際は %d", http.StatusOK, resp.StatusCode)
	}

	// JSON配列をパースしてスライスに展開
	var posts []map[string]any
	err = json.NewDecoder(resp.Body).Decode(&posts)
	if err != nil {
		t.Errorf("JSONパースエラー: %v", err)
	}

	// 取得した結果を整形して出力する（インデント文字列を階層ごとに繰り返す）
	postsJSON, err := json.MarshalIndent(posts, "", "  ")
	if err != nil {
		t.Errorf("JSON整形エラー: %v", err)
	} else {
		t.Logf("取得した投稿データ: \n%s", postsJSON)
	}

	// Todo: Getした内容がただしかを比較する処理を追加

}

// 投稿作成用APIのテスト
func TestCreatePostHandler(t *testing.T) {
	//テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	//テスト用のサーバーを作成する
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	//JWTトークンを発行
	token, err := handler.GenerateJWT(3)
	if err != nil {
		t.Fatal("Failed to generate token")
		return
	}

	//jsonデータを構築
	postJSON := `{"title": "テスト投稿", "content": "これはテスト用です"}`
	req, err := http.NewRequest(http.MethodPost, server.URL+"/posts", strings.NewReader(postJSON))
	if err != nil {
		t.Fatal("リクエスト生成エラー:", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	client := server.Client()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("HTTPリクエスト失敗:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusCreated, resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	t.Logf("Response body: %s", string(body))

	// Todo:取得してあってるか確認もやる
}
