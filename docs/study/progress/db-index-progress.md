# DB Index 学習進捗

## AIエージェントへの再開指示

学習を再開するときは、次の順番で確認してください。

1. `docs/study/plans/db-index-plan.md` を読み、全体のゴールとStepの依存関係を確認する。
2. このファイルの `現在地`、`未解決事項・ブロッカー`、`次回の開始地点` を確認する。
3. 完了済みStepの理解の証拠と、Plan変更履歴を確認する。
4. 学習者に現在地と今回行う1 Stepを短く説明する。
5. `次回最初に行うこと` から再開し、一度に複数Stepへ進まない。

---

## 基本情報

| 項目 | 内容 |
| --- | --- |
| 学習ID | `db-index` |
| 学習テーマ | PostgreSQLのインデックス設計とEXPLAIN / EXPLAIN ANALYZE |
| 対応するPlan | `docs/study/plans/db-index-plan.md` |
| 対象成果物 | `docs/architecture/database/index-design.md` |
| 学習開始日 | `2026-08-09` |
| 最終更新日 | `2026-08-09` |
| 全体ステータス | `進行中` |

## 現在地

| 項目 | 内容 |
| --- | --- |
| 現在のStep | `Step 1: インデックスの基礎と既存インデックス把握` |
| Stepの状態 | `完了（学習者確認待ち）` |
| 現在取り組んでいる内容 | `Step 1の完了条件を満たし、進捗を記録済み。Step 2は未開始` |
| 最後に完了したこと | `PK / UNIQUEによる自動インデックス、複合一意インデックス、B-treeの概要、BlogApiの既存インデックスを自分の言葉で説明` |
| 次回最初に行うこと | `学習者の確認後、明示的な開始指示があればStep 2: BlogApiのSQL分析を開始する` |

## Step進捗一覧

| Step | タイトル | 状態 | 完了日 | 完了根拠の要約 |
| --- | --- | --- | --- | --- |
| 1 | インデックスの基礎と既存インデックス把握 | `完了（確認待ち）` | `2026-08-09` | `PK / UNIQUEと自動indexを説明し、BlogApiの既存indexを正しく洗い出した` |
| 2 | BlogApiのSQL分析 | `未着手` | `-` | `-` |
| 3 | EXPLAIN基礎 | `未着手` | `-` | `-` |
| 4 | 検証用データとベースライン取得 | `未着手` | `-` | `-` |
| 5 | 単一列インデックス設計 | `未着手` | `-` | `-` |
| 6 | 複合インデックス設計 | `未着手` | `-` | `-` |
| 7 | likesと既存UNIQUE index分析 | `未着手` | `-` | `-` |
| 8 | migration実装とEXPLAIN ANALYZE比較 | `未着手` | `-` | `-` |
| 9 | コストとトレードオフ | `未着手` | `-` | `-` |
| 最終 | 復習・応用 | `未着手` | `-` | `-` |

## 現在のStepの詳細

### Step 1: インデックスの基礎と既存インデックス把握

#### 今回の目的

BlogApiに新しいインデックスを追加する前に、既に存在するインデックスを把握する。

#### 実施した作業

- `posts.id` がPRIMARY KEYであるケースから、同じ列への追加indexが不要かを検討した。
- `users.username` のUNIQUE制約によって一意indexが自動作成されることを確認した。
- `likes` の `UNIQUE(user_id, post_id)` が単一indexではなく複合一意indexを作ることを確認した。
- `(user_id, post_id)` と `post_id` 単独検索を比較し、複合indexでは列順が重要であることを確認した。
- `sql/init.sql` の定義を基に、BlogApiで自動作成されているindexをテーブルごとに洗い出した。
- B-treeを、本の索引のように目的値へ効率的に到達するための構造として学習し、値が順序付けされた木構造で管理される点を確認した。

#### 作成・変更した成果物

- 対象成果物 `docs/architecture/database/index-design.md` への反映はまだ行っていない。
- 今回は学習・理解確認のみを実施した。

#### 理解できたこと

