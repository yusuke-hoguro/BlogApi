import React, { useEffect, useState } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import client from '../api/client';

export default function PostDetail(){
    // useParams:URLに含まれるパラメータをオブジェクトとして返す
    const { id } = useParams();                                     // 投稿ID
    // useNavigate:「navigate」関数を取得して任意のパスに移動できる
    const navigate = useNavigate();                                 // navigate関数
    // useState:状態管理フック 変数の初期値を設定し、その変数を更新するための関数を返す
    const [post, setPost] = useState(null);                         // 投稿管理
    const [comments, setComments] = useState([]);                   // コメント一覧
    const [loading, setLoading] = useState(true);                   // 表示読み込み管理
    const [newComment, setNewComment] = useState('');               // 新規コメント
    const [submitting, setSubmitting] = useState(false);            // コメント送信管理
    const [successMsg, setSuccessMsg] = useState('');               // 送信完了時のメッセージ
    const [editingCommentId, setEditingCommentId] = useState(null); // 更新中のコメントID
    const [editingContent, setEditingContent] = useState('');       // 更新用コメント
    const [errorMsg, setErrorMsg] = useState('');                   // エラーメッセージ

    // 初回レンダリング時のみ実行
    useEffect(() => {
        fetchPostAndComments();
    }, [id]);

    // 投稿とコメントの取得を実施する
    async function fetchPostAndComments(){
        try{
            const [postRes, commentRes] = await Promise.all([
                client.get(`/api/posts/${id}`),
                client.get(`/api/posts/${id}/comments`)
            ]);
            setPost(postRes.data);
            setComments(commentRes.data || []);
        }catch(error){
            console.error('取得エラー:', error);
            setErrorMsg('投稿とコメントの取得に失敗しました。');
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
        setErrorMsg('');
        try{
            const token = localStorage.getItem('token');
            const response = await client.post(
                `/api/posts/${id}/comments`,
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
            // 送信成功メッセージ
            setSuccessMsg('コメントを送信しました！');
            // 3秒後にメッセージを消す
            setTimeout(() => setSuccessMsg(''), 3000);
        }catch(error){
            console.error('コメント投稿エラー:', error);
            setErrorMsg('コメント投稿でエラーが発生しました。' + error.message);
        }finally{
            setSubmitting(false);
        }
    }

    // コメント削除用関数
    async function handleDeleteComment(commentId) {
        if(!window.confirm("本当にコメントを削除しますか？")) return;

        try{
            const token = localStorage.getItem("token");
            await client.delete(`/api/comments/${commentId}`, {
                headers: { Authorization: token },
            });
            // 削除後にコメント一覧を再取得
            await fetchPostAndComments();
        } catch (error) {
            console.error("コメント削除エラー:", error);
            setErrorMsg('コメントの削除に失敗しました。');
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

    // 投稿削除を実行する
    async function handleDeletePost() {
        if (!window.confirm("本当に投稿を削除しますか？")) return;

        try {
            const token = localStorage.getItem("token");
            await client.delete(`/api/posts/${post.id}`, { headers: { Authorization: token } });
            navigate("/"); // 削除後は投稿一覧に戻る
        } catch (error) {
            console.error("投稿削除エラー:", error);
            setErrorMsg("投稿の削除に失敗しました。");
        }
    }

    // コメント更新用の関数
    async function handleUpdateComment(commentId) {
        if(!editingContent.trim()) return;

        try{
            const token = localStorage.getItem("token");
            await client.put(
                `/api/comments/${commentId}`,
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
            setErrorMsg('コメントの更新に失敗しました。');
        }
    }

    if(loading) return <p className="p-4">読み込み中...</p>;
    if(!post) return <p className="p-4 text-red-600">投稿が見つかりません</p>;

    return(
        <div className="p-4 max-w-3xl mx-auto w-full overflow-x-hidden box-border">
            {/* 戻るリンク &larr;は左向き矢印 */}
            <Link to="/" className="text-blue-600 hover:underline">
                &larr; 投稿一覧に戻る
            </Link>

            {/* 投稿タイトルと内容 */}
            <h1 data-testid="post-title" className='text-2xl font-bold mt-4'>{post.title}</h1>
            {/* mt:margin-top whitespace-pre-wrap:改行や連続スペースをそのまま表示しつつ、必要に応じで自動で折り返す */}
            <p data-testid="post-content" className='mt-2 text-gray-800 whitespace-pre-wrap break-all max-w-full'>{post.content}</p>


            {/* 自分の投稿なら削除ボタンを表示 */}
            {post.user_id === getCurrentUserId() && (
                <div className="mt-4">
                    <Link
                        to={`/post/${post.id}/edit`}
                        className="px-3 py-1 bg-yellow-500 text-white rounded hover:bg-yellow-600 transition-colors"
                    >
                        投稿編集
                    </Link>
                    <button
                        onClick={handleDeletePost}
                        className="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600 transition-colors"
                    >
                        投稿削除
                    </button>
                </div>
            )}

            {/* コメント一覧 */}
            <div className="mt-6">
                <h2 className='text-xl font-semibold mb-2'>コメント一覧</h2>
                {/* JSX内ではif文は使えないので三項演算子で記載 */}
                {!comments || comments.length === 0 ? (
                    <p className='text-gray-500'>コメントはまだありません。</p>
                ):(
                    <ul className='space-y-4'>
                        {comments.map(comment => (
                            // コメントカード全体
                            <li key={comment.id} data-testid="comment-item" className='border  rounded-lg shadow-sm bg-white p-4 max-w-full mx-auto min-w-0 overflow-x-hidden break-words'>
                                {/* 編集モードか表示モードかを切り替え */}
                                {editingCommentId === comment.id ? (
                                    <>
                                        {/* 編集用テキストエリア */}
                                        <textarea
                                            className="w-full border rounded p-2 focus:outline-none focus:ring-2 focus:ring-green-400 resize-none max-w-full"
                                            value={editingContent}
                                            onChange={(e) => setEditingContent(e.target.value)}
                                            // 高さを3行で揃える
                                            rows={3}
                                        />
                                        {/* 保存・キャンセルボタン */}
                                        <div className="mt-2 flex flex-wrap gap-2">
                                            <button
                                                onClick={() => handleUpdateComment(comment.id)}
                                                className="px-3 py-1 bg-green-500 text-white rounded hover:bg-green-600 transition-colors"
                                            >
                                                保存
                                            </button>
                                            <button
                                                onClick={() => setEditingCommentId(null)}
                                                className="px-3 py-1 bg-gray-400 text-white rounded hover:bg-gray-500 transition-colors"
                                            >
                                                キャンセル
                                            </button>
                                        </div>
                                    </>                                
                                ):(
                                    <>
                                        {/* コメント本文 */}
                                        <p className='text-gray-800 whitespace-pre-wrap break-all max-w-full overflow-x-hidden leading-relaxed'>
                                            {comment.content}
                                        </p>
                                        <div className="mt-2 flex flex-col sm:flex-row sm:justify-between sm:items-center gap-2">
                                            <div className="text-sm text-gray-500">
                                                <span>ユーザーID: {comment.user_id}</span> | 
                                                <span className="ml-1">{new Date(comment.created_at).toLocaleString()}</span>
                                            </div>
                                            {/* 編集・削除ボタン（自分のコメントのみ） */}
                                            {comment.user_id === getCurrentUserId() && (
                                                <div className="mt-2 sm:mt-0 flex flex-wrap gap-2">
                                                    <button
                                                        onClick={() => {
                                                            setEditingCommentId(comment.id);
                                                            setEditingContent(comment.content);
                                                        }}
                                                        className="px-3 py-1 bg-yellow-500 text-white rounded hover:bg-yellow-600 transition-colors"
                                                    >
                                                        編集
                                                    </button>
                                                    <button
                                                        onClick={() => handleDeleteComment(comment.id)}
                                                        className="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600 transition-colors"
                                                    >
                                                        削除
                                                    </button>
                                                </div>  
                                            )}
                                        </div>
                                    </>
                                )}
                            </li>
                        ))}
                    </ul>
                )}
            </div>

            <form onSubmit={handleSubmit} className="mt-6 border-t pt-4 space-y-2">
                <label htmlFor="new-comment" className="block font-medium">コメントを入力</label>
                <textarea 
                    id="new-comment" 
                    className='block w-full border rounded p-2 box-border max-w-full resize-none' 
                    rows="3" 
                    placeholder='コメントを入力…' 
                    value={newComment} 
                    disabled={submitting}   // 送信中は入力不可
                    onChange={(e) => {
                        setNewComment(e.target.value);
                        // 自動リサイズ
                        e.target.style.height = "auto";
                        e.target.style.height = e.target.scrollHeight + "px";
                    }}
                    maxLength={500} // 最大文字数を500文字に設定
                />
                {/* 文字数カウント表示 */}
                <p className="text-sm text-gray-500 text-right">{newComment.length} / 500文字</p>

                {/* エラー・成功メッセージ */}
                {errorMsg && <p className="!text-red-500 text-sm">{errorMsg}</p>}
                {successMsg && <p className="!text-green-500 text-sm">{successMsg}</p>}

                <div className="flex justify-end">
                    <button 
                        type="submit" 
                        disabled={submitting || !newComment.trim()} 
                        className='mt-2 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:bg-gray-400 box-border max-w-full'
                    >
                        {submitting ?  '送信中...' : 'コメント送信'}
                    </button>
                </div>
            </form>
        </div>
    );
}

