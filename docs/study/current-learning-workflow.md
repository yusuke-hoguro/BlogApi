# AI Learning Workflow（Terraform）

## 学習テーマ

Terraformを使用して、BlogAPIのAWSインフラを構築する。

---

# 学習目的

今回の目的は、

**TerraformでEC2インスタンスを作成することではない。**

Terraformの考え方、設計方法、レビュー方法、実装方法を理解し、自分の言葉で説明できるようになることを最優先とする。

---

# 基本方針

AIへ丸投げはしない。

以下の流れで、一つずつ理解しながら進める。

```text
設計
    ↓
設計レビュー
    ↓
実装
    ↓
コードレビュー
    ↓
動作確認
    ↓
振り返り
```

---

# 学習の進め方

## ① ゴールを小さく設定する

一度に大きなものは作らない。

例)

* EC2インスタンスを作成する
* Security Groupを追加する
* Elastic IPを追加する
* Dockerを導入する
* GitHub ActionsからTerraformを実行する

毎回、小さな単位で完成させる。

---

## ② まず設計を決める

コードを書く前に設計を行う。

ChatGPTと相談しながら、

* 何を作るのか
* なぜ必要なのか
* 他にどんな選択肢があるのか
* 今回なぜその方法を選ぶのか

を一つずつ決定する。

---

## ③ ChatGPTが設計書を作成する

話し合って決まった内容を設計書へまとめる。

設計書には以下を記載する。

* 目的
* 対象範囲
* 対象外
* 現在の構成
* Terraform導入後の構成
* 使用するAWSサービス
* ディレクトリ構成
* Terraformリソース一覧
* 変数設計
* Output
* 実行手順
* 削除手順
* セキュリティ
* 将来の拡張
* 設計判断理由

コードだけではなく、

**「なぜその設計なのか」**

を残すことを重視する。

---

## ④ 自分で設計書を読む

設計書をそのまま採用しない。

必ず一度読み、

* 分からない用語
* 理解できない設計
* なぜその設計なのか

をChatGPTへ質問する。

理解できるまで確認する。

---

## ⑤ Codexへ設計レビューを依頼する

まだ実装しない。

まず設計レビューだけを依頼する。

レビュー内容

* リポジトリ構成との整合性
* Terraform設計として問題ないか
* セキュリティ
* 保守性
* 初学者にも理解しやすいか

Codexからの指摘はChatGPTとも確認し、

理由を理解した上で採用する。

---

## ⑥ Codexへ実装を依頼する

設計確定後に実装する。

一度に全部ではなく、

* Provider
* Variables
* Security Group
* EC2
* Outputs

など小さく区切って実装する。

---

## ⑦ コードレビュー

生成されたTerraformをChatGPTと読む。

確認すること

* このresourceは何か
* この引数は何を意味するか
* この値はなぜ必要なのか
* 他の書き方はあるか
* 設計書と一致しているか

理解できるまで確認する。

---

## ⑧ Terraform実行

実行手順

```bash
terraform init
terraform fmt
terraform validate
terraform plan
terraform apply
```

特に

**terraform plan**

を理解してからapplyする。

---

## ⑨ 動作確認

AWS上で確認する。

* EC2が作成されたか
* Security Group
* SSH接続
* タグ
* Output

を確認する。

---

## ⑩ 振り返り

最後に

* 学んだこと
* 理解できなかったこと
* 次回改善すること

を整理する。

---

# AIの役割

## ChatGPT

担当

* 概念説明
* 設計相談
* 選択肢の提示
* 設計書作成
* コードレビュー
* 学習サポート
* Codexレビュー内容の確認

---

## Codex

担当

* 設計レビュー
* Terraform実装
* リファクタリング
* リポジトリへの反映

---

## 自分

担当

* 設計を判断する
* 提案を採用するか決める
* コードを読む
* Terraformを実行する
* AWSを確認する
* 疑問を質問する
* 最終的に説明できる状態にする

---

# 判断基準

AIが提案したから採用するのではなく、

**「なぜその設計なのか」を説明できること**

を採用条件とする。

AIで実装を高速化し、

**自分は設計・判断・レビューに集中する。**

---

# 今回のゴール

Terraformで

**EC2インスタンスを1台作成できること**

対象

* Provider
* EC2
* Security Group
* Variables
* Outputs

