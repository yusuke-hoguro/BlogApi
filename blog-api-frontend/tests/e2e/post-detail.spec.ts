import { test, expect } from '@playwright/test';
import { loginAsTestUser, createPost, deletePost } from './utils';
import { TEST_USERS } from './users';
import { POST_ITEM_TEST_ID, POST_TITLE_TEST_ID, POST_CONTENT_TEST_ID } from './constants/selectors';
import { CREAT_POST_TITLE, CREAT_POST_CONTENT } from './constants/posts';


test('投稿詳細ページへ遷移でき、詳細が表示される', async ({ page }) => {
    // テストユーザーでログイン
    const token = await loginAsTestUser(page, TEST_USERS.testuser)
    // APIを使用してテスト用の投稿を作成する
    const testTitle = CREAT_POST_TITLE + `${Date.now()}`;
    const testContent = CREAT_POST_CONTENT
    const post = await createPost(page, token, testTitle, testContent)
    // トップページ（投稿一覧表示）へ遷移する
    await page.goto('/');
    // 投稿一覧の中に新規作成した投稿があることを確認する
    await expect(page.getByTestId(POST_ITEM_TEST_ID).filter({ hasText: testTitle})).toHaveCount(1)
    // 新規追加した投稿の詳細画面を開く
    await page.getByRole('link', { name: testTitle }).click();
    // 新規追加した投稿の詳細画面に遷移できたかをチェック
    await expect(page.getByRole('heading', { name: testTitle })).toBeVisible();
    // 正しいURLに遷移しているかを確認する
    await expect(page).toHaveURL(/\/post\/\d+$/)
    // 投稿のタイトルが表示されていることを確認する
    await expect(page.getByTestId(POST_TITLE_TEST_ID)).toBeVisible()
    // 投稿の内容が表示されていることを確認する
    await expect(page.getByTestId(POST_CONTENT_TEST_ID)).toBeVisible()
    // APIを使用して投稿を削除する
    await deletePost(page, token, post.id)
    // トップページ（投稿一覧表示）へ遷移する
    await page.goto('/');
    // 投稿をすべて取得する
    const posts = page.getByTestId(POST_ITEM_TEST_ID);
    // テスト用の投稿が削除されたことを確認する
    await expect(posts.filter({ hasText: testTitle })).toHaveCount(0);

});
