import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import React from "react";
import { useCustomers, useCustomer, useCreateCustomer, useDeleteCustomer } from "@/hooks/use-customers";

vi.mock("@/lib/api", () => ({
  apiFetch: vi.fn(),
}));

import { apiFetch } from "@/lib/api";
const mockApiFetch = vi.mocked(apiFetch);

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  });
  return function Wrapper({ children }: { children: React.ReactNode }) {
    return React.createElement(QueryClientProvider, { client: queryClient }, children);
  };
}

describe("useCustomers", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("returns customer list", async () => {
    mockApiFetch.mockResolvedValueOnce({
      items: [
        { id: "1", name: "Client A" },
        { id: "2", name: "Client B" },
      ],
      total: 2,
    });

    const wrapper = createWrapper();
    const { result } = renderHook(() => useCustomers(), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toHaveLength(2);
    expect(result.current.data![0].name).toBe("Client A");
  });

  it("calls apiFetch with correct path", async () => {
    mockApiFetch.mockResolvedValueOnce({ items: [], total: 0 });

    const wrapper = createWrapper();
    renderHook(() => useCustomers(), { wrapper });

    await waitFor(() => expect(mockApiFetch).toHaveBeenCalledWith("/api/customers"));
  });
});

describe("useCustomer", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("fetches single customer by id", async () => {
    mockApiFetch.mockResolvedValueOnce({ id: "123", name: "Single" });

    const wrapper = createWrapper();
    const { result } = renderHook(() => useCustomer("123"), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.name).toBe("Single");
  });

  it("does not fetch when id is empty", () => {
    const wrapper = createWrapper();
    const { result } = renderHook(() => useCustomer(""), { wrapper });

    expect(result.current.fetchStatus).toBe("idle");
    expect(mockApiFetch).not.toHaveBeenCalled();
  });
});

describe("useCreateCustomer", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("calls POST on mutateAsync", async () => {
    mockApiFetch.mockResolvedValueOnce({ id: "new", name: "New Client" });

    const wrapper = createWrapper();
    const { result } = renderHook(() => useCreateCustomer(), { wrapper });

    await result.current.mutateAsync({ name: "New Client", tax_id: "", phone: "", email: "" } as any);

    expect(mockApiFetch).toHaveBeenCalledWith(
      "/api/customers",
      expect.objectContaining({ method: "POST" })
    );
  });
});

describe("useDeleteCustomer", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("calls DELETE on mutateAsync", async () => {
    mockApiFetch.mockResolvedValueOnce(undefined as any);

    const wrapper = createWrapper();
    const { result } = renderHook(() => useDeleteCustomer(), { wrapper });

    await result.current.mutateAsync("del-123");

    expect(mockApiFetch).toHaveBeenCalledWith(
      "/api/customers/del-123",
      expect.objectContaining({ method: "DELETE" })
    );
  });
});
