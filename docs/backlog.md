# BlogApi 実務レベル改善バックログ

作成日: 2026-07-01

このバックログは、既存の `docs/issue/TODO_LIST.md` を確認し、重複を整理して統合したものです。既存Todoの扱いは末尾の「既存Todo統合メモ」に記載しています。

## 現状採点

実務レベル評価: **62 / 100**

Go標準 `net/http`、Repository/Service/Handler分離、JWT認証、PostgreSQL、Docker Compose、Playwright E2E、GitHub Actions、EC2/Nginx/Certbot運用まで入っており、ポートフォリオとしての材料はかなり良いです。一方で、JWT秘密鍵のコード固定、ページングなし一覧API、DBインデックス/制約不足、本番composeにpgAdmin、構造化ログ/監視不足、マイグレーション運用の弱さ、テスト品質のばらつきがあり、実務評価ではここが減点になります。

## 80点以上を目指す優先方針

1. **セキュリティと設定管理を先に固める**: JWT秘密鍵、認証エラー、CORS、Rate Limit、Secrets、prod composeの公開範囲を直す。
2. **API/DBの実務品質を上げる**: ページング、ソート、検索、統一エラー、DB制約、インデックス、トランザクション整合性を整える。
3. **テストを仕様の証拠にする**: Handlerの正常/異常だけでなく、Repository/Service、認可、E2E、マイグレーション、CI品質ゲートを強化する。
4. **運用できる証拠を作る**: 構造化ログ、メトリクス、ヘルスチェック、バックアップ/リストア、デプロイ/ロールバック手順をドキュメントとCIで示す。

## 30分でできる改善

- JWT秘密鍵の環境変数化
- 認証エラーのJSON形式統一
- Delete APIのSwaggerメソッド誤記修正
- CORS許可Originの環境変数化
- `comments` テーブルの `IF NOT EXISTS` 追加
- 代表的なDBインデックス追加のマイグレーション作成
- CI権限の最小化
- 本番composeからpgAdminを外す
- `.env.example` 作成
- Docker/Goバージョン不整合修正

## 休日に腰を据えて取り組む改善

- ページング/検索/ソートを含む一覧API再設計
- 認証/認可のService層設計見直し
- DBマイグレーションのup/down対応と本番ロールバック設計
- ログ/メトリクス/アラートの運用設計
- E2Eの安定化とテストデータ管理
- Docker本番構成の分離、非root化、イメージ脆弱性スキャン
- バックアップ/リストア訓練
- OpenAPIを契約としてCI検証

## 改善バックログ

