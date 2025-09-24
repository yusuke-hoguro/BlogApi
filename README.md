# Go Blog API

Go言語とPostgreSQLを使ったブログAPIのポートフォリオです。  
投稿とコメント、いいねのCRUD機能を備え、JWTによる認証・認可も実装しています。  
開発・テスト環境はDocker Composeで統一し、GitHub ActionsによるCI/CDを構築しています。

---

## 使用技術

- 言語: Go (net/http, encoding/json, sql)
- データベース: PostgreSQL
- 認証: JWT（github.com/golang-jwt/jwt）
- コンテナ: Docker, docker-compose
- CI/CD: GitHub Actions
- テスト: GO `testing` パッケージ

---

## 環境構築

### 必要ツール

- Go: 1.24.2
- PostgreSQL: 15.13
- OS: Ubuntu 22.04（推奨）
- Docker:27.5.1
- docker compose V2

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

# DB接続用URL（Goのアプリで利用）
DATABASE_URL=postgres://postgres:yourpassword@db:5432/blog?sslmode=disable
```

---

## 実行方法(ローカル)

```bash
docker compose up --build
```

http://localhost:8080 にてAPIサーバーが起動します。  
pgAdmin も利用可能です。（http://localhost:5050）


---

## テスト実行方法(Docker環境)

毎回クリーンなDBでテストをするため、`down -v` を利用します。

```bash
# DBを初期化してテスト実行
docker compose -f docker-compose.test.yml up --build --abort-on-container-exit

# テスト後にコンテナ削除
docker compose -f docker-compose.test.yml down -v
```
---

## ログ確認方法

```bash
# DBログ
docker compose -f docker-compose.test.yml logs db

# アプリケーションログ
docker compose -f docker-compose.test.yml logs goapp
```

---

## CI/CD

GitHub Actionsにて `docker-compose.test.yml` を利用し、PR作成時やmainブランチ、developブランチにpush時に自動テストを実行。  
ローカルと同一の環境でテストすることで、再現性の高いCIを実現しています。

---

## API仕様

詳しくは[API設計書](docs/API_SPEC.md)を参照

---

## 今後の改善予定

- フロントエンド実装（React/Vueなど）
- Swagger/OpenAPIドキュメント生成
- ソーシャル認証連携（Googleなど）
- タグ機能
- CI/CDの本番デプロイ自動化

---

## 作者について

このプロジェクトはGoとバックエンド開発の理解を深める目的で作成したポートフォリオです。

---
