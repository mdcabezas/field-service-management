"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { VisitPhoto } from "@/types/api";

const VISIT_PHOTOS_KEY = (visitId: string) => ["visit-photos", visitId];

export function useVisitPhotos(visitId: string) {
  return useQuery({
    queryKey: VISIT_PHOTOS_KEY(visitId),
    queryFn: () => apiFetch<VisitPhoto[]>(`/api/visits/${visitId}/photos`),
    enabled: !!visitId,
  });
}

export function useCreateVisitPhoto(visitId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<VisitPhoto, "id" | "visit_id" | "created_at" | "updated_at">) =>
      apiFetch<VisitPhoto>("/api/visit-photos", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, visit_id: visitId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISIT_PHOTOS_KEY(visitId) });
    },
  });
}

export function useDeleteVisitPhoto(visitId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/visit-photos/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISIT_PHOTOS_KEY(visitId) });
    },
  });
}
