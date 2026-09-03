// ============================================================================
// Auth Store — Zustand store for authentication state + JWT tokens
// ============================================================================

import { create } from "zustand";
import { persist } from "zustand/middleware";

import type { UserRole } from "@/types/api";

interface AuthState {
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  user: {
    id: string;
    email: string;
    name: string;
    role: UserRole;
  } | null;
  setTokens: (accessToken: string, refreshToken: string) => void;
  clearTokens: () => void;
  login: (id: string, email: string, name: string, role: UserRole) => void;
  logout: () => void;
  updateUser: (user: { id: string; email: string; name: string; role: UserRole }) => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
      user: null,

      setTokens: (accessToken: string, refreshToken: string) =>
        set({ accessToken, refreshToken }),

      clearTokens: () =>
        set({ accessToken: null, refreshToken: null }),

      login: (id: string, email: string, name: string, role: UserRole) =>
        set({
          isAuthenticated: true,
          user: { id, email, name, role },
        }),

      logout: () =>
        set({
          accessToken: null,
          refreshToken: null,
          isAuthenticated: false,
          user: null,
        }),

      updateUser: (user) => set({ user }),
    }),
    {
      name: "auth-storage",
      partialize: (state) => ({
        accessToken: state.accessToken,
        refreshToken: state.refreshToken,
        isAuthenticated: state.isAuthenticated,
        user: state.user,
      }),
    }
  )
);
