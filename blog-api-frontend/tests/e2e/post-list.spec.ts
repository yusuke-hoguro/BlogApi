import { test, expect } from '@playwright/test';
import { loginAsTestUser, createPost, deletePost } from './utils';
import { WAIT_FOR_ELEMENT_TIMEOUT_MS } from './constants/config';
import { POST_ITEM_TEST_ID } from './constants/selectors';
import { TEST_USERS } from './users';
import { CREAT_POST_TITLE, CREAT_POST_CONTENT } from './constants/posts';

// PostListが正しくAPIを叩けて表示できるかを確認する
test('投稿一覧が表示される', async({ page }) => {
    // テストユーザーでログインする
    const token = await loginAsTestUser(page, TEST_USERS.testuser)
    await page.goto('/');
    // APIを使用してテスト用の投稿を作成する
    const testTitle = CREAT_POST_TITLE + `${Date.now()}`;
    const testContent = CREAT_POST_CONTENT
    const post = await createPost(page, token, testTitle, testContent)
    // 投稿作成後に一覧をリロード
    await page.goto('/');
    // 投稿リストが描画されるまで最大10秒待つ
    await page.waitForSelector(`[data-testid='${POST_ITEM_TEST_ID}']`, { timeout: WAIT_FOR_ELEMENT_TIMEOUT_MS });
    // テスト用に追加した投稿の要素を取得する
    const checkPost = page.getByTestId(POST_ITEM_TEST_ID).filter({ hasText: testTitle})
    // 投稿リストが表示されているかを確認する
    await expect(checkPost).toBeVisible();
    // 投稿内の見出し要素（h1など）を取得し、空でないことを確認する
    await expect(checkPost.getByRole('heading')).not.toBeEmpty();
    // APIを使用して投稿を削除する
    await deletePost(page, token, post.id)
    // トップページ（投稿一覧表示）へ遷移する
    await page.goto('/');
    // 投稿をすべて取得する
    const posts = page.getByTestId(POST_ITEM_TEST_ID);
    // テスト用の投稿が削除されたことを確認する
    await expect(posts.filter({ hasText: testTitle })).toHaveCount(0);
});
