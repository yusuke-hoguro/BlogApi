package handler_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"github.com/yusuke-hoguro/BlogApi/handler"
	"github.com/yusuke-hoguro/BlogApi/testutils"
)

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
		t.Fatal("JWTの生成に失敗:", err)
		return
	}

	//JSONデータを構築
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
}

// 投稿作成用API JWTトークンが無い場合のテストを実施する
func TestCreatePostHandlerUnauthorized(t *testing.T) {
	//テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	//テスト用のサーバーを作成する
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	//JSONデータを構築
	postJSON := `{"title": "テスト投稿", "content": "これはテスト用です"}`
	req, err := http.NewRequest(http.MethodPost, server.URL+"/posts", strings.NewReader(postJSON))
	if err != nil {
		t.Fatal("リクエスト生成エラー:", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := server.Client()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("HTTPリクエスト失敗:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusUnauthorized, resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	t.Logf("Response body: %s", string(body))
}

// 投稿作成用API 不正なリクエストのテストを実施する
func TestCreatePostHandlerInvalidJSON(t *testing.T) {
	//テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	//テスト用のサーバーを作成する
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	//JWTトークンを発行
	token, err := handler.GenerateJWT(3)
	if err != nil {
		t.Fatal("JWTの生成に失敗:", err)
		return
	}

	//不正なJSONデータを作成
	postJSON := `{"title": "テスト投稿", "content":}`
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

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusBadRequest, resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	t.Logf("Response body: %s", string(body))
}

// 投稿作成用API タイトルが空のテストを実施する
func TestCreatePostHandlerMissingTitle(t *testing.T) {
	//テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	//テスト用のサーバーを作成する
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	//JWTトークンを発行
	token, err := handler.GenerateJWT(3)
	if err != nil {
		t.Fatal("JWTの生成に失敗:", err)
		return
	}

	//タイトルが空のJSONデータを作成
	postJSON := `{"title": "", "content": "これはテスト用です"}`
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

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusBadRequest, resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	t.Logf("Response body: %s", string(body))
}

// 記事更新用ハンドラー関数のテスト
func TestUpdatePostHandler(t *testing.T) {
	// テスト用DBのセットアップ
	db := testutils.SetupTestDB(t)
	defer db.Close()

	// テスト用サーバーを作成する
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	// JWTトークンを発行
	token, err := handler.GenerateJWT(1)
	if err != nil {
		t.Fatal("JWTの生成に失敗:", err)
		return
	}

	// 更新用データのJSON
	updateJSON := `{"title": "更新されたタイトル", "content": "更新された内容"}`
	postID := 1
	url := fmt.Sprintf("%s/posts/%d", server.URL, postID)
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(updateJSON))
	if err != nil {
		t.Fatal("リクエスト生成エラー:", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	// HTTPリクエストを実行
	client := server.Client()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("HTTPリクエスト失敗:", err)
	}
	defer resp.Body.Close()

	// ステータスコードの確認
	if resp.StatusCode != http.StatusOK {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusOK, resp.StatusCode)
	}
	// レスポンスを表示
	body, _ := io.ReadAll(resp.Body)
	t.Logf("Response body: %s", string(body))
}

// 記事削除用ハンドラー関数のテスト
func TestDeletePostHandler(t *testing.T) {
	// テスト用DBのセットアップ
	db := testutils.SetupTestDB(t)
	defer db.Close()

	// テスト用サーバーのセットアップ
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	// JWTトークンを発行
	token, err := handler.GenerateJWT(2)
	if err != nil {
		t.Fatal("JWTの生成に失敗:", err)
	}

	// 削除対象のIDからURLを作成
	postID := 2
	url := fmt.Sprintf("%s/posts/%d", server.URL, postID)

	// リクエストの作成
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatal("リクエスト生成エラー:", err)
	}
	req.Header.Set("Authorization", token)

	// リクエスト実行
	client := server.Client()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("HTTPリクエスト失敗:", err)
	}
	defer resp.Body.Close()

	// ステータスコードの確認
	if resp.StatusCode != http.StatusOK {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusOK, resp.StatusCode)
	}

	// レスポンスを表示
	body, _ := io.ReadAll(resp.Body)
	t.Logf("Response body: %s", body)

}
