import { test, expect } from '@playwright/test';
import { loginAsTestUser, createPost, deletePost, logout } from './utils';
import { TEST_USERS } from './users';
import { POST_ITEM_TEST_ID, POST_TITLE_TEST_ID, POST_CONTENT_TEST_ID } from './constants/selectors';
import { CREATE_POST_TITLE, CREATE_POST_CONTENT } from './constants/posts';
import { BUTTON_DELETE_POST } from './constants/buttons';
import { PAGE_TITLE_POST_LIST } from './constants/pageTitles';
import { LABEL_EDIT_POST } from './constants/label';


test.describe('投稿詳細表示画面：正常系テスト', () => {

    test('投稿詳細ページへ遷移でき、詳細が表示される', async ({ page }) => {
        // テストユーザーでログイン
        const token = await loginAsTestUser(page, TEST_USERS.testuser)
        // APIを使用してテスト用の投稿を作成する
        const testTitle = CREATE_POST_TITLE + `${Date.now()}`;
        const testContent = CREATE_POST_CONTENT
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
});

test.describe('投稿詳細表示画面：異常系テスト', () => {
    test('存在しない投稿IDを選択した場合のテスト', async ({ page }) => {
        // テストユーザーでログイン
        const token = await loginAsTestUser(page, TEST_USERS.testuser)
        // トップページ（投稿一覧表示）へ遷移する
        await page.goto('/');
        // リクエストをフックして投稿が存在しないエラーを返す
        await page.route('**/api/posts/*', async route => {
            return route.fulfill({ 
                status: 404,
                body: JSON.stringify({ message: 'Post not found'}),
            });
        });
        await page.goto('/post/999999');
        // 画面に「投稿がない」ときの表示がただしく表示されているか確認
        await expect(page.getByText('投稿が見つかりません')).toBeVisible();
    });

    test('投稿取得APIが500の場合のテスト', async ({ page }) => {
        // テストユーザーでログイン
        const token = await loginAsTestUser(page, TEST_USERS.testuser)
        // APIを使用してテスト用の投稿を作成する
        const testTitle = CREATE_POST_TITLE + `${Date.now()}`;
        const testContent = CREATE_POST_CONTENT
        const post = await createPost(page, token, testTitle, testContent)
        // トップページ（投稿一覧表示）へ遷移する
        await page.goto('/');
        // 投稿一覧の中に新規作成した投稿があることを確認する
        await expect(page.getByTestId(POST_ITEM_TEST_ID).filter({ hasText: testTitle})).toHaveCount(1)
        // リクエストをフックして投稿が存在しないエラーを返す
        await page.route('**/api/posts/*', async route => {
            return route.fulfill({ 
                status: 500,
                body: JSON.stringify({ message: 'Database error'}),
            });
        });
        await page.goto(`/post/${post.id}`);
        // 画面に「投稿がない」ときの表示がただしく表示されているか確認
        await expect(page.getByText('投稿が見つかりません')).toBeVisible();
    });

    test('他ユーザーが作成した投稿に対して編集・削除ができないことをテスト', async ({ page }) => {
        // テストユーザーでログイン
        const token = await loginAsTestUser(page, TEST_USERS.testuser)
        // APIを使用してテスト用の投稿を作成する
        const testTitle = CREATE_POST_TITLE + `${Date.now()}`;
        const testContent = CREATE_POST_CONTENT
        const post = await createPost(page, token, testTitle, testContent)
        // トップページ（投稿一覧表示）へ遷移する
        await page.goto('/');
        // 投稿一覧の中に新規作成した投稿があることを確認する
        await expect(page.getByTestId(POST_ITEM_TEST_ID).filter({ hasText: testTitle})).toHaveCount(1)
        // ログアウト
        await logout(page);
        // 違うユーザーでログインする
        await loginAsTestUser(page, TEST_USERS.otheruser)
        // トップページ（投稿一覧表示）へ遷移する
        await page.goto('/');
        // 投稿一覧の中に新規作成した投稿があることを確認する
        await expect(page.getByTestId(POST_ITEM_TEST_ID).filter({ hasText: testTitle})).toHaveCount(1)
        // 新規追加した投稿の詳細画面を開く
        await page.getByRole('link', { name: testTitle }).click();
        // 新規追加した投稿の詳細画面に遷移できたかをチェック
        await expect(page.getByRole('heading', { name: testTitle })).toBeVisible();
        // 投稿編集が表示されていないことを確認
        await expect(page.getByRole('link', { name: LABEL_EDIT_POST })).toHaveCount(0);
        // 投稿削除ボタンが表示されていないことを確認
        await expect(page.getByRole('button', { name: BUTTON_DELETE_POST })).toHaveCount(0);
        // ログアウト
        await logout(page);
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
        // 削除ボタンをクリックしてコメントを削除
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

    test('削除済み投稿にアクセスした場合のテスト', async ({ page }) => {
        // テストユーザーでログイン
        const token = await loginAsTestUser(page, TEST_USERS.testuser)
        // トップページ（投稿一覧表示）へ遷移する
        await page.goto('/');
        // APIを使用してテスト用の投稿を作成する
        const testTitle = CREATE_POST_TITLE + `${Date.now()}`;
        const testContent = CREATE_POST_CONTENT
        const post = await createPost(page, token, testTitle, testContent)
        // トップページ（投稿一覧表示）へ遷移する
        await page.goto('/');
        // 投稿一覧の中に新規作成した投稿があることを確認する
        await expect(page.getByTestId(POST_ITEM_TEST_ID).filter({ hasText: testTitle})).toHaveCount(1)
        // テスト用に追加した投稿の詳細画面を開く
        await page.getByRole('link', { name: testTitle }).click();
        // テスト用に追加した投稿の詳細画面に遷移できたかをチェック
        await expect(page.getByRole('heading', { name: testTitle })).toBeVisible();
        // 削除ボタンをクリックしてコメントを削除
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
        await page.goto(`/post/${post.id}`);
        // 画面に「投稿がない」ときの表示がただしく表示されているか確認
        await expect(page.getByText('投稿が見つかりません')).toBeVisible();
    });
});
