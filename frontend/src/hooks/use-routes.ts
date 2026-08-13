"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { Route } from "@/types/api";

const ROUTES_KEY = ["routes"];

export function useRoutes() {
  return useQuery({
    queryKey: ROUTES_KEY,
    queryFn: () => apiFetch<{ items: Route[]; total: number }>("/api/routes").then(r => r.items),
  });
}

export function useRoute(id: string) {
  return useQuery({
    queryKey: [...ROUTES_KEY, id],
    queryFn: () => apiFetch<Route>(`/api/routes/${id}`),
    enabled: !!id,
  });
}

export function useCreateRoute() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<Route, "id" | "created_at" | "updated_at">) =>
      apiFetch<Route>("/api/routes", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ROUTES_KEY });
    },
  });
}

export function useUpdateRoute() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Route> }) =>
      apiFetch<Route>(`/api/routes/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ROUTES_KEY });
    },
  });
}

export function useDeleteRoute() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/routes/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ROUTES_KEY });
    },
  });
}
