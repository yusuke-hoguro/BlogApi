# API仕様書

BlogAPIで仕様するRESTAPIの仕様と利用方法を記載しています。

※一部作成途中

---

## API一覧

| メソッド   | エンドポイント                | 内容           |
| ------ | ---------------------- | ------------ |
| GET    | `/posts`               | 投稿一覧取得       |
| POST   | `/posts`               | 投稿作成（認証必要）   |
| GET    | `/posts/{id}`          | 投稿詳細取得       |
| PUT    | `/posts/{id}`          | 投稿更新（本人のみ）   |
| DELETE | `/posts/{id}`          | 投稿削除（本人のみ）   |
| GET    | `/posts/{id}/comments` | コメント一覧取得     |
| POST   | `/comments`            | コメント作成（認証必要） |
| PUT    | `/comments/{id}`       | コメント更新（本人のみ） |
| DELETE | `/comments/{id}`       | コメント削除（本人のみ） |
