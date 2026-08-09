# Terraform EC2インスタンス学習Plan

このPlanは、TerraformによるEC2インスタンスの設計、実装、検証、振り返りを、AIエージェントと1 Stepずつ進めるためのテーマ専用Planです。

## AIエージェントへの指示

### Plan作成・変更時

- `AGENTS.md`、`docs/agents/operations.md`、学習運用ガイド、進捗ファイル、EC2設計書を先に確認する。
- 実装・設定・既存ドキュメントから確認できる事実、一般的な推奨、推測を区別する。
- 未決定事項を推測で確定せず、該当Stepで学習者と一つずつ決定する。
- Planを変更する場合は、理由と後続Stepへの影響を説明し、学習者の確認後にPlanと進捗へ記録する。

### 学習進行時

- 一度に扱うのは、進捗ファイルが示す現在の1 Stepだけとする。
- Step開始時に、目的、論点、完了条件を説明する。
- 原則として、学習者より先に完成回答、完成した設計書、完成コードを提示しない。
- 最初に学習者へ、予想、説明、判断、調査または草案作成を促す。
- 学習者の回答後に、正しい点、不足、誤解と、その理由をフィードバックする。
- 成果物が作成されただけでは完了にせず、学習者が設計理由やコードの意味を説明できるか確認する。
- 学習者が完了を確認するまで次のStepへ進まない。
- Step完了時とセッション終了時に、進捗ファイルを更新する。
- AWS認証情報、秘密鍵、Terraform Stateなどの秘密情報をファイルや会話へ出力しない。

## 基本情報

| 項目 | 内容 |
| --- | --- |
| 学習ID | `terraform-ec2-instance` |
| 学習テーマ | Terraformを使用したBlogAPI学習・検証用EC2インスタンスの構築 |
| 対象成果物 | `docs/architecture/terraform/ec2-instance.md`、Terraformコード、検証記録、デモ、学習の振り返り |
| 対応する進捗ファイル | `docs/study/progress/terraform-ec2-instance-progress.md` |
| Plan作成日 | `2026-08-09` |
| Plan最終更新日 | `2026-08-09` |
| 全体ステータス | `レビュー待ち` |

## 入力情報

### 学習の背景

TerraformでEC2を作成することだけを目的にせず、Terraformの考え方、設計、レビュー、実装、実行結果の確認を理解し、自分の言葉で説明できるようになることを目指す。

旧Workflowで設計検討を開始しており、Terraform管理対象、AWSリージョン、OS、インスタンスタイプ、Security Groupの基本方針までは決定済みである。本Planでは、その現在地を引き継いで学習を継続する。

### 今回達成したいこと

- TerraformでEC2インスタンス1台とSecurity Group1個を作成・確認・削除できる。
- `terraform plan`の内容を理解してから`apply`を実行できる。
- 設計書とTerraformコードの対応関係を説明できる。
- 採用した設計、代替案、不採用理由、実務との差を説明できる。
- 学習内容をデモ、Pull Request、動画、Qiita記事などへ整理できる。

### 重点的に学びたいこと

- TerraformによるAWSリソース管理の基本
- 既存AWSリソースとTerraform管理リソースの境界
- EC2、Security Group、AMI、VPC、Subnet、Key Pairの関係
- Variables、Outputs、State、Providerの設計
- `init`、`fmt`、`validate`、`plan`、`apply`、`destroy`の役割
- セキュリティ上の簡略化と実務推奨構成の違い

### 制約・希望

- 一度に1 Stepずつ進める。
- 設計を確定し、内容を説明できるようになるまで実装しない。
- VPC、Subnet、Internet Gateway、Route Table、Key Pairは新規作成しない。
- 検証後は`terraform destroy`を実行し、課金対象を残さない。
- 未決定事項をAIが独自判断で実装しない。

## Plan作成前の理解度診断

### 自分の言葉で説明できること

- Terraform管理対象と既存AWSリソースを分ける必要性
- 学習範囲をEC2とSecurity Groupに限定する理由
- リージョン、OS、CPUアーキテクチャ、インスタンスタイプの採用理由
- 必要なポートだけを許可するSecurity Groupの基本
- SSHの`0.0.0.0/0`公開を学習・短時間のデモに限定すべき理由

### 名前や概要だけ知っていること

- TerraformのVariables、Outputs、State
- Data Sourceによる既存リソースやAMIの検索
- TerraformとAWS Providerのバージョン制約
- Security Groupルールの記述方法

