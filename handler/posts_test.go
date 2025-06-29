package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/yusuke-hoguro/BlogApi/middleware"
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
	//DBのセットアップを開始する
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("DB接続に失敗: %v", err)
	}
	defer db.Close()

	// HTTPリクエストの擬似オブジェクトの作成
	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	// ResponsWriterの擬似オブジェクト作成
	w := httptest.NewRecorder()

	// ハンドラー関数を取得して実行する
	handler := GetAllPostsHandler(db)
	handler(w, req)

	// 実行結果のレスポンスを取得
	resp := w.Result()
	defer resp.Body.Close()

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
}

// 投稿作成用APIのテスト
func TestCreatePostHandler(t *testing.T) {
	//DBのセットアップを開始する
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("DB接続に失敗: %v", err)
	}
	defer db.Close()

	//JWTトークンを発行
	token, err := generateJWT(3)
	if err != nil {
		t.Fatal("Failed to generate token")
		return
	}

	//jsonデータを構築
	postJSON := `{"title": "テスト投稿", "content": "これはテスト用です"}`
	req := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(postJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	// ResponsWriterの擬似オブジェクト作成
	w := httptest.NewRecorder()

	//ハンドラー関数を取得して実行
	handler := middleware.AuthMiddleware(CreatePostHandler(db))
	handler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("期待するテストコード %d, 実際は %d", http.StatusCreated, w.Code)
	}

	t.Logf("Response body: %s", w.Body.String())
}
