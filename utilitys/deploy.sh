#!/bin/bash
set -e

echo "=== Starting deployment ==="
cd ~/BlogApi
# 最新コードを取得
git pull origin main
# Docker 再起動
docker compose -f docker-compose.prod.yml down
docker compose -f docker-compose.prod.yml up -d --build
echo "=== Deployment finished ==="