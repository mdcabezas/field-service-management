"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { VisitSLATracking } from "@/types/api";

const ALL_VISIT_SLA_TRACKINGS_KEY = ["visit-sla-trackings"];

export function useAllVisitSLATrackings() {
  return useQuery({
    queryKey: ALL_VISIT_SLA_TRACKINGS_KEY,
    queryFn: () => apiFetch<{ items: VisitSLATracking[]; total: number }>("/api/visit-sla-trackings").then(r => r.items),
  });
}

export function useVisitSLATracking(id: string) {
  return useQuery({
    queryKey: [...ALL_VISIT_SLA_TRACKINGS_KEY, id],
    queryFn: () => apiFetch<VisitSLATracking>(`/api/visit-sla-trackings/${id}`),
    enabled: !!id,
  });
}

export function useCreateVisitSLATrackingGlobal() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<VisitSLATracking, "id" | "created_at" | "updated_at">) =>
      apiFetch<VisitSLATracking>("/api/visit-sla-trackings", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ALL_VISIT_SLA_TRACKINGS_KEY });
    },
  });
}

export function useUpdateVisitSLATracking() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<VisitSLATracking> }) =>
      apiFetch<VisitSLATracking>(`/api/visit-sla-trackings/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ALL_VISIT_SLA_TRACKINGS_KEY });
    },
  });
}

export function useDeleteVisitSLATrackingGlobal() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/visit-sla-trackings/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ALL_VISIT_SLA_TRACKINGS_KEY });
    },
  });
}
