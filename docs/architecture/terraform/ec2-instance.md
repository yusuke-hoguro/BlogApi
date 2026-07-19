# Terraform EC2インスタンス設計書

## 1. 目的

Terraformを使用して、BlogAPIの学習・検証用EC2インスタンスを1台構築する。

本フェーズでは、EC2を作成することだけでなく、以下を理解し、自分の言葉で説明できる状態を目指す。

- TerraformによるAWSリソース管理の基本
- `terraform plan`による変更内容の確認
- EC2とSecurity Groupの依存関係
- 既存AWSリソースをTerraformから参照する考え方
- 設計判断と実装内容の対応関係

本設計書は、CodexがTerraformコードを実装する前提で作成する。ただし、未決定事項はCodexが推測して実装せず、設計確定後に着手する。

---

## 2. フェーズ1のゴール

以下を一連のデモとして実行できることを完成条件とする。

1. Terraformコードを確認する
2. `terraform init`を実行する
3. `terraform fmt`を実行する
4. `terraform validate`を実行する
5. `terraform plan`で作成内容を確認する
6. `terraform apply`でリソースを作成する
7. AWSコンソールでEC2とSecurity Groupを確認する
8. 必要に応じてSSH接続を確認する
9. `terraform destroy`でリソースを削除する
10. AWSコンソールで削除を確認する

---

## 3. 対象範囲

### Terraformで管理するリソース

- EC2インスタンス 1台
- EC2用Security Group 1個

### Terraformで設定する項目

- AWS Provider
- AWSリージョン
- AMI
- EC2インスタンスタイプ
- 既存Subnetへの配置
- 既存VPCへのSecurity Group配置
- SSH用インバウンドルール
- アウトバウンドルール
- Variables
- Outputs
- リソース識別用タグ

---

## 4. 対象外

以下は本フェーズのTerraform管理対象に含めない。

- VPCの新規作成
- Subnetの新規作成
- Internet Gatewayの新規作成
- Route Tableの新規作成・変更
- Elastic IP
- Key Pairの新規作成・秘密鍵管理
- Dockerのインストール
- BlogAPIのデプロイ
- Nginx
- Route 53
- HTTPS・TLS証明書
- RDS
- ALB
- GitHub ActionsからのTerraform実行
- S3バックエンドなどを用いたリモートState管理

これらは、Terraformの基本操作とEC2作成を理解した後のフェーズで検討する。

---

## 5. 全体構成

### 既存AWSリソース

以下はAWS上に存在するリソースを利用し、Terraformでは作成・変更・削除しない。

- 既存VPC
- 既存Subnet
- 既存Internet Gateway
- 既存Route Table
- 既存Key Pair

### Terraform導入後の構成

```text
既存VPC
└── 既存Subnet
    └── EC2インスタンス（Terraform管理）
        └── Security Group（Terraform管理）

既存Key Pair
└── EC2インスタンスから参照
```

---

## 6. 設計方針

### 6.1 AWSリージョン

| 項目 | 設定値 |
|---|---|
| リージョン | `ap-northeast-1` |
| 表示名 | アジアパシフィック（東京） |

#### 採用理由

- 日本からの利用を前提としている
- BlogAPIの既存AWS環境と統一しやすい
- AWSおよびTerraformの日本語情報が多い
- RDSやALBなどを追加する場合も同一リージョンで拡張しやすい

リージョンはVariableとして定義し、初期値を`ap-northeast-1`とする方針とする。

### 6.2 AMI・OS

| 項目 | 設定値 |
|---|---|
| OS | Ubuntu Server 24.04 LTS |
| CPUアーキテクチャ | x86_64 |

#### 採用理由

- 現在のBlogAPI運用環境とOSを統一しやすい
- `apt`を利用した既存の知識を活用できる
- Docker、Nginx、Goなどの情報が豊富である
- LTS版であり、今後のフェーズでも継続利用しやすい
- ARM固有の差異を今回の学習範囲へ持ち込まない

#### 実装上の注意

AMI IDはリージョンごとに異なり、将来変更される可能性がある。

Codexは、特定のAMI IDを設計判断なしにハードコードしない。取得方法は今後、以下から選択する。

- Ubuntu公式AMIを条件指定してData Sourceから取得する
- 公開パラメータを利用する
- AMI IDをVariableとして渡す

現時点では、**Ubuntu Server 24.04 LTS・x86_64を使用することのみ確定**とする。

### 6.3 EC2インスタンスタイプ

| 項目 | 設定値 |
|---|---|
| インスタンスタイプ | `t3.micro` |

#### 採用理由

- TerraformおよびEC2の学習用途として十分な性能である
- コストを抑えやすい
- x86_64環境であり、既存のソフトウェアとの互換性を確保しやすい
- 必要になった場合、`t3.small`などへ変更しやすい

インスタンスタイプはVariableとして定義し、初期値を`t3.micro`とする。

### 6.4 VPC・Subnet

VPCおよびSubnetは新規作成せず、既存リソースを利用する。

