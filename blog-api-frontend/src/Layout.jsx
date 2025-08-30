import React from 'react';
import Header from './components/Header';
import { Outlet } from 'react-router-dom';

export default function Layout() {
  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      <main className="p-4">
        {/* 子ルートのコンポーネントが表示される */}
        <Outlet />
      </main>
    </div>
  );
}
