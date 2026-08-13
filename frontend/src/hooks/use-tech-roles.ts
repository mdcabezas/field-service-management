"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { TechRole } from "@/types/api";

const TECH_ROLES_KEY = ["tech-roles"];

export function useTechRoles() {
  return useQuery({
    queryKey: TECH_ROLES_KEY,
    queryFn: () => apiFetch<{ items: TechRole[]; total: number }>("/api/tech-roles").then(r => r.items),
  });
}

export function useTechRole(id: string) {
  return useQuery({
    queryKey: [...TECH_ROLES_KEY, id],
    queryFn: () => apiFetch<TechRole>(`/api/tech-roles/${id}`),
    enabled: !!id,
  });
}

export function useCreateTechRole() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<TechRole, "id" | "created_at">) =>
      apiFetch<TechRole>("/api/tech-roles", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: TECH_ROLES_KEY });
    },
  });
}

export function useUpdateTechRole() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<TechRole> }) =>
      apiFetch<TechRole>(`/api/tech-roles/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: TECH_ROLES_KEY });
    },
  });
}

export function useDeleteTechRole() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/tech-roles/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: TECH_ROLES_KEY });
    },
  });
}
