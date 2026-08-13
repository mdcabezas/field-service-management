"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { VisitReport } from "@/types/api";

const VISIT_REPORTS_KEY = (visitId: string) => ["visit-reports", visitId];

export function useVisitReports(visitId: string) {
  return useQuery({
    queryKey: VISIT_REPORTS_KEY(visitId),
    queryFn: () => apiFetch<VisitReport[]>(`/api/visits/${visitId}/reports`),
    enabled: !!visitId,
  });
}

export function useCreateVisitReport(visitId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<VisitReport, "id" | "visit_id" | "created_at" | "updated_at">) =>
      apiFetch<VisitReport>("/api/visit-reports", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, visit_id: visitId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISIT_REPORTS_KEY(visitId) });
    },
  });
}

export function useDeleteVisitReport(visitId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/visit-reports/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISIT_REPORTS_KEY(visitId) });
    },
  });
}
