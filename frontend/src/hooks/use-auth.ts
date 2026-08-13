"use client";

import { useRouter } from "next/navigation";
import { useAuthStore } from "@/stores/auth-store";
import { authApi } from "@/lib/api";

export function useAuth() {
  const router = useRouter();
  const { isAuthenticated, user, login, logout } = useAuthStore();

  const handleLogin = async (employeeNumber: string, password: string) => {
    await authApi.login(employeeNumber, password);
    const userInfo = await authApi.me();
    login(userInfo.employee_number, userInfo.role as "admin" | "manager" | "supervisor" | "technician");
    router.push("/dashboard");
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
