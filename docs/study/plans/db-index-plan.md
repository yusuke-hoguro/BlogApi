# DB Index 学習Plan

## 基本情報

| 項目 | 内容 |
| --- | --- |
| 学習ID | `db-index` |
| 学習テーマ | PostgreSQLのインデックス設計とEXPLAIN / EXPLAIN ANALYZE |
| 対象成果物 | `docs/architecture/database/index-design.md` |
| 対応する進捗ファイル | `docs/study/progress/db-index-progress.md` |
| Plan作成日 | `2026-08-09` |
| Plan最終更新日 | `2026-08-09` |
| 全体ステータス | `確定` |

## 入力情報

### 学習の背景

BlogApiの実際のrepository SQLを題材として、PostgreSQLのインデックス設計と実行計画の読み方を習得する。

現在のBackendは `handler -> service -> repository -> database/sql + PostgreSQL` の構成で、SQLはrepository層に閉じ込められている。

実際に次のような検索条件が存在する。

- `posts WHERE id = $1`
- `posts WHERE user_id = $1`
- `posts ORDER BY created_at DESC`
- `comments WHERE post_id = $1 ORDER BY created_at ASC`
- `likes WHERE user_id = $1 AND post_id = $2`
- `likes WHERE post_id = $1`
- `users WHERE username = $1`

### 今回達成したいこと

- SQLから必要なインデックス候補を自分で判断できる。
- PK / UNIQUE制約によって既に作られているインデックスを把握できる。
- 単一列インデックスと複合インデックスを使い分けられる。
- `EXPLAIN` / `EXPLAIN ANALYZE` の主要項目を読める。
- インデックス追加前後の実行計画を比較できる。
- 「なぜこのインデックスを追加したか」をポートフォリオで説明できる。
- 不要なインデックスを増やすデメリットも説明できる。

### 重点的に学びたいこと

- B-treeインデックスの基本
- WHERE / ORDER BY とインデックスの関係
- 複合インデックスの列順
- Seq Scan / Index Scan / Bitmap Heap Scan
- cost / rows / actual time
- `EXPLAIN ANALYZE`
- 少量データでPostgreSQLがSeq Scanを選ぶ理由
- インデックスの書き込みコスト

### 制約・希望

- BlogApiの実コードを教材にする。
- 一度に1 Stepだけ進める。
- 学習者が先に考えてからAIが解説する。
- 完成SQLや最適解を最初から提示しない。
- DB変更はmigrationを使い、`sql/init.sql` と `testdata/init_test.sql` の整合性を保つ。
- DB volume削除を前提にしない。

## Plan作成前の理解度診断

### 自分の言葉で説明できること

- 主キー
- 外部キー
- CRUD
- SQLの基本的なSELECT / INSERT / UPDATE / DELETE
- トランザクションの目的
- `BeginTx / Commit / Rollback`
- DB処理へcontextを渡す意味

### 名前や概要だけ知っていること

- インデックスは検索高速化に利用する
- EXPLAINという仕組みがある

### 未経験または理解が曖昧なこと

- PostgreSQLがインデックスを選択する仕組み
- B-treeの具体的な動作
- 複合インデックスの列順
- WHEREとORDER BYを同時に満たす設計
- Seq Scan / Index Scanの判断
- costの意味
- `EXPLAIN ANALYZE` の読み方
- インデックスが逆効果になるケース

### 実際に操作・実装したことがあるもの

- PostgreSQLでBlogApiのCRUDを実装
- repository層で`QueryContext`等を使用
- `posts`と`post_stats`をトランザクションで作成

### 診断用の質問と回答

| # | AIからの質問 | 学習者の回答 | 判定・補足 |
| --- | --- | --- | --- |
| 1 | インデックスは何のために使うか | 検索を高速化するためという概要は理解 | 要確認 |
| 2 | EXPLAINで何を見るか | 名前・概要レベル | 未経験 |
| 3 | 複合インデックスの列順を説明できるか | 理解が曖昧 | 未経験 |

