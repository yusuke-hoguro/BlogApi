-- 投稿統計のテーブル追加
CREATE TABLE IF NOT EXISTS post_stats(
    post_id INTEGER PRIMARY KEY REFERENCES posts(id) ON DELETE CASCADE,
    view_count INTEGER NOT NULL DEFAULT 0,
    like_count INTEGER NOT NULL DEFAULT 0,
    comment_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 投稿テーブルからIDを取得して存在しないIDは投稿統計テーブルに挿入する
INSERT INTO post_stats (post_id) SELECT id FROM posts ON CONFLICT (post_id) DO NOTHING;