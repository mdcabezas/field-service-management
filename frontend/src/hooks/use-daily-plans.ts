"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { DailyPlan } from "@/types/api";

const DAILY_PLANS_KEY = ["daily-plans"];

export function useDailyPlans() {
  return useQuery({
    queryKey: DAILY_PLANS_KEY,
    queryFn: () => apiFetch<{ items: DailyPlan[]; total: number }>("/api/daily-plans").then(r => r.items),
  });
}

export function useDailyPlan(id: string) {
  return useQuery({
    queryKey: [...DAILY_PLANS_KEY, id],
    queryFn: () => apiFetch<DailyPlan>(`/api/daily-plans/${id}`),
    enabled: !!id,
  });
}

export function useCreateDailyPlan() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<DailyPlan, "id" | "created_at" | "updated_at">) =>
      apiFetch<DailyPlan>("/api/daily-plans", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: DAILY_PLANS_KEY });
    },
  });
}

export function useUpdateDailyPlan() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<DailyPlan> }) =>
      apiFetch<DailyPlan>(`/api/daily-plans/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: DAILY_PLANS_KEY });
    },
  });
}

export function useDeleteDailyPlan() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/daily-plans/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: DAILY_PLANS_KEY });
    },
  });
}
