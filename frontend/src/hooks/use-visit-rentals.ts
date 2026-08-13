"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { VisitRental } from "@/types/api";

const VISIT_RENTALS_KEY = (visitId: string) => ["visit-rentals", visitId];

export function useVisitRentals(visitId: string) {
  return useQuery({
    queryKey: VISIT_RENTALS_KEY(visitId),
    queryFn: () => apiFetch<{ items: VisitRental[]; total: number }>(`/api/visits/${visitId}/rentals`).then(r => r.items),
    enabled: !!visitId,
  });
}

export function useCreateVisitRental(visitId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<VisitRental, "id" | "visit_id" | "created_at" | "updated_at">) =>
      apiFetch<VisitRental>("/api/visit-rentals", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, visit_id: visitId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISIT_RENTALS_KEY(visitId) });
    },
  });
}

export function useDeleteVisitRental(visitId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/visit-rentals/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISIT_RENTALS_KEY(visitId) });
    },
  });
}
