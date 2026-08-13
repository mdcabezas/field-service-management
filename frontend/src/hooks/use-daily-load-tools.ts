"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { DailyLoadTool } from "@/types/api";

const DAILY_LOAD_TOOLS_KEY = (dailyPlanId: string) => ["daily-load-tools", dailyPlanId];

export function useDailyLoadTools(dailyPlanId: string) {
  return useQuery({
    queryKey: DAILY_LOAD_TOOLS_KEY(dailyPlanId),
    queryFn: () => apiFetch<{ items: DailyLoadTool[]; total: number }>(`/api/daily-plans/${dailyPlanId}/tools`).then(r => r.items),
    enabled: !!dailyPlanId,
  });
}

export function useCreateDailyLoadTool(dailyPlanId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<DailyLoadTool, "id" | "daily_plan_id" | "created_at" | "updated_at">) =>
      apiFetch<DailyLoadTool>("/api/daily-load-tools", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, daily_plan_id: dailyPlanId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: DAILY_LOAD_TOOLS_KEY(dailyPlanId) });
    },
  });
}

export function useDeleteDailyLoadTool(dailyPlanId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/daily-load-tools/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: DAILY_LOAD_TOOLS_KEY(dailyPlanId) });
    },
  });
}