### 診断結果によるPlan調整

- インデックスの概念説明だけで終わらせず、実際のクエリ分析 → 仮説 → EXPLAIN → インデックス追加 → 再検証の順番に分割する。
- 既に理解しているSQL基礎やトランザクションの説明は最小限とする。
- 未経験のEXPLAIN、複合インデックス、planner判断は独立Stepとして細かく扱う。

## 学習ゴール

### 成果物のゴール

`docs/architecture/database/index-design.md` に以下を残す。

- BlogApiの主要クエリ一覧
- 既存インデックス一覧
- インデックス候補と採用理由
- 採用しなかった候補と理由
- EXPLAIN / EXPLAIN ANALYZEの比較結果
- インデックス追加によるメリット・デメリット
- 今後データ量が増えた場合の再評価方針

### 理解のゴール

- PK / UNIQUE制約がインデックスを作る理由を説明できる。
- WHERE句とORDER BYから候補列を考えられる。
- 単一列と複合インデックスを比較できる。
- 複合インデックスの先頭列が重要な理由を説明できる。
- PostgreSQLがSeq Scanを選ぶ場合がある理由を説明できる。
- EXPLAINの結果からインデックスが使われたか判断できる。
- インデックスを増やし過ぎる問題を説明できる。

### Plan全体の完了条件

- [ ] 対象成果物が完成している。
- [ ] 主要な設計判断と採用理由を、資料を見ずに説明できる。
- [ ] 代替案と不採用理由を説明できる。
- [ ] 前提条件が変わった場合の影響を説明できる。
- [ ] 代表的な障害の原因と切り分け方を説明できる。
- [ ] 復習・応用Stepを完了している。
- [ ] 未解決事項が、今後の学習候補として整理されている。

## 学習範囲

### 対象

- PostgreSQL B-treeインデックス
- PRIMARY KEY / UNIQUEと自動インデックス
- WHERE
- ORDER BY
- 単一列インデックス
- 複合インデックス
- EXPLAIN
- EXPLAIN ANALYZE
- Seq Scan
- Index Scan
- Bitmap Index / Bitmap Heap Scan
- migrationによるインデックス追加
- BlogApiのPost / Comment / Likeクエリ

### 対象外

- GIN
- GiST
- BRIN
- Hash index
- Full-text search
- パーティショニング
- PostgreSQL planner統計の高度なチューニング
- `pg_stat_statements`
- 大規模負荷試験

## 参照対象

### リポジトリ内

- `AGENTS.md`
- `docs/agents/overview.md`
- `docs/agents/database.md`
- `docs/agents/backend.md`
- `docs/agents/testing.md`
- `docs/agents/review.md`
- `internal/repository/post_repository.go`
- `internal/repository/comment_repository.go`
- `internal/repository/like_repository.go`
- `internal/repository/user_repository.go`
- `sql/init.sql`
- `sql/migrations/`
- `testdata/init_test.sql`

### 外部資料

- PostgreSQL公式ドキュメント（Indexes / Using EXPLAIN）

> 外部情報を使用する場合は、可能な限り一次情報・公式ドキュメントを優先し、参照日を記録する。

## Step一覧

