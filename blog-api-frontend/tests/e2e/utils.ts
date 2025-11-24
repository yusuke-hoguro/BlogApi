import { Page } from "@playwright/test";

/*
* テストユーザーでログインしてTokenを取得する
* (addInitScriptでLocalStorageにJWTを入れる)
*/ 
export async function loginAsTestUser(page: Page) {
    // テストユーザーでログインしてTokenを取得する
    const responce = await page.request.post('http://localhost:8080/api/login', {
        data:{
            username: "testuser2", 
            password: "validpassword"
        }
    });

    const body = await responce.json();
    const token = body.token;
    
    // addInitScriptは第2引数で渡した値をブラウザ側で実行される関数の第1引数として注入
    // 第1引数はブラウザ側で実行する関数
    await page.addInitScript(([jwt]) => {
        window.localStorage.setItem('token', jwt);
    }, [token])

    return token;
}
