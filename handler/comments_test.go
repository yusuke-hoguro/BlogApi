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
	token, err := handler.GenerateJWT(3)
	if err != nil {
		t.Fatal("JWTの生成に失敗", err)
		return
	}

	//コメント投稿用のJSONデータ作成
	commentJSON := `{"content":"テストコメント"}`
	postID := 3
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
	token, err := handler.GenerateJWT(2)
	if err != nil {
		t.Fatal("JWTの生成に失敗", err)
		return
	}

	//コメント削除用のJSONデータ作成
	commentID := 2
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

// コメント削除用APIで他人のコメント削除拒否テストを実施する
func TestDeleteCommentHandlerUnauthorized(t *testing.T) {
	//テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	//テスト用サーバーのセットアップ
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	//コメント投稿した人以外のユーザーIDを設定する
	token, err := handler.GenerateJWT(99)
	if err != nil {
		t.Fatal("JWTの生成に失敗", err)
		return
	}

	//コメント削除用のJSONデータ作成
	commentID := 2
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
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusForbidden, resp.StatusCode)
	}

	//ログに表示
	body, _ := io.ReadAll(resp.Body)
	t.Logf("レスポンス: %s", string(body))

}

// コメント削除用API 存在しないIDを指定した場合のテストを実施する
func TestDeleteCommentHandlerNotFound(t *testing.T) {
	//テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	//テスト用サーバーのセットアップ
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	//コメント投稿した人以外のユーザーIDを設定する
	token, err := handler.GenerateJWT(99)
	if err != nil {
		t.Fatal("JWTの生成に失敗", err)
		return
	}

	//存在しないコメントIDを指定してJSONデータ作成
	commentID := 9999
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
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusNotFound, resp.StatusCode)
	}

	//ログに表示
	body, _ := io.ReadAll(resp.Body)
	t.Logf("レスポンス: %s", string(body))

}

// コメント削除用API JWTトークンが無い場合のテストを実施する
func TestDeleteCommentHandlerNoAuthorization(t *testing.T) {
	//テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	//テスト用サーバーのセットアップ
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	//コメントIDを指定してJSONデータ作成
	commentID := 2
	url := fmt.Sprintf("%s/comments/%d", server.URL, commentID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatal("リクエスト生成失敗:", err)
	}

	//リクエスト送信
	client := server.Client()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("HTTPリクエスト失敗:", err)
	}
	defer resp.Body.Close()

	//ステータスコード確認
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusUnauthorized, resp.StatusCode)
	}
	//ログに表示
	body, _ := io.ReadAll(resp.Body)
	t.Logf("レスポンス: %s", string(body))

}

// 投稿のコメント取得用APIのテスト
func TestGetCommentsByPostIDHandler(t *testing.T) {
	// テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	// テスト用サーバーのセットアップ
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	// 投稿IDを設定
	postID := 1

	// リクエストの作成
	url := fmt.Sprintf("%s/posts/%d/comments", server.URL, postID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal("リクエスト生成失敗:", err)
	}

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

// コメント更新用APIのテスト
func TestUpdateCommentHandler(t *testing.T) {
	// テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	// テスト用サーバーのセットアップ
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	// テスト用のJWTトークン発行
	token, err := handler.GenerateJWT(1)
	if err != nil {
		t.Fatal("JWTの生成に失敗", err)
		return
	}

	// コメント投稿用のJSONデータ作成
	commentID := 3
	updateJSON := `{"content":"更新されたコメント内容"}`

	// リクエストを作成する
	url := fmt.Sprintf("%s/comments/%d", server.URL, commentID)
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(updateJSON))
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
	if resp.StatusCode != http.StatusOK {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusOK, resp.StatusCode)
	}

	//ログに表示
	body, _ := io.ReadAll(resp.Body)
	t.Logf("レスポンス: %s", string(body))

}

// コメント更新用API 存在しないコメントID選択時のテスト
func TestUpdateCommentHandlerNotFound(t *testing.T) {
	// テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	// テスト用サーバーのセットアップ
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	// テスト用のJWTトークン発行
	token, err := handler.GenerateJWT(1)
	if err != nil {
		t.Fatal("JWTの生成に失敗", err)
		return
	}

	// 存在しないコメントIDを設定して更新用JSONデータ作成
	commentID := 9999
	updateJSON := `{"content":"更新されたコメント内容"}`

	// リクエストを作成する
	url := fmt.Sprintf("%s/comments/%d", server.URL, commentID)
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(updateJSON))
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
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusNotFound, resp.StatusCode)
	}

	//ログに表示
	body, _ := io.ReadAll(resp.Body)
	t.Logf("レスポンス: %s", string(body))

}

// コメント更新用API 空のコンテンツを送った場合のテスト
func TestUpdateCommentHandlerEmptyContent(t *testing.T) {
	// テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	// テスト用サーバーのセットアップ
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	// テスト用のJWTトークン発行
	token, err := handler.GenerateJWT(2)
	if err != nil {
		t.Fatal("JWTの生成に失敗", err)
		return
	}

	// 空のコメントを設定して更新用JSONデータ作成
	commentID := 2
	updateJSON := `{"content":""}`

	// リクエストを作成する
	url := fmt.Sprintf("%s/comments/%d", server.URL, commentID)
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(updateJSON))
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
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusBadRequest, resp.StatusCode)
	}

	//ログに表示
	body, _ := io.ReadAll(resp.Body)
	t.Logf("レスポンス: %s", string(body))

}

// コメント更新用API JWTトークンが無い場合のテスト
func TestUpdateCommentHandlerNoAuthorization(t *testing.T) {
	// テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	// テスト用サーバーのセットアップ
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	// 空のコメントを設定して更新用JSONデータ作成
	commentID := 2
	updateJSON := `{"content":""}`

	// リクエストを作成する
	url := fmt.Sprintf("%s/comments/%d", server.URL, commentID)
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(updateJSON))
	if err != nil {
		t.Fatal("リクエスト生成失敗:", err)
	}
	req.Header.Set("Content-Type", "application/json")

	//リクエスト送信
	client := server.Client()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("HTTPリクエスト失敗:", err)
	}
	defer resp.Body.Close()

	//ステータスコード確認
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusUnauthorized, resp.StatusCode)
	}

	//ログに表示
	body, _ := io.ReadAll(resp.Body)
	t.Logf("レスポンス: %s", string(body))

}
