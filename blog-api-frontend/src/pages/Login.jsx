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
            const res = await client.post('/api/login',{ username, password });
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
        <>
            <header className="bg-gradient-to-r from-blue-600 to-blue-800 text-white shadow-lg">
                <div className="max-w-7xl mx-auto px-4 py-6 flex justify-center">
                    <h1 className="text-3xl font-bold">BlogAPI</h1>
                </div>
            </header>
            <div className='min-h-[calc(100vh-180px)] flex items-center justify-center bg-gray-50'>
                <div className='p-8 max-w-md bg-white rounded-lg shadow-lg overflow-x-hidden box-border'>
                    <form onSubmit={handleSubmit} className="space-y-6">
                        {error && <p className='text-red-500 text-center mb-4'>{error}</p>}
                        <input 
                            type="text" 
                            placeholder='ユーザー名'
                            value={username}
                            onChange={(e) => setUsername(e.target.value)} 
                            className='w-full rounded-lg p-8 bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-400 box-border text-2xl'
                        />
                        <input 
                            type="password" 
                            placeholder='パスワード'
                            value={password}
                            onChange={(e) => setPassword(e.target.value)} 
                            className='w-full rounded-lg p-8 bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-400 box-border text-2xl'
                        />
                        <button
                            type='submit'
                            disabled={loading}
                            className='w-full bg-blue-500 text-white py-3 rounded-lg hover:bg-blue-600 disabled:bg-gray-400 transition-colors box-border font-medium text-lg'
                        >
                            {loading ? "ログイン中..." : "ログイン"}
                        </button>
                    </form>
                </div>
            </div>
        </>
   );
}

