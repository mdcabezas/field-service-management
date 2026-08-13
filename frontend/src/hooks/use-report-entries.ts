"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { ReportEntry } from "@/types/api";

const REPORT_ENTRIES_KEY = (reportId: string) => ["report-entries", reportId];

export function useReportEntries(reportId: string) {
  return useQuery({
    queryKey: REPORT_ENTRIES_KEY(reportId),
    queryFn: () => apiFetch<{ items: ReportEntry[]; total: number }>(`/api/visit-reports/${reportId}/entries`).then(r => r.items),
    enabled: !!reportId,
  });
}

export function useCreateReportEntry(reportId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<ReportEntry, "id" | "report_id" | "created_at" | "updated_at">) =>
      apiFetch<ReportEntry>("/api/report-entries", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, report_id: reportId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: REPORT_ENTRIES_KEY(reportId) });
    },
  });
}

export function useDeleteReportEntry(reportId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/report-entries/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: REPORT_ENTRIES_KEY(reportId) });
    },
  });
}
