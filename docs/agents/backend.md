# Backend Guide

Backend API や Go 実装を変更する前に読む資料です。handler、service、repository、middleware、認証・認可、エラー処理、ログ処理のルールをまとめます。

## レイヤー間の依存関係

守るべき依存方向:

- `handler` は `service` に依存する。
- `service` は `repository` に依存する。
- `repository` は `models` と `apperror` に依存し、DB 操作を閉じ込める。
- `router` は handler と middleware を組み合わせる。
- `app` は repository と service の生成をまとめる。
- `models` は下位の共通データ構造として扱い、handler/service/repository へ逆依存しない。

避ける依存:

- `repository` から `handler` や `middleware` を import しない。
- `service` で HTTP request/response を扱わない。
- `handler` に SQL を直接書かない。
- フロントエンドから DB やバックエンド内部構造を前提にした処理を書かない。

## API / 認可マップ

公開 API:

- `GET /api/healthz` / `HEAD /api/healthz`
- `GET /api/posts`
- `GET /api/posts/{id}`
- `POST /api/signup`
- `POST /api/login`
- `GET /api/posts/{id}/comments`
- `GET /api/comments/{id}`
- `GET /api/posts/{id}/likes`
- `/swagger/` 配下の Swagger UI

認証必須 API:

- `POST /api/posts`
- `PUT /api/posts/{id}`
- `DELETE /api/posts/{id}`
- `GET /api/myposts`
- `POST /api/posts/{id}/comments`
- `PUT /api/comments/{id}`
- `DELETE /api/comments/{id}`
- `POST /api/posts/{id}/like`
- `DELETE /api/posts/{id}/like`

認可の境界:

- 投稿の更新・削除は `PostService.EnsurePostOwner` で投稿者本人のみ許可する。
- コメントの更新・削除は `CommentService.EnsureCommentOwner` でコメント作成者本人のみ許可する。
- いいね追加・削除はログインユーザー自身の `user_id` と対象 `post_id` の組み合わせで行う。投稿所有者チェックはしない。

## コーディング規約

- Go の標準フォーマットを使う。`gofmt` / `goimports` は `.golangci.yml` の formatter として有効。
- Lint は `golangci-lint`。有効 linters は `govet`, `errcheck`, `staticcheck`, `unused`, `gocritic`。
- context は `r.Context()` から受け取り、DB 呼び出しでは `QueryContext`, `QueryRowContext`, `ExecContext` を使う。
- HTTP handler は `func XxxHandler(service *service.XxxService, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc` の既存パターンに合わせる。
- JSON レスポンスは `respondJSON`、エラーレスポンスは `respondAppError` を使う。
- 入力検証は handler 層の `validation.go` か近い共通関数に集約する。
- SQL は repository 層に閉じ込め、プレースホルダ `$1`, `$2` を使う。
- 複数テーブル更新が必要な場合は `BeginTx` を使い、rollback defer と `sql.ErrTxDone` チェックの既存パターンに合わせる。

## エラーハンドリング方針

実装済みの方針:

- アプリケーションエラーは `apperror.NewAppError(type, message, cause)` で生成する。
- `TypeBadRequest`, `TypeUnauthorized`, `TypeForbidden`, `TypeNotFound`, `TypeConflict`, `TypeTimeout`, `TypeInternalServer`, `TypeMethodNotAllowed` を HTTP ステータスへ変換する。
- repository では `sql.ErrNoRows` を `TypeNotFound` に変換する。
- DB 由来などの内部エラーは `TypeInternalServer` とし、cause を `Err` に保持する。
- handler は `respondAppError` でログ出力し、クライアントへは `{"message": "..."}` 形式で返す。

注意:

- 認証 middleware は現在 `http.Error` を直接返している。新規実装では既存挙動との互換を優先し、全体統一を行う場合はテストとフロントエンドの期待値を確認する。
- クライアントに DB エラー詳細や stack trace を返さない。
- `context.Context` のキャンセル・タイムアウトを握りつぶさない。必要なら `apperror.TypeTimeout` などに変換する。

## ログ出力方針

実装済みの方針:

- Backend は標準 `log` パッケージを使う。
- サーバー起動、shutdown、監査イベント、アプリケーションエラー、予期しないエラーをログ出力する。
- 監査イベントは `workerpool.AuditWorkerPool` に非同期 enqueue し、queue full や closed はリクエスト失敗にせずログに残す。

推奨:

- パスワード、JWT、生の Authorization header、DB 接続情報などの秘密情報はログに出さない。
- 新しい監査対象操作を追加したら、成功レスポンス後に `enqueueAuditEvent` を呼ぶ。
- 監査イベント名は既存の `post_created`, `post_updated` のような snake_case の action 名に合わせる。

## 新規 Backend API 追加時のルール

1. `models` に必要な構造体を追加・更新する。
2. `repository` に SQL 操作を追加する。
3. `service` にユースケース・認可ロジックを追加する。
4. `handler` に request decode、validation、service 呼び出し、response を実装する。
5. `router.RegisterRoutes` に route を登録する。
6. 認証が必要な route は `middleware.AuthMiddleware` で包む。
7. 成功操作に必要な監査イベントを `enqueueAuditEvent` で追加する。
8. DB スキーマ変更がある場合は `docs/agents/database.md` を読み、migration と初期化 SQL の整合を取る。
9. handler テストを追加し、正常系と主要な異常系を確認する。
10. Swagger コメントと `docs/swagger.*` の更新が必要か確認する。
