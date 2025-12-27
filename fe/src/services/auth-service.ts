import { api } from '@/lib/api';
import { config } from '@/lib/config';
import type {
  RegisterPayload,
  LoginPayload,
  TokenResponse,
} from '@/types';

export const authService = {
  /**
   * Register a new user
   */
  register: async (payload: RegisterPayload): Promise<{ message: string }> => {
    const response = await api.post(config.endpoints.register, payload);
    return response.data;
  },

  /**
   * Login with email and password
   */
  login: async (payload: LoginPayload): Promise<TokenResponse> => {
    const response = await api.post(config.endpoints.login, payload);
    return response.data;
  },

  /**
   * Get Google OAuth URL
   */
  getGoogleOAuthUrl: (): string => {
    return `${config.apiUrl}${config.endpoints.googleOauth}`;
  },

  /**
   * Get Facebook OAuth URL
   */
  getFacebookOAuthUrl: (): string => {
    return `${config.apiUrl}${config.endpoints.facebookOauth}`;
  },

  /**
   * Get GitHub OAuth URL
   */
  getGithubOAuthUrl: (): string => {
    return `${config.apiUrl}${config.endpoints.githubOauth}`;
  },
};
