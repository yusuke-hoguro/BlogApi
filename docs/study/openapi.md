# OpneAPI 仕様書作成

## 手順の概要

1. swagツールの導入

- Go言語で書かれたAPIの仕様定義をコード内のコメントから自動で生成・文書化するためのツール

    ```bash
    go install github.com/swaggo/swag/cmd/swag@latest
    ```

2. 必要なライブラリの追加

- Swagger UIを表示するためのハンドラーを追加

    ```bash
    go get github.com/swaggo/http-swagger
    ```

3. コメントを記述

    ```go
    // @title Blog API
    // @version 1.0
    // @description This is a sample blog API built with Go net/http.
    // @host localhost:8080
    // @BasePath /
    ```

4. 各ハンドラーにコメントを追加

- `@Tags`, `@Param`, `@Success`, `@Failure` などを必要に応じて付与

5. Swagger Docs生成

- 下記のコマンドを実行してSwagger仕様を自動生成

    ```bash
    swag init -g main.go
    ```

6. Swagger UIをルーターに追加

    ```go
    import (
        "github.com/swaggo/http-swagger"
        _ "BlogApi/docs" // swag init で生成されたdocsパッケージをimport
    )

    func RegisterRoutes(r *mux.Router, db *sql.DB) {
        // Swagger UI
        r.Handle("/swagger/", httpSwagger.WrapHandler)

    }
    ```

7. 動作確認

    ```bash
    go run main.go
    ```
