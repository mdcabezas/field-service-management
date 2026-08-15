"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { VisitType } from "@/types/api";

const VISIT_TYPES_KEY = ["visit-types"];

export function useVisitTypes() {
  return useQuery({
    queryKey: VISIT_TYPES_KEY,
    queryFn: () => apiFetch<{ items: VisitType[]; total: number }>("/api/visit-types").then(r => r.items),
  });
}

export function useVisitType(id: string) {
  return useQuery({
    queryKey: [...VISIT_TYPES_KEY, id],
    queryFn: () => apiFetch<VisitType>(`/api/visit-types/${id}`),
    enabled: !!id,
  });
}

export function useCreateVisitType() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<VisitType, "id" | "created_at" | "updated_at">) =>
      apiFetch<VisitType>("/api/visit-types", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISIT_TYPES_KEY });
    },
  });
}

export function useUpdateVisitType() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<VisitType> }) =>
      apiFetch<VisitType>(`/api/visit-types/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISIT_TYPES_KEY });
    },
  });
}

export function useDeleteVisitType() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/visit-types/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISIT_TYPES_KEY });
    },
  });
}