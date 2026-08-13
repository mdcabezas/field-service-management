"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { VisitCheckpoint } from "@/types/api";

const VISIT_CHECKPOINTS_KEY = (visitId: string) => ["visit-checkpoints", visitId];

export function useVisitCheckpoints(visitId: string) {
  return useQuery({
    queryKey: VISIT_CHECKPOINTS_KEY(visitId),
    queryFn: () => apiFetch<VisitCheckpoint[]>(`/api/visits/${visitId}/checkpoints`),
    enabled: !!visitId,
  });
}

export function useCreateVisitCheckpoint(visitId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<VisitCheckpoint, "id" | "visit_id" | "created_at" | "updated_at">) =>
      apiFetch<VisitCheckpoint>("/api/visit-checkpoints", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, visit_id: visitId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISIT_CHECKPOINTS_KEY(visitId) });
    },
  });
}

export function useDeleteVisitCheckpoint(visitId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/visit-checkpoints/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISIT_CHECKPOINTS_KEY(visitId) });
    },
  });
}
