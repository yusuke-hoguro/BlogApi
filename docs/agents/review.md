# Review Guide

PR レビュー、AI 生成コードの確認、実装前後のセルフチェックで読む資料です。

## PRレビュー観点

API / Backend:

- handler/service/repository の責務が混ざっていないか。
- 認証が必要な route に `middleware.AuthMiddleware` が付いているか。
- 更新・削除・自分専用取得などで所有者チェックが抜けていないか。
- 入力検証、ID parse、JSON decode のエラーが適切な HTTP status になっているか。
- DB エラー、`sql.ErrNoRows`、`RowsAffected == 0` が `apperror` に変換されているか。
- 新しい SQL が SQL injection を避け、プレースホルダを使っているか。
- 複数テーブル更新でトランザクションが必要な箇所に `BeginTx` があるか。
- migration と `sql/init.sql` / `testdata/init_test.sql` の整合性が取れているか。

Frontend:

- API client の認証・エラー処理と矛盾していないか。
- ルーティング、ログイン状態、権限なし状態の UI が破綻しないか。
- E2E が role/name または安定した testid で対象を取れているか。
- ユーザー操作後の遷移・表示更新を適切に await しているか。

Tests / CI:

- 変更した振る舞いに対する Go test または Playwright test が追加・更新されているか。
- `make ci-test` の対象範囲を壊していないか。
- Docker Compose や port、env 変更が README/Makefile/CI と矛盾していないか。

Security / Ops:

- JWT secret、DB password、証明書、DuckDNS token などの秘密情報を commit していないか。
- 本番 Nginx、deploy、migration の順序を壊していないか。
- graceful shutdown や timeout の設計を壊していないか。

## このプロジェクトで避けるべき実装

- Go の大型 Web framework へ無断で置き換えること。現在の設計意図は `net/http` ベース。
- ORM を導入して repository 層を全面的に置き換えること。
- handler から直接 SQL を実行すること。
- service から HTTP response を書くこと。
- repository から HTTP status を直接扱うこと。HTTP status への変換は `apperror` と handler に寄せる。
- route ごとにバラバラの JSON エラー形式を返すこと。
- JWT や password をログ出力すること。
- DB volume を消す前提のスキーマ変更を行うこと。既存 DB 更新は migration で行う。
- E2E テストで不安定な CSS selector や固定データ名だけに依存すること。
- 依存関係の大幅更新、Docker/CI/デプロイ構成の変更を、アプリ機能変更と同じ PR に混ぜること。
- 生成物や大きなバイナリを不用意に追加すること。既にルートに `main` バイナリが存在するため、新たな build artifact の追加には注意する。

## AI がコードを生成・修正するときの注意事項

- まず `AGENTS.md` を読み、作業対象に応じて `docs/agents/*.md` を読む。
- README、Makefile、対象 package の既存実装、関連テストを読む。
- 既存の責務分離を保ち、最小範囲で変更する。
- 存在しないルールを「既存ルール」として書かない。迷う場合は「推奨」と明記する。
- Go 変更後は `gofmt` / `goimports` を前提に整形する。
- フロントエンド変更後は ESLint と Vite build の影響を確認する。
- DB スキーマ変更では、migration、初期化 SQL、テスト SQL、repository、tests の整合を必ず確認する。
- 認証・認可・所有者チェックの変更は、正常系よりも先に抜け漏れリスクを見る。
- エラー文言を変更すると backend tests、frontend 表示、E2E に影響する可能性がある。
- `AuthMiddleware` と `respondAppError` はエラー形式が完全には統一されていない。統一する場合は互換性とテストを含む明示的なリファクタとして扱う。
- `postIDFromRequest` は URL 文字列から ID を取り出す実装なので、route 形状を変える場合は `mux.Vars` 利用への移行も含めて影響を確認する。
- `workerpool.Enqueue` は非ブロッキングで、queue full はリクエスト失敗にしない設計。監査ログの失敗を本処理の失敗に変えない。
- `docs/backlog.md` にある課題は現行仕様ではなく、既知の改善候補として扱う。
