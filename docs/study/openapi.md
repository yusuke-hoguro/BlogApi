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
- 標準は`docs/`に生成されるが仕様書の混在を防ぐために`docs/swagger`に変更する

    ```bash
    swag init -g main.go -o docs/swagger
    ```

6. Swagger UIをルーターに追加

- 標準は`docs/`に生成されるが仕様書の混在を防ぐために`docs/swagger`に変更する

    ```go
    import (
        "github.com/swaggo/http-swagger"
        _ "BlogApi/docs/swagger" // swag init で生成されたdocsパッケージをimport
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

## Swaggerコメントの書き方ルール

- 各ハンドラー関数の直前に下記のような形式でコメントを追加する

    ```go
    // 関数名 godoc
    // @Summary 一言でAPIの概要
    // @Description 詳しい説明
    // @Tags グループ名
    // @Accept json (エンドポイントが受け取るリクエストの MIME タイプ)
    // @Produce json (エンドポイントが返すレスポンスの MIME タイプ) 
    // @Param パラメータ名 in (path|query|body) 型 必須説明 例: @Param <name> <in> <type> <required> "<description>"
    // @Success ステータスコード {object|array} モデル名 例: @Success <code> {object|array} <model(struct)>
    // @Failure ステータスコード {object} エラーレスポンスモデル　例: @Failure 401 {object} ErrorResponse
    // @Router パス [メソッド]　例: @Router /posts/{id} [get]
    ```

- 注意事項
    - モデル名(型)はパッケージ名を含めて指定 する

## トラブルシューティング
| No | 概要 | 内容 | 解決方法 |
|----|-------------|------|------|
| 1 | routerでのURI設定ミス | `HandleFunc`だとPathの完全一致になっており、`/swagger/index.html`や`/swagger/doc.json`にマッチしない | `PathPrefix`に変更し、接頭辞マッチに変更した。複数ファイルをまとめて使う場合は`PathPrefix`を使用する。|
| 2 | 仕様書が反映されない | 各Handler関数にコメントを記述したがOpenAPIの仕様書に反映されなかった | 自動で生成を組み込むまでは手動で`swag init -g main.go -o docs/swagger`を実行して再生成する必要がある |
