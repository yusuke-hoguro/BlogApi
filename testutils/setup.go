package testutils

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/yusuke-hoguro/BlogApi/handler"
	"github.com/yusuke-hoguro/BlogApi/middleware"
)

// 初期化処理
func init() {
	// 環境変数の読み込みを実施
	godotenv.Load("../.env")
}

// テスト用のDBを設定する
func SetupTestDB(t *testing.T) *sql.DB {
	// ヘルパー関数定義（呼び出し元のエラーを表示する）
	t.Helper()
	// DB接続設定
	dbHost := "localhost"
	dbPort := os.Getenv("DB_TEST_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_TEST_NAME")

	// Data Souce Nameの設定
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("DB接続失敗: %v", err)
	}
	//テスト用のSQLを読み込んで実行する
	loadTestSQL(t, db, getTestdataPath("init_test.sql"))

	return db
}

// テスト用のSQL文を読み込む
func loadTestSQL(t *testing.T, db *sql.DB, filepath string) {
	t.Helper()
	// 指定したファイルからSQLを読み込む
	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("SQLファイル読み込み失敗: %v", err)
	}

	// 読み込んだSQLをすべて実行する
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		t.Fatalf("SQL実行失敗: %v", err)
	}
}

// テスト用のサーバーを設定する
func SetupTestServer(db *sql.DB) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/posts", handler.GetAllPostsHandler(db)).Methods("GET")                                           //全投稿取得用
	r.HandleFunc("/posts", middleware.AuthMiddleware(handler.CreatePostHandler(db))).Methods("POST")                //個別投稿作成用
	r.HandleFunc("/posts/{id}", middleware.AuthMiddleware(handler.UpdatePostHandler(db))).Methods("PUT")            //個別投稿更新用
	r.HandleFunc("/posts/{id}", middleware.AuthMiddleware(handler.DeletePostHandler(db))).Methods("DELETE")         //個別投稿削除用
	r.HandleFunc("/posts/{id}/comments", middleware.AuthMiddleware(handler.PostCommentHandler(db))).Methods("POST") //コメント投稿
	r.HandleFunc("/comments/{id}", middleware.AuthMiddleware(handler.DeleteCommentHandler(db))).Methods("DELETE")   //コメントIDで削除
	r.HandleFunc("/posts/{id}/like", middleware.AuthMiddleware(handler.LikePostHandler(db))).Methods("POST")        //投稿にいいねをつける
	r.HandleFunc("/posts/{id}/likes", handler.GetLikesHandler(db)).Methods("GET")                                   //投稿のいいねを取得する
	return r
}

// テスト用データのパスを取得する
func getTestdataPath(filename string) string {
	_, b, _, _ := runtime.Caller(0) // 呼び出し元のファイルパスを取得
	base := filepath.Dir(b)         // 現在のファイルのディレクトリ
	return filepath.Join(base, "..", "testdata", filename)
}
