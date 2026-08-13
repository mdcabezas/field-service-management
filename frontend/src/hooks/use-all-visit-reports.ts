"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { VisitReport } from "@/types/api";

const ALL_VISIT_REPORTS_KEY = ["visit-reports"];

export function useAllVisitReports() {
  return useQuery({
    queryKey: ALL_VISIT_REPORTS_KEY,
    queryFn: () => apiFetch<{ items: VisitReport[]; total: number }>("/api/visit-reports").then(r => r.items),
  });
}

export function useVisitReport(id: string) {
  return useQuery({
    queryKey: [...ALL_VISIT_REPORTS_KEY, id],
    queryFn: () => apiFetch<VisitReport>(`/api/visit-reports/${id}`),
    enabled: !!id,
  });
}

export function useCreateVisitReportGlobal() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<VisitReport, "id" | "created_at" | "updated_at">) =>
      apiFetch<VisitReport>("/api/visit-reports", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ALL_VISIT_REPORTS_KEY });
    },
  });
}

export function useUpdateVisitReport() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<VisitReport> }) =>
      apiFetch<VisitReport>(`/api/visit-reports/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ALL_VISIT_REPORTS_KEY });
    },
  });
}

export function useDeleteVisitReportGlobal() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/visit-reports/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ALL_VISIT_REPORTS_KEY });
    },
  });
}
