package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

// Post構造体
type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// DBへの接続処理
func connectDB() (*sql.DB, error) {
	// DB接続を実施する
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Data Souce Nameの設定
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	// Postgress SQLに接続
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	// 接続確認
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// 記事一覧を取得する
func getPostsHandler(w http.ResponseWriter, r *http.Request) {
	// DB接続を実施
	db, err := connectDB()
	if err != nil {
		http.Error(w, "DB接続エラー", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// postsテーブルから指定したカラムのデータを取得する
	rows, err := db.Query("SELECT id, title, content FROM posts")
	if err != nil {
		http.Error(w, "DBクエリエラー", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		// DBから取得したデータをGoの変数に格納
		if err := rows.Scan(&p.ID, &p.Title, &p.Content); err != nil {
			http.Error(w, "データ取得エラー", http.StatusInternalServerError)
			return
		}
		// Post構造体のスライスに追加
		posts = append(posts, p)
	}
	// json形式に変更してレスポンスに書き込む
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, Blog API!")
}

func main() {
	// ポート取得
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// ハンドラー関数の設定
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/posts", getPostsHandler)
	// サーバー起動
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
	log.Println("Server started at :8080")
}
