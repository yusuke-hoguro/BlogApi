import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import client from '../api/client';
import { FiEye, FiEyeOff } from "react-icons/fi";

export default function Signup(){
    // useState:状態管理フック 変数の初期値を設定し、その変数を更新するための関数を返す
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [passwordConfirm, setPasswordConfirm] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);
    const [showPassword, setShowPassword] = useState(false);
    const [showPasswordConfirm, setShowPasswordConfirm] = useState(false);
    // useNavigate:「navigate」関数を取得して任意のパスに移動できる
    const navigate = useNavigate();

    async function handleSubmit(e) {
        // ブラウザのデフォルト動作を停止させる
        e.preventDefault();
        setError('');

        // ユーザー名のチェック
        if(!username.trim()) {
            setError('ユーザー名を入力してください。');
            return;
        }

        // パスワードのチェック        
        if(password.length < 8) {
            setError('パスワードは8文字以上で入力してください。');
            return;
        }

        // パスワードの確認
        if(password !== passwordConfirm) {
            setError('パスワードが一致しません。');
            return;
        }

        setLoading(true);
        try{
            const res = await client.post('/api/signup',{ username: username.trim(), password });
            // サインアップ成功時はログインページへ（登録完了後に来たことを値として渡す）
            navigate('/login', { state: { registered: true } });
        } catch(error){
            console.error("サインアップ失敗:", error)
            setError('サインアップに失敗しました。\nユーザー名とパスワードを確認してください。');
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
                    <h2 className="text-center text-lg font-semibold mb-6">新規登録</h2>
                    <form onSubmit={handleSubmit} className="space-y-6">
                        {error && <p className='text-red-500 text-center mb-4 whitespace-pre-line'>{error}</p>}
                        <input 
                            type="text" 
                            placeholder='ユーザー名'
                            value={username}
                            onChange={(e) => setUsername(e.target.value)} 
                            className='w-full rounded-lg p-5 bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-400 box-border text-xl'
                        />
                        <div className="relative">
                            <input 
                                type={showPassword ? "text" : "password"} 
                                placeholder='パスワード(8文字以上)'
                                value={password}
                                onChange={(e) => setPassword(e.target.value)} 
                                className='w-full rounded-lg p-5 pr-14 bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-400 box-border text-xl'
                            />
                            <button 
                                type="button" 
                                onClick={() => setShowPassword(!showPassword)}
                                className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 z-10"
                            >
                                {showPassword ? <FiEyeOff size={20} /> : <FiEye size={20} />}
                            </button>
                        </div>
                        <div className="relative">
                            <input 
                                type={showPasswordConfirm ? "text" : "password"} 
                                placeholder='パスワード確認'
                                value={passwordConfirm}
                                onChange={(e) => setPasswordConfirm(e.target.value)} 
                                className='w-full rounded-lg p-5 pr-14 bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-400 box-border text-xl'
                            />
                            <button 
                                type="button" 
                                onClick={() => setShowPasswordConfirm(!showPasswordConfirm)}
                                className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 z-10"
                            >
                                {showPasswordConfirm ? <FiEyeOff size={20} /> : <FiEye size={20} />}
                            </button>
                        </div>
                        <button
                            type='submit'
                            disabled={loading}
                            className='w-full bg-blue-500 text-white py-3 rounded-lg hover:bg-blue-600 disabled:bg-gray-400 transition-colors box-border font-medium text-lg'
                        >
                            {loading ? "新規登録中..." : "新規登録"}
                        </button>
                    </form>
                </div>
            </div>
            <div className="mt-6 text-center text-base text-gray-600">
                すでにアカウントをお持ちですか？ {' '}
                <Link to="/login" className="text-blue-600 hover:underline font-medium">
                    ログインはこちら
                </Link>
            </div>
        </>
   );
}

