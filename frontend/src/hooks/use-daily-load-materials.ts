"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { DailyLoadMaterial } from "@/types/api";

const DAILY_LOAD_MATERIALS_KEY = (dailyPlanId: string) => ["daily-load-materials", dailyPlanId];

export function useDailyLoadMaterials(dailyPlanId: string) {
  return useQuery({
    queryKey: DAILY_LOAD_MATERIALS_KEY(dailyPlanId),
    queryFn: () => apiFetch<{ items: DailyLoadMaterial[]; total: number }>(`/api/daily-plans/${dailyPlanId}/materials`).then(r => r.items),
    enabled: !!dailyPlanId,
  });
}

export function useCreateDailyLoadMaterial(dailyPlanId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<DailyLoadMaterial, "id" | "daily_plan_id" | "created_at" | "updated_at">) =>
      apiFetch<DailyLoadMaterial>("/api/daily-load-materials", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, daily_plan_id: dailyPlanId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: DAILY_LOAD_MATERIALS_KEY(dailyPlanId) });
    },
  });
}

export function useDeleteDailyLoadMaterial(dailyPlanId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/daily-load-materials/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: DAILY_LOAD_MATERIALS_KEY(dailyPlanId) });
    },
  });
}
