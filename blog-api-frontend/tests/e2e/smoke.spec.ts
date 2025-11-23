import { test, expect } from '@playwright/test';
import { loginAsTestUser } from './utils';

test('トップページが開ける', async ({ page }) => {
  // テストユーザーでログインする
  await loginAsTestUser(page)
  // トップページにアクセスする
  await page.goto('http://localhost:3000');
  // 投稿一覧の見出しが表示されるまで待つ
  await page.waitForSelector('h1', { timeout: 20000 });
  await expect(page).toHaveTitle(/Blog/);
  // トップページのタイトルを確認する
  await expect(page.getByRole('heading', { name: '投稿一覧' })).toBeVisible();
});

