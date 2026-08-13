"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { Rental } from "@/types/api";

const RENTALS_KEY = ["rentals"];

export function useRentals() {
  return useQuery({
    queryKey: RENTALS_KEY,
    queryFn: () => apiFetch<{ items: Rental[]; total: number }>("/api/rentals").then(r => r.items),
  });
}

export function useRental(id: string) {
  return useQuery({
    queryKey: [...RENTALS_KEY, id],
    queryFn: () => apiFetch<Rental>(`/api/rentals/${id}`),
    enabled: !!id,
  });
}

export function useCreateRental() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<Rental, "id" | "created_at" | "updated_at">) =>
      apiFetch<Rental>("/api/rentals", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: RENTALS_KEY });
    },
  });
}

export function useUpdateRental() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Rental> }) =>
      apiFetch<Rental>(`/api/rentals/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: RENTALS_KEY });
    },
  });
}

export function useDeleteRental() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/rentals/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: RENTALS_KEY });
    },
  });
}
