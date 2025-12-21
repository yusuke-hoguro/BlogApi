import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Layout from './Layout.jsx';
import PostCreate from './pages/PostCreate.jsx';
import PostList from './pages/PostList.jsx';
import PostDetail from './pages/PostDetail.jsx';
import PostEdit from './pages/PostEdit.jsx';
import Login from "./pages/Login.jsx";
import './index.css'
import RequireAuth from './components/RequireAuth.jsx';

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <BrowserRouter>
      <Routes>
        {/* Headerありの共通レイアウト */}
        <Route element={<RequireAuth/>}>
          <Route element={<Layout />}>
            <Route path="/" element={<PostList />} />
            <Route path="/post/create" element={<PostCreate />} />
            <Route path="/post/:id" element={<PostDetail />} />
            <Route path="/post/:id/edit" element={<PostEdit />} />
          </Route>
        </Route>

        {/* Headerなしの単独ルート */}
        <Route path="/login" element={<Login />} />
      </Routes>
    </BrowserRouter>
  </React.StrictMode>
)