#### 採用理由

- 今回はTerraformの基本とEC2作成の理解を優先する
- AWSネットワーク設計を同時に扱うと学習範囲が広がりすぎる
- ネットワーク構築を後続フェーズとして分離した方が、設計判断を説明しやすい

#### 実装上の注意

既存VPC IDおよびSubnet IDの指定方法は未決定である。

Codexは、特定のIDをコードへ直接ハードコードしない。今後、以下から選択する。

- VariableとしてIDを渡す
- タグなどを条件にData Sourceで検索する

EC2へパブリックIPv4アドレスを付与するかどうかも、既存Subnetの設定とSSH接続方法を確認して決定する。

---

## 7. Security Group設計

### 7.1 今回採用する設定

#### インバウンド

| 用途 | プロトコル | ポート | 接続元 |
|---|---|---:|---|
| SSH | TCP | 22 | `0.0.0.0/0` |

#### アウトバウンド

| 用途 | プロトコル | ポート | 接続先 |
|---|---|---|---|
| すべての外向き通信 | All | All | `0.0.0.0/0` |

#### 採用理由

- Terraformによる作成とSSH接続確認を優先する
- 接続元IPの変更によって検証が止まることを避ける
- HTTPおよびHTTPSは今回使用しないため開放しない
- 将来の`apt update`やDockerイメージ取得を考慮し、アウトバウンドは許可する

### 7.2 セキュリティ上の注意

SSHを`0.0.0.0/0`へ公開する設定は、学習・一時的なデモ用途に限定する。

以下を運用条件とする。

- 検証終了後は`terraform destroy`を実行する
- 長時間起動したまま放置しない
- 秘密鍵をGitリポジトリへコミットしない
- 不要なポートを追加しない

### 7.3 実務で推奨する設定

実務または継続運用する環境では、SSHの接続元をグローバルIPv4アドレスの`/32`に限定する。

```text
203.0.113.10/32
```

さらに、以下も検討する。

- AWS Systems Manager Session Managerを利用し、SSHポートを公開しない
- 踏み台サーバーまたはVPN経由に限定する
- HTTP/HTTPSはALBやリバースプロキシ経由で公開する
- Security Group間参照を利用する

今回の簡略化された設定と、実務での推奨構成を明確に分けて扱う。

---

## 8. Key Pair・SSH接続

既存Key Pairを利用する方針は確定している。

ただし、以下は未決定である。

- 使用する既存Key Pair名
- Key Pair名をVariableとして渡すか
- SSH接続をフェーズ1の必須完了条件にするか
- パブリックIPv4アドレスの付与方法
- SSHユーザー名

Codexは、Key Pairを新規作成したり、秘密鍵をTerraform StateまたはGitリポジトリで管理したりしない。

---

## 9. Terraformリソース一覧

| 種別 | Terraform上の想定 | 用途 |
|---|---|---|
| Terraform設定 | `terraform` block | TerraformとProviderのバージョン要件 |
| Provider | `provider "aws"` | 東京リージョンのAWS操作 |
| Data SourceまたはVariable | 未確定 | 既存VPCの参照 |
| Data SourceまたはVariable | 未確定 | 既存Subnetの参照 |
| Data SourceまたはVariable | 未確定 | Ubuntu 24.04 LTS AMIの取得 |
| Resource | `aws_security_group` | EC2用Security Group |
| Resource | `aws_instance` | EC2インスタンス1台 |
| Variable | 複数 | 環境依存値と変更可能値 |
| Output | 複数 | 作成結果の確認 |

Security Groupルールを`aws_security_group`内へ記述するか、個別リソースとして定義するかは未決定とする。

---

## 10. Variables設計

| 変数候補 | 初期値・例 | 状態 |
|---|---|---|
| AWSリージョン | `ap-northeast-1` | 確定 |
| EC2インスタンスタイプ | `t3.micro` | 確定 |
| VPC IDまたは検索条件 | 未定 | 未確定 |
| Subnet IDまたは検索条件 | 未定 | 未確定 |
| AMI IDまたは検索条件 | Ubuntu 24.04 LTS x86_64 | 方式未確定 |
| Key Pair名 | 未定 | 未確定 |
| SSH接続元CIDR | `0.0.0.0/0` | フェーズ1では確定 |
| EC2名・タグ用の名前 | 未定 | 未確定 |

### 変数設計方針

- 環境によって変わる値をコードへ直接埋め込まない
- 秘密情報をVariableへ含めない
- 秘密鍵をTerraformで管理しない
- 各Variableには型と説明を設定する
- 必要に応じて初期値や`validation`を設定する

---

## 11. Outputs設計

| Output候補 | 用途 | 状態 |
|---|---|---|
| EC2インスタンスID | AWSコンソールとの照合 | 採用予定 |
| EC2パブリックIPv4アドレス | SSH接続確認 | パブリックIP設計後に確定 |
| EC2プライベートIPv4アドレス | ネットワーク確認 | 採用予定 |
| Security Group ID | AWSコンソールとの照合 | 採用予定 |
| Availability Zone | 配置先確認 | 採用候補 |

