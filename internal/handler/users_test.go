package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yusuke-hoguro/BlogApi/testutils"
)

// ユーザー登録用APIの重複ユーザー確認テスト
func TestSignupHandlerDuplicateUsername(t *testing.T) {
	// テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	// テスト用サーバーのセットアップ
	h, cleanup := testutils.SetupTestServer(db)
	server := httptest.NewServer(h)
	defer server.Close()
	defer cleanup()

	// testdata/init_test.sql で作成済みのユーザー名を指定する
	signupJSON := `{"username":"testuser","password":"password123"}`
	req, err := http.NewRequest(http.MethodPost, server.URL+"/api/signup", strings.NewReader(signupJSON))
	if err != nil {
		t.Fatal("リクエスト生成失敗:", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// リクエスト送信
	resp, err := server.Client().Do(req)
	if err != nil {
		t.Fatal("HTTPリクエスト失敗:", err)
	}
	defer resp.Body.Close()

	// ステータスコード確認
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusConflict, resp.StatusCode)
	}
}
