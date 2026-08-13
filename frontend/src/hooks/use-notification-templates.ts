"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { NotificationTemplate } from "@/types/api";

const NOTIFICATION_TEMPLATES_KEY = ["notification-templates"];

export function useNotificationTemplates() {
  return useQuery({
    queryKey: NOTIFICATION_TEMPLATES_KEY,
    queryFn: () => apiFetch<{ items: NotificationTemplate[]; total: number }>("/api/notification-templates").then(r => r.items),
  });
}

export function useNotificationTemplate(id: string) {
  return useQuery({
    queryKey: [...NOTIFICATION_TEMPLATES_KEY, id],
    queryFn: () => apiFetch<NotificationTemplate>(`/api/notification-templates/${id}`),
    enabled: !!id,
  });
}

export function useCreateNotificationTemplate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<NotificationTemplate, "id" | "created_at" | "updated_at">) =>
      apiFetch<NotificationTemplate>("/api/notification-templates", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: NOTIFICATION_TEMPLATES_KEY });
    },
  });
}

export function useUpdateNotificationTemplate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<NotificationTemplate> }) =>
      apiFetch<NotificationTemplate>(`/api/notification-templates/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: NOTIFICATION_TEMPLATES_KEY });
    },
  });
}

export function useDeleteNotificationTemplate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/notification-templates/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: NOTIFICATION_TEMPLATES_KEY });
    },
  });
}
