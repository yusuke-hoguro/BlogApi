# Go Blog API

Go言語とPostgreSQLを使ったブログAPIのポートフォリオです。投稿とコメント、いいねのCRUD機能を備え、JWTによる認証・認可も実装しています。

---

## 使用技術

- Go (net/http, encoding/json, sql)
- PostgreSQL
- JWT 認証（github.com/golang-jwt/jwt）
- Docker（開発用PostgreSQL）
- `testing` パッケージによるユニットテスト

---

## 環境構築

### 必要ツール

- Go: 1.24.2
- PostgreSQL: 15.13
- OS: Ubuntu 22.04
- Docker:27.5.1

### クローンと初期化

```bash
git clone https://github.com/yusuke-hoguro/BlogApi.git
cd BlogApi
```

### .envの設定

ルートに `.env` ファイルを作成し、以下を記述：

```env
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=blog
DB_PORT=5432
DB_TEST_NAME=blog_test
DB_TEST_PORT=5433
DB_HOST=db
JWT_SECRET=your_jwt_secret
APP_PORT=8080
```

---

## 実行方法

```bash
go run main.go
```

http://localhost:8080 にてAPIサーバーが起動します。


---

## API仕様

詳しくは[API設計書](docs/API_SPEC.md)を参照

---

## テスト実行方法

```bash
go test ./... -v
```

テストは実装したハンドラー関数に対して正常系と異常系のテストを実施しています。※未対応のものは随時追加中

---

## 今後の改善予定

- フロントエンド実装（React/Vueなど）
- Swagger/OpenAPIドキュメント生成
- ソーシャル認証連携（Googleなど）
- タグ機能

---

## 作者について

このプロジェクトはGoとバックエンド開発の理解を深める目的で作成したポートフォリオです。

---