| Step | タイトル | 主な目的 | 成果物への反映 | 完了確認 | 状態 |
| --- | --- | --- | --- | --- | --- |
| 1 | インデックスの基礎と既存インデックス把握 | PK・UNIQUE・B-treeの基礎を理解する | 既存インデックス一覧 | 自動生成されるインデックスを説明 | `未着手` |
| 2 | BlogApiのSQL分析 | どの検索・ソートが候補か考える | クエリ一覧 | 自分で候補列を判断 | `未着手` |
| 3 | EXPLAIN基礎 | 実行計画の読み方を理解する | EXPLAIN結果 | Seq Scan等を説明 | `未着手` |
| 4 | 検証用データとベースライン取得 | インデックス追加前の計画を測る | Before結果 | 少量/大量データ差を説明 | `未着手` |
| 5 | 単一列インデックス設計 | posts検索・ソートを題材に設計する | 候補と比較 | 採否理由を説明 | `未着手` |
| 6 | 複合インデックス設計 | commentsを題材に列順を理解する | 複合index設計 | `(post_id, created_at)`等を説明 | `未着手` |
| 7 | likesと既存UNIQUE index分析 | 既存複合indexの利用可能範囲を理解する | Like設計 | 左端列ルールを説明 | `未着手` |
| 8 | migration実装とEXPLAIN ANALYZE比較 | 実装前後を実測する | migration / After結果 | plan差分を説明 | `未着手` |
| 9 | コストとトレードオフ | 読み込み高速化と書込みコストを理解する | 設計判断 | 不要indexを判断 | `未着手` |
| 最終 | 復習・応用 | 理解の定着と応用 | 成果物全体 | 条件変更・障害対応 | `未着手` |

---

## Step詳細

### Step 1: インデックスの基礎と既存インデックス把握

#### 目的

BlogApiに新しいインデックスを追加する前に、既に存在するインデックスを把握する。

#### 学ぶ内容・重要な論点

- インデックスとは何か
- B-treeの概要
- PKとUNIQUE
- PostgreSQLが自動作成するindex
- index追加前に既存indexを確認する重要性

#### 必要な前提知識

- 主キー、UNIQUE制約、基本的なSELECT

#### 調査対象・参照資料

- `sql/init.sql`
- PostgreSQL公式ドキュメント

#### 最初に考えること

- `posts WHERE id = $1` に自分でindexを追加する必要があるか。
- `users WHERE username = $1` はどうか。
- `likes WHERE user_id=$1 AND post_id=$2` はどうか。

#### 演習・作業

- BlogApiで既に存在するindexを洗い出す。
- PK / UNIQUEによって自動作成されるものと、自分で追加するものを分類する。

#### 成果物への反映

- 対象箇所: `docs/architecture/database/index-design.md` の既存インデックス一覧
- 反映内容: 既存indexと作成理由

#### AIの確認観点

- PK / UNIQUEとindexの関係を正しく説明できるか。
- 「検索対象だから追加する」ではなく既存indexを確認できているか。

#### 学習者が説明できるようになること

- PKとUNIQUEによって自動的にindexが作成される理由。

#### 完了条件と確認証拠

- [ ] PKとUNIQUEによる既存indexを説明できる。
- [ ] BlogApiで既にindexが存在する列を洗い出せる。
- 確認証拠: 自分の言葉での説明と既存index一覧。

#### 理解できなかった場合

- 補足する内容: B-treeと制約によるindex生成。
- 再確認方法: BlogApiの各PRIMARY KEY / UNIQUEについて再分類する。

### Step 2: BlogApiのSQL分析

#### 目的

「列にindexを貼る」ではなく、実際のクエリからindexを設計する考え方を身につける。

#### 学ぶ内容・重要な論点

- WHERE対象列
- ORDER BY対象列
- 既存indexとの重複
- index候補の洗い出し

#### 必要な前提知識

- Step 1

#### 調査対象・参照資料

- `internal/repository/post_repository.go`
- `internal/repository/comment_repository.go`
- `internal/repository/like_repository.go`
- `internal/repository/user_repository.go`

#### 最初に考えること

- 各SQLのWHERE対象、ORDER BY対象、既存index、新規index候補を自分で分類する。

#### 演習・作業

- repository内の主要SQLを一覧化する。
- クエリごとにindex候補を仮置きする。

#### 成果物への反映

- 対象箇所: 主要クエリ一覧
- 反映内容: SQL、用途、検索条件、ソート条件、既存index、候補index

#### AIの確認観点

- 実際のクエリを根拠に候補を出しているか。
- 不要なindexを無条件で追加していないか。

#### 学習者が説明できるようになること

- なぜクエリからindex設計を始めるのか。

#### 完了条件と確認証拠

