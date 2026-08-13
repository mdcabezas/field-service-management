"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { Technician } from "@/types/api";

const TECHNICIANS_KEY = (customerId: string) => ["technicians", customerId];
const ALL_TECHNICIANS_KEY = ["technicians", "all"];

export function useAllTechnicians() {
  return useQuery({
    queryKey: ALL_TECHNICIANS_KEY,
    queryFn: () => apiFetch<{ items: Technician[]; total: number }>("/api/technicians").then(r => r.items),
  });
}

export function useTechnicians(customerId: string) {
  return useQuery({
    queryKey: TECHNICIANS_KEY(customerId),
    queryFn: () => apiFetch<{ items: Technician[]; total: number }>(`/api/customers/${customerId}/technicians`).then(r => r.items),
    enabled: !!customerId,
  });
}

export function useCreateTechnician(customerId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<Technician, "id" | "customer_id" | "created_at" | "updated_at">) =>
      apiFetch<Technician>("/api/technicians", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, customer_id: customerId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: TECHNICIANS_KEY(customerId) });
    },
  });
}

export function useDeleteTechnician(customerId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/technicians/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: TECHNICIANS_KEY(customerId) });
    },
  });
}
