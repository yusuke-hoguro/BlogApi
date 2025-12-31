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
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// レスポンス：認証エラー時はログアウト
client.interceptors.response.use(
  res => res,
  err => {
    const status = err.response?.status;
    const requestURL = err.config?.url;
    const isLoginRequest = requestURL?.includes('/api/login');
    const hasToken = !!localStorage.getItem('token');
    switch (status) {
      case 401:
        if (!isLoginRequest && hasToken) {
          console.warn('認証エラーにより、ログアウト');
          localStorage.removeItem('token');
          window.location.href = '/login';
        }
        break;
      case 403:
        console.warn('権限がありません');
        break;
      default:
        break;
    }
    
    return Promise.reject(err)
  }
)

export default client;
