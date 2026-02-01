/* 
明日の方針を記載
1.入口で受信データのサイズチェック
2.C++で受信したバイト文字列のJSONデータを作成した関数にわたしてハッシュ化とプレビューにしてログに出力する
3.nlohmann/jsonライブラリを使用してJSONデータに変換する
4.変換したJSONデータを引数にして値バリデーション関数を呼び出す
5.問題なければ受信したバイト文字列は信用できるのでそのままTCPで別アプリに送信する
6.受信したアプリ側は１でつかった関数を使用してログに出力する
7.nlohmann/jsonライブラリを使用してJSONデータに変換する
8.変換したJSONデータを引数にして値バリデーション関数を呼び出す。3でつかったものと同じ関数を使用する
9.問題なければ受信したバイト文字列は信用できるのでそのまま処理を続行する

※JSONデータを構造体にデコードする関数は不要なのでバックログとして自分のC++ライブラリにいれておく

*/


// ログ出力のサンプルコード
inline void LogBytesPreview(const char* tag, const void* buf, size_t len) {
  uint64_t h = XXH3_64bits(buf, len);
  const size_t kMax = 256;
  size_t preview_len = std::min(len, kMax);
  const char* suffix = (len > kMax) ? "..." : "";

  LogD("%s len=%zu hash=%016llx preview=%.*s%s",
           tag, len, (unsigned long long)h,
           (int)preview_len, (const char*)buf, suffix);
}

// 以下はサイズチェックのサンプル

constexpr size_t MAX_JSON_BYTES = 64 * 1024; // 64KB

if (len == 0 || len > MAX_JSON_BYTES) {
  LOG_WARN("reject: invalid size len=%zu", len);
  return;
}
