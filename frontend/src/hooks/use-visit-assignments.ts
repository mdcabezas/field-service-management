"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { VisitAssignment } from "@/types/api";

const VISIT_ASSIGNMENTS_KEY = (visitId: string) => ["visit-assignments", visitId];

export function useVisitAssignments(visitId: string) {
  return useQuery({
    queryKey: VISIT_ASSIGNMENTS_KEY(visitId),
    queryFn: () => apiFetch<VisitAssignment[]>(`/api/visits/${visitId}/assignments`),
    enabled: !!visitId,
  });
}

export function useCreateVisitAssignment(visitId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<VisitAssignment, "id" | "visit_id" | "created_at" | "updated_at">) =>
      apiFetch<VisitAssignment>("/api/visit-assignments", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, visit_id: visitId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISIT_ASSIGNMENTS_KEY(visitId) });
    },
  });
}

export function useDeleteVisitAssignment(visitId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/visit-assignments/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISIT_ASSIGNMENTS_KEY(visitId) });
    },
  });
}
