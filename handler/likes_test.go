package handler_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yusuke-hoguro/BlogApi/handler"
	"github.com/yusuke-hoguro/BlogApi/internal/models"
	"github.com/yusuke-hoguro/BlogApi/testutils"
)

// いいね設定用APIのテスト
func TestLikePostHandler(t *testing.T) {
	// テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	// テスト用サーバーのセットアップ
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	// テスト用のJWTトークン発行
	userID := 1
	token, err := handler.GenerateJWT(userID)
	if err != nil {
		t.Fatal("JWTの生成に失敗", err)
		return
	}

	// いいね設定用
	postID := 1
	url := fmt.Sprintf("%s/api/posts/%d/like", server.URL, postID)

	// 同じリクエストを送信して重複登録されないことを確認する
	for i := 1; i <= 2; i++ {
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			t.Fatalf("[%d回目]リクエスト生成失敗: %v", i, err)
		}
		req.Header.Set("Authorization", "Bearer "+token)

		// リクエスト送信
		client := server.Client()
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("[%d回目]HTTPリクエスト失敗: %v", i, err)
		}
		defer resp.Body.Close()

		// ステータスコード確認
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("[%d回目] 期待するステータスコード %d, 実際は %d", i, http.StatusCreated, resp.StatusCode)
		}

		// ログに表示
		body, _ := io.ReadAll(resp.Body)
		t.Logf("[%d回目] レスポンス: %s", i, string(body))
	}

	// いいねが1個であることを確認する
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE user_id = $1 AND post_id = $2", userID, postID).Scan(&count)
	if err != nil {
		t.Fatal("DBからいいね件数の取得に失敗:", err)
	}
	if count != 1 {
		t.Errorf("重複登録されている可能性あり: いいね件数 = %d (期待値: 1)", count)
	}

}

// いいね取得用APIのテスト
func TestGetLikesHandler(t *testing.T) {
	// テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	// テスト用サーバーのセットアップ
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	// いいねを取得する投稿
	postID := 1
	url := fmt.Sprintf("%s/api/posts/%d/likes", server.URL, postID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal("リクエスト生成失敗:", err)
	}

	// リクエスト送信
	client := server.Client()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("HTTPリクエスト失敗:", err)
	}
	defer resp.Body.Close()

	// ステータスコード確認
	if resp.StatusCode != http.StatusOK {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusOK, resp.StatusCode)
	}

	// レスポンスをパースして検証
	var result models.LikesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("JSONのデコード失敗: %v", err)
	}

	// 結果を比較する
	if result.PostID != postID {
		t.Errorf("PostID が一致しない: get %d, want %d", result.PostID, postID)
	}
	if result.LikeCount != len(result.UserIDs) {
		t.Errorf("like_count (%d) と user_ids の数 (%d) が一致しません", result.LikeCount, len(result.UserIDs))
	}

	// ログに表示
	t.Logf("取得したいいねの情報: %+v", result)
}
