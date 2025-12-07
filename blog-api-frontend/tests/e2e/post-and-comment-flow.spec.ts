import { test, expect } from '@playwright/test';
import { loginAsTestUser } from './utils';
import { ASSERTION_TIMEOUT_MS, WAIT_FOR_ELEMENT_TIMEOUT_MS } from './constants/config';
import { TEST_USERS } from './users';
import { POST_ITEM_TEST_ID, POST_TITLE_TEST_ID, POST_CONTENT_TEST_ID, COMMENT_ITEM_TEST_ID } from './constants/selectors';
import { PAGE_TITLE_POST_LIST, PAGE_TITLE_POST_CREATE } from './constants/pageTitles';
import { CREAT_POST_TITLE, CREAT_POST_CONTENT, UPDATE_POST_TITLE, UPDATE_POST_CONTEXT } from './constants/posts';
import { BUTTON_UPDATE_POST, BUTTON_CREATE_POST, BUTTON_DELETE_POST, BUTTON_SEND_COMMENT, BUTTON_EDIT_COMMENT, BUTTON_SAVE_COMMENT, BUTTON_DELETE_COMMENT } from './constants/buttons';
import { LABEL_EDIT_POST, LABEL_POST_CREATE_TITLE, LABEL_POST_CREATE_CONTEXT } from './constants/label';
import { TEST_COMMENT } from './constants/comments';

test.describe('全体機能テスト:正常系', () => {

    test('新規投稿作成→コメント作成、編集、削除→投稿編集、削除のテストを実施', async ({ page }) => {
        // テストユーザーでログインする
        const token = await loginAsTestUser(page, TEST_USERS.testuser);
        // 投稿作成ページへ遷移
        await page.goto('/post/create');
        // 投稿作成ページが開けたかを確認する
        await expect(page.getByRole('heading', { name: PAGE_TITLE_POST_CREATE })).toBeVisible();
        // 新規投稿の作成を実施する
        const title = CREAT_POST_TITLE + `${Date.now()}`;
        const content = CREAT_POST_CONTENT;
        await page.getByLabel(LABEL_POST_CREATE_TITLE).fill(title);
        await page.getByLabel(LABEL_POST_CREATE_CONTEXT).fill(content);
        await page.getByRole('button', { name: BUTTON_CREATE_POST }).click();
        // 投稿一覧ページに遷移したかを確認する
        await expect(page.getByRole('heading', { name: PAGE_TITLE_POST_LIST })).toBeVisible();
        // 投稿一覧の中に新規作成した投稿があることを確認する
        await expect(page.getByTestId(POST_ITEM_TEST_ID).filter({ hasText: title})).toHaveCount(1)
        // 新規追加した投稿の詳細画面を開く
        await page.getByRole('link', { name: title }).click();
        // 新規追加した投稿の詳細画面に遷移できたかをチェック
        await expect(page.getByRole('heading', { name: title })).toBeVisible();
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
        // 編集ボタンを押す
        await page.getByRole('link', { name: LABEL_EDIT_POST }).click();
        // タイトルを編集する
        const updateTitle = UPDATE_POST_TITLE + ` ${Date.now()}`;
        await page.getByLabel(LABEL_POST_CREATE_TITLE).fill(updateTitle);
        // タイトルが変わったか確認
        await expect(page.getByLabel(LABEL_POST_CREATE_TITLE)).toHaveValue(updateTitle, { timeout: ASSERTION_TIMEOUT_MS })
        // 投稿内容を編集する
        const updateContext = UPDATE_POST_CONTEXT;
        await page.getByLabel(LABEL_POST_CREATE_CONTEXT).fill(updateContext);
        // 投稿内容が変わったか確認
        await expect(page.getByLabel(LABEL_POST_CREATE_CONTEXT)).toHaveValue(updateContext, { timeout: ASSERTION_TIMEOUT_MS })
        // 更新ボタンを押下する
        await page.getByRole('button', { name: BUTTON_UPDATE_POST }).click();
        // 投稿のタイトルが表示されていることを確認する
        await expect(page.getByTestId(POST_TITLE_TEST_ID)).toContainText(updateTitle, { timeout: ASSERTION_TIMEOUT_MS })
        // 投稿の内容が表示されていることを確認する
        await expect(page.getByTestId(POST_CONTENT_TEST_ID)).toContainText(updateContext, { timeout: ASSERTION_TIMEOUT_MS })
        // クリックによってconfirmがでることを想定してイベントハンドラをセット
        page.once('dialog', dialog => {
            console.log(`Dialog message: ${dialog.message()}`);
            dialog.accept();
        });
        // 投稿削除を実施
        await page.getByRole('button', { name: BUTTON_DELETE_POST }).click();
        // 投稿一覧に戻るまで待機
        await page.waitForURL('/');
        // 投稿一覧ページに遷移したかを確認する
        await expect(page.getByRole('heading', { name: PAGE_TITLE_POST_LIST })).toBeVisible();
        // 投稿をすべて取得する
        const posts = page.getByTestId(POST_ITEM_TEST_ID);
        // タイトル一致しないことを確認する
        await expect(posts.filter({ hasText: updateTitle })).toHaveCount(0);
    });

});
