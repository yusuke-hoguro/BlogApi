import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import client from "../api/client";

export default function PostCreate() {
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");
  const navigate = useNavigate();

  async function handleSubmit(e) {
    e.preventDefault();
    if (!title.trim() || !content.trim()) return;

    setSubmitting(true);
    setErrorMsg("");

    try {
      const token = localStorage.getItem("token");
      await client.post(
        "/api/posts",
        { title, content }
      );
      navigate("/"); // 投稿一覧へ戻る
    } catch (error) {
      console.error("投稿作成エラー:", error);
      setErrorMsg("投稿の作成に失敗しました。");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* 共通外枠 */}
      <div className="max-w-5xl mx-auto px-4 py-10">
        {/* フォームは少し細めに */}
        <div className="max-w-3xl mx-auto">
          <h1 className="text-3xl font-bold text-gray-900 mb-6">新規投稿作成</h1>

          {errorMsg && (
            <div className="mb-4 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-red-700">
              {errorMsg}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-5">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1" htmlFor="titleInput">
                タイトル
              </label>
              <input
                id="titleInput"
                type="text"
                placeholder="タイトル"
                className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500/40"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                disabled={submitting}
                maxLength={100}
              />
              <div className="mt-1 text-right text-xs text-gray-400">{title.length} / 100</div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1" htmlFor="contentInput">
                内容
              </label>
              <textarea
                id="contentInput"
                className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500/40 resize-none"
                rows={8}
                value={content}
                onChange={(e) => setContent(e.target.value)}
                disabled={submitting}
                maxLength={1000}
                placeholder="本文を入力..."
              />
              <div className="mt-1 text-right text-xs text-gray-400">{content.length} / 1000</div>
            </div>

            <button
              type="submit"
              disabled={submitting || !title.trim() || !content.trim()}
              className="inline-flex items-center justify-center rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-blue-700 disabled:bg-gray-400 transition focus:outline-none focus:ring-2 focus:ring-blue-500/40"
            >
              {submitting ? "送信中..." : "投稿作成"}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}
