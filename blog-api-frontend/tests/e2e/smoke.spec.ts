import { test, expect } from '@playwright/test';
import { loginAsTestUser } from './utils';
import { WAIT_FOR_ELEMENT_TIMEOUT_MS } from './constants';
import { TEST_USERS } from './users';

test('トップページが開ける', async ({ page }) => {
  // テストユーザーでログインする
  await loginAsTestUser(page, TEST_USERS.testuser)
  // トップページにアクセスする
  await page.goto('/');
  // 投稿一覧の見出しが表示されるまで待つ
  await page.waitForSelector('h1', { timeout: WAIT_FOR_ELEMENT_TIMEOUT_MS });
  await expect(page).toHaveTitle(/Blog/);
  // トップページのタイトルを確認する
  await expect(page.getByRole('heading', { name: '投稿一覧' })).toBeVisible();
});

