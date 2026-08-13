"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { Visit } from "@/types/api";

const VISITS_KEY = ["visits"];

export function useVisits() {
  return useQuery({
    queryKey: VISITS_KEY,
    queryFn: () => apiFetch<{ items: Visit[]; total: number }>("/api/visits").then(r => r.items),
  });
}

export function useVisit(id: string) {
  return useQuery({
    queryKey: [...VISITS_KEY, id],
    queryFn: () => apiFetch<Visit>(`/api/visits/${id}`),
    enabled: !!id,
  });
}

export function useCreateVisit() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<Visit, "id" | "created_at" | "updated_at">) =>
      apiFetch<Visit>("/api/visits", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISITS_KEY });
    },
  });
}

export function useUpdateVisit() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Visit> }) =>
      apiFetch<Visit>(`/api/visits/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISITS_KEY });
    },
  });
}

export function useDeleteVisit() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/visits/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISITS_KEY });
    },
  });
}
