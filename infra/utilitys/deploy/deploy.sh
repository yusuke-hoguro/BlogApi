#!/bin/bash
set -e

APP_DIR="$HOME/BlogApi"
COMPOSE_FILE="$APP_DIR/infra/docker-compose.prod.yml"
ENV_FILE="$APP_DIR/.env"
COMPOSE=(docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE")

# 指定パスの空き容量を取得する
function get_free_space() {
	df -BG "$1" | awk 'NR==2 {gsub(/G/,"",$4); print $4}'
}

# DBの起動完了を待つ関数
function wait_db(){
	echo "Waiting for DB..."
	local max_attempts=60
	local sleep_seconds=2
	local attempt=1

	# コマンドが成功するまで待つ
	until "${COMPOSE[@]}" exec -T db pg_isready -U postgres >/dev/null 2>&1; do
		if [ "$attempt" -ge "$max_attempts" ]; then
			echo "DB did not become ready after $((max_attempts * sleep_seconds)) seconds"
			return 1
		fi

		echo "  ...still waiting (${attempt}/${max_attempts})"
		attempt=$((attempt + 1))
		sleep "$sleep_seconds"
	done
	echo "DB is ready"
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

# BlogAPIの最新コードを取得する
function update_repository(){
	cd "$APP_DIR"
	git fetch origin
	git checkout -f main
	git pull origin main
}

# フロントエンドのビルドを実行する
function build_frontend(){
	echo "Building frontend start: $(date)"
	cd "$APP_DIR/blog-api-frontend"
	npm install
	npm run build
	echo "Building frontend end: $(date)"
	# Node modules削除して容量確保
	rm -rf node_modules
	npm cache clean --force || true
	cd "$APP_DIR"
}

# Dockerコンテナを停止する
function stop_containers() {
	echo "Stopping Docker containers"
	"${COMPOSE[@]}" down --remove-orphans
}

# appコンテナのDockerfileをビルドしてイメージ作成
function build_app_image() {
	echo "Building Docker images start: $(date)"
	"${COMPOSE[@]}" build app
	echo "Building Docker images end: $(date)"
}

# Postgress DBを起動する
function start_db(){
	"${COMPOSE[@]}" up -d db
	wait_db	
}

# DBマイグレーションを実行する
function run_migrations(){
	echo "Running DB migrations start: $(date)"
	"${COMPOSE[@]}" run --rm app ./migrate
	echo "Running DB migrations end: $(date)"
}

# Deploy後にBlogAPIを起動する
function start_services(){
	echo "Deploying Docker containers start: $(date)"
	"${COMPOSE[@]}" up -d
	echo "Deployment finished: $(date)"
}

# ビルド後のキャッシュを削除する
function cleanup_docker_cache(){
	echo "Post-clean (builder cache)"
	docker builder prune -af || true
	echo "Done"
}

# AWSへデプロイを実施する関数
function aws_deploy(){
	echo "Starting deployment"
	# デプロイ前に空き容量を確認
	ensure_disk_space
	# アプリケーションディレクトリへ移動して最新コードを取得
	update_repository
	# フロントエンドのビルドを実施
	build_frontend
	# Docker 再起動
	stop_containers
	# 空き容量を再チェック
	ensure_disk_space
	# App用のDockerfileをビルドする
	build_app_image
	# Postgres DB用のコンテナを起動する
	start_db
	# DBマイグレーションを実行する
	run_migrations
	# デプロイと起動を実行
	start_services
	# ビルドキャッシュだけ削除
	cleanup_docker_cache
}

# main
# AWSへデプロイを実行
aws_deploy
