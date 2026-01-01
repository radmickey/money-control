import axios from 'axios';

// Vite uses import.meta.env instead of process.env
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:9080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor for adding auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor for handling token refresh
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      const refreshToken = localStorage.getItem('refreshToken');
      if (refreshToken) {
        try {
          const response = await axios.post(`${API_BASE_URL}/auth/refresh`, {
            refreshToken,
          });

          const { accessToken, refreshToken: newRefreshToken } = response.data;
          localStorage.setItem('token', accessToken);
          localStorage.setItem('refreshToken', newRefreshToken);

          originalRequest.headers.Authorization = `Bearer ${accessToken}`;
          return api(originalRequest);
        } catch {
          localStorage.removeItem('token');
          localStorage.removeItem('refreshToken');
          window.location.href = '/login';
        }
      }
    }

    return Promise.reject(error);
  }
);

// Auth API
export const authAPI = {
  login: (credentials: { email: string; password: string }) =>
    api.post('/auth/login', credentials),
  register: (data: { email: string; password: string; firstName: string; lastName: string; baseCurrency?: string }) =>
    api.post('/auth/register', {
      email: data.email,
      password: data.password,
      first_name: data.firstName,
      last_name: data.lastName,
      base_currency: data.baseCurrency || 'USD',
    }),
  getProfile: () => api.get('/auth/profile'),
  refresh: (refreshToken: string) => api.post('/auth/refresh', { refresh_token: refreshToken }),
  googleAuth: (token: string) => api.post('/auth/google', { token }),
  // Telegram auth - sends initData from Telegram WebApp
  telegramAuth: (initData: string) => api.post('/auth/telegram', { init_data: initData }),
  // Get Google OAuth URL for redirect
  getGoogleAuthUrl: () => api.get('/auth/google/url'),
};

// Accounts API
export const accountsAPI = {
  list: () => api.get('/accounts'),
  get: (id: string) => api.get(`/accounts/${id}`),
  create: (data: any) => api.post('/accounts', data),
  update: (id: string, data: any) => api.put(`/accounts/${id}`, data),
  delete: (id: string) => api.delete(`/accounts/${id}`),
  // Sub-accounts
  listSubAccounts: (accountId: string) => api.get(`/accounts/${accountId}/sub-accounts`),
  createSubAccount: (accountId: string, data: any) =>
    api.post(`/accounts/${accountId}/sub-accounts`, data),
  updateSubAccount: (accountId: string, subId: string, data: any) =>
    api.put(`/accounts/${accountId}/sub-accounts/${subId}`, data),
  deleteSubAccount: (accountId: string, subId: string) =>
    api.delete(`/accounts/${accountId}/sub-accounts/${subId}`),
};

// Transactions API
export const transactionsAPI = {
  list: (params?: any) => api.get('/transactions', { params }),
  get: (id: string) => api.get(`/transactions/${id}`),
  create: (data: any) => api.post('/transactions', data),
  update: (id: string, data: any) => api.put(`/transactions/${id}`, data),
  delete: (id: string) => api.delete(`/transactions/${id}`),
};

// Assets API
export const assetsAPI = {
  list: () => api.get('/assets'),
  get: (id: string) => api.get(`/assets/${id}`),
  create: (data: any) => api.post('/assets', data),
  update: (id: string, data: any) => api.put(`/assets/${id}`, data),
  delete: (id: string) => api.delete(`/assets/${id}`),
  getPrice: (symbol: string, type: string) =>
    api.get(`/assets/price/${symbol}`, { params: { type } }),
};

// Currency API
export const currencyAPI = {
  getRate: (from: string, to: string) =>
    api.get('/currency/rate', { params: { from, to } }),
  convert: (from: string, to: string, amount: number) =>
    api.post('/currency/convert', { from, to, amount }),
};

// Insights API
export const insightsAPI = {
  getNetWorth: (baseCurrency: string) =>
    api.get('/insights/net-worth', { params: { baseCurrency } }),
  getAllocation: (baseCurrency: string) =>
    api.get('/insights/allocation', { params: { baseCurrency } }),
  getTrends: (baseCurrency: string, period: string) =>
    api.get('/insights/trends', { params: { baseCurrency, period } }),
};

export default api;
