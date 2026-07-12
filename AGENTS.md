# AGENTS.md

このファイルは、Codex、GitHub Copilot、Cursor、Claude Code などの AI エージェントが、このリポジトリで作業を始めるための入口です。

詳細な設計・実装ルールは `docs/agents/` 配下に分割しています。作業前に、このファイルを読んだうえで、対象領域に対応する詳細ドキュメントを必ず確認してください。

## プロジェクト概要

BlogApi は、Go 標準の `net/http` を中心に実装されたブログ API と、React/Vite 製フロントエンドを同居させたポートフォリオプロジェクトです。

- Backend: Go, `net/http`, `gorilla/mux`, `database/sql`, PostgreSQL
- Frontend: React, Vite, Tailwind CSS, Axios, React Router
- Auth: JWT による認証、投稿・コメントの所有者チェックによる認可
- Features: 投稿、コメント、いいね、ユーザー登録、ログイン、ヘルスチェック
- Tests: Go handler/integration tests, Playwright E2E tests
- Infra: Docker Compose, Nginx, PostgreSQL, pgAdmin, AWS EC2, Certbot
- Docs: Swagger/OpenAPI, README, study notes

## 最初に読む資料

作業内容に応じて、以下の詳細ドキュメントを読むこと。

| 作業内容 | 読む資料 |
| --- | --- |
| リポジトリ全体の把握 | `docs/agents/overview.md` |
| Backend API / Go 実装 | `docs/agents/backend.md` |
| Frontend / React 実装 | `docs/agents/frontend.md` |
| DB スキーマ / migration | `docs/agents/database.md` |
| テスト追加 / CI 確認 | `docs/agents/testing.md` |
| Docker / Swagger / deploy / 運用 | `docs/agents/operations.md` |
| PR レビュー / AI 生成時の注意 | `docs/agents/review.md` |

複数領域にまたがる変更では、該当する資料をすべて読むこと。たとえば、APIレスポンス変更は `backend.md`, `frontend.md`, `testing.md`, `review.md` を確認する。

## 最重要ルール

- 実装から確認できる事実を優先し、存在しないルールを前提にしない。迷う場合は「推奨」または「既知の改善候補」として扱う。
- Backend は `handler -> service -> repository -> database/sql + PostgreSQL` の依存方向を守る。
- `handler` に SQL を直接書かない。`service` で HTTP response を扱わない。`repository` から HTTP status を直接扱わない。
- エラーは原則 `internal/apperror.AppError` で表現し、HTTP レスポンスへの変換は handler 側に寄せる。
- 認証が必要な route は `middleware.AuthMiddleware` で包む。更新・削除では所有者チェックを確認する。
- DB スキーマ変更は `sql/migrations` を追加し、`sql/init.sql` と `testdata/init_test.sql` との整合を確認する。
- JWT、password、生の Authorization header、DB 接続情報などの秘密情報をログに出さない。
- E2E で触る UI は role/name または安定した `data-testid` で取得できるようにする。
- Go 変更後は `gofmt` / `goimports`、Frontend 変更後は ESLint / Vite build の影響を確認する。

## よく使うコマンド

```bash
make up-dev
make down-dev
make migrate-dev
make test-go
make go-lint
make fe-install
make fe-build
make test-e2e
make ci-test
```

## 注意

- `.env` / README には `JWT_SECRET` があるが、現行コードの JWT 署名・検証は `internal/middleware/jwt.go` の固定値 `JwtKey = []byte("your_secret_key")` を使っている。環境変数化は `docs/backlog.md` の改善候補として扱う。
- `AuthMiddleware` と `respondAppError` はエラー形式が完全には統一されていない。統一する場合は互換性とテストを含む明示的なリファクタとして扱う。
- `post_stats` は投稿作成時に初期行を作るが、現状ではコメント数・いいね数・閲覧数の集計更新には使われていない。
- `make test-go` と `make test-e2e` は Docker を使う。実行できない環境では、実行不能だった理由を明記する。
