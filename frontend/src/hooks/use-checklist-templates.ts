"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { ChecklistTemplate } from "@/types/api";

const CHECKLIST_TEMPLATES_KEY = ["checklist-templates"];

export function useChecklistTemplates() {
  return useQuery({
    queryKey: CHECKLIST_TEMPLATES_KEY,
    queryFn: () => apiFetch<{ items: ChecklistTemplate[]; total: number }>("/api/checklist-templates").then(r => r.items),
  });
}

export function useChecklistTemplate(id: string) {
  return useQuery({
    queryKey: [...CHECKLIST_TEMPLATES_KEY, id],
    queryFn: () => apiFetch<ChecklistTemplate>(`/api/checklist-templates/${id}`),
    enabled: !!id,
  });
}

export function useCreateChecklistTemplate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<ChecklistTemplate, "id" | "created_at" | "updated_at">) =>
      apiFetch<ChecklistTemplate>("/api/checklist-templates", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: CHECKLIST_TEMPLATES_KEY });
    },
  });
}

export function useUpdateChecklistTemplate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<ChecklistTemplate> }) =>
      apiFetch<ChecklistTemplate>(`/api/checklist-templates/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: CHECKLIST_TEMPLATES_KEY });
    },
  });
}

export function useDeleteChecklistTemplate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/checklist-templates/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: CHECKLIST_TEMPLATES_KEY });
    },
  });
}
