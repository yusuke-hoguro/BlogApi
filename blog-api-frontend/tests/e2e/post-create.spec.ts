import { test, expect } from '@playwright/test';
import { loginAsTestUser } from './utils';
import { ASSERTION_TIMEOUT_MS } from './constants';

test.describe('投稿機能:正常系テスト', () => {

    test('新規投稿作成→編集→削除のテストを実施', async ({ page }) => {
        // テストユーザーでログインする
        const token = await loginAsTestUser(page);
        // 投稿作成ページへ遷移
        await page.goto('/post/create');
        // 投稿作成ページが開けたかを確認する
        await expect(page.getByRole('heading', { name: '新規投稿作成' })).toBeVisible();
        // 新規投稿の作成を実施する
        const title = `E2E 投稿タイトル ${Date.now()}`;
        const content = 'E2E 投稿本文です。';
        await page.getByLabel('タイトル').fill(title);
        await page.getByLabel('内容').fill(content);
        await page.getByRole('button', { name: '投稿作成' }).click();
        // 投稿一覧ページに遷移したかを確認する
        await expect(page.getByRole('heading', { name: '投稿一覧' })).toBeVisible();
        // 投稿一覧の中に新規作成した投稿があることを確認する
        await expect(page.getByTestId('post-item').filter({ hasText: title})).toHaveCount(1)
        // 新規追加した投稿の詳細画面を開く
        await page.getByRole('link', { name: title }).click();
        // 新規追加した投稿の詳細画面に遷移できたかをチェック
        await expect(page.getByRole('heading', { name: title })).toBeVisible();
        // 編集ボタンを押す
        await page.getByRole('link', { name: '投稿編集' }).click();
        // タイトルを編集する
        const updateTitle = `E2E 投稿タイトル Update ${Date.now()}`;
        await page.getByLabel('タイトル').fill(updateTitle);
        // タイトルが変わったか確認
        await expect(page.getByLabel('タイトル')).toHaveValue(updateTitle, { timeout: ASSERTION_TIMEOUT_MS })
        // 投稿内容を編集する
        const updateContext = `E2E 投稿本文です。更新しました。`;
        await page.getByLabel('内容').fill(updateContext);
        // 投稿内容が変わったか確認
        await expect(page.getByLabel('内容')).toHaveValue(updateContext, { timeout: ASSERTION_TIMEOUT_MS })
        // 更新ボタンを押下する
        await page.getByRole('button', { name: '更新' }).click();
        // 投稿のタイトルが表示されていることを確認する
        await expect(page.getByTestId('post-title')).toContainText(updateTitle, { timeout: ASSERTION_TIMEOUT_MS })
        // 投稿の内容が表示されていることを確認する
        await expect(page.getByTestId('post-content')).toContainText(updateContext, { timeout: ASSERTION_TIMEOUT_MS })
        // クリックによってconfirmがでることを想定してイベントハンドラをセット
        page.once('dialog', dialog => {
            console.log(`Dialog message: ${dialog.message()}`);
            dialog.accept();
        });
        // 投稿削除を実施
        await page.getByRole('button', { name: '投稿削除' }).click();
        // 投稿一覧に戻るまで待機
        await page.waitForURL('http://localhost:3000/');
        // 投稿一覧ページに遷移したかを確認する
        await expect(page.getByRole('heading', { name: '投稿一覧' })).toBeVisible();
        // 投稿をすべて取得する
        const posts = page.getByTestId('post-item');
        // タイトル一致しないことを確認する
        await expect(posts.filter({ hasText: updateTitle })).toHaveCount(0);
    });

});
