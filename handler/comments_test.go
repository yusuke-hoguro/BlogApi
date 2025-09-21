package handler_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yusuke-hoguro/BlogApi/handler"
	"github.com/yusuke-hoguro/BlogApi/models"
	"github.com/yusuke-hoguro/BlogApi/testutils"
)

// コメント投稿用APIのテスト
func TestPostCommentHandle(t *testing.T) {
	// テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	// テスト用サーバーのセットアップ
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	// テスト用のJWTトークン発行
	token, err := handler.GenerateJWT(3)
	if err != nil {
		t.Fatal("JWTの生成に失敗", err)
		return
	}

	// コメント投稿用のJSONデータ作成
	commentJSON := `{"content":"テストコメント"}`
	postID := 3
	url := fmt.Sprintf("%s/posts/%d/comments", server.URL, postID)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(commentJSON))
	if err != nil {
		t.Fatal("リクエスト生成失敗:", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	// リクエスト送信
	client := server.Client()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("HTTPリクエスト失敗:", err)
	}
	defer resp.Body.Close()

	// ステータスコード確認
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusCreated, resp.StatusCode)
	}

	// ログに表示
	body, _ := io.ReadAll(resp.Body)
	t.Logf("レスポンス: %s", string(body))

}

// コメント投稿用APIのバリデーション確認用テスト
func TestPostCommentHandleValidation(t *testing.T) {
	// テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	// テスト用サーバーのセットアップ
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	// テスト用のJWTトークン発行
	token, err := handler.GenerateJWT(3)
	if err != nil {
		t.Fatal("JWTの生成に失敗", err)
		return
	}

	// 投稿は起動時にテスト用sqlで作成済み（postID=3）
	postID := 3
	// テスト実施用のテーブル作成
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"empty content", `{"content": ""}`, http.StatusBadRequest},
		{"content too long", fmt.Sprintf(`{"content": "%s"}`, strings.Repeat("a", handler.MaxCommentLength+1)), http.StatusBadRequest},
	}

	// サブテストを実行する
	for _, tt := range tests {
		// サブテストを作成（第1引数：サブテストの名前 第2引数：サブテストの処理）
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/posts/%d/comments", server.URL, postID)
			req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(tt.body))
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

// コメント削除用APIのテストを実施する
func TestDeleteCommentHandler(t *testing.T) {
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

	// コメント削除用のJSONデータ作成
	commentID := 2
	url := fmt.Sprintf("%s/comments/%d", server.URL, commentID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatal("リクエスト生成失敗:", err)
	}
	req.Header.Set("Authorization", token)

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

	// ログに表示
	body, _ := io.ReadAll(resp.Body)
	t.Logf("レスポンス: %s", string(body))

}

// コメント削除用APIで他人のコメント削除拒否テストを実施する
func TestDeleteCommentHandlerUnauthorized(t *testing.T) {
	// テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	// テスト用サーバーのセットアップ
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	// コメント投稿した人以外のユーザーIDを設定する
	token, err := handler.GenerateJWT(99)
	if err != nil {
		t.Fatal("JWTの生成に失敗", err)
		return
	}

	// コメント削除用のJSONデータ作成
	commentID := 2
	url := fmt.Sprintf("%s/comments/%d", server.URL, commentID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatal("リクエスト生成失敗:", err)
	}
	req.Header.Set("Authorization", token)

	// リクエスト送信
	client := server.Client()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("HTTPリクエスト失敗:", err)
	}
	defer resp.Body.Close()

	// ステータスコード確認
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusForbidden, resp.StatusCode)
	}

	// ログに表示
	body, _ := io.ReadAll(resp.Body)
	t.Logf("レスポンス: %s", string(body))

}

// コメント削除用API 存在しないIDを指定した場合のテストを実施する
func TestDeleteCommentHandlerNotFound(t *testing.T) {
	// テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	// テスト用サーバーのセットアップ
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	// JWTトークンを発行
	token, err := handler.GenerateJWT(2)
	if err != nil {
		t.Fatal("JWTの生成に失敗", err)
		return
	}

	// 存在しないコメントIDを指定してJSONデータ作成
	commentID := 9999
	url := fmt.Sprintf("%s/comments/%d", server.URL, commentID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatal("リクエスト生成失敗:", err)
	}
	req.Header.Set("Authorization", token)

	// リクエスト送信
	client := server.Client()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("HTTPリクエスト失敗:", err)
	}
	defer resp.Body.Close()

	// ステータスコード確認
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusNotFound, resp.StatusCode)
	}

	// ログに表示
	body, _ := io.ReadAll(resp.Body)
	t.Logf("レスポンス: %s", string(body))

}

