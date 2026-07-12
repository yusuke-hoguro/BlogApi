# Overview

このドキュメントは、AI エージェントが BlogApi の全体像を把握するための詳細資料です。作業を始める前に、対象領域が不明な場合やリポジトリ全体の前提を確認したい場合に読んでください。

## プロジェクト概要

BlogApi は、Go 標準の `net/http` を中心に実装されたブログ API と、React/Vite 製フロントエンドを同居させたポートフォリオプロジェクトです。

- Backend: Go, `net/http`, `gorilla/mux`, `database/sql`, PostgreSQL
- Frontend: React, Vite, Tailwind CSS, Axios, React Router
- Auth: JWT による認証、投稿・コメントの所有者チェックによる認可
- Features: 投稿、コメント、いいね、ユーザー登録、ログイン、ヘルスチェック
- Tests: Go handler/integration tests, Playwright E2E tests
- Infra: Docker Compose, Nginx, PostgreSQL, pgAdmin, AWS EC2, Certbot
- Docs: Swagger/OpenAPI, README, study notes

## ディレクトリ構成

```text
.
├── cmd/
│   ├── api/                 # API サーバーのエントリポイント
│   ├── migrate/             # DB マイグレーション実行コマンド
│   └── examples/            # 学習・検証用サンプル
├── internal/
│   ├── app/                 # repository/service の組み立て
│   ├── apperror/            # アプリケーションエラー型と HTTP ステータス変換
│   ├── config/              # アプリ内定数
│   ├── db/                  # DB 接続、マイグレーション
│   ├── handler/             # HTTP handler、入力検証、レスポンス共通処理
│   ├── middleware/          # JWT 認証、CORS、timeout
│   ├── models/              # API/DB で使う構造体
│   ├── repository/          # SQL 実行、永続化
│   ├── router/              # ルーティング登録
│   ├── service/             # ユースケース、認可などの業務ロジック
│   └── workerpool/          # 監査イベント処理用 worker pool
├── blog-api-frontend/
│   ├── src/                 # React アプリ本体
│   └── tests/e2e/           # Playwright E2E
├── sql/
│   ├── init.sql             # 新規 DB 初期化用 SQL
│   └── migrations/          # 既存 DB 更新用マイグレーション SQL
├── testdata/                # Go テスト用 SQL
├── testutils/               # Go テスト共通セットアップ
├── infra/                   # Docker Compose、Nginx、デプロイ関連
├── docs/                    # Swagger/OpenAPI、メモ、学習ドキュメント
├── Makefile                 # 開発・テスト・マイグレーション用コマンド
└── README.md                # セットアップ、運用、CI、デプロイ説明
```

## アーキテクチャ

Backend は明確なレイヤー分割を持つシンプルな構成です。

```text
cmd/api
  -> router
    -> middleware
    -> handler
      -> service
        -> repository
          -> database/sql + PostgreSQL
```

主な流れ:

1. `cmd/api/main.go` が DB 接続、worker pool、router、middleware、HTTP server を組み立てる。
2. `internal/router/routers.go` が URL と handler を対応付ける。
3. `handler` が HTTP リクエストの解析、JSON decode、入力検証、レスポンス生成を担当する。
4. `service` が所有者チェックなどの業務ロジックを担当する。
5. `repository` が SQL と DB エラー変換を担当する。
6. エラーは原則 `internal/apperror.AppError` として上位へ返し、handler で HTTP レスポンスへ変換する。

Frontend は `blog-api-frontend/src/api/client.ts` の Axios client を API 通信の中心にしています。JWT は `localStorage` の `token` を request interceptor で `Authorization: Bearer ...` として付与し、401/403 は response interceptor で処理します。

## 実行環境・設定の注意

- 開発環境では API は `http://localhost:8080`、フロントエンドは `http://localhost:3000`、pgAdmin は `http://localhost:5050` で公開される。
- Backend の DB 接続は `.env` の `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` を `internal/db/db.go` が読み込む。
- `.env` / README には `JWT_SECRET` があるが、現行コードの JWT 署名・検証は `internal/middleware/jwt.go` の固定値 `JwtKey = []byte("your_secret_key")` を使っている。環境変数化は既知の改善候補であり、現行仕様として扱わない。
- Frontend の API 接続先は Vite の `VITE_API_BASE_URL` で決まる。開発用 `.env.development` では `http://localhost:8080`。
- `.env`、Vite 環境変数、Docker Compose、Nginx 設定を変更する場合は backend/frontend/Docker の接続経路をセットで確認する。
- E2E やテストでは既存の fixtures / utils / constants を優先して再利用する。