### 未経験または理解が曖昧なこと

- 既存Key PairとSSH接続の具体的な設計
- パブリックIPv4アドレスの付与方法
- Ubuntu 24.04 LTS AMIの取得方法
- 既存VPC・Subnetの指定方法
- Terraformファイルの分割とローカルState管理
- 実装後の一連のTerraform操作とAWS上での確認

### 実際に操作・実装したことがあるもの

- 既存AWS環境とBlogAPIのEC2運用に関する検討
- EC2とSecurity Groupの設計書の叩き台作成
- Terraform実装と実行は未着手

### 診断結果によるPlan調整

- 決定済みの対象範囲と基本設計は完了済みStepとして扱う。
- 次回は、未決定のKey Pair・SSH・パブリックIPv4設計から再開する。
- 設計確定、設計レビュー、実装、コードレビュー、AWS検証を別Stepに分ける。

## 学習ゴール

### 成果物のゴール

- `docs/architecture/terraform/ec2-instance.md`の未決定事項が解消され、実装可能な設計書になっている。
- 設計書に一致したTerraformコードが作成されている。
- `terraform fmt -check`と`terraform validate`が成功している。
- `terraform plan`、`apply`、AWS確認、`destroy`の結果が記録されている。
- デモと振り返りを実施し、必要に応じてPull Request、動画、Qiita記事へ整理されている。

### 理解のゴール

- 各Terraformブロックとリソースの役割を説明できる。
- Terraform管理対象と既存AWSリソースの境界を説明できる。
- 採用した設計と代替案を比較できる。
- `terraform plan`の作成・変更・削除内容を読み取れる。
- セキュリティやネットワークの前提が変わった場合の設計変更を説明できる。
- エラー発生時に、設定、認証、ネットワーク、Stateの観点から切り分けられる。

### Plan全体の完了条件

- [ ] EC2設計書が完成し、未決定事項が解消されている。
- [ ] 設計レビューの指摘を理解し、必要な修正が反映されている。
- [ ] Terraformコードと設計書が一致している。
- [ ] 主要なTerraform設定を自分の言葉で説明できる。
- [ ] `fmt`、`validate`、`plan`、`apply`、AWS確認、`destroy`を完了している。
- [ ] EC2とSecurity Groupが削除されたことを確認している。
- [ ] 採用案、代替案、実務向け改善案を説明できる。
- [ ] デモと振り返りを完了している。
- [ ] 残った課題を今後の学習候補として整理している。

## 学習範囲

### 対象

- AWS Provider、EC2インスタンス1台、Security Group1個
- 既存VPC、Subnet、Key Pairの参照
- Ubuntu Server 24.04 LTS x86_64 AMI
- Variables、Outputs、タグ、ローカルState
- Terraformの初期化、検証、計画、適用、削除
- 設計レビュー、コードレビュー、AWSコンソール確認、振り返り

### 対象外

- VPC、Subnet、Internet Gateway、Route Table、Key Pairの新規作成・変更
- Elastic IP、Docker、BlogAPIデプロイ、Nginx、Route 53、HTTPS、RDS、ALB
- GitHub ActionsからのTerraform実行
- S3などを使用したリモートState管理

## 参照対象

### リポジトリ内

- `AGENTS.md`
- `docs/agents/operations.md`
- `docs/study/learning-workflow.md`
- `docs/study/templates/learning-plan-template.md`
- `docs/study/templates/learning-progress-template.md`
- `docs/study/progress/terraform-ec2-instance-progress.md`
- `docs/architecture/terraform/ec2-instance.md`
- `.gitignore`
- Terraform実装時に決定する配置先

### 外部資料

- HashiCorp Terraform公式ドキュメント
- HashiCorp AWS Provider公式ドキュメント
- AWS EC2、VPC、Security Group、AMI、Key Pairの公式ドキュメント
- Ubuntu公式のAWS AMI情報

外部情報を使用するときは一次情報を優先し、参照日を記録する。

## Step一覧