- PostgreSQLではPRIMARY KEYを定義すると、その一意性を保証する一意indexが自動作成される。
- UNIQUE制約でも一意indexが自動作成される。
- 複数列のUNIQUE制約では、その列の組み合わせに対する複合一意indexが作成される。
- `(user_id, post_id)` のような複合indexでは列順が重要で、先頭列を条件にしない `post_id` 単独検索には効率よく利用しにくい。
- B-tree indexは値を順序付けされた木構造で管理するため、全行を順番に確認せず目的の値へ効率よく到達できる。
- 既存のPK / UNIQUE indexを確認せず同じindexを追加すると、重複indexを作る可能性がある。

#### 自分の言葉で説明できたこと

- `posts.id` はPRIMARY KEYなので、PostgreSQLが自動でindexを付与するため追加indexは不要と説明できた。
- PRIMARY KEYはNULLを許容せず、一意であることを説明できた。
- `users.username` はUNIQUEなので一意indexが自動作成され、別途indexを作成する必要がないと説明できた。
- `UNIQUE(user_id, post_id)` が複合一意indexになることを、指摘後に理解した。
- `(user_id, post_id)` は `post_id` 単独検索には使いにくく、列順が逆なら `post_id` を先頭にした検索に利用しやすいと説明できた。
- B-tree indexを「本の索引のように、指定値がどこにあるかを効率よく把握する仕組み」というイメージで説明できた。
- BlogApiで自動作成されるindexを、`posts.id`、`users.id`、`users.username`、`comments.id`、`likes.id`、`likes(user_id, post_id)`、`post_stats.post_id` と洗い出せた。

#### AIからの指摘と修正内容

- 指摘: 当初 `UNIQUE(user_id, post_id)` は単独indexであり、別途複合indexが必要と認識していた。
- 修正: 複数列UNIQUEは、その列の組み合わせに対する複合一意indexを自動作成することを確認した。
- 指摘: B-treeの説明は「場所を把握している」という索引のイメージに留まっていた。
- 修正: 値を大小関係に基づいて順序付けされた木構造で管理し、探索対象外の範囲を除外しながら目的値へ到達する点を補足した。

#### まだ曖昧なこと

- B-tree内部構造の詳細は今回のStepでは概要理解に留めた。
- PostgreSQLが実際にSeq Scan / Index Scanのどちらを選択するかは未学習で、後続のEXPLAIN Stepで確認する。
- 複合indexの列順について基本的な感覚は得たが、WHERE + ORDER BYを含む具体的な設計判断は後続Stepで深める。

#### 完了条件の確認

- [x] PKとUNIQUEによる既存indexを説明できる。
- [x] BlogApiで既にindexが存在する列を洗い出せる。
- 確認方法: 資料の完成回答を先に提示せず、質問形式で学習者自身に理由と既存index一覧を説明してもらった。
- 完了根拠: `posts.id` と `users.username` への追加indexが不要な理由を説明でき、複数列UNIQUEの誤解を修正したうえで、BlogApiの既存index一覧を正しく列挙できた。B-treeについても本の索引を用いて自分の言葉で概要を説明できた。
- 学習者の確認: `確認待ち`

#### 次のアクション

`学習者がこの進捗記録を確認する。Step 2は開始しない。学習者から明示的にStep 2開始の指示があった場合のみ、Step 2: BlogApiのSQL分析を開始する。`

## 完了Stepの記録

### Step 1: インデックスの基礎と既存インデックス把握

- 実施日: `2026-08-09`
- 状態: `完了（学習者確認待ち）`
- 完了を確認した方法: PK / UNIQUE / 複合UNIQUE / B-tree / BlogApi既存indexについて質問し、学習者自身の言葉で回答してもらった。
- 完了根拠: Planの2つの完了条件を満たし、途中の誤解も対話内で修正できた。
- 成果物への反映: `docs/architecture/database/index-design.md` への反映はまだ行っていない。

## 理解の記録

### 説明できるようになったこと

- PRIMARY KEYとUNIQUE制約によってPostgreSQLが一意indexを自動作成すること。
- 複数列UNIQUEは複合一意indexになること。
- 複合indexでは列順が重要であること。
- B-tree indexを本の索引に例え、全件走査を避けて目的値へ効率的に到達するイメージ。
- BlogApiですでに存在するPK / UNIQUE由来のindex一覧。

