import { test, expect } from '@playwright/test';
import { loginAsTestUser } from './utils';
import { WAIT_FOR_ELEMENT_TIMEOUT_MS } from './constants/config';
import { POST_ITEM_TEST_ID } from './constants/selectors';
import { TEST_USERS } from './users';

// PostListが正しくAPIを叩けて表示できるかを確認する
test('投稿一覧が表示される', async({ page }) => {
    // テストユーザーでログインする
    await loginAsTestUser(page, TEST_USERS.testuser)
    await page.goto('/');
    // 投稿リストが描画されるまで最大10秒待つ
    await page.waitForSelector(`[data-testid='${POST_ITEM_TEST_ID}']`, { timeout: WAIT_FOR_ELEMENT_TIMEOUT_MS });
    // 要素のうち、data-testid="post-item" が付いたものを取得し最初の要素を選択
    const firstPost = page.getByTestId(POST_ITEM_TEST_ID).first();
    // 投稿リストが表示されているかを確認する
    await expect(firstPost).toBeVisible();
    // 投稿内の見出し要素（h1など）を取得し、空でないことを確認する
    await expect(firstPost.getByRole('heading')).not.toBeEmpty();
});
