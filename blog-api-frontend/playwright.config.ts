import { defineConfig } from '@playwright/test';

export default defineConfig({
    testDir: './tests',
    use: {
        headless: true,
        baseURL: 'http://localhost:3000',
    },
    webServer: {
        command: 'docker compose -f ../docker-compose.yml up --build frontend',
        url: 'http://localhost:3000',
        reuseExistingServer: true,
        timeout: 120_000,
    },
});
