"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { EPPItem } from "@/types/api";

const EPP_ITEMS_KEY = ["epp-items"];

export function useEPPItems() {
  return useQuery({
    queryKey: EPP_ITEMS_KEY,
    queryFn: () => apiFetch<{ items: EPPItem[]; total: number }>("/api/epp-items").then(r => r.items),
  });
}

export function useEPPItem(id: string) {
  return useQuery({
    queryKey: [...EPP_ITEMS_KEY, id],
    queryFn: () => apiFetch<EPPItem>(`/api/epp-items/${id}`),
    enabled: !!id,
  });
}

export function useCreateEPPItem() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<EPPItem, "id" | "created_at" | "updated_at">) =>
      apiFetch<EPPItem>("/api/epp-items", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: EPP_ITEMS_KEY });
    },
  });
}

export function useUpdateEPPItem() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<EPPItem> }) =>
      apiFetch<EPPItem>(`/api/epp-items/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: EPP_ITEMS_KEY });
    },
  });
}

export function useDeleteEPPItem() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/epp-items/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: EPP_ITEMS_KEY });
    },
  });
}
