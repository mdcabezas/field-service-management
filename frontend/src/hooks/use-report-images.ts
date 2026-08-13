"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { ReportImage } from "@/types/api";

const REPORT_IMAGES_KEY = (reportId: string) => ["report-images", reportId];

export function useReportImages(reportId: string) {
  return useQuery({
    queryKey: REPORT_IMAGES_KEY(reportId),
    queryFn: () => apiFetch<{ items: ReportImage[]; total: number }>(`/api/visit-reports/${reportId}/images`).then(r => r.items),
    enabled: !!reportId,
  });
}

export function useCreateReportImage(reportId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<ReportImage, "id" | "report_id" | "created_at" | "updated_at">) =>
      apiFetch<ReportImage>("/api/report-images", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, report_id: reportId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: REPORT_IMAGES_KEY(reportId) });
    },
  });
}

export function useDeleteReportImage(reportId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/report-images/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: REPORT_IMAGES_KEY(reportId) });
    },
  });
}
