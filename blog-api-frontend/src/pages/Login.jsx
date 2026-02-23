import React, { useState } from 'react';
import { useNavigate, useLocation, Link } from 'react-router-dom';
import client from '../api/client';
import { FiEye, FiEyeOff } from "react-icons/fi";

export default function Login(){
    // useState:状態管理フック 変数の初期値を設定し、その変数を更新するための関数を返す
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);
    const [showPassword, setShowPassword] = useState(false);
    // useNavigate:「navigate」関数を取得して任意のパスに移動できる
    const navigate = useNavigate();
    // useLocation: 現在のURLの場所に関する情報を取得
    const location = useLocation();
    // 元のアクセス先があれば戻す
    const from = location.state?.from?.pathname || '/';
    // 新規登録完了後に来た場合はメッセージを表示
    const registered = location.state?.registered;

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
            setError('ログインに失敗しました。\nユーザー名とパスワードを確認してください。');
        } finally{
            setLoading(false);
        }
    }

    return(
        <>
            <header className="bg-gradient-to-r from-blue-600 to-blue-800 text-white shadow-lg">
                <div className="max-w-7xl mx-auto px-4 py-6 flex justify-center">
                    <h1 data-testid="login-title" className="text-3xl font-bold">BlogAPI</h1>
                </div>
            </header>
            <div className='min-h-[calc(100vh-180px)] flex items-center justify-center bg-gray-50'>
                <div className='p-8 max-w-md bg-white rounded-lg shadow-lg overflow-x-hidden box-border'>
                    <form onSubmit={handleSubmit} className="space-y-6">
                        {registered && <p className='text-green-500 text-center mb-4'>登録が完了しました。ログインしてください。</p>}
                        {error && <p className='text-red-500 text-center mb-4 whitespace-pre-line'>{error}</p>}
                        <input 
                            type="text" 
                            placeholder='ユーザー名'
                            value={username}
                            onChange={(e) => setUsername(e.target.value)} 
                            className='w-full rounded-lg p-6 bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-400 box-border text-xl'
                        />
                        <div className="relative">
                            <input 
                                type={showPassword ? "text" : "password"} 
                                placeholder='パスワード'
                                value={password}
                                onChange={(e) => setPassword(e.target.value)} 
                                className='w-full rounded-lg p-6 pr-14 bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-400 box-border text-xl'
                            />
                            <button 
                                type="button" 
                                onClick={() => setShowPassword(!showPassword)}
                                className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 z-10"
                            >
                                {showPassword ? <FiEyeOff size={20} /> : <FiEye size={20} />}
                            </button>
                        </div>
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
            <div className="mt-6 text-center text-base text-gray-600">
                アカウントをお持ちでないですか？ {' '}
                <Link to="/signup" className="text-blue-600 hover:underline font-medium">
                    新規登録はこちら
                </Link>
            </div>

        </>
   );
}