| No | タイトル | 内容 | なぜ実務で評価されるのか | 優先度 | 想定工数 | ステータス | 実装順 |
|---:|---|---|---|---|---|---|---:|
| 1 | JWT秘密鍵を環境変数化する | `middleware.JwtKey` の固定値を廃止し、起動時に `JWT_SECRET` 必須チェックを行う。 | 秘密情報をコードに置かない基本ができていることを示せる。 | High | 30分 | Todo | 1 |
| 2 | JWT署名方式と期限を検証する | `alg` がHS256以外なら拒否し、claims構造体で `exp` を明示検証する。 | JWTの典型的な実装ミスを避けられる。 | High | 2時間 | Todo | 2 |
| 3 | 認証エラーをJSON形式に統一する | `http.Error` をやめ、既存の `AppError` / `respondAppError` に統一する。 | API利用者が一貫してエラー処理できる。 | High | 30分 | Todo | 3 |
| 4 | CORSを環境別に制御する | `AllowedOrigins` を環境変数化し、本番で不要Originを許可しない。 | 本番運用を意識したセキュリティ設計として評価される。 | High | 30分 | Todo | 4 |
| 5 | 本番composeからpgAdminを除外する | `docker-compose.prod.yml` からpgAdminを削除し、dev専用に限定する。 | 本番公開面を減らす判断ができることを示せる。 | High | 30分 | Todo | 5 |
| 6 | API一覧にページングを導入する | `GET /api/posts` と `GET /api/myposts` に `limit` / `cursor` を追加する。 | データ増加時に破綻しないAPI設計になる。 | High | 半日 | Todo | 6 |
| 7 | 投稿一覧に安定したソートを導入する | `created_at DESC, id DESC` のように同時刻でも順序が揺れないSQLにする。 | E2Eや本番UIで再現性のある一覧表示になる。 | High | 30分 | Todo | 7 |
| 8 | DBインデックス最適化 | `posts(created_at,id)`, `posts(user_id,created_at)`, `comments(post_id,created_at)`, `likes(post_id)` などを追加する。 | パフォーマンス劣化をDB設計で防げる。既存Todoから統合。 | High | 2時間 | Todo | 8 |
| 9 | posts.user_idに外部キーを追加する | `posts.user_id REFERENCES users(id)` を追加し、孤児投稿を防ぐ。 | データ整合性をDBで担保できる。 | High | 2時間 | Todo | 9 |
| 10 | post_statsの集計整合性を直す | コメント/いいね作成削除時に `post_stats` を同一トランザクションで更新する。 | 集計値のズレを防ぎ、業務データの信頼性が上がる。 | High | 半日 | Todo | 10 |
| 11 | いいね作成を冪等かつ結果が分かる設計にする | `ON CONFLICT DO NOTHING` の影響行数を見て、作成済みか新規作成かを返す。 | クライアントが状態を正確に扱える。 | High | 2時間 | Todo | 11 |
| 12 | コメント作成時に投稿存在確認を明示する | FK違反を500にせず、存在しない投稿は404へ変換する。 | DBエラーをAPI仕様に正しくマッピングできる。 | High | 2時間 | Todo | 12 |
| 13 | 更新APIのレスポンスをDB結果から返す | `UPDATE ... RETURNING` を使い、更新後の値を返す。 | APIレスポンスと永続化状態のズレを防ぐ。 | High | 2時間 | Todo | 13 |
| 14 | `updated_at` を全主要テーブルに追加する | posts/comments/usersなどに `updated_at` を追加し、更新時に反映する。 | 実務で必要な監査・同期・差分検知に使える。 | High | 2時間 | Todo | 14 |
| 15 | マイグレーションをup/down対応にする | 独自migrationにロールバック方針を追加するか、goose/migrate等を採用する。 | 本番変更の戻し方を説明できる。 | High | 1日以上 | Todo | 15 |
| 16 | 初期化SQLとmigrationの責務を分離する | `sql/init.sql` と `sql/migrations` の二重管理を整理し、スキーマの正をmigrationに寄せる。 | 環境差分による事故を減らせる。 | High | 半日 | Todo | 16 |
| 17 | 設定管理をConfig構造体に集約する | DB、JWT、CORS、worker、timeout、portを `config.Load()` で読み込み検証する。 | 起動時に設定不備を検知でき、保守性が上がる。 | High | 半日 | Todo | 17 |
| 18 | DB接続プールを設定する | `SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime` を環境値で設定する。 | 負荷時のDB接続枯渇を予防できる。 | High | 30分 | Todo | 18 |
| 19 | DB接続にPing timeoutを設定する | `PingContext` と起動timeoutを使う。 | 起動ハングを防ぎ、障害検知が明確になる。 | High | 30分 | Todo | 19 |
| 20 | Rate Limitを導入する | login/signup/comment/post作成にIPまたはユーザー単位の制限を入れる。 | ブルートフォースやスパム対策として評価される。 | High | 半日 | Todo | 20 |
| 21 | ログを構造化する | 標準 `log` から `log/slog` に移行し、request_id, method, path, status, latencyを出す。 | 運用時に追跡可能なログになる。 | High | 半日 | Todo | 21 |
| 22 | Request IDミドルウェアを追加する | リクエストごとにIDを発行/伝播し、レスポンスヘッダとログへ出す。 | 障害調査の基本ができていることを示せる。 | High | 2時間 | Todo | 22 |
| 23 | アクセスログミドルウェアを追加する | 全APIのステータス、処理時間、ユーザーIDを記録する。 | API運用に必要な可観測性が上がる。 | High | 2時間 | Todo | 23 |
| 24 | ヘルスチェックをreadiness/livenessに分ける | DB接続を含むreadinessとプロセス生存のlivenessを分離する。 | コンテナ運用で正しい再起動判断ができる。 | High | 2時間 | Todo | 24 |
| 25 | OpenAPIと実装のズレを修正する | DeletePostの `@Router` がputになっているなどSwagger注釈を棚卸しする。 | API仕様書を信頼できる成果物にできる。 | High | 30分 | Todo | 25 |
| 26 | OpenAPIをCIで検証する | Swagger生成後の差分チェック、lint、破壊的変更検知をCIに追加する。 | API契約を守る開発フローを示せる。 | High | 半日 | Todo | 26 |
| 27 | リクエスト/レスポンスDTOを分離する | `models.User` を入力/出力に流用せず、passwordをレスポンスへ出さない設計にする。 | 情報漏洩防止とAPI契約の明確化につながる。 | High | 半日 | Todo | 27 |
| 28 | Signupレスポンスからpasswordを除外する | 登録後にハッシュ済み/平文passwordが返らないようDTOを使う。 | セキュリティ観点で重要。 | High | 30分 | Todo | 28 |
| 29 | ログに機密情報を出さない | usernameや原因エラーの出し方を見直し、password/tokenは絶対に出さないルールにする。 | 個人情報・秘密情報管理の意識を示せる。 | High | 2時間 | Todo | 29 |
| 30 | 入力JSONの厳格化 | `DisallowUnknownFields`、Content-Type検証、最大body sizeを追加する。 | 想定外入力を早期に拒否できる。 | High | 2時間 | Todo | 30 |
| 31 | パスワードポリシーを明確化する | 長さ上限、空白、bcryptコスト、ログイン失敗時レスポンスを仕様化する。 | 認証機能の現実性が上がる。 | Medium | 2時間 | Todo | 31 |
| 32 | トランザクションヘルパーを導入する | `WithTx(ctx, func(tx) error)` でrollback/commitを共通化する。 | 複数Repository操作の整合性を保ちやすくなる。 | Medium | 半日 | Todo | 32 |
| 33 | RepositoryのDB依存型を揃える | トランザクション対応しやすいよう `DBExecutor` 利用方針を統一する。 | トランザクション境界をServiceで制御しやすい。 | Medium | 2時間 | Todo | 33 |
| 34 | Handlerをinterface依存にする | Handlerは具象Serviceではなく小さなinterfaceに依存する。 | Goらしい小さなinterface設計を示せる。 | Medium | 半日 | Todo | 34 |
| 35 | Handlerのメソッドチェック重複を削除する | muxのMethodsに任せ、Signup/Login内の重複チェックを整理する。 | 責務が明確になり保守性が上がる。 | Medium | 30分 | Todo | 35 |
| 36 | N+1を避けた投稿一覧レスポンスを設計する | 投稿一覧にコメント数/いいね数を返す場合はJOIN/Statsを使う。 | 一覧性能を実務目線で説明できる。 | Medium | 半日 | Todo | 36 |
| 37 | 検索APIを追加する | title/contentの検索、将来的には全文検索indexを検討する。 | API機能として実用性が上がる。 | Medium | 半日 | Todo | 37 |
| 38 | ソフトデリート方針を決める | posts/comments/usersを物理削除のままにするか、`deleted_at` を入れるか決める。 | 業務データ保持の設計判断を示せる。 | Medium | 半日 | Todo | 38 |
| 39 | 監査ログを永続化する | workerpoolのaudit eventをDBまたはログ基盤へ保存する。 | 非同期処理を実務的な監査機能に育てられる。 | Medium | 半日 | Todo | 39 |
| 40 | workerpoolの停止時drainを保証する | shutdown時にキュー内イベントを処理してから終了する設計にする。 | graceful shutdownの品質が上がる。 | Medium | 2時間 | Todo | 40 |
| 41 | workerpoolのバックプレッシャー方針を明記する | キュー満杯時にdropするか同期処理へ切り替えるかを実装/ログ化する。 | 障害時の挙動を説明できる。 | Medium | 2時間 | Todo | 41 |
| 42 | Repository単体テストを追加する | posts/comments/likes/usersのSQL結果、not found、unique違反をテストする。 | DB境界の品質を証明できる。 | High | 半日 | Todo | 42 |
| 43 | Service単体テストを追加する | 認可、トランザクション、エラー変換をmock repositoryでテストする。 | ビジネスロジックの仕様が明確になる。 | High | 半日 | Todo | 43 |
| 44 | Middleware単体テストを追加する | JWT正常、期限切れ、署名不正、alg不正、Authorization形式不正をテストする。 | セキュリティ修正の退行を防げる。 | High | 2時間 | Todo | 44 |
| 45 | Handlerテストをtable-drivenに整理する | 重複が多いテストを共通ヘルパー化し、可読性を上げる。 | テスト追加コストが下がる。 | Medium | 半日 | Todo | 45 |
| 46 | テストでレスポンスbodyを検証する | ステータスだけでなくJSON schema/主要フィールド/エラー形式を確認する。 | API仕様の証拠として強くなる。 | High | 半日 | Todo | 46 |
| 47 | テストDB初期化をmigrationベースにする | `testdata/init_test.sql` と本番schemaの差分をなくす。 | テストだけ通るスキーマ事故を防げる。 | High | 半日 | Todo | 47 |
| 48 | E2Eにサインアップ画面を追加する | 既存Todoの「新規登録画面E2E」を追加し、成功/重複/validationを確認する。 | ユーザー導線の品質を示せる。 | Medium | 2時間 | Todo | 48 |
| 49 | E2Eにいいね機能を追加する | UI操作からAPI反映まで確認する。既存Todoから統合。 | フロント/バックの結合品質が上がる。 | Medium | 半日 | Todo | 49 |
| 50 | E2Eテストデータをテストごとに独立させる | global setup依存を減らし、並列実行でも壊れないfixtureにする。 | CIの flaky を減らせる。 | Medium | 半日 | Todo | 50 |
| 51 | CIにrace testを追加する | `go test -race ./...` を追加し、workerpool等の競合を検出する。 | Go実務で重要な並行処理品質を示せる。 | High | 2時間 | Todo | 51 |
| 52 | CIにcoverage閾値を追加する | カバレッジをartifactに置くだけでなく、最低ラインを設定する。 | 品質ゲートとして機能する。 | Medium | 2時間 | Todo | 52 |
| 53 | CI権限を最小化する | backend-testsの `contents: write` を必要なjobだけに限定する。 | GitHub Actionsのセキュリティ意識を示せる。 | High | 30分 | Todo | 53 |
| 54 | CIのDocker layer cacheを導入する | BuildKit/GHA cacheでDocker buildを高速化する。 | 開発体験とCIコストを改善できる。 | Medium | 2時間 | Todo | 54 |
| 55 | 脆弱性スキャンを追加する | Trivy等でGo依存、npm依存、Docker imageをCIスキャンする。 | サプライチェーン対策として評価される。 | High | 半日 | Todo | 55 |
| 56 | PostgreSQL 18移行検証を行う | volume mount、backup/restore、CI/E2E、ロールバックを検証する。既存Todoから統合。 | DBメジャー更新の現実的な運用力を示せる。 | Medium | 1日以上 | Todo | 56 |
| 57 | Dockerfileを非root実行にする | runtime imageで専用ユーザーを作り、バイナリを非rootで起動する。 | コンテナセキュリティの基本。 | High | 2時間 | Todo | 57 |
| 58 | Docker imageのタグを固定する | `nginx:latest` などを具体的なversion/digestへ固定する。 | 予期しない本番差分を防げる。 | High | 30分 | Todo | 58 |
| 59 | DockerfileのGoバージョン不整合を直す | `go.mod` は1.25.0、Dockerfileは1.26-alpineなので方針を揃える。 | ビルド再現性が上がる。 | High | 30分 | Todo | 59 |
| 60 | `.dockerignore` を強化する | `.git`, coverage, frontend node_modules, binary `main` などを除外する。 | build contextを小さくし、秘密情報混入を防ぐ。 | Medium | 30分 | Todo | 60 |
| 61 | `.env.example` を作成する | `.env` 実ファイル前提をやめ、必要な環境変数をサンプル化する。 | セットアップしやすく、秘密情報管理も明確になる。 | High | 30分 | Todo | 61 |
| 62 | AWS容量不足問題の運用対策を文書化する | Docker prune、log rotate、EBS監視、CloudWatch alarmを手順化する。既存Todoから統合。 | 小規模本番運用の現実的な課題に対応できる。 | Medium | 2時間 | Todo | 62 |
| 63 | バックアップ/リストア手順を作る | PostgreSQL dump、復元、定期バックアップ、復旧テストをdocsに追加する。 | データを守る運用設計として高評価。 | High | 半日 | Todo | 63 |
| 64 | デプロイのロールバック手順を整える | 前回イメージ、DB migration失敗時、nginx切替の戻し方を明記する。 | 本番運用で必須の事故対応力を示せる。 | High | 半日 | Todo | 64 |
| 65 | prod deployで成果物を明確化する | CIでimageをbuild/pushし、EC2はpullして起動する形を検討する。 | デプロイ再現性が上がる。 | Medium | 1日以上 | Todo | 65 |
| 66 | フロントのAPIクライアントを型安全にする | OpenAPIから型生成、またはレスポンス型を整備する。 | フロント/バック間の契約ミスを減らせる。 | Medium | 1日以上 | Todo | 66 |
| 67 | エラーコード体系を導入する | `message` だけでなく `code`, `request_id`, `details` を返す。 | クライアント実装と運用調査がしやすい。 | Medium | 半日 | Todo | 67 |
| 68 | APIバージョニング方針を決める | `/api/v1/...` へ移行するか、現状維持の理由をREADMEに書く。 | APIの長期運用を考えられることを示せる。 | Medium | 2時間 | Todo | 68 |
| 69 | 404/405の共通ハンドラを設定する | muxのNotFoundHandler/MethodNotAllowedHandlerでJSONを返す。 | API全体の一貫性が上がる。 | Medium | 30分 | Todo | 69 |
| 70 | Nginxセキュリティヘッダを追加する | HSTS, X-Content-Type-Options, Referrer-Policyなどを本番confに入れる。 | Web公開時の基本防御になる。 | Medium | 30分 | Todo | 70 |
| 71 | TLS更新の監視を追加する | certbot renew成功/失敗と証明書期限を通知する。 | HTTPS運用の継続性を担保できる。 | Medium | 2時間 | Todo | 71 |
| 72 | 設計判断ADRを追加する | net/http採用、database/sql採用、EC2構成、JWT方式などの理由を書く。 | 面接で設計意図を説明しやすくなる。 | Medium | 半日 | Todo | 72 |
| 73 | READMEを採用担当者向けに再構成する | 機能、技術選定、工夫、運用、テスト、今後の改善を短く整理する。 | ポートフォリオとして伝わりやすくなる。 | Medium | 半日 | Todo | 73 |
| 74 | コメントを「なぜ」に寄せて整理する | コードの逐語説明コメントを減らし、設計意図や注意点を残す。 | 読みやすい実務コードに近づく。 | Low | 半日 | Todo | 74 |
| 75 | package責務をdocsに図示する | handler/service/repository/app/config/middlewareの責務を簡潔に書く。 | アーキテクチャ理解がしやすい。 | Low | 2時間 | Todo | 75 |
| 76 | フロントのいいねUIを実装する | API実装済みのいいねをPostList/PostDetailに表示・操作できるようにする。既存Todoから統合。 | APIだけでなくユーザー体験まで完成する。 | Medium | 半日 | Todo | 76 |
| 77 | Goルーチン活用を実務的な非同期処理にする | 既存Todoの「Goルーチン強化」はaudit workerの永続化/安全停止/監視へ統合する。 | ただ並行処理を使うのでなく、運用価値のある用途にできる。 | Medium | 半日 | Todo | 77 |
| 78 | `/api/healthz` のレスポンス仕様を明確化する | status, db, version, build_shaなどを必要範囲で返す。 | デプロイ確認と監視に使いやすい。 | Medium | 2時間 | Todo | 78 |
| 79 | バイナリ成果物をgit管理対象から外す | ルートの `main` のようなビルド成果物を `.gitignore` に追加する。 | リポジトリが軽くなり、差分レビューが健全になる。 | Medium | 30分 | Todo | 79 |
| 80 | 本番運用Runbookを作る | 起動停止、ログ確認、DB接続、デプロイ失敗、容量不足、証明書更新を手順化する。 | 実務運用を想定した完成度になる。 | Medium | 半日 | Todo | 80 |

## 既存Todo統合メモ

| 既存Todo | 統合先 |
|---|---|
| いいね機能（フロント実装） | No.49, No.76 |
| Goルーチンを使用した実装強化 | No.39, No.40, No.41, No.77 |
| 新規登録画面のE2Eテストを追加する | No.48 |
| AWSサーバーでの容量不足問題を解決する | No.62, No.80 |
| DBインデックス最適化 | No.8 |
| PostgreSQL 18への移行検証を行う | No.56 |
| graceful shutdown対応 | 対応済みとして評価済み。追加改善はNo.40 |
| Dependabot運用調整 | 対応済みとして評価済み。追加改善はNo.56 |
