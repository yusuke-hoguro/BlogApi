import axios from 'axios';

console.log("VITE_API_BASE_URL =", import.meta.env.VITE_API_BASE_URL);
const client = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL, // Go APIのベースURL
  headers: {
    'Content-Type': 'application/json',
  },
});

export default client;
