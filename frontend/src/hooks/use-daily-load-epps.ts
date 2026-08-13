"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { DailyLoadEPP } from "@/types/api";

const DAILY_LOAD_EPPS_KEY = (dailyPlanId: string) => ["daily-load-epps", dailyPlanId];

export function useDailyLoadEPPs(dailyPlanId: string) {
  return useQuery({
    queryKey: DAILY_LOAD_EPPS_KEY(dailyPlanId),
    queryFn: () => apiFetch<{ items: DailyLoadEPP[]; total: number }>(`/api/daily-plans/${dailyPlanId}/epps`).then(r => r.items),
    enabled: !!dailyPlanId,
  });
}

export function useCreateDailyLoadEPP(dailyPlanId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<DailyLoadEPP, "id" | "daily_plan_id" | "created_at" | "updated_at">) =>
      apiFetch<DailyLoadEPP>("/api/daily-load-epps", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, daily_plan_id: dailyPlanId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: DAILY_LOAD_EPPS_KEY(dailyPlanId) });
    },
  });
}

export function useDeleteDailyLoadEPP(dailyPlanId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/daily-load-epps/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: DAILY_LOAD_EPPS_KEY(dailyPlanId) });
    },
  });
}
