import React, { useState, useEffect } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import client from "../api/client";


export default function PostEdit() {
    // useParams:URLに含まれるパラメータをオブジェクトとして返す
    const { id } = useParams();
    // useNavigate:「navigate」関数を取得して任意のパスに移動できる
    const navigate = useNavigate();
    // useState:状態管理フック 変数の初期値を設定し、その変数を更新するための関数を返す
    const [title, setTitle] = useState("");
    const [content, setContent] = useState("");
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [errorMsg, setErrorMsg] = useState("");

    // 初期表示（元データを取得）:初回レンダリング時のみ実行
    useEffect(() => {
        async function fetchPost(){
            try {
                // APIで投稿を取得
                const res = await client.get(`/api/posts/${id}`);
                setTitle(res.data.title);
                setContent(res.data.content);                
            } catch(error){
                console.error("投稿取得エラー：", error);
                setErrorMsg("投稿の取得に失敗しました。")
            } finally{
                setLoading(false);
            }
        }
        fetchPost();
    }, [id]);

    // 投稿更新処理
    async function handleUpdate(e){
        // ブラウザのデフォルト動作を停止させる
        e.preventDefault();
        setSaving(true);
        setErrorMsg("");

        try{
            const token = localStorage.getItem("token");
            // 投稿更新用のAPIを送信
            await client.put(
                `/api/posts/${id}`,
                { title, content }
            );
            // 更新が成功したので詳細ページに遷移
            navigate(`/post/${id}`);
        } catch(error){
            console.error("投稿更新エラー:", error);
            setErrorMsg("投稿の更新に失敗しました。");
        } finally{
            setSaving(false);
        }
    }

    if(loading) return <p className="p-4">読み込み中...</p>;

    return (
        <div className="p-4 max-w-3xl mx-auto">
            {/* 戻るリンク &larr;は左向き矢印 */}
            <Link to={`/post/${id}`} className="text-blue-600 hover:underline">
                &larr; 投稿詳細に戻る
            </Link>

            <h1 className="text-2xl font-bold mt-4 mb-4">投稿編集</h1>

            {errorMsg && <p className="text-red-600 mb-2">{errorMsg}</p>}

            {/* 更新ボタンが押されたらhandleUpdateを実行する */}
            <form onSubmit={handleUpdate} className="space-y-4">
                <div>
                    <label htmlFor="post-title" className="block font-medium mb-1">タイトル</label>
                    <input
                        id="post-title"
                        type="text"
                        className="w-full border rounded p-2"
                        value={title}
                        disabled={saving}
                        onChange={(e) => setTitle(e.target.value)}
                        maxLength={100}
                        required
                    />
                </div>

                <div>
                    <label htmlFor="post-content" className="block font-medium mb-1">内容</label>
                    <textarea
                        id="post-content"
                        className="w-full border rounded p-2 resize-none"
                        rows="8"
                        value={content}
                        disabled={saving}
                        onChange={(e) => setContent(e.target.value)}
                        maxLength={1000}
                        required
                    />
                </div>

                <button
                    type="submit"
                    disabled={saving}
                    className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 disabled:bg-gray-400"
                >
                    {saving ? "保存中..." : "更新"}
                </button>
            </form>
            

        </div>

       

        
    )
}
