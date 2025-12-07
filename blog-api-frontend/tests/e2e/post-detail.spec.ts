import { test, expect } from '@playwright/test';
import { loginAsTestUser } from './utils';
import { TEST_USERS } from './users';
import { POST_ITEM_TEST_ID, POST_TITLE_TEST_ID, POST_CONTENT_TEST_ID } from './constants/selectors';


test('投稿詳細ページへ遷移でき、詳細が表示される', async ({ page }) => {
    // テストユーザーでログイン
    await loginAsTestUser(page, TEST_USERS.testuser)
    // トップページ（投稿一覧表示）へ遷移する
    await page.goto('/');
    // 最初の投稿のリンクを取得してクリックし、詳細ページへ遷移する
    const firstPost = page.getByTestId(POST_ITEM_TEST_ID).first().locator('a', { hasText: /./ });
    await firstPost.click();
    // 正しいURLに遷移しているかを確認する
    await expect(page).toHaveURL(/\/post\/\d+$/)
    // 投稿のタイトルが表示されていることを確認する
    await expect(page.getByTestId(POST_TITLE_TEST_ID)).toBeVisible()
    // 投稿の内容が表示されていることを確認する
    await expect(page.getByTestId(POST_CONTENT_TEST_ID)).toBeVisible()
});
