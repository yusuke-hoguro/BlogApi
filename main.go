package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
)

// DBのポインタを保管
var db *sql.DB

// JWT認証用の秘密鍵（最終的には環境変数から）
var jwtKey = []byte("your_secret_key")

// Post構造体
type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// User登録用の構造体
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
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

// 記事用のハンドラー関数
func postHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getPostsHandler(w, r)
	case http.MethodPost:
		createPostHandler(w, r)
	case http.MethodPut:
		updatePostHandler(w, r)
	case http.MethodDelete:
		deletePostHandler(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// 記事一覧取得用のハンドラー関数
func getPostsHandler(w http.ResponseWriter, r *http.Request) {
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

// 記事作成用のハンドラー関数
func createPostHandler(w http.ResponseWriter, r *http.Request) {
	// Post型の構造体にデコードして格納
	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// INSERT実行
	err := db.QueryRow("INSERT INTO posts (title, content) VALUES ($1, $2) RETURNING id", post.Title, post.Content).Scan(&post.ID)
	if err != nil {
		http.Error(w, "Failed to insert post", http.StatusInternalServerError)
		return
	}

	// 作成した記事IDをJSONで返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

// 記事更新のハンドラー関数
func updatePostHandler(w http.ResponseWriter, r *http.Request) {
	// URLからIDを取得する
	idStr := strings.TrimPrefix(r.URL.Path, "/posts/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Post型の構造体にデコードして格納
	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// UPDATE実行
	result, err := db.Exec("UPDATE posts SET title = $1, content = $2 WHERE id = $3", post.Title, post.Content, id)
	if err != nil {
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	// 更新行数の確認
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, "Post nor found or no changes", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Post update successfully!")
}

// 記事削除のハンドラー関数
func deletePostHandler(w http.ResponseWriter, r *http.Request) {
	// URLからIDを取得する
	idStr := strings.TrimPrefix(r.URL.Path, "/posts/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// DELETE実行
	result, err := db.Exec("DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	// 削除行数の確認
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to confirm deletion", http.StatusInternalServerError)
		return
	} else if rowsAffected == 0 {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ユーザー登録用のハンドラー関数
func signupHandler(w http.ResponseWriter, r *http.Request) {
	// Postであるかをチェックする
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// GOの構造体にデコード
	var userData User
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// バリデーション
	if userData.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	if len(userData.Password) < 8 {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// INSERT実行
	err = db.QueryRow("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", userData.Username, string(hashedPassword)).Scan(&userData.ID)
	if err != nil {
		http.Error(w, "Failed to insert post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userData)
}

// ログイン機能用用のハンドラー関数
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Postであるかをチェックする
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// GOの構造体にデコード
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// バリデーション
	if creds.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	if creds.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	// ユーザーをDBから検索
	var id int
	var hashedPassword string
	err = db.QueryRow("SELECT id, password FROM users WHERE username = $1", creds.Username).Scan(&id, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	// パスワード照合
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// JWT生成
	token, err := generateJWT(id)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// レスポンス
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// JWTの検証を実施するミドルウェア
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストヘッダーの確認
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// JWTの解析
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// 引数で指定されたハンドラー関数を実行
		next(w, r)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, Blog API!")
}

// JWTトークンを発行する
func generateJWT(userID int) (string, error) {
	// payloadの生成
	claims := &jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), //トークンの有効時間(24時間後)
	}

	// JWTを生成する
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 署名付きトークン生成
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func main() {
	var err error
	// DB接続を実施
	db, err = connectDB()
	if err != nil {
		log.Fatalf("DB接続エラー: %v", err)
		return
	}
	defer db.Close()

	// ポート取得
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// ハンドラー関数の設定
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/posts", authMiddleware(postHandler))
	http.HandleFunc("/posts/", authMiddleware(postHandler))
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/login", loginHandler)

	// サーバー起動
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
	log.Println("Server started at :8080")
}
