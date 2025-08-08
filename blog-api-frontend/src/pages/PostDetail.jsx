import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import client from '../api/client';

export default function PostDetail(){
    // useParams:URLに含まれるパラメータをオブジェクトとして返す
    const { id } = useParams();
    // useState:状態管理フック 変数の初期値を設定し、その変数を更新するための関数を返す
    const [post, setPost] = useState(null);
    const [comments, setComments] = useState([]);
    const [loading, setLoading] = useState(true);

    /*
    * 初回レンダリング時のみ実行
    * useEffectはasyncを直接渡せないので内部でasync関数を定義して呼び出す
    */
    useEffect(() => {
        async function fetchPostAndComments(){
            try{
                const [postRes, commentRes] = await Promise.all([
                    client.get(`/posts/${id}`),
                    client.get(`/posts/${id}/comments`)
                ]);
                setPost(postRes.data);
                setComments(commentRes.data || []);
            }catch(error){
                console.error('取得エラー:', error);
            }finally{
                setLoading(false);
            }
        }
        fetchPostAndComments();
    }, [id]);

    if(loading) return <p className="p-4">読み込み中...</p>;
    if(!post) return <p className="p-4 text-red-600">投稿が見つかりません</p>;

    return(
        <div className='p-4'>
            {/* &larr;は左向き矢印 */}
            <Link to="/" className="text-blue-600 hover:underline">&larr; 投稿一覧に戻る</Link>
            <h1 className='text-2xl font-bold mt-4'>{post.titile}</h1>
            {/* mt:margin-top whitespace-pre-wrap:改行や連続スペースをそのまま表示しつつ、必要に応じで自動で折り返す */}
            <p className='mt-2 text-gray-800 whitespace-pre-wrap'>{post.content}</p>

            <div className="mt-6">
                <h2 className='text-xl font-semibold mb-2'>コメント一覧</h2>
                {/* JSX内ではif文は使えないので三項演算子で記載 */}
                {!comments || comments.length === 0 ? (
                    <p className='text-gray-500'>コメントはまだありません。</p>
                ):(
                    <ul className='space-y-3'>
                        {comments.map(comment => (
                            <li key={comment.id} className='border p-3 rounded'>
                                <p className='text-gray-700'>{comment.content}</p>
                                <p className='text-sm text-gray-400'>ユーザーID:{comment.user_id}</p>
                            </li>
                        ))}
                    </ul>
                )}
            </div>
        </div>
    );
}

