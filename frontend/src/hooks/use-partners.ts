"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { Partner } from "@/types/api";

const PARTNERS_KEY = ["partners"];

export function usePartners() {
  return useQuery({
    queryKey: PARTNERS_KEY,
    queryFn: () => apiFetch<{ items: Partner[]; total: number }>("/api/partners").then(r => r.items),
  });
}

export function usePartner(id: string) {
  return useQuery({
    queryKey: [...PARTNERS_KEY, id],
    queryFn: () => apiFetch<Partner>(`/api/partners/${id}`),
    enabled: !!id,
  });
}

export function useCreatePartner() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<Partner, "id" | "created_at" | "updated_at">) =>
      apiFetch<Partner>("/api/partners", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PARTNERS_KEY });
    },
  });
}

export function useUpdatePartner() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Partner> }) =>
      apiFetch<Partner>(`/api/partners/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PARTNERS_KEY });
    },
  });
}

export function useDeletePartner() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/partners/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PARTNERS_KEY });
    },
  });
}
