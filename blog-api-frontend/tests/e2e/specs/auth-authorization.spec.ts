import { test, expect } from '@playwright/test';
import { TEST_USERS } from '@e2e/fixtures/users';
import { loginAsTestUser, createPost } from '@e2e/utils/utils';
import { BUTTON_LOGOUT, BUTTON_DELETE_POST, BUTTON_CREATE_POST } from '@e2e/constants/buttons';
import { CREATE_POST_TITLE, CREATE_POST_CONTENT } from '@e2e/constants/posts';
import { POST_ITEM_TEST_ID, } from '@e2e/constants/selectors';
import { PAGE_TITLE_POST_LIST, PAGE_TITLE_POST_CREATE } from '@e2e/constants/pageTitles';
import { LABEL_POST_CREATE_TITLE, LABEL_POST_CREATE_CONTEXT } from '@e2e/constants/label';


test.describe('認証・認可のテスト', () => {
    test('未ログインで投稿作成ページへアクセスするテスト', async ({ page }) => {
        // 未ログインで投稿作成ページへ遷移
        await page.goto('/post/create');
        // URLを確認する
        await expect(page).toHaveURL('/login');
        // ログイン後、ログアウトでも同様かを確認する
        const token = await loginAsTestUser(page, TEST_USERS.testuser)
        await page.goto('/');
        // ログアウト
        await page.getByRole('button', { name: BUTTON_LOGOUT }).click();
        // ログアウトの完了を待つ（token 削除まで待機）
        await page.waitForFunction(() => {
            return localStorage.getItem('token') === null;
        });
        // 未ログインで投稿作成ページへ遷移
        await page.goto('/post/create');
        // URLを確認する
        await expect(page).toHaveURL('/login');
    });

    test('未ログインで投稿詳細画面へアクセスするテスト', async ({ page }) => {
        // テストユーザーでログイン
        const token = await loginAsTestUser(page, TEST_USERS.testuser)
        // トップページ（投稿一覧表示）へ遷移する
        await page.goto('/');
        // APIを使用してテスト用の投稿を作成する
        const testTitle = CREATE_POST_TITLE + `${Date.now()}`;
        const testContent = CREATE_POST_CONTENT
        const post = await createPost(page, token, testTitle, testContent)
        // ログアウト
        await page.getByRole('button', { name: BUTTON_LOGOUT }).click();
        // ログアウトの完了を待つ（token 削除まで待機）
        await page.waitForFunction(() => {
            return localStorage.getItem('token') === null;
        });
        // 未ログインで投稿詳細ページへ遷移
        await page.goto(`/post/${post.id}`);
        // URLを確認する
        await expect(page).toHaveURL('/login');
        // 投稿を削除するために投稿作成ユーザーで再度ログイン
        await loginAsTestUser(page, TEST_USERS.testuser)
        // トップページに遷移
        await page.goto('/');
        // 投稿一覧の中に新規作成した投稿があることを確認する
        await expect(page.getByTestId(POST_ITEM_TEST_ID).filter({ hasText: testTitle})).toHaveCount(1)
        // テスト用に追加した投稿の詳細画面を開く
        await page.getByRole('link', { name: testTitle }).click();
        // テスト用に追加した投稿の詳細画面に遷移できたかをチェック
        await expect(page.getByRole('heading', { name: testTitle })).toBeVisible();
        // 削除ボタンをクリックして投稿を削除
        const button = page.getByRole('button', { name: BUTTON_DELETE_POST });
        // クリックによってconfirmがでることを想定してイベントハンドラをセットしておく
        page.once('dialog', async dialog =>{
            console.log(`Dialog message: ${dialog.message()}`); //デバッグ用ログ
            await dialog.accept();
        });
        await button.click();
        // 投稿一覧ページに遷移したかを確認する
        await expect(page.getByRole('heading', { name: PAGE_TITLE_POST_LIST })).toBeVisible();
        // 投稿をすべて取得する
        const posts = page.getByTestId(POST_ITEM_TEST_ID);
        // テスト用の投稿が削除されたことを確認する
        await expect(posts.filter({ hasText: testTitle })).toHaveCount(0);
    });

    test('トークンが期限切れになった場合のテスト', async ({ page }) => {
        // テストユーザーでログイン
        const token = await loginAsTestUser(page, TEST_USERS.testuser)
        // トップページにアクセスする
        await page.goto('/');
        // トップページのタイトルを確認する
        await expect(page.getByRole('heading', { name: PAGE_TITLE_POST_LIST })).toBeVisible();        
        // 不正なTokenに書き換える
        await page.evaluate(() => {
            localStorage.setItem('token', 'invalid.token.value');
        });
        // 新規投稿ページへ遷移
        await page.goto('/post/create');
        // 投稿作成ページが開けたかを確認する
        await expect(page.getByRole('heading', { name: PAGE_TITLE_POST_CREATE })).toBeVisible();
        // 新規投稿の作成を実施する
        const title = CREATE_POST_TITLE + `${Date.now()}`;
        const content = CREATE_POST_CONTENT;
        await page.getByLabel(LABEL_POST_CREATE_TITLE).fill(title);
        await page.getByLabel(LABEL_POST_CREATE_CONTEXT).fill(content);
        await page.getByRole('button', { name: BUTTON_CREATE_POST }).click();
        // ログインへ戻される
        await expect(page).toHaveURL('/login');
    });
});
