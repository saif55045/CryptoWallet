import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token to requests
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Auth API
export const authAPI = {
  signup: (data) => api.post('/auth/signup', data),
  login: (data) => api.post('/auth/login', data),
  verifyOTP: (data) => api.post('/auth/verify-otp', data),
  resendOTP: (data) => api.post('/auth/resend-otp', data),
  getProfile: () => api.get('/auth/profile'),
};

// Wallet API
export const walletAPI = {
  generateWallet: () => api.post('/wallet/generate'),
  getMyWallet: () => api.get('/wallet/my-wallet'),
  validateWalletId: (walletId) => api.get(`/wallet/validate/${walletId}`),
  getWalletInfo: (walletId) => api.get(`/wallet/info/${walletId}`),
  exportPrivateKey: () => api.get('/wallet/export-key'),
  getBeneficiaries: () => api.get('/wallet/beneficiaries'),
  addBeneficiary: (data) => api.post('/wallet/beneficiaries', data),
  deleteBeneficiary: (id) => api.delete(`/wallet/beneficiaries/${id}`),
};

// UTXO API
export const utxoAPI = {
  getMyBalance: () => api.get('/utxo/my-balance'),
  getBalance: (walletId) => api.get(`/utxo/balance/${walletId}`),
  getMyUTXOs: (includeSpent = false) => api.get(`/utxo/my-utxos?includeSpent=${includeSpent}`),
  getUTXOs: (walletId, includeSpent = false) => api.get(`/utxo/list/${walletId}?includeSpent=${includeSpent}`),
  getStats: () => api.get('/utxo/stats'),
  createCoinbase: (data) => api.post('/utxo/coinbase', data),
};

// Transaction API
export const transactionAPI = {
  send: (data) => api.post('/transaction/send', data),
  create: (data) => api.post('/transaction/create', data),
  broadcast: (data) => api.post('/transaction/broadcast', data),
  getMyTransactions: () => api.get('/transaction/my-transactions'),
  getTransaction: (txId) => api.get(`/transaction/${txId}`),
  getStats: () => api.get('/transaction/stats'),
};

// Blockchain API
export const blockchainAPI = {
  getStats: () => api.get('/blockchain/stats'),
  getBlocks: (page = 1, limit = 10) => api.get(`/blockchain/blocks?page=${page}&limit=${limit}`),
  getBlock: (identifier) => api.get(`/blockchain/block/${identifier}`),
  getLatestBlock: () => api.get('/blockchain/latest'),
  getMiningStatus: () => api.get('/blockchain/mining-status'),
  validate: () => api.get('/blockchain/validate'),
  createGenesis: (token) => api.post('/blockchain/genesis', {}, {
    headers: { Authorization: `Bearer ${token}` }
  }),
  mine: (token) => api.post('/blockchain/mine', {}, {
    headers: { Authorization: `Bearer ${token}` }
  }),
  getMyBlocks: (token) => api.get('/blockchain/my-blocks', {
    headers: { Authorization: `Bearer ${token}` }
  }),
};

// Zakat API
export const zakatAPI = {
  getSettings: () => api.get('/zakat/settings'),
  getRecipients: () => api.get('/zakat/recipients'),
  getSummary: (token) => api.get('/zakat/summary', {
    headers: { Authorization: `Bearer ${token}` }
  }),
  calculate: (token) => api.post('/zakat/calculate', {}, {
    headers: { Authorization: `Bearer ${token}` }
  }),
  pay: (token, data) => api.post('/zakat/pay', data, {
    headers: { Authorization: `Bearer ${token}` }
  }),
  getHistory: (token) => api.get('/zakat/history', {
    headers: { Authorization: `Bearer ${token}` }
  }),
};

// Activity Logs API
export const logsAPI = {
  getActivityLogs: (token, params = {}) => {
    const queryParams = new URLSearchParams(params).toString();
    return api.get(`/logs/?${queryParams}`, {
      headers: { Authorization: `Bearer ${token}` }
    });
  },
  getActivityStats: (token) => api.get('/logs/stats', {
    headers: { Authorization: `Bearer ${token}` }
  }),
};

// Reports API
export const reportsAPI = {
  generateTransactionReport: (token, data) => api.post('/reports/transactions', data, {
    headers: { Authorization: `Bearer ${token}` }
  }),
  getWalletReport: (token) => api.get('/reports/wallet', {
    headers: { Authorization: `Bearer ${token}` }
  }),
};

// Admin API
export const adminAPI = {
  getStats: (token) => api.get('/admin/stats', {
    headers: { Authorization: `Bearer ${token}` }
  }),
  getUsers: (token) => api.get('/admin/users', {
    headers: { Authorization: `Bearer ${token}` }
  }),
  getTransactions: (token) => api.get('/admin/transactions', {
    headers: { Authorization: `Bearer ${token}` }
  }),
  getBlocks: (token) => api.get('/admin/blocks', {
    headers: { Authorization: `Bearer ${token}` }
  }),
  getLogs: (token) => api.get('/admin/logs', {
    headers: { Authorization: `Bearer ${token}` }
  }),
  toggleAdmin: (token, userId) => api.put(`/admin/users/${userId}/toggle-admin`, {}, {
    headers: { Authorization: `Bearer ${token}` }
  }),
  deleteUser: (token, userId) => api.delete(`/admin/users/${userId}`, {
    headers: { Authorization: `Bearer ${token}` }
  }),
  googleAuth: (token) => api.post('/auth/google', { token }),
};

// Combined API object for convenient imports
const apiService = {
  auth: authAPI,
  wallet: walletAPI,
  utxo: utxoAPI,
  transaction: transactionAPI,
  blockchain: blockchainAPI,
  zakat: zakatAPI,
  logs: logsAPI,
  reports: reportsAPI,
  admin: adminAPI,
};

export default apiService;
