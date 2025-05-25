package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

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

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, Blog API!")
}

func main() {

	//Postgress SQLに接続
	db, err := connectDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()
	fmt.Println("Connected to DB successfully!")

	// ポート取得
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// ハンドラー関数の設定
	http.HandleFunc("/", helloHandler)
	fmt.Println("Starting server on :8080")

	// サーバー起動
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}
