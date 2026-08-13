"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { DailyPlanAssignment } from "@/types/api";

const ALL_ASSIGNMENTS_KEY = ["daily-plan-assignments"];
const ASSIGNMENTS_BY_PLAN_KEY = (dailyPlanId: string) => ["daily-plan-assignments", "plan", dailyPlanId];

export function useAllDailyPlanAssignments() {
  return useQuery({
    queryKey: ALL_ASSIGNMENTS_KEY,
    queryFn: () => apiFetch<{ items: DailyPlanAssignment[]; total: number }>("/api/daily-plan-assignments").then(r => r.items),
  });
}

export function useDailyPlanAssignments(dailyPlanId: string) {
  return useQuery({
    queryKey: ASSIGNMENTS_BY_PLAN_KEY(dailyPlanId),
    queryFn: () => apiFetch<{ items: DailyPlanAssignment[]; total: number }>(`/api/daily-plans/${dailyPlanId}/assignments`).then(r => r.items),
    enabled: !!dailyPlanId,
  });
}

export function useCreateDailyPlanAssignment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<DailyPlanAssignment, "id" | "created_at" | "updated_at">) =>
      apiFetch<DailyPlanAssignment>("/api/daily-plan-assignments", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ALL_ASSIGNMENTS_KEY });
    },
  });
}

export function useUpdateDailyPlanAssignment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<DailyPlanAssignment> }) =>
      apiFetch<DailyPlanAssignment>(`/api/daily-plan-assignments/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ALL_ASSIGNMENTS_KEY });
    },
  });
}

export function useDeleteDailyPlanAssignment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/daily-plan-assignments/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ALL_ASSIGNMENTS_KEY });
    },
  });
}
