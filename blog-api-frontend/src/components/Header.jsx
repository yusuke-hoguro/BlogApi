import React from "react";
import { Link, useNavigate } from "react-router-dom";

export default function Header(){
    const navigate = useNavigate();

    // ログアウト用の関数
    function handleLogout(){
        // Webブラウザに保存してあるトークンを削除
        localStorage.removeItem("token");
        navigate("/login");
    }

    // トークンを取得
    const token = localStorage.getItem("token");

    return(
        <header className="flex justify-between items-center p-4 bg-gray-100 border-b shadow">
            <Link to="/" className="text-xl font-bold">
                My Blog
            </Link>
            <nav>
                {token ? (
                    <button
                        onClick={handleLogout}
                        className="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600 transition-colors"
                    >
                        ログアウト
                    </button>
                ):(
                    <Link
                        to="/login"
                        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors"
                    >
                        ログイン
                    </Link>
                )}
            </nav>
        </header>
    );
}