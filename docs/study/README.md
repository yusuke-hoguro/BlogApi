# 学習ドキュメント

`docs/study` は、AIエージェントと進める学習の運用資料、テーマごとのPlan・進捗、学習中に得た技術メモを管理するディレクトリです。

## 学習を始める

新しい学習を始めるときは、最初に[学習Planの運用方法](./learning-workflow.md)を確認してください。

運用に使用する原本は次の2ファイルです。テンプレートへ直接記入せず、テーマごとのファイルへ複製して使用します。

- [学習Planテンプレート](./templates/learning-plan-template.md)
- [学習進捗管理テンプレート](./templates/learning-progress-template.md)

## ディレクトリ構成

```text
docs/study/
├── README.md
├── learning-workflow.md
├── templates/
│   ├── learning-plan-template.md
│   └── learning-progress-template.md
├── plans/                    # テーマ専用の学習Plan
├── progress/                 # テーマ専用の進捗
└── notes/                    # 分野別の技術メモ
    ├── go/
    ├── api/
    ├── testing/
    └── infrastructure/
```

`plans/` と `progress/` は、最初のテーマ専用ファイルを作成するときに追加します。

## ファイルの役割と配置

| 種類 | 役割 | 配置 |
| --- | --- | --- |
| 学習運用 | AIとの学習手順を定義する | `docs/study/learning-workflow.md` |
| テンプレート | Planと進捗の再利用可能な原本 | `docs/study/templates/` |
| 学習Plan | 何を、なぜ、どの順番で学ぶかを定義する | `docs/study/plans/<learning-id>-plan.md` |
| 学習進捗 | 現在地、理解の証拠、次の行動を記録する | `docs/study/progress/<learning-id>-progress.md` |
| 技術メモ | 学習中に確認した知識や手順を分野別に残す | `docs/study/notes/<分野>/` |
| 設計成果物 | 学習結果として確定したシステム設計を残す | `docs/architecture/<領域>/` |

## 技術メモ

### Go

- [ゴルーチンとチャネルの安全性](./notes/go/concurrency.md)

### API

- [OpenAPI仕様書作成](./notes/api/openapi.md)

### Testing

- [Playwright環境構築手順](./notes/testing/playwright.md)

### Infrastructure

- [GitHubへのEC2設定](./notes/infrastructure/ec2-setting.md)
- [Nginxでのリバースプロキシ設定](./notes/infrastructure/reverse-proxy.md)

## 新しいドキュメントの追加ルール

- 学習IDとファイル名には、小文字のkebab-caseを使用する。
- 1つの成果物に対して、同じ学習IDのPlanと進捗ファイルを作成する。
- 技術メモは `notes/go`、`notes/api`、`notes/testing`、`notes/infrastructure` の該当分野へ追加する。
- 学習途中の状態は進捗ファイルに記録し、完成した設計書には混在させない。
- 新しい技術メモを追加したら、このREADMEの索引も更新する。