- [ ] SQLごとのindex候補を自分で挙げられる。
- [ ] 「とりあえず全列に貼る」が不適切な理由を説明できる。
- 確認証拠: クエリ分析表と説明。

#### 理解できなかった場合

- 補足する内容: WHERE / ORDER BYとindex候補の関係。
- 再確認方法: postsの3クエリだけに絞って再分析する。

### Step 3: EXPLAIN基礎

#### 目的

index追加前に、PostgreSQLが現在どう実行しているかを読めるようにする。

#### 学ぶ内容・重要な論点

- `EXPLAIN`
- Seq Scan
- Index Scan
- Bitmap Heap Scan
- Sort
- cost
- rows
- width

#### 必要な前提知識

- Step 1〜2

#### 調査対象・参照資料

- PostgreSQL公式ドキュメント（Using EXPLAIN）

#### 最初に考えること

- EXPLAIN結果のどこを見れば、index利用の有無と処理の大まかな流れが分かるか予想する。

#### 演習・作業

- 主要クエリへ`EXPLAIN`を付け、実行計画を読む。

#### 成果物への反映

- 対象箇所: EXPLAINの読み方
- 反映内容: 主要nodeと読み方

#### AIの確認観点

- Seq ScanとIndex Scanの違いを説明できるか。
- costを実測時間と誤解していないか。

#### 学習者が説明できるようになること

- EXPLAINが何を示しているか。

#### 完了条件と確認証拠

- [ ] Seq ScanとIndex Scanの違いを説明できる。
- [ ] costが実測時間そのものではないことを説明できる。
- 確認証拠: EXPLAIN結果を自分の言葉で読み解いた説明。

#### 理解できなかった場合

- 補足する内容: plannerと実行計画node。
- 再確認方法: PK検索と全件検索の2本だけ比較する。

### Step 4: 検証用データとベースライン取得

#### 目的

少量データではindexがあってもSeq Scanが選ばれる場合があることを体験し、比較可能なBefore結果を作る。

#### 学ぶ内容・重要な論点

- plannerがSeq Scanを選ぶ理由
- データ件数
- 選択性
- `EXPLAIN ANALYZE`

#### 必要な前提知識

- Step 3

#### 調査対象・参照資料

- 開発DB
- repositoryの主要SELECT

#### 最初に考えること

- なぜ数件しかないテーブルではindexより全件走査の方が安い場合があるか。

#### 演習・作業

- 検証データを十分用意する。
- index追加前の`EXPLAIN` / `EXPLAIN ANALYZE`を取得する。

#### 成果物への反映

- 対象箇所: Before検証結果
- 反映内容: データ件数、実行計画、主要値

#### AIの確認観点

- 少量データと大量データでplanner判断が変わり得ることを理解しているか。

#### 学習者が説明できるようになること

- indexが存在しても必ず使われるわけではない理由。

#### 完了条件と確認証拠

- [ ] 少量データでSeq Scanが選ばれる理由を説明できる。
- [ ] Beforeの実行計画を取得できる。
- 確認証拠: Before結果と解説。

#### 理解できなかった場合

- 補足する内容: plannerのコスト比較。
- 再確認方法: 同一クエリをデータ量を変えて比較する。

### Step 5: 単一列インデックス設計

#### 目的

postsの検索・ソートを題材に、単一列indexの採否を判断する。

#### 学ぶ内容・重要な論点

- 単一列index
- 選択性
- ORDER BYとB-tree
- ASC / DESC

#### 必要な前提知識

- Step 1〜4

#### 調査対象・参照資料

- `internal/repository/post_repository.go`

#### 最初に考えること

- `posts(user_id)` と `posts(created_at)` は本当に必要か。

#### 演習・作業

- 候補indexごとに期待する効果とコストを整理する。

#### 成果物への反映

- 対象箇所: postsのindex設計
- 反映内容: 候補、採用/不採用、理由

#### AIの確認観点

- クエリとデータ特性を根拠に判断しているか。

#### 学習者が説明できるようになること