// コメント削除用API JWTトークンが無い場合のテストを実施する
func TestDeleteCommentHandlerNoAuthorization(t *testing.T) {
	// テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	// テスト用サーバーのセットアップ
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	// コメントIDを指定してJSONデータ作成
	commentID := 2
	url := fmt.Sprintf("%s/comments/%d", server.URL, commentID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
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
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusUnauthorized, resp.StatusCode)
	}
	// ログに表示
	body, _ := io.ReadAll(resp.Body)
	t.Logf("レスポンス: %s", string(body))

}

// コメント削除用API  無効なトークンを送信した場合のテスト
func TestDeleteCommentHandlerInvalidToken(t *testing.T) {
	// テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	// テスト用サーバーのセットアップ
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	// 改ざんされたトークンを用意
	invalidToken := "Bearer invalid.jwt.token"

	// コメントIDを指定してJSONデータ作成
	commentID := 2
	url := fmt.Sprintf("%s/comments/%d", server.URL, commentID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatal("リクエスト生成失敗:", err)
	}
	req.Header.Set("Authorization", invalidToken)

	// リクエスト送信
	client := server.Client()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("HTTPリクエスト失敗:", err)
	}
	defer resp.Body.Close()

	// ステータスコード確認
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusUnauthorized, resp.StatusCode)
	}
	// ログに表示
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
	postID := 2

	// リクエストの作成
	url := fmt.Sprintf("%s/posts/%d/comments", server.URL, postID)
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

	// ログに表示
	body, _ := io.ReadAll(resp.Body)
	var comments []models.Comment
	if err := json.Unmarshal(body, &comments); err != nil {
		t.Fatal("JSON Unmarshal失敗:", err)
	}

	// コメント件数の確認
	if len(comments) != 1 {
		t.Fatalf("期待件数 1, 実際 %d", len(comments))
	}

	// 取得した内容をチェック
	if comments[0].PostID != postID {
		t.Errorf("期待するPostID %d, 実際は %d", postID, comments[0].PostID)
	}
	if comments[0].Content == "" {
		t.Errorf("コメント内容が空です")
	}
	if comments[0].UserID != 2 {
		t.Errorf("期待するUserID 1, 実際は %d", comments[0].UserID)
	}

	t.Logf("レスポンス: %s", string(body))

}

// 投稿のコメント取得用APIのテスト
func TestGetCommentsByPostIDHandlerMultiple(t *testing.T) {
	// テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	// テスト用サーバーのセットアップ
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	// 投稿IDを設定(init_test.sqlで複数コメント登録済み)
	postID := 1

	// リクエストの作成
	url := fmt.Sprintf("%s/posts/%d/comments", server.URL, postID)
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

	// ログに表示
	body, _ := io.ReadAll(resp.Body)
	var comments []models.Comment
	if err := json.Unmarshal(body, &comments); err != nil {
		t.Fatal("JSON Unmarshal失敗:", err)
	}
	if len(comments) != 2 {
		t.Errorf("期待件数 2, 実際 %d", len(comments))
	}

	// 取得内容のチェック(PostID)
	if comments[0].PostID != postID {
		t.Errorf("期待するPostID %d, 実際は %d", postID, comments[0].PostID)
	}
	if comments[1].PostID != postID {
		t.Errorf("期待するPostID %d, 実際は %d", postID, comments[0].PostID)
	}
	// 取得内容のチェック(Content)
	if comments[0].Content == "" || comments[1].Content == "" {
		t.Errorf("コメント内容が空です")
	}
	// 取得内容のチェック(UserID)
	if comments[0].UserID != 1 || comments[1].UserID != 1 {
		t.Errorf("UserIDが不正です")
	}
	t.Logf("レスポンス: %s", string(body))
}

// 投稿のコメント取得用APIのテスト
func TestGetCommentsByPostIDHandlerEmpty(t *testing.T) {
	// テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	// テスト用サーバーのセットアップ
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	// 投稿IDを設定(init_test.sqlで複数コメント登録済み)
	postID := 3

	// リクエストの作成
	url := fmt.Sprintf("%s/posts/%d/comments", server.URL, postID)
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

	// 取得したコメントを格納
	body, _ := io.ReadAll(resp.Body)
	var comments []models.Comment
	if err := json.Unmarshal(body, &comments); err != nil {
		t.Fatal("JSON Unmarshal失敗:", err)
	}
	if len(comments) != 0 {
		t.Errorf("期待件数 0, 実際 %d", len(comments))
	}
	t.Logf("レスポンス: %s", string(body))
}

