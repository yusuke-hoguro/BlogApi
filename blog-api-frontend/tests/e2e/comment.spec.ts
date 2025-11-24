import { test, expect } from '@playwright/test';
import { loginAsTestUser } from './utils';

test.describe('コメント機能', () => {
    
    test('UI画面でコメントの作成→表示→編集→削除のテストを実施する', async({ page }) => {

        // テストユーザーでログイン
        await loginAsTestUser(page)

        // トップページ（投稿一覧表示）へ遷移する
        await page.goto('http://localhost:3000/');

        // 最初の投稿のリンクを取得してクリックし、詳細ページへ遷移する
        const firstPost = page.getByTestId('post-item').first().locator('a', { hasText: /./ });
        await firstPost.click();

        // テスト用コメント
        const testComment = `E2Eテスト用コメント ${Date.now()}`;

        // コメント入力
        await page.getByPlaceholder('コメントを入力').fill(testComment);    
        // 送信ボタンをクリック
        await page.getByRole('button', { name: 'コメント送信' }).click();

        // 作成したコメントが画面に存在することを確認する
        const commentLocator = page.getByTestId('comment-item').filter({ hasText: testComment });
        await expect(commentLocator).toBeVisible();

        // 編集ボタンをクリック
        const editButton = commentLocator.getByRole('button', { name: '編集' });
        await editButton.click();

        // テキストエリアに新しいコメントを入力する
        const newComment = testComment + ' - 更新';
        const textarea = commentLocator.locator('textarea');
        await textarea.fill(newComment, { timeout: 5000 });

        // 保存する
        const saveButton = commentLocator.getByRole('button', { name: '保存' });
        await saveButton.click();

        // 更新内容が反映されていることを確認する
        await expect(commentLocator).toContainText(newComment, { timeout: 5000 })

        // 削除ボタンをクリックしてコメントを削除
        const deleteButton = commentLocator.getByRole('button', { name: '削除' });
        // クリックによってconfirmがでることを想定してイベントハンドラをセットしておく
        page.once('dialog', async dialog =>{
            console.log(`Dialog message: ${dialog.message()}`); //デバッグ用ログ
            await dialog.accept();
        });
        await deleteButton.click();

        // UIからも消えていることを確認する
        await expect(commentLocator).toHaveCount(0, { timeout: 5000 });
    });

});
