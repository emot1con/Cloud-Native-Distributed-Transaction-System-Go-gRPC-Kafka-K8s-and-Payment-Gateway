import axios, { AxiosError, AxiosInstance, InternalAxiosRequestConfig } from 'axios';
import Cookies from 'js-cookie';
import { config } from './config';
import type { ApiError } from '@/types';

// Create axios instance with default config
const api: AxiosInstance = axios.create({
  baseURL: config.apiUrl,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor - add auth token
api.interceptors.request.use(
  (requestConfig: InternalAxiosRequestConfig) => {
    const token = Cookies.get(config.storageKeys.accessToken);
    if (token && requestConfig.headers) {
      // Token from backend already includes "Bearer " prefix
      // If it starts with "Bearer ", use as-is; otherwise add prefix
      requestConfig.headers.Authorization = token.startsWith('Bearer ') ? token : `Bearer ${token}`;
    }
    return requestConfig;
  },
  (error: AxiosError) => Promise.reject(error)
);

// Response interceptor - handle errors and token refresh
api.interceptors.response.use(
  (response: import('axios').AxiosResponse) => response,
  async (error: AxiosError<ApiError>) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };
    
    // Handle 401 Unauthorized - try to refresh token
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      
      try {
        const refreshToken = Cookies.get(config.storageKeys.refreshToken);
        if (refreshToken) {
          const response = await axios.get(`${config.apiUrl}${config.endpoints.refreshToken}`, {
            headers: {
              Authorization: `Bearer ${refreshToken}`,
            },
          });
          
          const { token, refresh_token, expired_at, refresh_token_expired_at } = response.data;
          
          // Update tokens in cookies
          Cookies.set(config.storageKeys.accessToken, token, { 
            expires: new Date(expired_at),
            sameSite: 'strict',
          });
          Cookies.set(config.storageKeys.refreshToken, refresh_token, { 
            expires: new Date(refresh_token_expired_at),
            sameSite: 'strict',
          });
          
          // Retry original request with new token
          if (originalRequest.headers) {
            originalRequest.headers.Authorization = `Bearer ${token}`;
          }
          return api(originalRequest);
        }
      } catch (refreshError) {
        // Refresh failed - clear tokens and redirect to login
        Cookies.remove(config.storageKeys.accessToken);
        Cookies.remove(config.storageKeys.refreshToken);
        
        if (typeof window !== 'undefined') {
          window.location.href = '/login';
        }
        return Promise.reject(refreshError);
      }
    }
    
    // Extract error message
    const message = error.response?.data?.error || error.message || 'An error occurred';
    return Promise.reject(new Error(message));
  }
);

export { api };
