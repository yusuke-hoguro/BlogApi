import React, { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import client from '../api/client';

export default function Login(){
    // useState:状態管理フック 変数の初期値を設定し、その変数を更新するための関数を返す
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);
    // useNavigate:「navigate」関数を取得して任意のパスに移動できる
    const navigate = useNavigate();
    // useLocation: 現在のURLの場所に関する情報を取得
    const location = useLocation();
    // 元のアクセス先があれば戻す
    const from = location.state?.from?.pathname || '/';

    async function handleSubmit(e) {
        // ブラウザのデフォルト動作を停止させる
        e.preventDefault();
        setError('');
        setLoading(true);
        try{
            const res = await client.post('/login',{ username, password });
            // 取得したトークンを保存する
            localStorage.setItem('token', res.data.token);
            // ログイン成功時はトップページへ
            navigate(from, { replace: true });
        } catch(error){
            console.error("ログイン失敗:", error)
            setError('ログインに失敗しました。ユーザー名とパスワードを確認してください。');
        } finally{
            setLoading(false);
        }
    }

    return(
        <div className='p-4 max-w-md mx-auto mt-10 bg-gray-50 border rounded shadow  box-border'>
            <h1 className='text-2xl font-bold mb-4 text-center'>ログイン</h1>
            {error && <p className='text-red-500 mb-4 text-center'>{error}</p>}
            <form onSubmit={handleSubmit} className="space-y-4 px-4">
                <input 
                    type="text" 
                    placeholder='ユーザー名'
                    value={username}
                    onChange={(e) => setUsername(e.target.value)} 
                    className='w-full border rounded p-2 focus:outline-none focus:ring-2 focus:ring-blue-400 box-border'
                />
                <input 
                    type="password" 
                    placeholder='パスワード'
                    value={password}
                    onChange={(e) => setPassword(e.target.value)} 
                    className='w-full border rounded p-2 focus:outline-none focus:ring-2 focus:ring-blue-400 box-border'
                />
                <button
                    type='submit'
                    disabled={loading}
                    className='w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600 disabled:bg-gray-400 transition-colors  box-border'
                >
                    {loading ? "ログイン中..." : "ログイン"}
                </button>
            </form>
        </div>
    );
}

