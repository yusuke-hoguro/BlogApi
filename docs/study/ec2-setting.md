
# GitHubにEC2を設定する

## 設定手順

1. GitHub 設定するリポジトリーにアクセスする
1. リポジトリーのSettingsを選択する
1. Secrets and variables → Actions → New repository secret の順で選択する
1. Secrets を登録
1. Name（環境変数名）Secret（値）を入力して、1 個ずつ登録する

## 設定例
- EC2_HOST → blog-api.duckdns.org
- EC2_USER → ubuntu（Amazon Linux なら ec2-user）
- EC2_SSH_KEY → EC2 鍵ペアの秘密鍵（id_rsa の中身）