// 投稿のコメント取得用API 存在しない投稿IDのテスト
func TestGetCommentsByPostIDHandlerPostNotFound(t *testing.T) {
	// テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	// テスト用サーバーのセットアップ
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	// 存在しない投稿IDを設定
	postID := 9999

	// リクエストの作成
	url := fmt.Sprintf("%s/posts/%d/comments", server.URL, postID)
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

	// 取得したコメントを格納する
	body, _ := io.ReadAll(resp.Body)
	var comments []models.Comment
	if err := json.Unmarshal(body, &comments); err != nil {
		t.Fatal("JSON Unmarshal失敗:", err)
	}

	// 存在しない投稿の場合は空想定
	if len(comments) != 0 {
		t.Errorf("期待件数 0, 実際 %d", len(comments))
	}

	t.Logf("レスポンス: %s", string(body))
}

// 投稿のコメント取得用API 数値でない投稿IDのテスト
func TestGetCommentsByPostIDHandlerInvalidID(t *testing.T) {
	// テスト用DBのセットアップを開始する
	db := testutils.SetupTestDB(t)
	defer db.Close()

	// テスト用サーバーのセットアップ
	server := httptest.NewServer(testutils.SetupTestServer(db))
	defer server.Close()

	// リクエストの作成
	url := fmt.Sprintf("%s/posts/ddd/comments", server.URL)
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
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusBadRequest, resp.StatusCode)
	}

	// 取得したコメントを格納
	body, _ := io.ReadAll(resp.Body)

	// エラーメッセージが含まれていることを確認する
	if !strings.Contains(string(body), "Invalid") {
		t.Errorf("レスポンスにエラーメッセージが含まれていません: %s", string(body))
	}
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

	// ログに表示
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

	// リクエスト送信
	client := server.Client()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("HTTPリクエスト失敗:", err)
	}
	defer resp.Body.Close()

	// ステータスコード確認
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusNotFound, resp.StatusCode)
	}

	// ログに表示
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

	// リクエスト送信
	client := server.Client()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("HTTPリクエスト失敗:", err)
	}
	defer resp.Body.Close()

	// ステータスコード確認
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusBadRequest, resp.StatusCode)
	}

	// ログに表示
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

	// リクエスト送信
	client := server.Client()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("HTTPリクエスト失敗:", err)
	}
	defer resp.Body.Close()

	// ステータスコード確認
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusUnauthorized, resp.StatusCode)
	}

	// ログに表示
	body, _ := io.ReadAll(resp.Body)
	t.Logf("レスポンス: %s", string(body))

}

// コメント更新用API 他人のコメントを更新しようとした場合のテスト
func TestUpdateCommentHandlerForbidden(t *testing.T) {
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

	// テスト用のコメントはinit_test.sqlで作成済み
	commentID := 3
	updateJSON := `{"content":"他人のコメントを不正に更新しようとする"}`

	// リクエストを作成する
	url := fmt.Sprintf("%s/comments/%d", server.URL, commentID)
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(updateJSON))
	if err != nil {
		t.Fatal("リクエスト生成失敗:", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	// リクエスト送信
	client := server.Client()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("HTTPリクエスト失敗:", err)
	}
	defer resp.Body.Close()

	// ステータスコード確認
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("期待するステータスコード %d, 実際は %d", http.StatusForbidden, resp.StatusCode)
	}

	// ログに表示
	body, _ := io.ReadAll(resp.Body)
	t.Logf("レスポンス: %s", string(body))

}

// コメント更新用APIのテスト
func TestUpdateCommentHandlerValidation(t *testing.T) {
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

	// init_test.sqlでテスト用のコメント作成済み
	commentID := 3

	// テスト実施用のテーブル作成
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"empty content", `{"content": ""}`, http.StatusBadRequest},
		{"content too long", fmt.Sprintf(`{"content": "%s"}`, strings.Repeat("a", handler.MaxCommentLength+1)), http.StatusBadRequest},
	}

	// サブテストを実行する
	for _, tt := range tests {
		// サブテストを作成（第1引数：サブテストの名前 第2引数：サブテストの処理）
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/comments/%d", server.URL, commentID)
			req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(tt.body))
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
