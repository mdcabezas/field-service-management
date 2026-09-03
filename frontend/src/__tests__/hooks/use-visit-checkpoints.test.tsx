import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import React from "react";
import { useVisitCheckpoints, useCreateVisitCheckpoint } from "@/hooks/use-visit-checkpoints";

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

describe("useVisitCheckpoints", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("fetches checkpoints for a visit", async () => {
    mockApiFetch.mockResolvedValueOnce([{ id: "c1", visit_id: "v1", type: "gps" }]);

    const wrapper = createWrapper();
    const { result } = renderHook(() => useVisitCheckpoints("v1"), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toHaveLength(1);
    expect(mockApiFetch).toHaveBeenCalledWith("/api/visits/v1/checkpoints");
  });

  it("does not fetch when visitId is empty", () => {
    const wrapper = createWrapper();
    const { result } = renderHook(() => useVisitCheckpoints(""), { wrapper });

    expect(result.current.fetchStatus).toBe("idle");
    expect(mockApiFetch).not.toHaveBeenCalled();
  });
});

describe("useCreateVisitCheckpoint", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("calls POST with visit_id injected", async () => {
    mockApiFetch.mockResolvedValueOnce({ id: "new-cp", visit_id: "v1", type: "gps" });

    const wrapper = createWrapper();
    const { result } = renderHook(() => useCreateVisitCheckpoint("v1"), { wrapper });

    await result.current.mutateAsync({ type: "gps" } as any);

    expect(mockApiFetch).toHaveBeenCalledWith(
      "/api/visit-checkpoints",
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({ type: "gps", visit_id: "v1" }),
      })
    );
  });
});