### 重要な設計判断

| 判断 | 採用理由 | 前提・制約 | 参照先 |
| --- | --- | --- | --- |
| PK / UNIQUE由来の既存indexと同じindexを重複して追加しない | PostgreSQLが制約を保証するindexを既に自動作成しているため | PostgreSQLを利用 | `sql/init.sql` |
| 新規index設計前に既存indexを確認する | 重複indexを避け、書き込み・容量コストを不要に増やさないため | 実際のクエリ分析はStep 2以降 | `sql/init.sql` |

### 検討した代替案

| 代替案 | メリット | デメリット | 不採用理由 |
| --- | --- | --- | --- |
| `posts(id)` を追加する | id検索用indexを明示できる | PRIMARY KEY由来のindexと重複する | 既存indexで目的を満たすため |
| `users(username)` を追加する | username検索用indexを明示できる | UNIQUE由来のindexと重複する | 既存indexで目的を満たすため |
| `likes(user_id, post_id)` を追加する | 複合検索に利用できる | UNIQUE制約由来の同じ複合indexと重複する | 既存indexで目的を満たすため |

### 条件が変わった場合の影響

| 条件変更 | 影響する設計 | 対応方針 |
| --- | --- | --- |
| `likes` を `post_id` 単独で頻繁に検索する | `(user_id, post_id)` の既存indexだけでは効率的に絞り込みにくい | Step 2以降で実クエリを分析し、新規index候補として検討する |

## 未解決事項・ブロッカー

| 種別 | 内容 | 解決に必要なこと | 状態 |
| --- | --- | --- | --- |
| `今後の学習` | PostgreSQLがSeq Scan / Index Scanをどう選ぶか | Step 3以降でEXPLAINを使って確認 | `未解決` |
| `今後の学習` | WHERE / ORDER BYを踏まえた複合indexの具体的な列順 | Step 2とStep 6で実クエリを題材に検討 | `未解決` |
| `理解深化` | B-tree内部構造の詳細 | 必要に応じて後続Stepで補足。現Stepの完了には不要 | `保留` |

## Plan変更履歴

| 日付 | 変更内容 | 変更理由 | 後続Stepへの影響 | 学習者確認 |
| --- | --- | --- | --- | --- |
| `2026-08-09` | Planを確定 | 学習順序・範囲について学習者が承認 | なし | `確認済み` |

## 復習・応用の記録

### 資料を見ない説明

- 実施日: `-`
- 説明できたこと: `未実施`
- 説明できなかったこと: `未実施`

### 代替案との比較

- 問い: `未実施`
- 回答の要約: `未実施`
- 再確認事項: `未実施`

### 条件変更問題

- シナリオ: `未実施`
- 回答の要約: `未実施`
- 再確認事項: `未実施`

### 障害対応問題

- 想定障害: `未実施`
- 切り分け方: `未実施`
- 再確認事項: `未実施`

### AIによる確認問題

| # | 問題 | 回答の要約 | 判定 | 補足 |
| --- | --- | --- | --- | --- |
| - | `未実施` | `未実施` | `-` | `-` |

## 今後の学習候補

- GIN / GiST / BRIN / Hash index
- Full-text search
- パーティショニング
- PostgreSQL planner統計の高度なチューニング
- `pg_stat_statements`
- 大規模負荷試験

## 次回の開始地点

### 次回最初に行うこと

`学習者がStep 1の進捗記録を確認する。明示的にStep 2開始の指示があった場合のみ、Step 2: BlogApiのSQL分析の目的・学ぶ論点・完了条件を確認して開始する。`

### 再開前に確認するもの

- `docs/study/plans/db-index-plan.md` の Step 2
- `docs/study/progress/db-index-progress.md`
- BlogApi repository層の実SQL

### AIに引き継ぐ注意事項

- 一度に1 Stepだけ進める。
- 学習者が考える前に完成回答や完成した成果物を提示しない。
- Step 1は完了条件を満たしているが、学習者確認待ち。
- Step 2はまだ開始しない。学習者から明示的な開始指示があるまで待つ。
