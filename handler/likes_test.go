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
	userID := 1
	token, err := handler.GenerateJWT(userID)
	if err != nil {
		t.Fatal("JWTの生成に失敗", err)
		return
	}

	//いいね設定用
	postID := 1
	url := fmt.Sprintf("%s/posts/%d/like", server.URL, postID)

	//おなじリクエストを送信して重複登録されないことを確認する
	for i := 1; i <= 2; i++ {
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			t.Fatalf("[%d回目]リクエスト生成失敗: %v", i, err)
		}
		req.Header.Set("Authorization", token)

		//リクエスト送信
		client := server.Client()
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("[%d回目]HTTPリクエスト失敗: %v", i, err)
		}
		defer resp.Body.Close()

		//ステータスコード確認
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("[%d回目] 期待するステータスコード %d, 実際は %d", i, http.StatusCreated, resp.StatusCode)
		}

		//ログに表示
		body, _ := io.ReadAll(resp.Body)
		t.Logf("[%d回目] レスポンス: %s", i, string(body))
	}

	//いいねが1個であることを確認する
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE user_id = $1 AND post_id = $2", userID, postID).Scan(&count)
	if err != nil {
		t.Fatal("DBからいいね件数の取得に失敗:", err)
	}
	if count != 1 {
		t.Errorf("重複登録されている可能性あり: いいね件数 = %d (期待値: 1)", count)
	}

}
