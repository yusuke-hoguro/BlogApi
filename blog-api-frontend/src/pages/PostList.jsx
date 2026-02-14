import { Link } from "react-router-dom";
import React, {useEffect, useState} from "react";
import client from "../api/client";

export default function PostList(){
    const [posts, setPosts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [errorMsg, setErrorMsg] = useState('');

    useEffect(() => {
        client.get('/api/posts')
            .then(response => {
                const data = response.data;
                // 投稿が配列であることを確認。配列でない場合は空配列を設定。（違うデータを設定すると例外で画面が真っ白になる）
                setPosts(Array.isArray(data) ? data : []);
            })
            .catch(error => {
                console.error('投稿取得エラー:', error);
                setErrorMsg('投稿の取得に失敗しました。');
            })
            .finally(() => setLoading(false))
    }, []);
    
    if(loading) return <p className="p-4">読み込み中...</p>

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-5xl mx-auto px-4 py-10">
        {/* タイトル行 + 作成ボタン */}
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-3xl font-bold text-gray-800">投稿一覧</h1>
          <Link
            to="/post/create"
            className="inline-flex items-center px-4 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 transition"
          >
            新規投稿作成
          </Link>
        </div>

        {/* エラー */}
        {errorMsg && (
          <p
            data-testid="post-fetch-error"
            className="text-red-600 bg-red-50 border border-red-200 rounded-lg p-3 mb-4"
          >
            {errorMsg}
          </p>
        )}

        {/* 空状態 */}
        {!errorMsg && posts.length === 0 && (
          <p
            data-testid="post-empty"
            className="text-gray-600 bg-white border border-gray-200 rounded-lg p-4"
          >
            投稿がありません
          </p>
        )}

        {/* 一覧 */}
        <ul className="space-y-4">
          {posts.map((post) => (
            <li
              key={post.id}
              data-testid="post-item"  // テスト用に追加
              className="bg-white border border-gray-200 rounded-xl p-5 shadow-sm hover:shadow-md transition break-words"
            >
              <h2 className="text-lg font-semibold">
                <Link
                  to={`/post/${post.id}`}
                  className="text-blue-700 hover:underline break-words"
                >
                  {post.title}
                </Link>
              </h2>

              <p className="text-gray-700 mt-2 break-words">
                {post.content?.slice(0, 100)}
                {post.content?.length > 100 ? "..." : ""}
              </p>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