対象外

* Docker
* BlogAPIデプロイ
* Nginx
* GitHub Actions
* Route53
* HTTPS
* RDS

---

# 成果物

今回の成果物は以下とする。

* Terraform設計書
* Terraformコード
* GitHubコミット
* Pull Request
* デモ動画
* Qiita記事

設計段階から、

**「最終的に何を説明するのか」**

を意識して進める。

---

# デモ完成条件

以下の流れをデモできる状態を完成とする。

1. Terraformコードを確認
2. terraform init
3. terraform fmt
4. terraform validate
5. terraform plan
6. terraform apply
7. AWSコンソールでEC2確認
8. terraform destroy
9. AWSコンソールで削除確認

---

# 毎回の作業開始時

このMarkdownを読み込み、以下を確認してから作業を開始する。

確認対象

- 学習テーマ
- 学習目的
- 今回のゴール
- 現在の進捗
- 本日完了した作業
- 新たに決定した内容
- 現在の疑問点
- 次回実施する作業

内容を確認したうえで、現在地を把握し、前回の続きから作業を開始する。

---

# 毎回の作業終了時

作業終了後は、このMarkdownを更新する。

更新内容

## 現在の進捗

Terraformフェーズ1は、実装前のEC2設計段階にある。

`docs/architecture/terraform/ec2-instance.md` に設計書の叩き台を作成し、管理対象、リージョン、OS、インスタンスタイプ、Security Groupまで決定・記録した。

まだ設計は完成しておらず、Codexへの設計レビューおよびTerraform実装は開始していない。次回は、設計書に残した未決定事項を一つずつ検討する。

---

## 本日完了した作業

実施日：2026年7月19日

* `current-learning-workflow.md` を読み、前回の現在地と作業方針を確認した
* フェーズ1でTerraform管理対象とするAWSリソースの範囲を決定した
* TerraformではEC2インスタンス1台とSecurity Group1個を管理する方針とした
* VPC、Subnet、Internet Gateway、Route Table、Key Pairは既存AWSリソースを利用する方針とした
* AWSリージョンを `ap-northeast-1`（東京）に決定した
* EC2のOSをUbuntu Server 24.04 LTS、CPUアーキテクチャをx86_64に決定した
* EC2インスタンスタイプを `t3.micro` に決定した
* Security GroupではSSH用の22/TCPのみをインバウンドで許可し、学習用として接続元を `0.0.0.0/0` とする方針を決定した
* アウトバウンドはすべての通信を許可する方針を決定した
* 学習用の簡略化された構成と、実務で推奨されるセキュリティ構成を設計書で分けて記載する方針とした
* `docs/architecture/terraform/ec2-instance.md` に、Codexが実装前提で参照できる設計書の叩き台を作成した
* Codexが未決定事項を独自判断で実装しないよう、未確定項目と実装制約を設計書に明記した

---

## 今日学んだこと

* 初回のTerraform学習では、VPCなどのネットワーク構築を分離し、EC2とSecurity Groupに範囲を限定すると理解しやすい
* 既存AWSリソースを参照するものと、Terraformが作成・削除するものを明確に区別する必要がある
* AMI IDはリージョンや時間によって変わり得るため、OSを決めることと具体的なAMI IDの取得方法を決めることは別の設計判断である
* `t3.micro` は今回の学習用途に十分であり、必要になった段階で変数の値を変更して拡張できる
* Security Groupでは必要なポートだけを許可することが基本である
* SSHを `0.0.0.0/0` に公開する構成は学習・短時間のデモ用途に限定し、実務では接続元IPの制限やSession Managerを検討すべきである
* 設計書では「今回なぜ簡略化したか」と「実務ではどうするか」を分けて書くことで、設計意図を説明しやすくなる
* Codexへ実装を依頼する前に、未決定事項を明示して独自判断による実装を防ぐ必要がある

---

## 新たに決定した内容

