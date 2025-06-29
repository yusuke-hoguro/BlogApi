DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS users;

-- ユーザー用のテーブル作成
CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

-- 投稿用のテーブル作成
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 初期データ投入
INSERT INTO posts (user_id, title, content) VALUES
  (1, 'テストタイトル1', 'テスト内容1'),
  (2, 'テストタイトル2', 'テスト内容2');