| Step | タイトル | 主な目的 | 成果物への反映 | 完了確認 | 状態 |
| --- | --- | --- | --- | --- | --- |
| 1 | 管理対象と既存リソースの整理 | 学習範囲と管理境界を決める | 設計書2〜5章 | 理由を説明できる | `完了` |
| 2 | EC2とSecurity Groupの基本設計 | 基本構成とセキュリティ方針を決める | 設計書6〜7章 | 採用理由と注意点を説明できる | `完了` |
| 3 | Key Pair・SSH・パブリックIPv4設計 | EC2への接続方法を決める | 設計書8章 | 接続経路と安全上の制約を説明できる | `未着手` |
| 4 | AMI・VPC・Subnetの参照設計 | 既存リソースとAMIの取得方法を決める | 設計書6章、9〜10章 | VariableとData Sourceを比較できる | `未着手` |
| 5 | Terraform構成要素の詳細設計 | Variables、Outputs、命名、Stateなどを確定する | 設計書9〜13章 | 各設計判断を説明できる | `未着手` |
| 6 | 設計書の完成と自己レビュー | 未決定事項をなくし、全体を説明する | 設計書全体 | 資料を見ず主要判断を説明できる | `未着手` |
| 7 | Codexによる設計レビュー | 実装前に整合性と安全性を確認する | 設計書の修正 | 指摘の採否と理由を説明できる | `未着手` |
| 8 | Terraformの段階的実装 | 設計に沿って小さく実装する | Terraformコード | 各ブロックを説明できる | `未着手` |
| 9 | コードレビュー | 実装と設計の一致を確認する | コード・設計書 | 差分と判断理由を説明できる | `未着手` |
| 10 | Terraform実行とAWS確認 | 作成・確認・削除を安全に実施する | 検証記録 | planとAWS状態を照合できる | `未着手` |
| 11 | デモと成果整理 | 学習内容を他者へ説明可能にする | デモ、PR、動画、記事 | 一連の流れを説明できる | `未着手` |
| 最終 | 復習・応用 | 理解の定着と応用力を確認する | 振り返り | 説明・条件変更・障害対応 | `未着手` |

状態は `未着手 / 進行中 / 完了 / 保留` のいずれかを使用する。

## Step詳細

### Step 1: 管理対象と既存リソースの整理

- 目的: Terraformが作成・削除するものと、既存AWSリソースとして参照するものを明確にする。
- 学ぶ内容: Terraform管理境界、学習範囲の分割、依存関係。
- 学習者が説明すること: EC2とSecurity Groupだけを管理し、ネットワークを対象外にした理由。
- 成果物への反映: EC2設計書の目的、対象範囲、対象外、全体構成。
- 完了条件: 管理対象と既存リソースを区別し、採用理由を説明できる。
- 確認証拠: 旧Workflowの決定記録とEC2設計書2〜5章。

### Step 2: EC2とSecurity Groupの基本設計

- 目的: EC2の基本仕様と、学習用Security Groupの方針を決める。
- 学ぶ内容: リージョン、OS、CPU、インスタンスタイプ、インバウンド・アウトバウンド。
- 学習者が説明すること: `ap-northeast-1`、Ubuntu 24.04 LTS x86_64、`t3.micro`、SSH 22/TCPを採用した理由。
- 成果物への反映: EC2設計書6〜7章。
- 完了条件: 学習用の`0.0.0.0/0`と実務推奨構成の違いを説明できる。
- 確認証拠: 旧Workflowの決定記録とEC2設計書6〜7章。

### Step 3: Key Pair・SSH・パブリックIPv4設計

- 目的: 既存Key Pairを使用したEC2への接続方法と完了条件を決める。
- 最初に考えること: Key Pair名をコードへ固定せず渡す方法、SSH確認を必須にする利点と注意点。
- 演習・作業: 既存Key Pair名、Variable化、SSH必須条件、パブリックIPv4付与、SSHユーザー名を一つずつ検討する。
- 成果物への反映: EC2設計書8章と関連する完了条件。
- 完了条件: 接続経路、必要な設定、秘密鍵の管理境界、安全上の制約を説明できる。
- 理解できなかった場合: Key Pairと秘密鍵、Subnet設定とパブリックIPの関係を図または具体例で再確認する。

### Step 4: AMI・VPC・Subnetの参照設計

- 目的: AMIと既存ネットワークリソースを再現可能かつ安全に指定する方法を決める。
- 最初に考えること: VariableでIDを渡す方法とData Sourceで検索する方法の違い。
- 演習・作業: AMI取得、VPC指定、Subnet指定の候補を比較し、今回の採用方法を決める。
- 成果物への反映: EC2設計書6.2、6.4、9〜10章。
- 完了条件: 採用方法、代替案、環境差や検索結果が複数になるリスクを説明できる。

