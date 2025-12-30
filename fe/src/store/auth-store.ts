import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import Cookies from 'js-cookie';
import type { TokenResponse, User } from '@/types';
import { config } from '@/lib/config';

interface AuthState {
  isAuthenticated: boolean;
  token: string | null;
  refreshToken: string | null;
  tokenExpiry: string | null;
  role: string | null;
  user: User | null;
  
  // Actions
  login: (tokenResponse: TokenResponse, user?: User) => void;
  setUser: (user: User) => void;
  logout: () => void;
  checkAuth: () => boolean;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      isAuthenticated: false,
      token: null,
      refreshToken: null,
      tokenExpiry: null,
      role: null,
      user: null,
      
      login: (tokenResponse: TokenResponse, user?: User) => {
        const { token, refresh_token, expired_at, refresh_token_expired_at, role } = tokenResponse;
        
        // Store tokens in cookies for HTTP-only like behavior
        Cookies.set(config.storageKeys.accessToken, token, {
          expires: new Date(expired_at),
          sameSite: 'strict',
          secure: process.env.NODE_ENV === 'production',
        });
        Cookies.set(config.storageKeys.refreshToken, refresh_token, {
          expires: new Date(refresh_token_expired_at),
          sameSite: 'strict',
          secure: process.env.NODE_ENV === 'production',
        });
        
        set({
          isAuthenticated: true,
          token,
          refreshToken: refresh_token,
          tokenExpiry: expired_at,
          role,
          user: user || null,
        });
      },
      
      setUser: (user: User) => {
        set({ user });
      },
      
      logout: () => {
        // Clear cookies
        Cookies.remove(config.storageKeys.accessToken);
        Cookies.remove(config.storageKeys.refreshToken);
        
        set({
          isAuthenticated: false,
          token: null,
          refreshToken: null,
          tokenExpiry: null,
          role: null,
          user: null,
        });
      },
      
      checkAuth: () => {
        const token = Cookies.get(config.storageKeys.accessToken);
        const isAuth = !!token;
        
        if (!isAuth && get().isAuthenticated) {
          get().logout();
        }
        
        return isAuth;
      },
    }),
    {
      name: 'auth-storage',
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({
        isAuthenticated: state.isAuthenticated,
        role: state.role,
        tokenExpiry: state.tokenExpiry,
        user: state.user, // Persist user info
      }),
    }
  )
);
