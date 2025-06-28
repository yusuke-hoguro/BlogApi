package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"
)

// テスト用のDBを設定
func setupTestDB() (*sql.DB, error) {
	return sql.Open("postgres", "host=localhost port=5432 user=postgres password=yourpassword dbname=blog sslmode=disable")
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

	//次回はまず、この関数内の処理を理解するところから！
}
