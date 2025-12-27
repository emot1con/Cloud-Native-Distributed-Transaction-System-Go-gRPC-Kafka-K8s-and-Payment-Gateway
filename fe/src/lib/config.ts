// Environment configuration with type safety
// All environment variables should be defined here

export const config = {
  // API Configuration - hardcoded for now due to env issues
  apiUrl: process.env.NEXT_PUBLIC_API_URL || 'http://192.168.49.2:31123',
  
  // OAuth Configuration
  oauthRedirectUrl: process.env.NEXT_PUBLIC_OAUTH_REDIRECT_URL || 'http://localhost:3000',
  
  // App Configuration
  appName: process.env.NEXT_PUBLIC_APP_NAME || 'Microservice Store',
  appDescription: process.env.NEXT_PUBLIC_APP_DESCRIPTION || 'E-Commerce Platform',
  
  // API Endpoints
  endpoints: {
    // Auth
    register: '/api/v1/auth/register',
    login: '/api/v1/auth/login',
    refreshToken: '/api/v1/auth/refresh-token',
    googleOauth: '/api/v1/auth/google',
    facebookOauth: '/api/v1/auth/facebook',
    githubOauth: '/api/v1/auth/github',
    
    // Products (trailing slash required)
    products: '/product/',
    
    // Orders (trailing slash required)
    orders: '/order/',
    
    // Payments
    paymentTransaction: '/payment/transaction',
  },
  
  // Storage keys
  storageKeys: {
    accessToken: 'access_token',
    refreshToken: 'refresh_token',
    tokenExpiry: 'token_expiry',
    cart: 'cart',
  },
} as const;

export type Config = typeof config;
