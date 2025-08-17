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
    const [newComment, setNewComment] = useState('');
    const [submitting, setSubmitting] = useState(false);
    // コメント編集用
    const [editingCommentId, setEditingCommentId] = useState(null);
    const [editingContent, setEditingContent] = useState('');

    // 初回レンダリング時のみ実行
    useEffect(() => {
        fetchPostAndComments();
    }, [id]);

    // 投稿とコメントの取得を実施する
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

    // コメント送信処理
    async function handleSubmit(e) {
        // ブラウザのデフォルト動作を停止させる
        e.preventDefault();
        if(!newComment.trim()) return;

        setSubmitting(true);
        try{
            const token = localStorage.getItem('token');
            const response = await client.post(
                `/posts/${id}/comments`,
                { content: newComment },
                {
                    headers:{
                        Authorization: token
                    }
                }
            );
            console.log('投稿APIレスポンス:', response);
            setNewComment('');
            await fetchPostAndComments();
        }catch(error){
            console.error('コメント投稿エラー:', error);
            alert('コメント投稿でエラーが発生しました: ' + error.message);
        }finally{
            setSubmitting(false);
        }
    }

    // コメント削除用関数
    async function handleDeleteComment(commentId) {
        if(!window.confirm("本当にコメントを削除しますか？")) return;

        try{
            const token = localStorage.getItem("token");
            await client.delete(`/comments/${commentId}`, {
                headers: { Authorization: token },
            });
            // 削除後にコメント一覧を再取得
            await fetchPostAndComments();
        } catch (error) {
            console.error("コメント削除エラー:", error);
            alert("コメントの削除に失敗しました。");
        }
    }

    // JWTからログイン中のユーザーIDを取得
    function getCurrentUserId(){
        const token = localStorage.getItem("token");
        if (!token) return null;

        try{
            const payload = JSON.parse(atob(token.split('.')[1]))
            return payload.user_id;
        }catch{
            return null;
        }
    }

    // コメント更新用の関数
    async function handleUpdateComment(commentId) {
        if(!editingContent.trim()) return;

        try{
            const token = localStorage.getItem("token");
            await client.put(
                `/comments/${commentId}`,
                { content: editingContent },
                { headers: { Authorization: token } }
            );
            // 編集終了のためリセットする
            setEditingCommentId(null);
            setEditingContent('');
            // コメント更新後に投稿とコメントを再取得
            await fetchPostAndComments();
        } catch(error){
            console.error("コメント更新エラー:", error);
            alert("コメントの更新に失敗しました。");            
        }
    }

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
                                {editingCommentId === comment.id ? (
                                    <>
                                        <textarea
                                            className="w-full border rounded p-2"
                                            value={editingContent}
                                            onChange={(e) => setEditingContent(e.target.value)}
                                        />
                                        <div className="mt-2 flex gap-2">
                                            <button
                                                onClick={() => handleUpdateComment(comment.id)}
                                                className="px-2 py-1 bg-green-500 text-white rounded hover:bg-green-600"
                                            >
                                                保存
                                            </button>
                                            <button
                                                onClick={() => setEditingCommentId(null)}
                                                className="px-2 py-1 bg-gray-400 text-white rounded hover:bg-gray-500"
                                            >
                                                キャンセル
                                            </button>
                                        </div>
                                    </>                                
                                ):(
                                    <>
                                        <p className='text-gray-700'>{comment.content}</p>
                                        <p className='text-sm text-gray-400'>ユーザーID:{comment.user_id}</p>

                                        {/* 自分のコメントのみ削除ボタン表示 */}
                                        {comment.user_id === getCurrentUserId() && (
                                            <div className="mt-2 flex gap-2">
                                                <button
                                                    onClick={() => {
                                                        setEditingCommentId(comment.id);
                                                        setEditingContent(comment.content);
                                                    }}
                                                    className="px-2 py-1 bg-yellow-500 text-white rounded hover:bg-yellow-600"
                                                >
                                                    編集
                                                </button>
                                                <button
                                                    onClick={() => handleDeleteComment(comment.id)}
                                                    className="mt-2 px-2 py-1 bg-red-500 text-white rounded hover:bg-red-600"
                                                >
                                                    削除
                                                </button>
                                            </div>  
                                        )}
                                    </>
                                )}
                            </li>
                        ))}
                    </ul>
                )}
            </div>

            <form onSubmit={handleSubmit} className="mt-6 border-t pt-4">
                <textarea className='w-full border rounded p-2' rows="3" placeholder='コメントを入力…' value={newComment} onChange={(e) => setNewComment(e.target.value)}/>
                <button type="submit" disabled={submitting} className='mt-2 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:bg-gray-400'>
                    {submitting ?  '送信中...' : 'コメント送信'}
                </button>
            </form>
        </div>
    );
}

