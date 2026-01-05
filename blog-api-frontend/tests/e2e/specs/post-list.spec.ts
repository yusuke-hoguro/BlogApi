import { test, expect } from '@playwright/test';
import { loginAsTestUser, createPost, deletePost } from '@e2e/utils/utils';
import { WAIT_FOR_ELEMENT_TIMEOUT_MS } from '@e2e/constants/config';
import { POST_ITEM_TEST_ID, POST_FETCH_ERROR_TEST_ID, POST_EMPTY_TEST_ID } from '@e2e/constants/selectors';
import { TEST_USERS } from '@e2e/fixtures/users';
import { CREATE_POST_TITLE, CREATE_POST_CONTENT } from '@e2e/constants/posts';
import { BUTTON_LOGOUT } from '@e2e/constants/buttons';
import { PAGE_TITLE_LOGIN } from '@e2e/constants/pageTitles';

test.describe('投稿一覧表示画面：正常系テスト', () => {
    // PostListが正しくAPIを叩けて表示できるかを確認する
    test('投稿一覧が表示される', async({ page }) => {
        // テストユーザーでログインする
        const token = await loginAsTestUser(page, TEST_USERS.testuser)
        await page.goto('/');
        // APIを使用してテスト用の投稿を作成する
        const testTitle = CREATE_POST_TITLE + `${Date.now()}`;
        const testContent = CREATE_POST_CONTENT
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
});

test.describe('投稿一覧表示画面：異常系テスト', () => {
    // PostListが正しくAPIを叩けて表示できるかを確認する
    test('API がサーバーエラー(500)を返した場合の表示確認', async({ page }) => {
        // リクエストをフックする
        await page.route('**/api/posts', async route => {
            // サーバーにはリクエストを送らずに500のレスポンスを返す 
            return route.fulfill({
                status: 500,
                body: JSON.stringify({ message: 'Failed to fetch posts'}),
            });
        });
        // テストユーザーでログインする
        await loginAsTestUser(page, TEST_USERS.testuser)
        await page.goto('/');
        await expect(page.getByTestId(POST_FETCH_ERROR_TEST_ID)).toBeVisible();
    });

    test('投稿が0件の場合に表示される文言確認', async({ page }) => {
        // リクエストをフックして投稿が1件も無いレスポンスを返す
        await page.route('**/api/posts', async route => {
            return route.fulfill({
                status: 200,
                body: JSON.stringify([]),
            })
        });
        // テストユーザーでログインする
        await loginAsTestUser(page, TEST_USERS.testuser)
        await page.goto('/');
        await expect(page.getByTestId(POST_EMPTY_TEST_ID)).toBeVisible();
    });

    test('未ログイン時のページ表示テスト', async ({ page }) => {
        // トップページへ遷移
        await page.goto('/');
        // URLを確認する
        await expect(page).toHaveURL('/login');
        // ログインページが開いているかを確認する
        await expect(page.getByRole('heading', { name: PAGE_TITLE_LOGIN })).toBeVisible();
        // テストユーザーでログインする
        await loginAsTestUser(page, TEST_USERS.testuser)
        await page.goto('/');
        // ログアウト
        await page.getByRole('button', { name: BUTTON_LOGOUT }).click();
        // URLを確認する
        await expect(page).toHaveURL('/login');
        // ログインページが開いているかを確認する
        await expect(page.getByRole('heading', { name: PAGE_TITLE_LOGIN })).toBeVisible();
    });    
});
