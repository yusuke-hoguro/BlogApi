# Database Guide

DB スキーマ、migration、SQL、テストデータを変更する前に読む資料です。

## DB / Migration

- 新規 DB 初期化は `sql/init.sql`、既存 DB 更新は `sql/migrations/*.sql` を使う。
- マイグレーションは `cmd/migrate` から `internal/db.RunMigrations` を呼び、`schema_migrations` テーブルで適用済みファイルを管理する。
- マイグレーション SQL はファイル名順に実行され、各ファイルはトランザクション内で適用される。
- `sql/init.sql` は新規 DB 初期化用、`sql/migrations` は既存 DB 更新用として扱う。
- テスト DB 初期化は `testdata/init_test.sql` を使う。

## post_stats の現状

- `post_stats` は投稿作成時に `posts` と同一トランザクションで初期行を作る。
- 現状ではコメント数・いいね数・閲覧数の集計更新には使われていない。
- 集計ロジックを追加する場合は既知の改善候補として別途設計し、コメント/いいね作成削除と同一トランザクションで整合性を保つ方針を検討する。

## スキーマ変更時のルール

1. `sql/migrations` に新しい `.sql` を追加する。
2. 新規環境用に `sql/init.sql` も整合させる。
3. テスト用に `testdata/init_test.sql` も整合させる。
4. repository の SQL と models の JSON/DB 構造を確認する。
5. handler/integration test を追加・更新する。

## SQL 実装ルール

- SQL は repository 層に閉じ込める。
- プレースホルダ `$1`, `$2` を使い、文字列連結で SQL を組み立てない。
- 複数テーブル更新が必要な場合は `BeginTx` を使う。
- rollback defer と `sql.ErrTxDone` チェックの既存パターンに合わせる。
- `sql.ErrNoRows` は `apperror.TypeNotFound` に変換する。
- `RowsAffected == 0` は not found として扱う既存 helper に合わせる。

## 避けるべきこと

- Docker volume を消す前提のスキーマ変更を行うこと。
- `sql/init.sql` だけを変えて migration を追加しないこと。
- `testdata/init_test.sql` と本番スキーマの差分を放置すること。
- handler から直接 SQL を実行すること。
