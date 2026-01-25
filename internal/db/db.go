package db

import (
	"database/sql"
	"fmt"
	"os"
)

// DBのポインタを保管
var DB *sql.DB

// DBへの接続処理
func ConnectDB() (*sql.DB, error) {
	// DB接続を実施する
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Data Souce Nameの設定
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	// Postgress SQLに接続
	DB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	// 接続確認
	err = DB.Ping()
	if err != nil {
		return nil, err
	}

	return DB, nil
}
