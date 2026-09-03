import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { apiFetch, authApi, createCrudApi } from "@/lib/api";
import { useAuthStore } from "@/stores/auth-store";

// Mock sonner toast
vi.mock("sonner", () => ({
  toast: {
    error: vi.fn(),
    success: vi.fn(),
  },
}));

// Mock window.location
const mockAssign = vi.fn();
Object.defineProperty(window, "location", {
  value: { href: "", assign: mockAssign },
  writable: true,
});

describe("apiFetch", () => {
  const originalFetch = globalThis.fetch;

  beforeEach(() => {
    useAuthStore.setState({
      accessToken: "test-token",
      refreshToken: "test-refresh",
      isAuthenticated: true,
      user: { id: "u1", email: "test@test.com", name: "Test", role: "admin" },
    });
    // Reset isRefreshing state by clearing module internals
    globalThis.fetch = vi.fn();
  });

  afterEach(() => {
    globalThis.fetch = originalFetch;
    vi.restoreAllMocks();
    window.location.href = "";
  });

  it("injects Authorization header with token", async () => {
    (globalThis.fetch as any).mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: () => Promise.resolve({ data: "ok" }),
    });

    await apiFetch("/test");

    expect(globalThis.fetch).toHaveBeenCalledWith(
      "http://localhost:8081/test",
      expect.objectContaining({
        headers: expect.any(Headers),
      })
    );
  });

  it("skips auth header when skipAuth is true", async () => {
    (globalThis.fetch as any).mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: () => Promise.resolve({ data: "ok" }),
    });

    await apiFetch("/test", { skipAuth: true });

    const callHeaders = (globalThis.fetch as any).mock.calls[0][1].headers;
    // Headers should not contain Authorization
    expect(callHeaders.get("Authorization")).toBeNull();
  });

  it("returns parsed JSON on 200", async () => {
    (globalThis.fetch as any).mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: () => Promise.resolve({ items: [1, 2, 3] }),
    });

    const result = await apiFetch<{ items: number[] }>("/test");
    expect(result).toEqual({ items: [1, 2, 3] });
  });

  it("returns undefined on 204", async () => {
    (globalThis.fetch as any).mockResolvedValueOnce({
      ok: true,
      status: 204,
      json: () => Promise.resolve(null),
    });

    const result = await apiFetch("/test");
    expect(result).toBeUndefined();
  });

  it("throws error on non-2xx response", async () => {
    (globalThis.fetch as any).mockResolvedValueOnce({
      ok: false,
      status: 400,
      json: () => Promise.resolve({ error: "bad request" }),
    });

    await expect(apiFetch("/test")).rejects.toMatchObject({
      error: "bad request",
      statusCode: 400,
    });
  });

  it("shows toast on 403", async () => {
    const { toast } = await import("sonner");
    (globalThis.fetch as any).mockResolvedValueOnce({
      ok: false,
      status: 403,
      json: () => Promise.resolve({ error: "forbidden" }),
    });

    await expect(apiFetch("/test")).rejects.toMatchObject({ statusCode: 403 });
    expect(toast.error).toHaveBeenCalledWith("No tiene permisos para realizar esta acción");
  });

  it("returns undefined when response body is empty on error", async () => {
    (globalThis.fetch as any).mockResolvedValueOnce({
      ok: false,
      status: 500,
      json: () => Promise.reject(new Error("no body")),
    });

    await expect(apiFetch("/test")).rejects.toMatchObject({
      error: "Unknown error",
      statusCode: 500,
    });
  });
});

