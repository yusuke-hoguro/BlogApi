import { test, expect } from '@playwright/test';
import { loginAsTestUser, logout } from './utils';
import { ASSERTION_TIMEOUT_MS, WAIT_FOR_ELEMENT_TIMEOUT_MS } from './constants/config';
import { TEST_COMMENT, TEST_COMMENT_LONG, TEST_COMMENT_TOO_LONG } from './constants/comments';
import { BUTTON_SEND_COMMENT, BUTTON_EDIT_COMMENT, BUTTON_SAVE_COMMENT, BUTTON_DELETE_COMMENT } from './constants/buttons';
import { COMMENT_ITEM_TEST_ID, POST_ITEM_TEST_ID } from './constants/selectors';
import { TEST_USERS } from './users';

test.describe('コメント機能：正常系テスト', () => {
    
    test('UI画面でコメントの作成→表示→編集→削除のテストを実施する', async({ page }) => {
        // テストユーザーでログイン
        await loginAsTestUser(page, TEST_USERS.testuser)
        // トップページ（投稿一覧表示）へ遷移する
        await page.goto('/');
        // 最初の投稿のリンクを取得してクリックし、詳細ページへ遷移する
        const firstPost = page.getByTestId(POST_ITEM_TEST_ID).first().locator('a', { hasText: /./ });
        await firstPost.click();
        // テスト用コメント
        const testComment = TEST_COMMENT + ` ${Date.now()}`;
        // コメント入力
        await page.getByPlaceholder('コメントを入力').fill(testComment);    
        // 送信ボタンをクリック
        await page.getByRole('button', { name: BUTTON_SEND_COMMENT }).click();
        // 作成したコメントが画面に存在することを確認する
        const commentLocator = page.getByTestId(COMMENT_ITEM_TEST_ID).filter({ hasText: testComment });
        // コメントが表示されるまで最大10秒待つ
        await page.waitForSelector(`[data-testid='${COMMENT_ITEM_TEST_ID}']`, { timeout: WAIT_FOR_ELEMENT_TIMEOUT_MS });
        await expect(commentLocator).toBeVisible();
        // 編集ボタンをクリック
        const editButton = commentLocator.getByRole('button', { name: BUTTON_EDIT_COMMENT });
        await editButton.click();
        // テキストエリアに新しいコメントを入力する
        const newComment = testComment + ' - 更新';
        const textarea = commentLocator.locator('textarea');
        await textarea.fill(newComment, { timeout: ASSERTION_TIMEOUT_MS });
        // 保存する
        const saveButton = commentLocator.getByRole('button', { name: BUTTON_SAVE_COMMENT });
        await saveButton.click();
        // 更新内容が反映されていることを確認する
        await expect(commentLocator).toContainText(newComment, { timeout: ASSERTION_TIMEOUT_MS })
        // 削除ボタンをクリックしてコメントを削除
        const deleteButton = commentLocator.getByRole('button', { name: BUTTON_DELETE_COMMENT });
        // クリックによってconfirmがでることを想定してイベントハンドラをセットしておく
        page.once('dialog', async dialog =>{
            console.log(`Dialog message: ${dialog.message()}`); //デバッグ用ログ
            await dialog.accept();
        });
        await deleteButton.click();
        // UIからも消えていることを確認する
        await expect(commentLocator).toHaveCount(0, { timeout: ASSERTION_TIMEOUT_MS });
    });

});

test.describe('コメント機能：異常系テスト', () => {

    test('空コメント、文字数オーバー、他ユーザーのコメント編集削除不可のテスト', async({ page }) => {
        // テストユーザーでログイン
        await loginAsTestUser(page, TEST_USERS.testuser)
        // トップページ（投稿一覧表示）へ遷移する
        await page.goto('/');
        // 最初の投稿のリンクを取得してクリックし、詳細ページへ遷移する
        const firstPost = page.getByTestId(POST_ITEM_TEST_ID).first().locator('a', { hasText: /./ });
        await firstPost.click();
        // コメント入力欄を取得
        const commentInput = page.getByPlaceholder('コメントを入力');
        // 送信ボタンを取得
        const sendButton = page.getByRole('button', { name: BUTTON_SEND_COMMENT });
        // 空コメント設定
        await commentInput.fill('');
        // ボタンが disabled であることを確認
        await expect(sendButton).toBeDisabled();
        // 500文字の文字列を入力
        await commentInput.fill(TEST_COMMENT_LONG);
        // 入力した値が500文字になっていることを確認
        await expect(commentInput).toHaveValue(TEST_COMMENT_LONG);
        // 一旦コメント削除
        await commentInput.fill('');
        await expect(commentInput).toHaveValue('');
        // 501文字の文字列を作成
        await commentInput.fill(TEST_COMMENT_TOO_LONG);
        // 最大500文字しかはいらないことを確認する
        await expect(commentInput).toHaveValue(TEST_COMMENT_LONG);
        // 仮のコメントを投稿する
        await commentInput.fill(TEST_COMMENT_LONG);
        // 送信ボタンをクリック
        await page.getByRole('button', { name: BUTTON_SEND_COMMENT }).click();
        // ログアウト
        await logout(page);
        // 違うユーザーでログインする
        await loginAsTestUser(page, TEST_USERS.otheruser)
        // トップページ（投稿一覧表示）へ遷移する
        await page.goto('/');
        // 最初の投稿のリンクを取得してクリックし、詳細ページへ遷移する
        const firstPostSecondLogin = page.getByTestId(POST_ITEM_TEST_ID).first().locator('a', { hasText: /./ });
        await firstPostSecondLogin.click();
        // 仮投稿したコメントが画面に存在することを確認する
        const commentLocator = page.getByTestId(COMMENT_ITEM_TEST_ID).filter({ hasText: TEST_COMMENT_LONG });
        await expect(commentLocator).toBeVisible();
        // 編集と削除ボタンの取得を実施
        const editButton = commentLocator.getByRole('button', { name: BUTTON_EDIT_COMMENT });
        const deleteButton = commentLocator.getByRole('button', { name: BUTTON_DELETE_COMMENT });
        // 存在しないことを確認する
        await expect(editButton).toHaveCount(0);
        await expect(deleteButton).toHaveCount(0);
        // 後処理で追加したコメントを削除する
        await loginAsTestUser(page, TEST_USERS.testuser)
        await page.goto('/');
        const firstPostLastLogin = page.getByTestId(POST_ITEM_TEST_ID).first().locator('a', { hasText: /./ });
        await firstPostLastLogin.click();
        // 削除ボタンをクリックしてコメントを削除
        const button = commentLocator.getByRole('button', { name: BUTTON_DELETE_COMMENT });
        // クリックによってconfirmがでることを想定してイベントハンドラをセットしておく
        page.once('dialog', async dialog =>{
            console.log(`Dialog message: ${dialog.message()}`); //デバッグ用ログ
            await dialog.accept();
        });
        await button.click();
        // UIからも消えていることを確認する
        await expect(commentLocator).toHaveCount(0, { timeout: ASSERTION_TIMEOUT_MS });
    });
});
    