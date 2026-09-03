import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, waitFor, act } from "@testing-library/react";
import React from "react";
import { useAuth } from "@/hooks/use-auth";
import { useAuthStore } from "@/stores/auth-store";

vi.mock("next/navigation", () => ({
  useRouter: () => ({
    push: vi.fn(),
  }),
}));

vi.mock("@/lib/api", () => ({
  authApi: {
    login: vi.fn(),
    me: vi.fn(),
    logout: vi.fn(),
  },
}));

import { authApi } from "@/lib/api";
const mockAuthApi = vi.mocked(authApi);

describe("useAuth", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    useAuthStore.setState({
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
      user: null,
    });
  });

  it("returns isAuthenticated and user from store", () => {
    useAuthStore.setState({
      isAuthenticated: true,
      user: { id: "u1", email: "a@b.com", name: "A", role: "admin" },
    });

    const { result } = renderHook(() => useAuth());
    expect(result.current.isAuthenticated).toBe(true);
    expect(result.current.user?.email).toBe("a@b.com");
  });

  it("login calls authApi.login then authApi.me then sets store", async () => {
    mockAuthApi.login.mockResolvedValueOnce({ access_token: "a", refresh_token: "b" });
    mockAuthApi.me.mockResolvedValueOnce({ id: "u1", email: "t@t.com", name: "T", role: "admin" });

    const { result } = renderHook(() => useAuth());

    await act(async () => {
      await result.current.login("t@t.com", "pass");
    });

    expect(mockAuthApi.login).toHaveBeenCalledWith("t@t.com", "pass");
    expect(mockAuthApi.me).toHaveBeenCalled();

    const state = useAuthStore.getState();
    expect(state.isAuthenticated).toBe(true);
    expect(state.user?.email).toBe("t@t.com");
  });

  it("logout calls authApi.logout then clears store", async () => {
    mockAuthApi.logout.mockResolvedValueOnce(undefined as any);
    useAuthStore.setState({
      isAuthenticated: true,
      user: { id: "u1", email: "a@b.com", name: "A", role: "admin" },
    });

    const { result } = renderHook(() => useAuth());

    await act(async () => {
      await result.current.logout();
    });

    expect(mockAuthApi.logout).toHaveBeenCalled();
    const state = useAuthStore.getState();
    expect(state.isAuthenticated).toBe(false);
    expect(state.user).toBeNull();
  });

  it("logout clears state even if API fails", async () => {
    mockAuthApi.logout.mockRejectedValueOnce(new Error("network"));
    useAuthStore.setState({
      isAuthenticated: true,
      user: { id: "u1", email: "a@b.com", name: "A", role: "admin" },
    });

    const { result } = renderHook(() => useAuth());

    await act(async () => {
      await result.current.logout();
    });

    const state = useAuthStore.getState();
    expect(state.isAuthenticated).toBe(false);
  });
});
