# Playwrite(プレイライト) 環境構築手順

## 概要

## Seleniumとの比較

### メリット

- 専用のブラウザバイナリを使用するので環境構築が用意
- Playwrightライブラリと専用ブラウザのバージョンが常に同期されているため、ブラウザの自動更新で動かないという問題が起こりにくい
- テストの再現性が高く、どの環境で実行しても同じバージョンのブラウザでテストが実施される

### 理由
- Seleniumは基本的にユーザーが普段使ってる既存のブラウザとそれに対応するドライバーを組み合わせて動作する
    - 動作の仕組み:WebDriverがブラウザのAPIと通信し、ブラウザを操作する
    - テストを実施するにはPCにインストール済みのブラウザとそのブラウザのバージョンに一致するWebDriverの実行ファイルが必要
    - ブラウザが自動更新されると、WebDriverのバージョンも合わせる必要があるので環境構築やメンテナンスの手間がかかる
- Playwrightはテスト実行のために最適化された専用のブラウザバイナリを使用する
    - `npx playwright install` コマンドにより、Playwrightのバージョンと完全に互換性のある特定のバージョンのブラウザがDLされる
    - 通常のブラウザとは別の場所にインストールされ、テスト実行時のみに使用される
- ロケーターはWebページ上のUI要素を一意に特定し、操作するための中心的な仕組み
    - Seleniumなどでは要素を特定するためにCSSセレクタやXPathといったDOM構造に依存した手法が主流
    - これだとHTMLの構造変更に弱くテストが壊れやすい
    - ロケーターはユーザーの視点に近い方法で要素を探す
    - ロケーターはWebページ上の特定の要素を見つけ出すための仕組み
- Playwrite環境の強み、自動待機メカニズム
    - WebDriverWaitのような明示的な待機処理はほとんど不要
    - Playwriteのアクション（クリックなど）は実行前に以下を自動チェックする
        - 要素がDOMにアタッチされているか
        - 要素が表示されているか (Visible)
        - 要素が安定しているか (Stable) (アニメーションが終了しているかなど)
        - 要素がイベントを受け付けられるか (Enabled)(他の要素に隠れていないか、無効化されていないかなど) 

## 手順

1. フロントエンドに`Playwright`をインストール

- フロントエンドのディレクトリに移動して`Playwright`をインストールする

    ```bash
    cd blog-api-frontend
    npm install --save-dev @playwright/test
    npx playwright install #「Playwrightを使ったテストを実行できるように、必要な専用ブラウザを準備してください」という指示を出すためのコマンド
    ```

2. Playwright の設定ファイルを作る

- Playwright は playwright.config.ts があれば動くので、手動で作る。
- フロントエンドのプロジェクト直下に作成する

    ```ts
    import { defineConfig } from '@playwright/test';

    export default defineConfig({
    testDir: './tests',
    use: {
        headless: true,
        baseURL: 'http://localhost:3000',
    },
    webServer: {
        command: 'npm run dev',
        port: 3000,
        reuseExistingServer: true,
    },
    });
    ```
- 最新のPlaywrightはwebServerの自動起動が標準
    - `npx playwrite test`でフロントが自動起動される
- baseURLを設定しておくと`page.goto('/')`が使える

3. E2E テスト用ディレクトリを作成

    ```markdown
    tests/
        e2e/
            smoke.spec.ts
    ```

4. とりあえずは動く簡単なE2Eテストを書く

- 仮として`smoke.spec.ts`を以下のソースコードで作成する

    ```ts
    import { test, expect } from '@playwright/test';

    test('トップページが開ける', async ({ page }) => {
        await page.goto('/');
        await expect(page).toHaveTitle(/.*/);
    });
    ```

5. テストを実行する

- 通常実行

```bash
npx playwright test
```
- GUI で見たい場合

```bash
npx playwright test --ui
```
- 録画しながら操作を保存する Playwright Codegen

```bash
npx playwright codegen http://localhost:3000
```
