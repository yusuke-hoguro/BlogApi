import { test, expect } from '@playwright/test';

test('トップページが開ける', async ({ page }) => {
    await page.goto('/');
    await expect(page).toHaveTitle(/.*/);
});
