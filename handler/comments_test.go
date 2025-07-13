package handler_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yusuke-hoguro/BlogApi/handler"
	"github.com/yusuke-hoguro/BlogApi/testutils"
)

// コメント投稿用APIのテスト
func TestPostCommentHandle(t *testing.T) {
	//テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	//テスト用サーバーのセットアップ
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	//テスト用のJWTトークン発行
	token, err := handler.GenerateJWT(2)
	if err != nil {
		t.Fatal("JWTの生成に失敗", err)
		return
	}

	//コメント投稿用のJSONデータ作成
	commentJSON := `{"content":"テストコメント"}`
	postID := 2
	url := fmt.Sprintf("%s/posts/%d/comments", server.URL, postID)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(commentJSON))
	if err != nil {
		t.Fatal("リクエスト生成失敗:", err)
	}
	req.Header.Set("Content-Type", "application/json")
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

// コメント削除用APIのテストを実施する
func TestDeleteCommentHandler(t *testing.T) {
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

	//コメント削除用のJSONデータ作成
	commentID := 1
	url := fmt.Sprintf("%s/comments/%d", server.URL, commentID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
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
	if resp.StatusCode != http.StatusOK {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusOK, resp.StatusCode)
	}

	//ログに表示
	body, _ := io.ReadAll(resp.Body)
	t.Logf("レスポンス: %s", string(body))

}
