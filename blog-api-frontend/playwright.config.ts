import { defineConfig, devices } from '@playwright/test';
import { register } from 'tsconfig-paths';

register({
    baseUrl: '.',
    paths: {
        '@e2e/*': ['tests/e2e/*'],
    },
});

export default defineConfig({
    testDir: './tests',
    projects: [
        {
            name: 'chromium',
            use: { ...devices['Desktop Chrome'] },
        },
        {
            name: 'firefox',
            use: { ...devices['Desktop Firefox'] },
        },
        {
            name: 'webkit',
            use: { ...devices['Desktop Safari'] },
        },
    ],
    use: {
        headless: true,
        // フロントエンドはポート3000で動作するdocker-composeのnginxコンテナを通じて提供
        baseURL: 'http://localhost:3000',
    },
    webServer: {
        // フロントエンドのE2Eテスト用にdocker-composeでコンテナを起動
        command: 'docker compose -f ../infra/docker-compose.yml --env-file ../.env up --build frontend',
        url: 'http://localhost:3000',
        reuseExistingServer: true,
        timeout: 120_000,
    },
    // テスト開始時に初期設定を実行する
    globalSetup: './tests/e2e/setup/global-setup',
});
