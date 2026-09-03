// ============================================================================
// API Wrapper — Centralized fetch with JWT handling and auto-refresh
// ============================================================================

import { useAuthStore } from "@/stores/auth-store";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081";

// ---- Cookie helpers (for middleware compatibility) ----

function setCookie(name: string, value: string, days: number) {
  const expires = new Date(Date.now() + days * 864e5).toUTCString();
  document.cookie = `${name}=${encodeURIComponent(value)}; expires=${expires}; path=/; SameSite=Lax`;
}

function deleteCookie(name: string) {
  document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/`;
}

// ---- Request queue for concurrent 401s ----

let isRefreshing = false;
let failedQueue: Array<{
  resolve: (value: unknown) => void;
  reject: (reason?: unknown) => void;
}> = [];

const processQueue = (error: unknown, token: string | null = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  failedQueue = [];
};

// ---- Core fetch wrapper ----

interface ApiOptions extends RequestInit {
  skipAuth?: boolean;
}

export async function apiFetch<T>(
  path: string,
  options: ApiOptions = {}
): Promise<T> {
  const { skipAuth = false, ...fetchOptions } = options;

  const headers = new Headers(fetchOptions.headers);

  if (!skipAuth) {
    const { accessToken } = useAuthStore.getState();
    if (accessToken) {
      headers.set("Authorization", `Bearer ${accessToken}`);
    }
  }

  const res = await fetch(`${API_URL}${path}`, {
    ...fetchOptions,
    headers,
  });

  if (res.status === 401 && !skipAuth) {
    if (isRefreshing) {
      return new Promise((resolve, reject) => {
        failedQueue.push({ resolve, reject });
      }).then(() => {
        return apiFetch<T>(path, options);
      });
    }

    isRefreshing = true;

    try {
      const { refreshToken } = useAuthStore.getState();

      if (!refreshToken) {
        processQueue(new Error("No refresh token"));
        deleteCookie("access_token");
        window.location.href = "/login";
        throw new Error("Session expired");
      }

      const refreshRes = await fetch(`${API_URL}/auth/refresh`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ refresh_token: refreshToken }),
      });

      if (refreshRes.ok) {
        const data = await refreshRes.json();
        useAuthStore.getState().setTokens(data.access_token, data.refresh_token);
        setCookie("access_token", data.access_token, 1);
        processQueue(null);
        return apiFetch<T>(path, options);
      } else {
        processQueue(new Error("Refresh failed"));
        useAuthStore.getState().logout();
        deleteCookie("access_token");
        window.location.href = "/login";
        throw new Error("Session expired");
      }
    } catch (error) {
      processQueue(error);
      useAuthStore.getState().logout();
      deleteCookie("access_token");
      window.location.href = "/login";
      throw error;
    } finally {
      isRefreshing = false;
    }
  }

  if (!res.ok) {
    const error = await res.json().catch(() => ({
      error: "Unknown error",
      statusCode: res.status,
    }));

    if (res.status === 403) {
      const { toast } = await import("sonner");
      toast.error("No tiene permisos para realizar esta acción");
    }

    throw { ...error, statusCode: res.status };
  }

  if (res.status === 204) {
    return undefined as T;
  }

  return res.json();
}

// ---- Auth-specific API calls ----

export const authApi = {
  login: (email: string, password: string) =>
    apiFetch<{ access_token: string; refresh_token: string }>("/auth/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
      skipAuth: true,
    }),

  /** Call after login to persist tokens in both Zustand and cookie */
  setAuthTokens: (accessToken: string, refreshToken: string) => {
    useAuthStore.getState().setTokens(accessToken, refreshToken);
    setCookie("access_token", accessToken, 1);
  },

  refresh: (refreshToken: string) =>
    apiFetch<{ access_token: string; refresh_token: string }>("/auth/refresh", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: refreshToken }),
      skipAuth: true,
    }),

  logout: () =>
    apiFetch<void>("/auth/logout", {
      method: "POST",
    }).finally(() => {
      useAuthStore.getState().logout();
      deleteCookie("access_token");
    }),

  me: () =>
    apiFetch<{ id: string; email: string; name: string; role: string }>("/auth/me"),
};

// ---- Generic CRUD factory ----

export function createCrudApi<T extends { id: string }>(basePath: string) {
  return {
    list: (params?: Record<string, string>) => {
      const searchParams = params ? new URLSearchParams(params) : undefined;
      const url = searchParams ? `${basePath}?${searchParams}` : basePath;
      return apiFetch<T[]>(url);
    },

    get: (id: string) => apiFetch<T>(`${basePath}/${id}`),

    create: (data: Omit<T, "id" | "created_at" | "updated_at">) =>
      apiFetch<T>(basePath, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),

    update: (id: string, data: Partial<T>) =>
      apiFetch<T>(`${basePath}/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),

    delete: (id: string) =>
      apiFetch<void>(`${basePath}/${id}`, {
        method: "DELETE",
      }),
  };
}