機密情報はOutputへ出力しない。

---

## 12. タグ設計

EC2とSecurity Groupには、AWSコンソール上で識別できるタグを設定する。

| タグ | 値の例 | 目的 |
|---|---|---|
| `Name` | 未定 | リソース名の識別 |
| `Project` | `BlogApi` | 対象プロジェクトの識別 |
| `ManagedBy` | `Terraform` | Terraform管理対象の明示 |
| `Environment` | `learning`または`dev` | 用途・環境の識別 |

正式な値は、既存リポジトリやAWSリソースの命名規則を確認して決定する。

---

## 13. Terraform State管理

フェーズ1ではローカルStateを利用する。

Stateファイルにはリソース情報が含まれるため、Git管理しない。

実装時には既存の`.gitignore`を確認し、少なくとも以下が除外されていることを確認する。

```gitignore
.terraform/
*.tfstate
*.tfstate.*
*.tfplan
```

`.terraform.lock.hcl`は、原則として再現可能性を確保するためGit管理対象とする。

チーム開発や継続運用へ移行する場合は、リモートState管理を後続フェーズで検討する。

---

## 14. Codexへの実装制約

Codexは以下を守る。

1. 本設計書で未確定の値を推測して実装しない
2. 未決定事項が実装に必要な場合は、実装前に指摘する
3. VPC、Subnet、Internet Gateway、Route Table、Key Pairを新規作成しない
4. 既存ネットワークリソースを変更・削除しない
5. Ubuntu Server 24.04 LTSのx86_64 AMIを使用する
6. EC2インスタンスタイプは`t3.micro`とする
7. インバウンドはSSHの22/TCPのみとする
8. HTTPの80/TCPおよびHTTPSの443/TCPを追加しない
9. 秘密鍵やAWS認証情報をコード、Variable、tfvars、Outputへ含めない
10. Terraform StateをGit管理しない
11. `terraform fmt`と`terraform validate`を通す
12. 設計書と実装内容が一致することを確認する

---

## 15. 実行・確認手順

設計確定および実装完了後、以下の順番で確認する。

```bash
terraform init
terraform fmt -check
terraform validate
terraform plan
terraform apply
```

AWSコンソールで以下を確認する。

- EC2インスタンスが1台作成されている
- OSがUbuntu Server 24.04 LTSである
- インスタンスタイプが`t3.micro`である
- 想定した既存Subnetへ配置されている
- Terraformで作成したSecurity Groupが関連付けられている
- インバウンドがSSHの22/TCPのみである
- HTTPおよびHTTPSが許可されていない
- タグが付与されている
- OutputとAWSコンソールの値が一致している

確認後、以下を実行する。

```bash
terraform destroy
```

EC2インスタンスとSecurity Groupが削除されたことをAWSコンソールで確認する。

---

## 16. 未決定事項

次回以降、以下を一つずつ検討する。

1. 既存Key Pairの利用方法
2. SSH接続を必須完了条件に含めるか
3. EC2へパブリックIPv4アドレスを付与する方法
4. 既存VPCの指定方法
5. 既存Subnetの指定方法
6. Ubuntu 24.04 LTS AMIの取得方法
7. Terraformファイルの配置先とファイル分割
8. TerraformおよびAWS Providerのバージョン制約
9. Security Group Ruleの記述方法
10. Variableの正式名称、型、説明、初期値
11. Outputの正式名称
12. タグとリソースの命名規則

未決定事項を解消した後に本設計書を更新し、Codexへ設計レビューを依頼する。

---

## 17. 将来の拡張候補

- SSH接続元CIDRの制限
- Systems Manager Session Managerの導入
- Elastic IP
- VPC・Subnetなどのネットワークリソース管理
- Docker導入
- BlogAPIデプロイ
- Nginx
- Route 53
- HTTPS
- RDS
- ALB
- GitHub ActionsによるTerraformの検証・実行
- リモートState管理

一度に追加せず、各フェーズで目的、対象範囲、設計判断を明確にする。

---

## 18. 現時点の設計判断まとめ

| 項目 | 決定内容 |
|---|---|
| Terraform管理対象 | EC2 1台、Security Group 1個 |
| ネットワーク | 既存VPC・既存Subnetを利用 |
| リージョン | `ap-northeast-1` |
| OS | Ubuntu Server 24.04 LTS |
| CPUアーキテクチャ | x86_64 |
| インスタンスタイプ | `t3.micro` |
| インバウンド | SSH 22/TCPのみ |
| SSH接続元 | `0.0.0.0/0`（学習・一時的なデモ用途） |
| アウトバウンド | すべて許可 |
| HTTP・HTTPS | 今回は許可しない |
| Key Pair | 既存Key Pairを利用するが詳細未決定 |
| State | ローカルStateを利用し、Git管理しない |

この設計判断を基準とし、未決定事項を確定させてからTerraform実装へ進む。
