import { Link } from "react-router-dom";
import React, {useEffect, useState} from "react";
import client from "../api/client";

export default function PostList(){
    const [posts, setPosts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [errorMsg, setErrorMsg] = useState('');

    useEffect(() => {
        client.get('/posts')
            .then(response => {
                setPosts(response.data)
            })
            .catch(error => {
                console.error('投稿取得エラー:', error);
                setErrorMsg('投稿の取得に失敗しました。');
            })
            .finally(() => setLoading(false))
    }, []);
    
    if(loading) return <p className="p-4">読み込み中...</p>

    return(
        <div className="p-4 max-w-3xl mx-auto">
            <h1 className="text-2xl font-bold mb-4" >投稿一覧</h1>
            <Link
                to="/post/create"
                className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 transition-colors"
            >
                新規投稿作成
            </Link>
            {errorMsg && <p className="text-red-500 mb-4">{errorMsg}</p>}
            <ul className="space-y-4">
                {posts.map(post => (
                    <li 
                        key={post.id} 
                        className="p-4 border rounded shadow bg-white hover:shadow-md transition-shadow break-words" 
                    >
                        <Link 
                            to={`/post/${post.id}`} 
                            className="text-xl font-semibold text-blue-600 hover:underline break-words"
                        >
                            {post.title}
                        </Link>
                        <p className="text-gray-700 mt-2 break-words">
                            {post.content?.slice(0, 100)}{post.content?.length > 100 ? '...' : ''}
                        </p>
                    </li>
                ))}
            </ul>
        </div>
    );
}