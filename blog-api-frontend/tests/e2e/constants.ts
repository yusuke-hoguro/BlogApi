// アサーション用タイムアウト時間
export const ASSERTION_TIMEOUT_MS = 5000;
// 要素が表示されるのを待つためのタイムアウト
export const WAIT_FOR_ELEMENT_TIMEOUT_MS = 10000; 
// コメント関連
export const COMMENT_MAX_LENGTH = 500;
// テスト用コメント文字列
export const TEST_COMMENT = 'これはE2Eテスト用コメントです';
export const TEST_COMMENT_LONG = 'a'.repeat(COMMENT_MAX_LENGTH);
export const TEST_COMMENT_TOO_LONG = 'a'.repeat(COMMENT_MAX_LENGTH + 1);
