"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { MaintenanceSchedule } from "@/types/api";

const MAINTENANCE_SCHEDULES_KEY = ["maintenance-schedules"];

export function useMaintenanceSchedules() {
  return useQuery({
    queryKey: MAINTENANCE_SCHEDULES_KEY,
    queryFn: () => apiFetch<{ items: MaintenanceSchedule[]; total: number }>("/api/maintenance-schedules").then(r => r.items),
  });
}

export function useMaintenanceSchedule(id: string) {
  return useQuery({
    queryKey: [...MAINTENANCE_SCHEDULES_KEY, id],
    queryFn: () => apiFetch<MaintenanceSchedule>(`/api/maintenance-schedules/${id}`),
    enabled: !!id,
  });
}

export function useCreateMaintenanceSchedule() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<MaintenanceSchedule, "id" | "created_at" | "updated_at">) =>
      apiFetch<MaintenanceSchedule>("/api/maintenance-schedules", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: MAINTENANCE_SCHEDULES_KEY });
    },
  });
}

export function useUpdateMaintenanceSchedule() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<MaintenanceSchedule> }) =>
      apiFetch<MaintenanceSchedule>(`/api/maintenance-schedules/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: MAINTENANCE_SCHEDULES_KEY });
    },
  });
}

export function useDeleteMaintenanceSchedule() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/maintenance-schedules/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: MAINTENANCE_SCHEDULES_KEY });
    },
  });
}
