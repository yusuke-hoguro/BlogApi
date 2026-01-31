#!/usr/bin/env bash
set -euo pipefail

echo "=== Starting domain certificate renewal ==="
cd ~/BlogApi
# certbotコンテナで証明書更新を実行
docker compose -f docker-compose.prod.yml --env-file ./.env run --rm certbot renew --quiet
# nginxコンテナに証明書更新を通知
docker compose -f docker-compose.prod.yml exec -T nginx nginx -s reload
# 更新された証明書の情報を表示
docker compose -f docker-compose.prod.yml --env-file ./.env run --rm certbot certificates
echo "=== Domain certificate renewal finished ==="