### Step 5: Terraform構成要素の詳細設計

- 目的: 実装に必要な残りの設計判断をすべて確定する。
- 学ぶ内容: Terraform・Providerバージョン、ファイル配置、Variables、Outputs、タグ、命名、State、Security Group Ruleの記述方法。
- 演習・作業: 設計書16章の未決定事項を順番に解消する。
- 成果物への反映: EC2設計書9〜16章。
- 完了条件: すべての入力、出力、リソース名、ファイル、State管理方法を説明できる。

### Step 6: 設計書の完成と自己レビュー

- 目的: 設計書を実装可能な状態にし、学習者自身が全体を説明する。
- 演習・作業: 未決定表現、矛盾、対象外への逸脱、秘密情報の扱いを確認する。
- 成果物への反映: EC2設計書全体。
- 完了条件: 資料を見ずに管理境界、接続、AMI、ネットワーク、変数、State、セキュリティを説明できる。

### Step 7: Codexによる設計レビュー

- 目的: 実装前にリポジトリ整合性、Terraform設計、セキュリティ、保守性、学びやすさを確認する。
- 演習・作業: 指摘ごとに理由を理解し、採用・不採用を判断して設計書へ反映する。
- 完了条件: すべての指摘に対応方針があり、実装を妨げる未決定事項が残っていない。

### Step 8: Terraformの段階的実装

- 目的: 確定した設計に沿って、小さな単位でTerraformコードを作成する。
- 実装順序: Terraform・Provider設定、Variables、Data Source、Security Group、EC2、Outputs。
- AIの確認観点: 未確定値の推測、秘密情報、対象外リソースの作成、設計との差異がないこと。
- 完了条件: 各ブロックの役割と主要引数を学習者が説明でき、`terraform fmt -check`と`terraform validate`が成功する。

### Step 9: コードレビュー

- 目的: 生成されたコードを読み、設計と実装の対応を確認する。
- 演習・作業: resource、data、variable、output、依存関係、代替記法を説明する。
- 完了条件: 設計書との差異がなく、差異がある場合は理由と両方の修正が記録されている。

### Step 10: Terraform実行とAWS確認

- 目的: `plan`を理解したうえで安全にリソースを作成し、削除まで確認する。
- 実行順序: `init`、`fmt -check`、`validate`、`plan`、学習者確認、`apply`、AWS確認、必要に応じたSSH確認、`destroy`、削除確認。
- 完了条件: planとAWSコンソールの状態を照合でき、EC2とSecurity Groupが削除されている。
- 失敗時: 認証、入力値、AMI、既存ネットワーク、Security Group、Stateの順に切り分け、強引な再実行をしない。

### Step 11: デモと成果整理

- 目的: 学んだ設計と操作を第三者へ説明できる形にする。
- 演習・作業: Terraformコード確認からdestroyまでのデモ、Pull Request、動画、Qiita記事に必要な説明を整理する。
- 完了条件: コマンドの羅列ではなく、各操作の目的、設計判断、注意点を説明できる。

## 最終Step: 復習・応用

### 資料を見ない説明

- Terraform管理対象と既存AWSリソースの境界を説明する。
- EC2作成からdestroyまでの処理と確認点を説明する。
- 設計書とTerraformコードの対応を説明する。

### 代替案との比較

- VariableとData Sourceによる既存リソース指定を比較する。
- SSH公開とSession Managerを比較する。
- ローカルStateとリモートStateを比較する。

### 条件変更問題

- 継続運用する本番環境になった場合、SSH、State、ネットワーク、可用性をどう変更するか考える。
- HTTP/HTTPS、Docker、RDSを追加する場合、Planをどの成果物単位に分けるか考える。

### 障害対応問題

- AMIが見つからない、SSH接続できない、planが想定外の置換を示す、destroyに失敗する場合の切り分けを説明する。

### 最終確認

- [ ] 成果物の内容と説明に矛盾がない。
- [ ] 主要な判断を自分の言葉で説明できる。
- [ ] 条件変更と障害対応の問いに回答できる。
- [ ] 残った課題を今後の学習候補へ記録した。

## Plan変更ルール

- Planは学習中の理解度や新しく判明した前提に応じて変更してよい。
- AIエージェントは、変更前に理由と後続Stepへの影響を提示する。
- 学習者の確認後、Planと進捗ファイルの変更履歴を更新する。
- 完了済み記録は削除せず、再確認が必要な場合は理由を残す。
