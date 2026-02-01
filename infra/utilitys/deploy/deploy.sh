#!/bin/bash
set -e

APP_DIR="$HOME/BlogApi"
COMPOSE_FILE="infra/docker-compose.prod.yml"

# 指定パスの空き容量を取得する
function get_free_space() {
	df -BG "$1" | awk 'NR==2 {gsub(/G/,"",$4); print $4}'
}

# デプロイの前に空き容量を確認する関数
function ensure_disk_space(){
	# 指定パスの空き容量（GB）
	local free_space
	free_space=$(get_free_space "$APP_DIR")
	# 必要な空き容量（GB）
	local required_space=1  

	echo "Free disk space: ${free_space}GB (need >= ${required_space}GB)"

	if [ "$free_space" -lt "$required_space" ]; then
		echo "Low disk. Cleaning docker caches before build"
		# ビルドキャッシュ削除
		docker builder prune -af || true
		# 未使用のイメージ削除
		docker image prune -af || true
		# 未使用のネットワーク削除
		docker network prune -f || true
		# まだ容量が足りない場合はsystem pruneも実施
		free_space=$(get_free_space "$APP_DIR")
		if [ "$free_space" -lt "$required_space" ]; then
			echo "Still low disk. Performing full docker system prune"
			docker system prune -af || true
		fi

		# 再度空き容量を確認
		free_space=$(get_free_space "$APP_DIR")
		echo "Free disk space after cleanup: ${free_space}GB"
	fi
}

# AWSへデプロイを実施する関数
function aws_deploy(){
	echo "Starting deployment"
	# デプロイ前に空き容量を確認
	ensure_disk_space
	# アプリケーションディレクトリへ移動して最新コードを取得
	cd ~/BlogApi
	git fetch origin
	git checkout -f main
	git pull origin main

	# フロントエンドのビルドを実施
	echo "Building frontend start: $(date)"
	cd blog-api-frontend
	npm install
	npm run build
	echo "Building frontend end: $(date)"
	# Node modules削除して容量確保
	rm -rf node_modules
	npm cache clean --force || true
	cd ..

	# Docker 再起動
	echo "Restarting Docker containers"
	# 既存のコンテナを停止・削除
	docker compose -f "$COMPOSE_FILE" --env-file ./.env down --remove-orphans
	# 空き容量を再チェック
	ensure_disk_space
	echo "Deploying Docker containers start: $(date)"
	# Docker イメージのビルドとコンテナの起動
	docker compose -f "$COMPOSE_FILE" --env-file ./.env up -d --build
	echo "Deployment finished: $(date)"

	# ビルドキャッシュだけ削除
	echo "Post-clean (builder cache)"
	docker builder prune -af || true
	echo "Done"
}

# main
# AWSへデプロイを実行
aws_deploy
