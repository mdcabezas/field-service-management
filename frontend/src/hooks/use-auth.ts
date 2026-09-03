"use client";

import { useRouter } from "next/navigation";
import { useAuthStore } from "@/stores/auth-store";
import { authApi } from "@/lib/api";

export function useAuth() {
  const router = useRouter();
  const { isAuthenticated, user, login, logout } = useAuthStore();

  const handleLogin = async (email: string, password: string) => {
    await authApi.login(email, password);
    const userInfo = await authApi.me();
    login(userInfo.id, userInfo.email, userInfo.name, userInfo.role as "admin" | "manager" | "supervisor" | "technician");
    router.push("/");
  };

  const handleLogout = async () => {
    try {
      await authApi.logout();
    } catch (error) {
      // Even if logout fails, clear local state
    } finally {
      logout();
      router.push("/login");
    }
  };

  return {
    isAuthenticated,
    user,
    login: handleLogin,
    logout: handleLogout,
  };
}
