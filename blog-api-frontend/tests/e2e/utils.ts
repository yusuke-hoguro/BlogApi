import { Page } from "@playwright/test";
import { TestUser } from './users';

/*
* テストユーザーでログインしてTokenを取得する
* (addInitScriptでLocalStorageにJWTを入れる)
*/ 
export async function loginAsTestUser(page: Page, user: TestUser) {
    // テストユーザーでログインしてTokenを取得する
    const response = await page.request.post('http://localhost:8080/api/login', {
        data:{
            username: user.username, 
            password: user.password
        },
    });

    if (!response.ok()) {
        throw new Error(`ログイン失敗: ${response.status()} ${await response.text()}`);
    }

    const body = await response.json();
    const token = body.token;
    
    // addInitScriptは第2引数で渡した値をブラウザ側で実行される関数の第1引数として注入
    // 第1引数はブラウザ側で実行する関数
    await page.addInitScript(([jwt]) => {
        window.localStorage.setItem('token', jwt);
    }, [token])

    return token;
}

/**
 * ログアウト処理
 * localStorage の token を削除してトップページに遷移する
 */
export async function logout(page: Page) {
    await page.evaluate(() => {
        localStorage.removeItem('token');
    });
    await page.goto('/');
}

// RESTAPIで投稿を作成する
export async function createPost(page: Page, token: string, title: string, content: string){
    // Playwrite環境用のAPI呼び出し
    const res = await page.request.post('http://localhost:8080/api/posts',{
        data: { title, content },
        headers: {
            Authorization: `${token}`,
            'Content-Type': 'application/json',
        },
    });

    if(!res.ok()){
        throw new Error(`投稿作成に失敗: ${res.status()}`);
    }
    return await res.json();
}

// RESTAPIで投稿を削除する
export async function deletePost(page: Page, token: string, postId: number) {
    const res = await page.request.delete(`http://localhost:8080/api/posts/${postId}`, {
        headers: {
            Authorization: `${token}`,
        },
    });

    if(!res.ok()){
        throw new Error(`投稿削除に失敗: ${res.status()}`);
    }
}
