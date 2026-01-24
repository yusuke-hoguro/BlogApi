# Nginxでのリバースプロキシ設定

## 1.背景・目的

HTTPS対応のためにNginxをリバースプロキシとして利用する方法を学ぶ。

## 2.問題

### 問題1: リバースプロキシを入れてHTTPS化する理由

- #### 状況:
    - リバースプロキシを入れてHTTPS化する理由が理解できていない
- #### 解決策:
    - リバースプロキシを導入する理由を調べる
- #### 結果:
    - ToDo：調べて記載すること

---

### 問題2: CertbotでHTTPS証明書取得に失敗

- #### 状況:
    - 証明書取得のために作成したCertbotコンテナで証明書を作成しようとして失敗
- #### 実施手順:
    - 下記のコマンドを実行して証明書取得を試みて失敗した

    ```bash
    docker run -it --rm \
    -v $(pwd)/certbot/conf:/etc/letsencrypt \
    -v $(pwd)/certbot/www:/var/www/certbot \
    certbot/certbot certonly \
    --webroot --webroot-path=/var/www/certbot \
    --email your-email@example.com \
    --agree-tos --no-eff-email \
    -d example.com -d www.example.com
    ```
- #### 原因:
    - ドメインがないので証明書を作成することができない
- #### 解決策:
    - DuckDNSで無料のドメインを取得する
- #### 理由:
    - 自己署名証明書だと下記の問題がある
- #### 結果: 
    - 

---

### 問題3:EC2で試す前にローカルPCでも確認できるか？

- #### 状況: 
    - リバースプロキシの設定と動作確認がEC2上ではなく、ローカルPCでも実施できるかを確認
- #### 解決策:
    - Docker Compose はそのまま使える
    - hosts ファイルでローカル IP にマッピングすると

    ```text
    127.0.0.1 blog-api.duckdns.org
    ```
    - HTTPS 証明書についてはローカルだと自己署名証明書を作って `Nginx` に設定が必要
    - ローカルでは 80/443 が別プロセスに使われている場合は変更が必要

- #### 結果:
    - Let’s Encryptの証明書取得にEC2の公開IPとドメインが必要だったのでEC2上で確認をすすめた

---

### 問題4:DuckDNSの更新は自動化スクリプトで定期更新が必要

- #### 状況: 
    - DuckDNSは無料のダイナミックDNSサービスなので、EC2や自宅PCのグローバルIPが変わるたびにDNSレコードを更新する必要がある
- #### 作成と自動更新手順:
    - 更新スクリプト作成
    ```bash
    #!/bin/bash
    # DuckDNSトークンとサブドメイン
    TOKEN="あなたのトークン"
    DOMAIN="blog-api"

    # 更新実行
    curl -s "https://www.duckdns.org/update?domains=$DOMAIN&token=$TOKEN&ip="
    ```
    - 実行権限を付与
    ```bash
    chmod +x duck.sh
    ```
    - cronで定期実行
    ```bash
    crontab -e
    ```
    ```bash
    */5 * * * * /home/ubuntu/duckdns/duck.sh >> /home/ubuntu/duckdns/duck.log 2>&1
    ```
- #### 結果:
    - 現在は未対応だが、DuckDNSを使う場合は自動更新が必要なので上記の仕組みを追加する

---

### 問題5:ローカルPCでの動作確認やテストに影響はないのか

- #### 状況:
    - httpsでアクセスできるようにnginxを追加したが、ローカルPCでの確認に影響はないのかが気になった。
- #### 結果:
    - nginx を追加したのは サーバー（EC2 など）上 で、HTTPS 化のため
    - ローカルPCでは `Go` の `net/http` サーバー を直接 `localhost:8080` で立ち上げている状態
    - ローカルで http://localhost:8080 にアクセスする場合、nginx は介在しない
    - よってローカルPCでの確認に影響はなし

---

### 問題6:nginx が SSL証明書ファイルを読み込めず起動に失敗

- #### 状況:
    - nginx が SSL証明書ファイル `/etc/letsencrypt/live/blog-api.duckdns.org/fullchain.pem` を読み込もうとしたが失敗
- #### 原因:
    - Certbot で証明書がまだ取得できていない。ファイルがない段階でコンテナを起動しようとしているため失敗。
- #### 解決策:
    - まず Certbot で証明書を取得する
        - HTTP-01 チャレンジが通る環境 を作成
        - ホスト nginx を停止する（ポート80を空ける）
        - Docker nginx をポート80で起動
        - webroot パスを正しくマウントして Certbot を実行
    - 証明書取得後に nginx を SSL モードで起動
- #### 手順:
    1. ホストのポート80を空ける
    ```bash
    sudo systemctl stop nginx
    sudo systemctl disable nginx
    ```
    1. Docker nginx を HTTP だけで起動する
        - nginx 設定ファイルを一時的に SSL設定なしのHTTP-only にする
    ```nginx
    server {
        listen 80;
        server_name blog-api.duckdns.org;

        location / {
            proxy_pass http://go_app:8080;
        }

        location /.well-known/acme-challenge/ {
            root /var/www/certbot;
        }
    }
    ```
    1. Certbot で証明書を取得
        - webroot をマウントしていることを確認
    ```yml
    volumes:
        - ./certbot/www:/var/www/certbot
    ```
    1. 証明書取得コマンドを実行
    ```bash
    docker compose -f docker-compose.prod.yml run --rm certbot certonly \
    --webroot --webroot-path=/var/www/certbot \
    --agree-tos --no-eff-email \
    -m hoguro4649@yahoo.co.jp \
    -d blog-api.duckdns.org
    ```
    1. nginx を SSL設定に切り替えて再起動
        - 元の SSL 設定（`ssl_certificate` / `ssl_certificate_key`）に戻す
- #### 結果:
    - 上記の手順を実行することで`nginx`を無事に起動できた

---

### 問題7:SSL証明書の有効期限が切れていた

- #### 状況:
    - AWS上でBlogAPIをしたあとブラウザでアクセスすると警告が表示された
- #### 原因:
    - SSL証明書（httpsで通信するための鍵）は有効期限が約90日
    - 期限切れになっていたので「このサイトは危険である」と警告が表示された
- #### 解決策:
    - SSL証明書を更新する
- #### 手順:
    1. Certbotを使用して証明書を更新するシェルを作成する
        - certbotコンテナ起動時に更新を要求する（単純な起動では更新処理は実行されないので`renew`を指定して起動する）
        - プロキシーサーバーと共有箇所にある証明書が更新されるのでプロキシーサーバーを再起動して読み込む
        - ログとして更新した証明書情報を出力
        - 上記をシェル化すること（`utilitys/certbot/renew_cert.sh`）
    1. GitHub Actions のWorkflowでスケジュール設定してそのシェルがAWS内で実行されるようにする



## 4.学び・今後に活かすこと
