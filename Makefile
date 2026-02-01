# 環境変数の読み込み
include .env
export

# Makefile
# compose fileのパス
COMPOSE_DEV_FILE=infra/docker-compose.yml
COMPOSE_PROD_FILE=infra/docker-compose.prod.yml
COMPOSE_TEST_FILE=infra/docker-compose.test.yml
ENV_FILE=.env

# 共通変数を定義し、使いやすくする
DC=docker compose --env-file $(ENV_FILE)
DEV=$(DC) -f $(COMPOSE_DEV_FILE)
PROD=$(DC) -f $(COMPOSE_PROD_FILE)
TEST=$(DC) -f $(COMPOSE_TEST_FILE)


# ターゲット定義
.PHONY: help \
		up-dev down-dev down-volumes-dev restart-dev logs-dev build-dev rebuild-dev ps-dev \
		up-prod down-prod restart-prod logs-prod build-prod rebuild-prod ps-prod \
		up-test down-test down-volumes-test restart-test logs-test build-test rebuild-test ps-test \
		test-go test-e2e ci-test wait-test-db

help:
	@echo "Makefile commands:"
	@echo "  make up-dev               - Start development environment"
	@echo "  make down-dev             - Stop development environment"
	@echo "  make down-volumes-dev     - Stop development environment and remove volumes"
	@echo "  make restart-dev          - Restart development environment"
	@echo "  make logs-dev             - View logs of development environment"
	@echo "  make build-dev            - Build images for development environment"
	@echo "  make rebuild-dev          - Rebuild images for development environment without cache"
	@echo "  make ps-dev               - Show status of containers in development environment"
	@echo ""
	@echo "  make up-prod              - Start production environment"
	@echo "  make down-prod            - Stop production environment"
	@echo "  make restart-prod         - Restart production environment"
	@echo "  make logs-prod            - View logs of production environment"
	@echo "  make build-prod           - Build images for production environment"
	@echo "  make rebuild-prod         - Rebuild images for production environment without cache"
	@echo "  make ps-prod              - Show status of containers in production environment"
	@echo ""
	@echo "  make up-test              - Start test environment"
	@echo "  make down-test            - Stop test environment"
	@echo "  make down-volumes-test    - Stop test environment and remove volumes"
	@echo "  make restart-test         - Restart test environment"
	@echo "  make logs-test            - View logs of test environment"
	@echo "  make build-test           - Build images for test environment"
	@echo "  make rebuild-test         - Rebuild images for test environment without cache"
	@echo "  make ps-test              - Show status of containers in test environment"
	@echo ""
	@echo "  make test-go              - Run backend handler function tests"
	@echo "  make test-e2e             - Run E2E tests"
	@echo "  make ci-test              - Run all CI tests"	

# 開発環境用
# 起動
up-dev:
	$(DEV) up -d
# 停止
down-dev:
	$(DEV) down
# 停止時にボリュームも削除
down-volumes-dev:
	$(DEV) down -v
# 再起動
restart-dev: down-dev up-dev
# ログの確認
logs-dev:
	$(DEV) logs -f
# イメージのビルド
build-dev:
	$(DEV) build
# キャッシュを使わずにイメージをビルド
rebuild-dev:
	$(DEV) build --no-cache
# コンテナの状態確認
ps-dev:
	$(DEV) ps

# 本番環境用
# 起動
up-prod:
	$(PROD) up -d
# 停止
down-prod:
	$(PROD) down
# 再起動
restart-prod: down-prod up-prod
# ログの確認
logs-prod:
	$(PROD) logs -f
# イメージのビルド
build-prod:
	$(PROD) build
# キャッシュを使わずにイメージをビルド
rebuild-prod:
	$(PROD) build --no-cache
# コンテナの状態確認
ps-prod:
	$(PROD) ps

# テスト環境用
# 起動
up-test:
	$(TEST) up -d
# 停止
down-test:
	$(TEST) down
# 停止時にボリュームも削除
down-volumes-test:
	$(TEST) down -v
# 再起動
restart-test: down-test up-test
# ログの確認
logs-test:
	$(TEST) logs -f
# イメージのビルド
build-test:
	$(TEST) build
# キャッシュを使わずにイメージをビルド
rebuild-test:
	$(TEST) build --no-cache
# コンテナの状態確認
ps-test:
	$(TEST) ps

# バックンド ハンドラー関数のテスト実行
test-go:
	@set -e; \
	trap '$(MAKE) down-volumes-test' EXIT; \
	$(MAKE) down-volumes-test; \
	$(MAKE) build-test;	\
	$(MAKE) up-test; \
	$(MAKE) wait-test-db; \
	go test ./internal/handler/... -v

# E2Eテスト実行
test-e2e:
	@set -e; \
	trap 'cd $(CURDIR) && $(MAKE) down-volumes-dev' EXIT; \
	$(MAKE) down-volumes-dev; \
	$(MAKE) build-dev; \
	$(MAKE) up-dev; \
	$(MAKE) wait-dev-db; \
	cd blog-api-frontend && npx playwright test

# CIのテストをまとめて実行
ci-test: test-go test-e2e

# テスト用DBの起動待ち（5秒待機）
wait-test-db:
	@echo "Waiting for test DB..."
	@until docker exec postgres_test pg_isready -U  $(DB_USER) >/dev/null 2>&1; do \
		echo "  ...still waiting"; \
		sleep 1; \
	done

wait-dev-db:
	@echo "Waiting for dev DB..."
	@until docker exec postgres pg_isready -U  $(DB_USER) >/dev/null 2>&1; do \
		echo "  ...still waiting"; \
		sleep 1; \
	done
