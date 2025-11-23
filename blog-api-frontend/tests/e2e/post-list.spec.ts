import { test, expect } from '@playwright/test';
import { loginAsTestUser } from './utils';

// PostListが正しくAPIを叩けて表示できるかを確認する
test('投稿一覧が表示される', async({ page }) => {
    // テストユーザーでログインする
    await loginAsTestUser(page)
    await page.goto('http://localhost:3000/');
    // 投稿リストが描画されるまで最大10秒待つ
    await page.waitForSelector('[data-testid="post-item"]', { timeout: 10000 });
    // 要素のうち、data-testid="post-item" が付いたものを取得し最初の要素を選択
    const firstPost = page.getByTestId('post-item').first();
    // 投稿リストが表示されているかを確認する
    await expect(firstPost).toBeVisible();
    // 投稿内の見出し要素（h1など）を取得し、空でないことを確認する
    await expect(firstPost.getByRole('heading')).not.toBeEmpty();
});
