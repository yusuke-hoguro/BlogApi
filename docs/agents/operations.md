# Operations Guide

Docker、Nginx、Swagger/OpenAPI、deploy、運用系の変更前に読む資料です。

## Docker / Compose

- 開発環境は `infra/docker-compose.yml` を使う。
- 本番相当環境は `infra/docker-compose.prod.yml` を使う。
- テスト環境は `infra/docker-compose.test.yml` を使う。
- Makefile は Docker Compose の主要操作をラップしているため、基本的には Makefile のターゲットを使う。

## 実行環境

- 開発環境では API は `http://localhost:8080`、フロントエンドは `http://localhost:3000`、pgAdmin は `http://localhost:5050` で公開される。
- Backend の DB 接続は `.env` の `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` を `internal/db/db.go` が読み込む。
- Frontend の API 接続先は Vite の `VITE_API_BASE_URL` で決まる。
- `.env`、Vite 環境変数、Docker Compose、Nginx 設定を変更する場合は backend/frontend/Docker の接続経路をセットで確認する。

## Swagger/OpenAPI

- API handler には swaggo 用コメントがあり、`/swagger/` で Swagger UI を提供する。
- API 仕様を変えた場合は `docs/docs.go`, `docs/swagger.json`, `docs/swagger.yaml` の更新が必要か確認する。
- 生成コマンドは GitHub Actions の `generate-swagger.yml` と同じ `swag init -g ./cmd/api/main.go -o docs/` を使う。
- Swagger 生成物は `docs/` 配下にあるため、手動更新する場合は生成差分だけを確認し、無関係なドキュメント変更と混ぜない。

## Deployment / Production

- 本番デプロイは AWS EC2 + Docker Compose + Nginx + Certbot を想定している。
- 本番 Nginx、deploy script、migration 実行順序を変更する場合は README と workflow の説明も確認する。
- デプロイでは DB 起動確認後に migration を実行してからアプリケーションを起動する流れを壊さない。
- TLS 証明書、DuckDNS、GitHub Actions secrets などの秘密情報を commit しない。

## 避けるべきこと

- Docker/CI/deploy 構成変更を、無関係なアプリ機能変更と同じ PR に混ぜること。
- 生成物や大きなバイナリを不用意に追加すること。
- 本番 compose や Nginx 設定の変更時に README / workflow / deploy script の整合を確認しないこと。
