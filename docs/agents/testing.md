# Testing Guide

テスト追加、CI確認、E2E修正を行う前に読む資料です。

## Backend Tests

- 主要コマンドは `make test-go`。
- `make test-go` は `infra/docker-compose.test.yml` の PostgreSQL を起動し、DB の readiness を待ってからホスト側で `go test ./internal/handler/... -v` を実行する。
- テスト DB 初期化は `testutils.SetupTestDB(t)` と `testdata/init_test.sql` を使う。
- HTTP handler テストは `httptest.NewServer` と `testutils.SetupTestServer(db)` の既存パターンに合わせる。
- 認証が必要なテストでは `handler.GenerateJWT(userID)` を使う。
- 正常系だけでなく、400/401/403/404/409/500 相当の異常系を追加する。

## Frontend / E2E

- 主要コマンドは `make test-e2e` または `make pw-test`。
- Playwright は Chromium / Firefox / WebKit の 3 project で動く。
- `playwright.config.ts` が Docker Compose で frontend を起動し、Playwright の `baseURL` は `http://localhost:3000`。
- E2E の画面操作は frontend/nginx 経由で行い、API 呼び出しはフロントエンド build 時の `VITE_API_BASE_URL=http://localhost:8080` と Docker/Nginx 構成の整合を見る。
- テストユーザーや selector 文字列は `tests/e2e/fixtures` と `tests/e2e/constants` に寄せる。
- テストデータは `createUniqueText` など既存 utility を使い、並列実行や再実行で衝突しにくくする。

## CI 相当

```bash
make ci-test
```

必要に応じて個別に:

```bash
make go-lint
make test-go
make fe-build
make test-e2e
```

## 変更別の目安

- Backend handler/service/repository を変えたら、該当 handler test を追加・更新する。
- 認証・認可を変えたら、401/403 と所有者チェックのテストを追加・更新する。
- DB スキーマを変えたら、migration、`sql/init.sql`、`testdata/init_test.sql`、handler test を確認する。
- Frontend の主要導線を変えたら、Playwright E2E を追加・更新する。
- UI 文言、button label、testid を変えたら、E2E constants と selector を確認する。

## 注意

- `make test-go` と `make test-e2e` は Docker を使う。実行できない環境では、実行不能だった理由を明記する。
- テストのためだけに本番コードの責務分離を崩さない。
- E2E テストで不安定な CSS selector や固定データ名だけに依存しない。
