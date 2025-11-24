import { test, expect } from '@playwright/test';
import { loginAsTestUser } from './utils';

test('コメント一覧が表示される', async({ page }) => {
    // テストユーザーでログイン
    await loginAsTestUser(page)
    // トップページ（投稿一覧表示）へ遷移する
    await page.goto('http://localhost:3000/');
    // 最初の投稿のリンクを取得してクリックし、詳細ページへ遷移する
    const firstPost = page.getByTestId('post-item').first().locator('a', { hasText: /./ });
    await firstPost.click();
    // 投稿のコメントを取得する
    const comments = page.getByTestId('comment-item');
    await expect(comments.first()).toBeVisible();
});
