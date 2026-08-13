"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { Material } from "@/types/api";

const MATERIALS_KEY = ["materials"];

export function useMaterials() {
  return useQuery({
    queryKey: MATERIALS_KEY,
    queryFn: () => apiFetch<{ items: Material[]; total: number }>("/api/materials").then(r => r.items),
  });
}

export function useMaterial(id: string) {
  return useQuery({
    queryKey: [...MATERIALS_KEY, id],
    queryFn: () => apiFetch<Material>(`/api/materials/${id}`),
    enabled: !!id,
  });
}

export function useCreateMaterial() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<Material, "id" | "created_at" | "updated_at">) =>
      apiFetch<Material>("/api/materials", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: MATERIALS_KEY });
    },
  });
}

export function useUpdateMaterial() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Material> }) =>
      apiFetch<Material>(`/api/materials/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: MATERIALS_KEY });
    },
  });
}

export function useDeleteMaterial() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/materials/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: MATERIALS_KEY });
    },
  });
}