* フェーズ1でTerraform管理対象とするのはEC2インスタンス1台とSecurity Group1個のみとする
* VPC、Subnet、Internet Gateway、Route Tableは既存AWSリソースを利用し、Terraformでは作成・変更・削除しない
* 既存Key Pairを利用し、TerraformではKey Pairや秘密鍵を作成・管理しない
* AWSリージョンは `ap-northeast-1`（東京）とする
* EC2のOSはUbuntu Server 24.04 LTS、CPUアーキテクチャはx86_64とする
* EC2インスタンスタイプは `t3.micro` とし、Variableで変更可能な設計を基本とする
* Security GroupのインバウンドはSSH用の22/TCPのみとする
* 学習・デモの簡便さを優先し、SSH接続元は一時的に `0.0.0.0/0` とする
* Security Groupのアウトバウンドはすべての通信を許可する
* HTTP、HTTPSおよびElastic IPはフェーズ1の対象外とする
* 検証後は `terraform destroy` を実行し、EC2とSecurity Groupを残さない
* 今回採用する学習用構成と、実務で推奨するセキュリティ構成を設計書で分けて説明する
* 設計書は `docs/architecture/terraform/ec2-instance.md` で管理する
* 未決定事項が残っている間はCodexへ実装を依頼しない

---

## 現在の疑問点

以下は未決定であり、次回以降に一つずつ検討する。

* 使用する既存Key Pair名
* Key Pair名をVariableとして外部から渡すか
* SSH接続をフェーズ1の必須完了条件にするか
* EC2へパブリックIPv4アドレスを付与する方法
* SSHユーザー名を `ubuntu` として設計書へ明記するか
* Ubuntu 24.04 LTSのAMI取得方法
  * AWS ProviderのData Sourceで検索する
  * SSM Parameter Storeの公開パラメータを利用する
  * AMI IDをVariableで渡す
* 既存VPC IDとSubnet IDの指定方法
  * VariableでIDを渡す
  * タグなどを条件にData Sourceで検索する
* Variablesとして定義する値とdefault値
* Outputsとして公開する値
* リソース名およびタグの命名規則
* Terraformファイルの配置ディレクトリと分割方法
* ローカルStateの具体的な管理方法と `.gitignore` の確認
* AWS ProviderとTerraform本体のバージョン制約

---

## 次回実施する作業

1. `current-learning-workflow.md` を読み、現在地と作業方針を確認する
2. `docs/architecture/terraform/ec2-instance.md` を読み、決定済み事項と未決定事項を確認する
3. **既存Key PairとSSH接続の設計**を最初に検討する
4. 使用する既存Key Pair名と、Key Pair名をVariableとして渡すかを決定する
5. SSH接続をフェーズ1の必須完了条件に含めるかを決定する
6. パブリックIPv4アドレスの付与方法とSSHユーザー名を決定する
7. 決定内容と理由を `docs/architecture/terraform/ec2-instance.md` に反映する
8. 続いて、AMI取得方法、VPC・Subnetの指定方法、Variables、Outputs、ディレクトリ構成、State管理を順番に検討する
9. すべての設計項目を決定した後、自分で設計書を読み、不明点をChatGPTへ質問する
10. 設計を理解できた段階でCodexへ設計レビューのみを依頼する

次回最初に検討する項目は、**既存Key PairとSSH接続の設計**とする。

現時点では、CodexへTerraform実装を依頼しない。

---

このMarkdownを**今回のTerraform学習の最新状態**として管理する。

次回作業時は、このMarkdownをChatGPTへ読み込ませ、前回の続きから作業を開始する。

Terraformの学習が終了したら、このMarkdownをベースに次の学習テーマ（Kubernetes、AWS、gRPCなど）へ更新し、同じ運用で継続する。

# ChatGPTへのプロンプト例

## 作業開始時

以下のプロンプトを使用する。

```text
docs/study/current-learning-workflow.md を読み込み、「毎回の作業開始時」に記載された内容を実行してください。

現在地を把握し、今回実施すべき作業を提案してください。

作業中は、このWorkflowに従って進めてください。
```

---

## 作業終了時

以下のプロンプトを使用する。

```text
今日の作業はここで終了します。

docs/study/current-learning-workflow.md の「毎回の作業終了時」に従い、更新対象を現在の内容へ更新してください。

更新が完了したら、今回の作業内容を簡単にまとめてください。
```

---

## 新しいチャットを開始するとき

以下のプロンプトを使用する。

```text
docs/study/current-learning-workflow.md を読み込み、この学習テーマの現在地を把握してください。

その後、「毎回の作業開始時」に従って作業を開始してください。
```

---

このプロンプトを基本形とし、必要に応じて今回実施したい内容を追加して指示する。