- 単一列indexを追加する判断基準。

#### 完了条件と確認証拠

- [ ] 各indexの採用/不採用理由を説明できる。
- 確認証拠: 設計判断と理由。

#### 理解できなかった場合

- 補足する内容: 選択性とソートコスト。
- 再確認方法: user_id検索とcreated_atソートを個別に再検討する。

### Step 6: 複合インデックス設計

#### 目的

commentsの`WHERE post_id = $1 ORDER BY created_at ASC`を題材に、複合indexと列順を理解する。

#### 学ぶ内容・重要な論点

- 複合index
- 列順
- WHERE + ORDER BY
- `(post_id, created_at)`
- 単一index2本との違い

#### 必要な前提知識

- Step 5

#### 調査対象・参照資料

- `internal/repository/comment_repository.go`

#### 最初に考えること

- `(post_id)`、`(created_at)`、`(post_id, created_at)`、`(created_at, post_id)`を比較する。

#### 演習・作業

- 各候補がどの部分を高速化できるか整理する。

#### 成果物への反映

- 対象箇所: commentsの複合index設計
- 反映内容: 列順と採用理由

#### AIの確認観点

- WHERE条件とORDER BYの両方を考慮できているか。

#### 学習者が説明できるようになること

- 複合indexの列順が重要な理由。

#### 完了条件と確認証拠

- [ ] 複合indexの列順を説明できる。
- [ ] 単一index2本との違いを説明できる。
- 確認証拠: 候補比較と説明。

#### 理解できなかった場合

- 補足する内容: B-treeの並び順と左端列。
- 再確認方法: 4候補の利用可能クエリを表にする。

### Step 7: likesと既存UNIQUE index分析

#### 目的

既存の`UNIQUE(user_id, post_id)`を題材に、複合indexの利用可能範囲と重複index回避を理解する。

#### 学ぶ内容・重要な論点

- 複合indexの左端列
- `(user_id, post_id)`が使える検索
- `post_id`単独検索との関係
- 重複index

#### 必要な前提知識

- Step 6

#### 調査対象・参照資料

- `internal/repository/like_repository.go`
- `sql/init.sql`

#### 最初に考えること

- `WHERE user_id=$1 AND post_id=$2`と`WHERE post_id=$1`の両方を既存UNIQUE indexだけで十分に処理できるか。

#### 演習・作業

- 既存UNIQUE indexで利用できるクエリを分類する。

#### 成果物への反映

- 対象箇所: likesのindex設計
- 反映内容: 既存indexの利用範囲と不足箇所

#### AIの確認観点

- 左端列ルールを丸暗記ではなくクエリと結び付けて説明できるか。

#### 学習者が説明できるようになること

- 既存複合indexで十分な検索と不十分な検索。

#### 完了条件と確認証拠

- [ ] 既存UNIQUE indexで十分なクエリと、不十分なクエリを説明できる。
- 確認証拠: Likeクエリの分類。

#### 理解できなかった場合

- 補足する内容: 複合B-treeの左端列。
- 再確認方法: `(user_id, post_id)`で利用可能なWHERE条件を列挙する。

### Step 8: migration実装とEXPLAIN ANALYZE比較

#### 目的

採用したindexをBlogApiへ導入し、実装前後を実測する。

#### 学ぶ内容・重要な論点

- migrationによるindex追加
- `sql/init.sql`との整合
- `testdata/init_test.sql`との整合
- Before / After比較

#### 必要な前提知識

- Step 1〜7

#### 調査対象・参照資料

- `sql/migrations/`
- `sql/init.sql`
- `testdata/init_test.sql`
- `docs/agents/database.md`

#### 最初に考えること

- 採用したindexを既存DB、新規DB、テストDBへどう整合させるか。

#### 演習・作業

- migrationを追加する。
- init.sqlとtest SQLを更新する。
- `EXPLAIN ANALYZE`のBefore / Afterを比較する。

#### 成果物への反映

