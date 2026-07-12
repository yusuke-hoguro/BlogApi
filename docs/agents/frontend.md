# Frontend Guide

Frontend や React/Vite 実装を変更する前に読む資料です。API client、認証ルート、E2E を意識した UI 実装ルールをまとめます。

## 技術スタック

- React
- Vite
- Tailwind CSS
- Axios
- React Router
- Playwright E2E

## API 通信

- API 通信は原則 `blog-api-frontend/src/api/client.ts` の Axios client を経由する。
- Axios client は `VITE_API_BASE_URL` を `baseURL` として使う。
- 開発用 `.env.development` では `VITE_API_BASE_URL=http://localhost:8080`。
- request interceptor で `localStorage` の `token` を `Authorization: Bearer ...` として付与する。
- response interceptor では 401 の一部ケースで token を削除し、`/login` へ遷移する。

## ルーティングと認証

現行ルーティングでは、以下が `RequireAuth` 配下です。

- `/`
- `/post/create`
- `/post/:id`
- `/post/:id/edit`

以下は非認証画面です。

- `/signup`
- `/login`

ログイン必須ページを追加する場合は、既存の `RequireAuth` と `Layout` の構造を確認して合わせる。

## 実装ルール

- React コンポーネントは既存の JSX スタイルに合わせる。
- API 通信は `src/api/client.ts` の Axios client を経由する。
- 認証が必要な画面は既存の `RequireAuth` パターンを確認して合わせる。
- E2E で触る UI には、role/name で取れるラベルや必要に応じて `data-testid` を維持する。
- 画面で使う文言、button label、testid は E2E の constants と整合させる。
- `npm run lint` と `npm run build` が通る変更にする。

## 新規 Frontend 機能追加時のルール

1. API 呼び出しは `src/api/client.ts` を経由する。
2. ログイン必須ページは `RequireAuth` と既存 routing を確認する。
3. 画面で使う文言、button label、testid は E2E の constants と整合させる。
4. 主要ユーザーフローは Playwright に追加する。

## 注意

- フロントエンドから DB やバックエンド内部構造を前提にした処理を書かない。
- API のエラー形式は backend の `respondAppError` と `AuthMiddleware` で完全には統一されていないため、エラー表示を変える場合は backend/frontend/E2E の期待値を確認する。
- E2E のセレクタを壊す UI 変更では、テスト側の constants や fixtures も確認する。
