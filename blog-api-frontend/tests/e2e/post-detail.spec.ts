import { test, expect } from '@playwright/test';
import { loginAsTestUser } from './utils';

test('投稿詳細ページへ遷移でき、詳細が表示される', async ({ page }) => {
    // テストユーザーでログイン
    await loginAsTestUser(page)
    // トップページ（投稿一覧表示）へ遷移する
    await page.goto('http://localhost:3000/');
    // 最初の投稿のリンクを取得してクリックし、詳細ページへ遷移する
    const firstPost = page.getByTestId('post-item').first().locator('a', { hasText: /./ });
    await firstPost.click();
    // 正しいURLに遷移しているかを確認する
    await expect(page).toHaveURL(/\/post\/\d+$/)
    // 投稿のタイトルが表示されていることを確認する
    await expect(page.getByTestId('post-title')).toBeVisible()
    // 投稿の内容が表示されていることを確認する
    await expect(page.getByTestId('post-content')).toBeVisible()
});
