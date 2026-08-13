"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { SLA } from "@/types/api";

const SLAS_KEY = (partnerId: string) => ["slas", partnerId];

export function useSLAs(partnerId: string) {
  return useQuery({
    queryKey: SLAS_KEY(partnerId),
    queryFn: () => apiFetch<{ items: SLA[]; total: number }>(`/api/partners/${partnerId}/slas`).then(r => r.items),
    enabled: !!partnerId,
  });
}

export function useCreateSLA(partnerId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<SLA, "id" | "partner_id" | "created_at" | "updated_at">) =>
      apiFetch<SLA>("/api/slas", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, partner_id: partnerId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: SLAS_KEY(partnerId) });
    },
  });
}

export function useDeleteSLA(partnerId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/slas/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: SLAS_KEY(partnerId) });
    },
  });
}
