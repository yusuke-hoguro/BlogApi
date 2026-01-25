// @title Blog API
// @version 1.0
// @description This is a sample blog API built with Go net/http.
// @host localhost:8080
// @BasePath /
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/yusuke-hoguro/BlogApi/internal/db"
	"github.com/yusuke-hoguro/BlogApi/internal/middleware"
	"github.com/yusuke-hoguro/BlogApi/router"

	_ "github.com/lib/pq"
)

func main() {
	if err := runServer(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func runServer() error {
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

	r := mux.NewRouter()
	router.RegisterRoutes(r, db.DB)
	// CORSミドルウェアを適用
	handler := middleware.CorsMiddleware(r)
	log.Println("Server started at :" + port)

	// サーバー起動
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}
