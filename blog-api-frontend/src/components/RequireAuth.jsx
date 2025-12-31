import React from "react";
import { Navigate, replace, Outlet, useLocation } from "react-router-dom";

// トークンがなければログインページへリダイレクト
export default function RequireAuth(){
    const token = localStorage.getItem('token');
    // 現在のURLの場所に関する情報を取得
    const location = useLocation();

    if(!token){
        // トークンがなければログインページへリダイレクト
        // 元のアクセス先をstateに保存しておく
        return <Navigate to="/login" state={{from: location}} replace />;
    }
    return <Outlet />;
}
