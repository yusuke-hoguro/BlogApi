DROP TABLE IF EXISTS likes;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS posts;

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

-- コメントのテーブル作成
CREATE TABLE comments(
    id SERIAL PRIMARY KEY,
    post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- いいねのテーブル作成
CREATE TABLE IF NOT EXISTS likes(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, post_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);

-- 初期データ投入
INSERT INTO posts (user_id, title, content) VALUES
  (1, 'テストタイトル1', 'テスト内容1'),
  (2, 'テストタイトル2', 'テスト内容2'),
  (3, 'コメントテスト用', 'コメント追加テスト');

INSERT INTO users (id, username, password) VALUES
  (1, 'testuser', 'pass'),
  (2, 'testuser2', 'pass2'),
  (3, 'testuser3', 'pass3');

INSERT INTO comments (id, post_id, user_id, content) VALUES
  (3, 1, 1, 'Default Comment'),
  (4, 1, 1, 'Default Comment'),
  (2, 2, 2, 'Delete Comment');

INSERT INTO likes (id, user_id, post_id) VALUES
  (1, 1, 1),
  (2, 2, 1),
  (3, 3, 1);
