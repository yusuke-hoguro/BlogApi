import { request } from "@playwright/test";
import { TEST_USERS } from '@e2e/fixtures/users';

const API_BASE_URL = process.env.E2E_API_BASE_URL ?? 'http://localhost:8080';

// APIが起動するまで待機する関数
async function waitForApi(apiBaseURL: string) {
  const healthPaths = '/api/posts';
  const maxTries = 30;
  const sleepMs = 1000;

  const ctx = await request.newContext({ baseURL: apiBaseURL });

  for (let i = 1; i <= maxTries; i++) {
    try {
        const res = await ctx.get(healthPaths);
        if (res.ok()) {
            console.log('API is ready');
            await ctx.dispose();
            return;
        }
    } catch {
        console.log(`API not ready yet, attempt ${i}/${maxTries}`);
    }
    // 非同期で指定した時間待機
    await new Promise(r => setTimeout(r, sleepMs));
  }

  await ctx.dispose();
  throw new Error(`API is not ready at ${apiBaseURL} after ${maxTries} seconds`);
}

// グローバルセットアップ関数
export default async function globalSetup(){
    // APIが起動するまで待機
    await waitForApi(API_BASE_URL);
    // PlaywrightのAPIRequestContextを作成
    const ctx = await request.newContext({
        baseURL: API_BASE_URL,
        extraHTTPHeaders: {
            'Content-Type': 'application/json',
        },
    })

    // テスト用ユーザーを登録する
    for(const user of Object.values(TEST_USERS)){
        // API経由でユーザー登録を実行
        const res = await ctx.post('/api/signup', {data: user});
        // 実行結果を判定する
        if (![200, 201, 409].includes(res.status())){
            const body = await res.text().catch(() => '');
            throw new Error(`Failed to create test user ${user.username}. Status: ${res.status()}. Body: ${body}`);
        }
    }
    // 作成したリクエストコンテキストを破棄する
    await ctx.dispose();
    console.log('Global setup completed: Test users are created.');
}

