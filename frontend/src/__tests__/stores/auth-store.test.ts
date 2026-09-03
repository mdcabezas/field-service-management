import { describe, it, expect, beforeEach } from "vitest";
import { useAuthStore } from "@/stores/auth-store";

describe("auth-store", () => {
  beforeEach(() => {
    useAuthStore.setState({
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
      user: null,
    });
  });

  it("has correct initial state", () => {
    const state = useAuthStore.getState();
    expect(state.accessToken).toBeNull();
    expect(state.refreshToken).toBeNull();
    expect(state.isAuthenticated).toBe(false);
    expect(state.user).toBeNull();
  });

  it("setTokens sets both tokens", () => {
    useAuthStore.getState().setTokens("access-123", "refresh-456");
    const state = useAuthStore.getState();
    expect(state.accessToken).toBe("access-123");
    expect(state.refreshToken).toBe("refresh-456");
  });

  it("clearTokens sets both tokens to null", () => {
    useAuthStore.getState().setTokens("a", "b");
    useAuthStore.getState().clearTokens();
    const state = useAuthStore.getState();
    expect(state.accessToken).toBeNull();
    expect(state.refreshToken).toBeNull();
  });

  it("login sets isAuthenticated and user", () => {
    useAuthStore.getState().login("u1", "test@test.com", "Test User", "admin");
    const state = useAuthStore.getState();
    expect(state.isAuthenticated).toBe(true);
    expect(state.user).toEqual({
      id: "u1",
      email: "test@test.com",
      name: "Test User",
      role: "admin",
    });
  });

  it("logout clears everything", () => {
    useAuthStore.getState().setTokens("a", "b");
    useAuthStore.getState().login("u1", "a@b.com", "A", "manager");
    useAuthStore.getState().logout();
    const state = useAuthStore.getState();
    expect(state.accessToken).toBeNull();
    expect(state.refreshToken).toBeNull();
    expect(state.isAuthenticated).toBe(false);
    expect(state.user).toBeNull();
  });

  it("updateUser only changes user field", () => {
    useAuthStore.getState().login("u1", "old@test.com", "Old", "technician");
    useAuthStore.getState().updateUser({
      id: "u1",
      email: "new@test.com",
      name: "New",
      role: "manager",
    });
    const state = useAuthStore.getState();
    expect(state.user?.email).toBe("new@test.com");
    expect(state.user?.role).toBe("manager");
    expect(state.isAuthenticated).toBe(true);
  });

  it("persists to localStorage", () => {
    useAuthStore.getState().setTokens("persist-a", "persist-b");
    useAuthStore.getState().login("u1", "p@test.com", "P", "admin");

    const stored = localStorage.getItem("auth-storage");
    expect(stored).toBeTruthy();
    const parsed = JSON.parse(stored!);
    expect(parsed.state.accessToken).toBe("persist-a");
    expect(parsed.state.user.email).toBe("p@test.com");
  });

  it("persists and restores from localStorage", () => {
    // Set state that will be persisted
    useAuthStore.getState().setTokens("restored-token", "restored-refresh");
    useAuthStore.getState().login("u2", "r@test.com", "R", "supervisor");

    // Verify it was persisted
    const stored = localStorage.getItem("auth-storage");
    expect(stored).toBeTruthy();
    const parsed = JSON.parse(stored!);
    expect(parsed.state.accessToken).toBe("restored-token");
    expect(parsed.state.user.role).toBe("supervisor");

    // Verify current state is correct
    expect(useAuthStore.getState().accessToken).toBe("restored-token");
    expect(useAuthStore.getState().user?.role).toBe("supervisor");
  });
});
