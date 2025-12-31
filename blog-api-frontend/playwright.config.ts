import { defineConfig } from '@playwright/test';

export default defineConfig({
    testDir: './tests',
    use: {
        headless: true,
        // フロントエンドはポート3000で動作するdocker-composeのnginxコンテナを通じて提供
        baseURL: 'http://localhost:3000',
    },
    webServer: {
        // フロントエンドのE2Eテスト用にdocker-composeでコンテナを起動
        command: 'docker compose -f ../docker-compose.yml up --build frontend',
        url: 'http://localhost:3000',
        reuseExistingServer: true,
        timeout: 120_000,
    },
});
