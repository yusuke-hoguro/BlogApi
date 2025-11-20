import { test, expect } from '@playwright/test';

test('トップページが開ける', async ({ page }) => {
  await page.goto('http://localhost:3000/');
  await expect(page).toHaveTitle(/Blog/);
  await expect(page.getByRole('heading', { name: 'Posts' })).toBeVisible();
});
