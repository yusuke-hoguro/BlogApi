import React from "react";
import { Link, useNavigate } from "react-router-dom";

export default function Header(){
    const navigate = useNavigate();

    // ログアウト用の関数
    function handleLogout(){
        // Webブラウザに保存してあるトークンを削除
        localStorage.removeItem("token");
        navigate("/login", { replace: true });
    }

    // トークンを取得
    const token = localStorage.getItem("token");

    return(
        <header className="bg-gradient-to-r from-blue-600 to-blue-800 text-white shadow-lg">
            <div className="px-6 py-5 flex justify-between items-center">
                <Link to="/" className="text-3xl font-semibold text-white/95 visited:text-white hover:text-white transition hover:opacity-90">
                    BlogAPI
                </Link>
                <nav>
                    {token ? (
                        <button
                            onClick={handleLogout}
                            className="px-4 py-2 text-white/80 hover:text-white rounded transition-colors"
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
            </div>
        </header>
    );
}
