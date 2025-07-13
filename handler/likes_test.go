package handler_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yusuke-hoguro/BlogApi/handler"
	"github.com/yusuke-hoguro/BlogApi/testutils"
)

// いいね設定用APIのテスト
func TestLikePostHandler(t *testing.T) {
	//テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	//テスト用サーバーのセットアップ
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	//テスト用のJWTトークン発行
	token, err := handler.GenerateJWT(1)
	if err != nil {
		t.Fatal("JWTの生成に失敗", err)
		return
	}

	//いいね設定用
	postID := 1
	url := fmt.Sprintf("%s/posts/%d/like", server.URL, postID)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		t.Fatal("リクエスト生成失敗:", err)
	}
	req.Header.Set("Authorization", token)

	//リクエスト送信
	client := server.Client()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("HTTPリクエスト失敗:", err)
	}
	defer resp.Body.Close()

	//ステータスコード確認
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusCreated, resp.StatusCode)
	}

	//ログに表示
	body, _ := io.ReadAll(resp.Body)
	t.Logf("レスポンス: %s", string(body))

}
