# Go Blog API

- Go言語とPostgreSQLを使ったブログAPIのポートフォリオです。  
- 投稿とコメント、いいねのCRUD機能を備え、JWTによる認証・認可も実装しています。  
- 開発・テスト環境はDocker Composeで統一し、GitHub ActionsによるCI/CDを構築しています。

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

- ルートに `.env` ファイルを作成し、以下を記述：

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

## EC2 デプロイ（ポートフォリオ用）

ブログAPIをAWS EC2上にデプロイし、DuckDNS と Nginx を使って HTTPS 化する手順を記載します。

### 前提

- AWSアカウントを保有している必要があります。
- EC2 インスタンス（Ubuntu 22.04推奨）を作成する必要があります。
- セキュリティグループで 80/tcp と 443/tcp を開放してください。
- Docker, docker-compose がAWS EC2上にインストールをしておいてください。
- DuckDNS で取得したサブドメインを持っている必要があります。（例: `blog-api.duckdns.org`）

### ディレクトリ構成（EC2用例）

```text
~/BlogApi
├── docker-compose.yml        # 本番用Docker Compose
├── nginx/
│   └── conf.d/
│       └── blogapi.conf      # Nginxリバースプロキシ設定
└── certbot/
    ├── conf/                 # 証明書保存用
    └── www/                  # Let's Encrypt認証用
```

### 環境変数

EC2でも `.env` をルートに配置します。

```env
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=blog
DB_PORT=5432
DB_HOST=db
JWT_SECRET=your_jwt_secret
APP_PORT=8080
DATABASE_URL=postgres://postgres:yourpassword@db:5432/blog?sslmode=disable
DUCKDNS_DOMAIN=blog-api.duckdns.org
```

### Docker Compose 起動

```bash
docker compose up --build -d
```

- `app` コンテナで Go API が動作します。
- `db` コンテナで PostgreSQL が動作します。
- `nginx` コンテナで HTTPS リバースプロキシが動作
- `certbot` コンテナで HTTPS リバースプロキシが動作

### Nginx + HTTPS 設定例

`nginx/conf.d/blogapi.conf`:

```nginx
# HTTP から HTTPS へのリダイレクト
server {
    listen 80;
    server_name blog-api.duckdns.org;

    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        return 301 https://$host$request_uri;
    }
}

# HTTPS
server {
    listen 443 ssl;
    server_name blog-api.duckdns.org;

    ssl_certificate /etc/letsencrypt/live/blog-api.duckdns.org/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/blog-api.duckdns.org/privkey.pem;

    location / {
        proxy_pass http://app:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

### Certbot（Let's Encrypt）証明書取得

```bash
sudo certbot certonly --webroot \
  -w ./certbot/www \
  -d blog-api.duckdns.org \
  --email <メールアドレス> --agree-tos --non-interactive
```


### 動作確認

```bash
curl -I https://blog-api.duckdns.org/api/posts
```

- ブラウザでもアクセス可能

### 注意点

-  本番用パスワードや JWT シークレットは絶対に公開しないこと。
-  DuckDNS の更新は自動化スクリプトで定期更新を実施すること。
-  EC2 再起動後も `docker compose up -d` で復旧可能です。

---


## テスト実行方法(Docker環境)

- 毎回クリーンなDBでテストをするため、`down -v` を利用します。

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

- GitHub Actionsにて `docker-compose.test.yml` を利用し、PR作成時やmainブランチ、developブランチにpush時に自動テストを実行。  
- ローカルと同一の環境でテストすることで、再現性の高いCIを実現しています。

---

## API仕様

- 詳しくは [API設計書 (Swagger UI)](https://yusuke-hoguro.github.io/BlogApi/) を参照

---

## 今後の改善予定

- 課題や改善については[TODO List](docs/issue/TODO_LIST.md)を参照

---

## 作者について

- このプロジェクトはGoとバックエンド開発の理解を深める目的で作成したポートフォリオです。

---
