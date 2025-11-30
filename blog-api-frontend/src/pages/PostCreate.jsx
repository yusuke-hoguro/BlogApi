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
        { title, content },
        { headers: { Authorization: token } }
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
    <div className="p-4 max-w-3xl mx-auto">
      <h1 className="text-2xl font-bold mb-4">新規投稿作成</h1>
      {errorMsg && <p className="text-red-500 mb-2">{errorMsg}</p>}

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block font-medium mb-1" htmlFor="titleInput">タイトル</label>
          <input
            id="titleInput"
            type="text"
            placeholder="タイトル"
            className="w-full border rounded p-2"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            disabled={submitting}
            maxLength={100}
          />
        </div>

        <div>
          <label className="block font-medium mb-1" htmlFor="contentInput">内容</label>
          <textarea
            id="contentInput"
            className="w-full border rounded p-2 resize-none"
            rows={6}
            value={content}
            onChange={(e) => setContent(e.target.value)}
            disabled={submitting}
            maxLength={1000}
          />
        </div>

        <button
          type="submit"
          disabled={submitting || !title.trim() || !content.trim()}
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:bg-gray-400"
        >
          {submitting ? "送信中..." : "投稿作成"}
        </button>
      </form>
    </div>
  );
}
