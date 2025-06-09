package main

import (
	"log"
	"net/http"
	"os"

	"github.com/yusuke-hoguro/BlogApi/db"
	"github.com/yusuke-hoguro/BlogApi/middleware"
	"github.com/yusuke-hoguro/BlogApi/router"

	_ "github.com/lib/pq"
)

func main() {
	// DB接続を実施
	var err error
	db.DB, err = db.ConnectDB()
	if err != nil {
		log.Fatal("DB接続失敗:", err)
	}
	defer db.DB.Close()

	// ポート取得
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	router.RegisterRoutes(mux, db.DB)
	// CORSミドルウェアを適用
	handler := middleware.CorsMiddleware(mux)

	// サーバー起動
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
	log.Println("Server started at :8080")
}
