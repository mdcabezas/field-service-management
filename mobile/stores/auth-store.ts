import { create } from 'zustand';
import { login as apiLogin, logout as apiLogout, getAccessToken, getMe } from '../lib/auth';

interface AuthState {
  isAuthenticated: boolean;
  user: { id: string; email: string; name: string } | null;
  isLoading: boolean;
  error: string | null;
  login: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  checkAuth: () => Promise<void>;
}

export const useAuthStore = create<AuthState>((set) => ({
  isAuthenticated: false,
  user: null,
  isLoading: true,
  error: null,

  login: async (email: string, password: string) => {
    set({ isLoading: true, error: null });
    try {
      await apiLogin(email, password);
      const userInfo = await getMe();
      set({ isAuthenticated: true, user: { id: userInfo.id, email: userInfo.email, name: userInfo.name }, isLoading: false });
    } catch (error) {
      set({ error: String(error), isLoading: false });
      throw error;
    }
  },

  logout: async () => {
    await apiLogout();
    set({ isAuthenticated: false, user: null });
  },

  checkAuth: async () => {
    set({ isLoading: true });
    try {
      const token = await getAccessToken();
      if (token) {
        const userInfo = await getMe();
        set({ isAuthenticated: true, user: { id: userInfo.id, email: userInfo.email, name: userInfo.name }, isLoading: false });
      } else {
        set({ isAuthenticated: false, isLoading: false });
      }
    } catch {
      set({ isAuthenticated: false, isLoading: false });
    }
  },
}));
