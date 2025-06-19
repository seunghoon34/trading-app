// API Configuration for Mobile Frontend
// This file centralizes all API endpoint configurations

const API_CONFIG = {
  // Base URL - Update this for your backend
  BASE_URL: __DEV__ 
    ? 'http://localhost:3000'  // Development - Local backend
    : 'https://your-production-api.com',  // Production - Update this URL

  // Timeout settings
  TIMEOUT: 10000, // 10 seconds

  // Endpoints
  ENDPOINTS: {
    // Authentication
    LOGIN: '/api/v1/auth/login',
    REGISTER: '/api/v1/auth/register',
    LOGOUT: '/api/v1/auth/logout',
    ME: '/api/v1/user/me',
    
    // Portfolio
    PORTFOLIO_PERFORMANCE: '/api/v1/portfolio/performance',
    PORTFOLIO_POSITIONS: '/api/v1/portfolio/positions',
    PORTFOLIO_SUMMARY: '/api/v1/portfolio/summary',
    
    // Trading
    ORDERS: '/api/v1/trading/orders',
    PLACE_ORDER: '/api/v1/trading/order',
    ORDER_HISTORY: '/api/v1/trading/history',
    
    // Market Data
    MARKET_DATA: '/api/v1/market/data',
    STOCK_QUOTE: '/api/v1/market/quote',
    TRENDING: '/api/v1/market/trending',
    GAINERS: '/api/v1/market/gainers',
    LOSERS: '/api/v1/market/losers',
    
    // Chatbot
    CHAT: '/api/v1/chat/message',
    CHAT_HISTORY: '/api/v1/chat/history',
  },

  // Headers
  getHeaders: (token = null) => {
    const headers = {
      'Content-Type': 'application/json',
    };
    
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }
    
    return headers;
  },

  // Full URL builder
  getFullUrl: (endpoint) => {
    return `${API_CONFIG.BASE_URL}${endpoint}`;
  },
};

export default API_CONFIG; 