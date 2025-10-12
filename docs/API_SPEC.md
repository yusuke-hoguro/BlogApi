# API仕様書

BlogAPIで仕様するRESTAPIの仕様と利用方法を記載しています。

※一部作成途中

---

## API一覧

| メソッド   | エンドポイント                | 内容           |
| ------ | ---------------------- | ------------ |
| GET    | `/api/posts`               | 投稿一覧取得       |
| POST   | `/api/posts`               | 投稿作成（認証必要）   |
| GET    | `/api/posts/{id}`          | 投稿詳細取得       |
| PUT    | `/api/posts/{id}`          | 投稿更新（本人のみ）   |
| DELETE | `/api/posts/{id}`          | 投稿削除（本人のみ）   |
| GET    | `/api/posts/{id}/comments` | コメント一覧取得     |
| POST   | `/api/comments`            | コメント作成（認証必要） |
| PUT    | `/api/comments/{id}`       | コメント更新（本人のみ） |
| DELETE | `/api/comments/{id}`       | コメント削除（本人のみ） |
