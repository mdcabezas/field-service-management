"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { VisitMeasurement } from "@/types/api";

const VISIT_MEASUREMENTS_KEY = (visitId: string) => ["visit-measurements", visitId];

export function useVisitMeasurements(visitId: string) {
  return useQuery({
    queryKey: VISIT_MEASUREMENTS_KEY(visitId),
    queryFn: () => apiFetch<VisitMeasurement[]>(`/api/visits/${visitId}/measurements`),
    enabled: !!visitId,
  });
}

export function useCreateVisitMeasurement(visitId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<VisitMeasurement, "id" | "visit_id" | "created_at" | "updated_at">) =>
      apiFetch<VisitMeasurement>("/api/visit-measurements", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, visit_id: visitId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISIT_MEASUREMENTS_KEY(visitId) });
    },
  });
}

export function useDeleteVisitMeasurement(visitId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/visit-measurements/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISIT_MEASUREMENTS_KEY(visitId) });
    },
  });
}