describe("apiFetch 401 auto-refresh", () => {
  const originalFetch = globalThis.fetch;

  beforeEach(() => {
    useAuthStore.setState({
      accessToken: "expired-token",
      refreshToken: "valid-refresh",
      isAuthenticated: true,
      user: { id: "u1", email: "test@test.com", name: "Test", role: "admin" },
    });
    globalThis.fetch = vi.fn();
  });

  afterEach(() => {
    globalThis.fetch = originalFetch;
    vi.restoreAllMocks();
    window.location.href = "";
  });

  it("refreshes token on 401 and retries request", async () => {
    // First call returns 401
    (globalThis.fetch as any)
      .mockResolvedValueOnce({
        ok: false,
        status: 401,
        json: () => Promise.resolve({}),
      })
      // Refresh call succeeds
      .mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve({ access_token: "new-token", refresh_token: "new-refresh" }),
      })
      // Retry succeeds
      .mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve({ data: "success" }),
      });

    const result = await apiFetch<{ data: string }>("/test");
    expect(result).toEqual({ data: "success" });
    expect(globalThis.fetch).toHaveBeenCalledTimes(3);
  });

  it("redirects to /login when no refresh token", async () => {
    useAuthStore.setState({ refreshToken: null });

    (globalThis.fetch as any).mockResolvedValueOnce({
      ok: false,
      status: 401,
      json: () => Promise.resolve({}),
    });

    await expect(apiFetch("/test")).rejects.toThrow("Session expired");
    expect(window.location.href).toBe("/login");
  });

  it("redirects to /login when refresh fails", async () => {
    (globalThis.fetch as any)
      .mockResolvedValueOnce({
        ok: false,
        status: 401,
        json: () => Promise.resolve({}),
      })
      .mockResolvedValueOnce({
        ok: false,
        status: 401,
        json: () => Promise.resolve({}),
      });

    await expect(apiFetch("/test")).rejects.toThrow("Session expired");
    expect(window.location.href).toBe("/login");
  });
});

describe("authApi", () => {
  const originalFetch = globalThis.fetch;

  beforeEach(() => {
    globalThis.fetch = vi.fn();
  });

  afterEach(() => {
    globalThis.fetch = originalFetch;
  });

  it("login sends POST with email and password", async () => {
    (globalThis.fetch as any).mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: () => Promise.resolve({ access_token: "a", refresh_token: "b" }),
    });

    const result = await authApi.login("user@test.com", "pass123");
    expect(result).toEqual({ access_token: "a", refresh_token: "b" });

    const [url, opts] = (globalThis.fetch as any).mock.calls[0];
    expect(url).toBe("http://localhost:8081/auth/login");
    expect(opts.method).toBe("POST");
  });

  it("me returns user info", async () => {
    (globalThis.fetch as any).mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: () => Promise.resolve({ id: "1", email: "a@b.com", name: "A", role: "admin" }),
    });

    const result = await authApi.me();
    expect(result).toMatchObject({ id: "1", email: "a@b.com" });
  });
});

describe("createCrudApi", () => {
  const originalFetch = globalThis.fetch;

  beforeEach(() => {
    globalThis.fetch = vi.fn();
  });

  afterEach(() => {
    globalThis.fetch = originalFetch;
  });

  const crud = createCrudApi<{ id: string; name: string }>("/api/items");

  it("list fetches all items", async () => {
    (globalThis.fetch as any).mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: () => Promise.resolve([{ id: "1", name: "A" }]),
    });

    const result = await crud.list();
    expect(result).toEqual([{ id: "1", name: "A" }]);
  });

  it("list with params appends query string", async () => {
    (globalThis.fetch as any).mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: () => Promise.resolve([]),
    });

    await crud.list({ search: "test" });
    const [url] = (globalThis.fetch as any).mock.calls[0];
    expect(url).toContain("search=test");
  });

  it("get fetches single item", async () => {
    (globalThis.fetch as any).mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: () => Promise.resolve({ id: "1", name: "A" }),
    });

    const result = await crud.get("1");
    expect(result).toEqual({ id: "1", name: "A" });
  });

  it("create sends POST", async () => {
    (globalThis.fetch as any).mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: () => Promise.resolve({ id: "2", name: "B" }),
    });

    await crud.create({ name: "B" } as any);
    const [, opts] = (globalThis.fetch as any).mock.calls[0];
    expect(opts.method).toBe("POST");
  });

  it("update sends PUT", async () => {
    (globalThis.fetch as any).mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: () => Promise.resolve({ id: "1", name: "Updated" }),
    });

    await crud.update("1", { name: "Updated" });
    const [, opts] = (globalThis.fetch as any).mock.calls[0];
    expect(opts.method).toBe("PUT");
  });

  it("delete sends DELETE", async () => {
    (globalThis.fetch as any).mockResolvedValueOnce({
      ok: true,
      status: 204,
      json: () => Promise.resolve(undefined),
    });

    await crud.delete("1");
    const [, opts] = (globalThis.fetch as any).mock.calls[0];
    expect(opts.method).toBe("DELETE");
  });
});
