"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { ReportTemplate } from "@/types/api";

const REPORT_TEMPLATES_KEY = ["report-templates"];

export function useReportTemplates() {
  return useQuery({
    queryKey: REPORT_TEMPLATES_KEY,
    queryFn: () => apiFetch<{ items: ReportTemplate[]; total: number }>("/api/report-templates").then(r => r.items),
  });
}

export function useReportTemplate(id: string) {
  return useQuery({
    queryKey: [...REPORT_TEMPLATES_KEY, id],
    queryFn: () => apiFetch<ReportTemplate>(`/api/report-templates/${id}`),
    enabled: !!id,
  });
}

export function useCreateReportTemplate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<ReportTemplate, "id" | "created_at" | "updated_at">) =>
      apiFetch<ReportTemplate>("/api/report-templates", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: REPORT_TEMPLATES_KEY });
    },
  });
}

export function useUpdateReportTemplate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<ReportTemplate> }) =>
      apiFetch<ReportTemplate>(`/api/report-templates/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: REPORT_TEMPLATES_KEY });
    },
  });
}

export function useDeleteReportTemplate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/report-templates/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: REPORT_TEMPLATES_KEY });
    },
  });
}
