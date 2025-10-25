#!/bin/bash
set -e

echo "=== Starting deployment ==="
cd ~/BlogApi
git fetch origin
git checkout -f main
# 最新コードを取得
git pull origin main

# フロントエンドのビルドを実施
cd blog-api-frontend
npm install
npm run build
cd ..

# Docker 再起動
docker compose -f docker-compose.prod.yml down
docker compose -f docker-compose.prod.yml up -d --build
echo "=== Deployment finished ==="

# Docker環境のクリーン（容量問題で暫定追加）
echo "=== Cleaning Docker environment ==="
docker system prune -af || true
echo "=== Cleaned ==="
