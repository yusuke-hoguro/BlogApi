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

// 投稿作成用API バリデーション確認用のテストを実施する
func TestCreatePostHandlerValidation(t *testing.T) {
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

	// テスト実施用のテーブル（関数ごとに書いたほうがわかりやすいのでグローバルにしない）
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		// タイトルが空のテスト
		{
			name:       "empty title",
			body:       `{"title": "", "content": "本文"}`,
			wantStatus: http.StatusBadRequest,
		},
		// タイトルが文字数オーバーのテスト
		{
			name:       "title too long",
			body:       fmt.Sprintf(`{"title": "%s", "content": "本文"}`, strings.Repeat("a", handler.MaxTitleLength+1)),
			wantStatus: http.StatusBadRequest,
		},
		// 投稿内容が空の場合のテスト
		{
			name:       "empty content",
			body:       `{"title": "タイトル", "content": ""}`,
			wantStatus: http.StatusBadRequest,
		},
		// 投稿内容が文字数オーバーのテスト
		{
			name:       "content too long",
			body:       fmt.Sprintf(`{"title": "タイトル", "content": "%s"}`, strings.Repeat("a", handler.MaxContentLength+1)),
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		// サブテストを作成（第1引数：サブテストの名前 第2引数：サブテストの処理）
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, server.URL+"/posts", strings.NewReader(tt.body))
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

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("[%s] 期待するステータスコード %d, 実際は %d", tt.name, tt.wantStatus, resp.StatusCode)
			}

			body, _ := io.ReadAll(resp.Body)
			t.Logf("[%s] Response body: %s", tt.name, string(body))
		})
	}
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

// 投稿取得用APIのテスト
func TestGetPostsByIDHandler(t *testing.T) {
	//テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	//テスト用のサーバーを作成する
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	//投稿IDを指定してJSONデータ作成
	postID := 1
	url := fmt.Sprintf("%s/posts/%d", server.URL, postID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal("リクエスト生成失敗:", err)
	}

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

	body, _ := io.ReadAll(resp.Body)
	t.Logf("Response body: %s", string(body))
}

// 投稿取得用API 存在しない投稿IDのテストを実施
func TestGetPostsByIDHandlerNotFound(t *testing.T) {
	//テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	//テスト用のサーバーを作成する
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	//存在しない投稿IDを指定してJSONデータ作成
	postID := 9999
	url := fmt.Sprintf("%s/posts/%d", server.URL, postID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal("リクエスト生成失敗:", err)
	}

	client := server.Client()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("HTTPリクエスト失敗:", err)
	}
	defer resp.Body.Close()

	// ステータスコードの確認
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusNotFound, resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	t.Logf("Response body: %s", string(body))
}

// 投稿取得用API IDが数値でないテストを実施
func TestGetPostsByIDHandlerInvalidID(t *testing.T) {
	//テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	//テスト用のサーバーを作成する
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	//数値以外のIDを指定してJSONデータ作成
	url := fmt.Sprintf("%s/posts/aaa", server.URL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal("リクエスト生成失敗:", err)
	}

	client := server.Client()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("HTTPリクエスト失敗:", err)
	}
	defer resp.Body.Close()

	// ステータスコードの確認
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusBadRequest, resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	t.Logf("Response body: %s", string(body))
}
