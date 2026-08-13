"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { VisitSLATracking } from "@/types/api";

const VISIT_SLA_TRACKINGS_KEY = (visitId: string) => ["visit-sla-trackings", visitId];

export function useVisitSLATrackings(visitId: string) {
  return useQuery({
    queryKey: VISIT_SLA_TRACKINGS_KEY(visitId),
    queryFn: () => apiFetch<VisitSLATracking[]>(`/api/visits/${visitId}/sla-trackings`),
    enabled: !!visitId,
  });
}

export function useCreateVisitSLATracking(visitId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<VisitSLATracking, "id" | "visit_id" | "created_at" | "updated_at">) =>
      apiFetch<VisitSLATracking>("/api/visit-sla-trackings", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, visit_id: visitId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISIT_SLA_TRACKINGS_KEY(visitId) });
    },
  });
}

export function useDeleteVisitSLATracking(visitId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/visit-sla-trackings/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISIT_SLA_TRACKINGS_KEY(visitId) });
    },
  });
}
