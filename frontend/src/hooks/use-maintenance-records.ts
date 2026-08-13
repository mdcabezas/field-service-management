"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { MaintenanceRecord } from "@/types/api";

const MAINTENANCE_RECORDS_KEY = ["maintenance-records"];

export function useMaintenanceRecords() {
  return useQuery({
    queryKey: MAINTENANCE_RECORDS_KEY,
    queryFn: () => apiFetch<{ items: MaintenanceRecord[]; total: number }>("/api/maintenance-records").then(r => r.items),
  });
}

export function useMaintenanceRecord(id: string) {
  return useQuery({
    queryKey: [...MAINTENANCE_RECORDS_KEY, id],
    queryFn: () => apiFetch<MaintenanceRecord>(`/api/maintenance-records/${id}`),
    enabled: !!id,
  });
}

export function useUpdateMaintenanceRecord() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<MaintenanceRecord> }) =>
      apiFetch<MaintenanceRecord>(`/api/maintenance-records/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: MAINTENANCE_RECORDS_KEY });
    },
  });
}

export function useCreateMaintenanceRecord() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<MaintenanceRecord, "id" | "created_at" | "updated_at">) =>
      apiFetch<MaintenanceRecord>("/api/maintenance-records", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: MAINTENANCE_RECORDS_KEY });
    },
  });
}

export function useDeleteMaintenanceRecord() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/maintenance-records/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: MAINTENANCE_RECORDS_KEY });
    },
  });
}
