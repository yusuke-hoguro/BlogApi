import React, {useEffect, useState} from "react";
import client from "../api/client";

export default function PostList(){
    const [posts, setPosts] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        client.get('/posts')
            .then(response => {
                setPosts(response.data)
            })
            .catch(error => {
                console.error('投稿取得エラー:', error);
            })
            .finally(() => setLoading(false))
    }, []);
    
    if(loading) return <p className="p-4">読み込み中...</p>

    return(
        <div className="p-4">
            <hi className="text-2xl font-bold mb-4" >投稿一覧</hi>
            <ul>
                {posts.map(post => (
                    <li key={post.id} className="mb-4 p-4 border rounded shadow" >
                        <h2 className="text-xl font-semibold">{post.title}</h2>
                        <p className="text-gray-700">{post.content?.slice(0,100)}...</p>
                    </li>
                ))}
            </ul>
        </div>
    );
}