- 対象箇所: migration、After検証結果
- 反映内容: 実装内容と比較結果

#### AIの確認観点

- DB volume削除に頼らずmigrationで更新できているか。
- plannerがindexを利用したか正しく判定できるか。

#### 学習者が説明できるようになること

- 実測結果を根拠にindex追加の効果を説明する方法。

#### 完了条件と確認証拠

- [ ] migrationでindexを追加できる。
- [ ] Before / Afterを比較できる。
- [ ] plannerがindexを使ったか判断できる。
- 確認証拠: SQL差分とEXPLAIN ANALYZE結果。

#### 理解できなかった場合

- 補足する内容: migration運用と実行計画比較。
- 再確認方法: index1本だけに絞ってBefore / Afterを再比較する。

### Step 9: インデックスのコストと設計判断

#### 目的

indexを追加しない判断も含めて、読み込み性能と書き込みコストのtrade-offを理解する。

#### 学ぶ内容・重要な論点

- INSERTコスト
- UPDATEコスト
- DELETEコスト
- disk容量
- maintenanceコスト
- 重複index

#### 必要な前提知識

- Step 1〜8

#### 調査対象・参照資料

- 今回追加したindex一覧
- PostgreSQL公式ドキュメント

#### 最初に考えること

- 「検索が速くなるなら全部の列にindexを貼ればよい」がなぜ誤りか。

#### 演習・作業

- 採用indexごとに得られる利益と増えるコストを整理する。
- 不採用候補の理由を明文化する。

#### 成果物への反映

- 対象箇所: 設計判断とtrade-off
- 反映内容: 採用・不採用理由

#### AIの確認観点

- 読み取り性能だけを見ず、書き込みと容量への影響まで説明できるか。

#### 学習者が説明できるようになること

- indexを追加しない判断。

#### 完了条件と確認証拠

- [ ] indexを追加しない判断も説明できる。
- [ ] 読み取り性能と書き込み性能のtrade-offを説明できる。
- 確認証拠: 採用/不採用判断表。

#### 理解できなかった場合

- 補足する内容: index更新コスト。
- 再確認方法: postsへのINSERT時に更新対象となるindexを列挙する。

---

## 最終Step: 復習・応用

### 資料を見ない説明

- インデックスとは何か。
- BlogApiに追加したindexは何か。
- なぜその列・列順にしたのか。
- なぜ追加しなかったindexがあるのか。
- EXPLAINで何を確認したか。

### 代替案との比較

- `comments(post_id)` + `comments(created_at)` と `comments(post_id, created_at)` を比較する。

### 条件変更問題

- 投稿一覧に `WHERE user_id = ? ORDER BY created_at DESC` が追加されたらどう設計を見直すか。
- コメントを最新100件だけ取得するようになったらどうするか。
- likes検索が `post_id` 中心から `user_id` 中心に変わったらどうするか。

### 障害対応問題

「indexを追加したのにEXPLAINがSeq Scanのまま」の場合に、次の観点で切り分ける。

- データ件数
- 選択性
- planner統計
- query条件
- index列順

### AIによる確認問題

- PK検索へ追加indexが不要な理由を説明する。
- `comments(post_id, created_at)`の列順を逆にした場合の影響を説明する。
- index追加後もSeq Scanになるケースを説明する。
- indexを削除する判断基準を説明する。

### 最終確認

- [ ] 成果物の内容と説明に矛盾がない。
- [ ] 主要なindex設計を資料なしで説明できる。
- [ ] EXPLAIN結果からアクセス方法を判断できる。
- [ ] 条件変更時にindex設計を再検討できる。
- [ ] 未解決事項を今後の学習候補として整理した。

## Plan変更ルール

- Planは学習中の理解度や新しく判明した前提に応じて変更してよい。
- AIエージェントは、変更前に理由と後続Stepへの影響を提示する。
- 学習者の確認後、Planと進捗ファイルの変更履歴を更新する。
- 過去の完了記録は削除せず、再確認が必要な場合はその理由を残す。
