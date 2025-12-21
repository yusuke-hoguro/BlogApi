import axios from 'axios';

console.log("VITE_API_BASE_URL =", import.meta.env.VITE_API_BASE_URL);

const client = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL, // Go APIのベースURL
  headers: {
    'Content-Type': 'application/json',
  },
});

// リクエスト：tokenを自動付与
client.interceptors.request.use(config => {
  const token = localStorage.getItem('token');
  if(token){
    config.headers.Authorization = `${token}`;
  }
  return config;
});

// レスポンス：認証エラー時はログアウト
client.interceptors.response.use(
  res => res,
  err => {
    const status = err.response?.status;
    if(status === 401 || status === 403) {
      console.warn('認証エラーにより、ログアウト');
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(err)
  }
)

export default client;
