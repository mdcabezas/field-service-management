"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { Tool } from "@/types/api";

const TOOLS_KEY = ["tools"];

export function useTools() {
  return useQuery({
    queryKey: TOOLS_KEY,
    queryFn: () => apiFetch<{ items: Tool[]; total: number }>("/api/tools").then(r => r.items),
  });
}

export function useTool(id: string) {
  return useQuery({
    queryKey: [...TOOLS_KEY, id],
    queryFn: () => apiFetch<Tool>(`/api/tools/${id}`),
    enabled: !!id,
  });
}

export function useCreateTool() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<Tool, "id" | "created_at" | "updated_at">) =>
      apiFetch<Tool>("/api/tools", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: TOOLS_KEY });
    },
  });
}

export function useUpdateTool() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Tool> }) =>
      apiFetch<Tool>(`/api/tools/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: TOOLS_KEY });
    },
  });
}

export function useDeleteTool() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/tools/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: TOOLS_KEY });
    },
  });
}
