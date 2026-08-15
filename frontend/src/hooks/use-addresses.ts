"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { Address } from "@/types/api";

const ADDRESSES_KEY = ["addresses"];

export function useAddresses(params?: { limit?: number; offset?: number }) {
  const searchParams = new URLSearchParams();
  if (params?.limit) searchParams.set("limit", params.limit.toString());
  if (params?.offset) searchParams.set("offset", params.offset.toString());
  
  const queryString = searchParams.toString();
  const url = `/api/addresses${queryString ? `?${queryString}` : ""}`;

  return useQuery({
    queryKey: [...ADDRESSES_KEY, params],
    queryFn: () => apiFetch<{ items: Address[]; total: number }>(url).then(r => r.items),
  });
}

export function useAddress(id: string) {
  return useQuery({
    queryKey: [...ADDRESSES_KEY, id],
    queryFn: () => apiFetch<Address>(`/api/addresses/${id}`),
    enabled: !!id,
  });
}

export function useCreateAddress() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<Address, "id" | "created_at">) =>
      apiFetch<Address>("/api/addresses", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ADDRESSES_KEY });
    },
  });
}

export function useUpdateAddress() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Address> }) =>
      apiFetch<Address>(`/api/addresses/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ADDRESSES_KEY });
    },
  });
}

export function useDeleteAddress() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/addresses/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ADDRESSES_KEY });
    },
  });
}