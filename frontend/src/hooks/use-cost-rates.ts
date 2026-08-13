"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { CostRate } from "@/types/api";

const COST_RATES_KEY = ["cost-rates"];

export function useCostRates() {
  return useQuery({
    queryKey: COST_RATES_KEY,
    queryFn: () => apiFetch<{ items: CostRate[]; total: number }>("/api/cost-rates").then(r => r.items),
  });
}

export function useCostRate(id: string) {
  return useQuery({
    queryKey: [...COST_RATES_KEY, id],
    queryFn: () => apiFetch<CostRate>(`/api/cost-rates/${id}`),
    enabled: !!id,
  });
}

export function useCreateCostRate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<CostRate, "id" | "created_at" | "updated_at">) =>
      apiFetch<CostRate>("/api/cost-rates", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: COST_RATES_KEY });
    },
  });
}

export function useUpdateCostRate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<CostRate> }) =>
      apiFetch<CostRate>(`/api/cost-rates/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: COST_RATES_KEY });
    },
  });
}

export function useDeleteCostRate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/cost-rates/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: COST_RATES_KEY });
    },
  });
}
