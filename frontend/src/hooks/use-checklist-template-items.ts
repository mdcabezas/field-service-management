"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { ChecklistTemplateMaterial, ChecklistTemplateTool, ChecklistTemplateEPP } from "@/types/api";

// Materials
const CHECKLIST_TEMPLATE_MATERIALS_KEY = (templateId: string) => ["checklist-template-materials", templateId];

export function useChecklistTemplateMaterials(templateId: string) {
  return useQuery({
    queryKey: CHECKLIST_TEMPLATE_MATERIALS_KEY(templateId),
    queryFn: () => apiFetch<{ items: ChecklistTemplateMaterial[]; total: number }>(`/api/checklist-templates/${templateId}/materials`).then(r => r.items),
    enabled: !!templateId,
  });
}

export function useCreateChecklistTemplateMaterial(templateId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<ChecklistTemplateMaterial, "id" | "template_id" | "created_at" | "updated_at">) =>
      apiFetch<ChecklistTemplateMaterial>("/api/checklist-template-materials", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, template_id: templateId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: CHECKLIST_TEMPLATE_MATERIALS_KEY(templateId) });
    },
  });
}

export function useDeleteChecklistTemplateMaterial(templateId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/checklist-template-materials/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: CHECKLIST_TEMPLATE_MATERIALS_KEY(templateId) });
    },
  });
}

// Tools
const CHECKLIST_TEMPLATE_TOOLS_KEY = (templateId: string) => ["checklist-template-tools", templateId];

export function useChecklistTemplateTools(templateId: string) {
  return useQuery({
    queryKey: CHECKLIST_TEMPLATE_TOOLS_KEY(templateId),
    queryFn: () => apiFetch<{ items: ChecklistTemplateTool[]; total: number }>(`/api/checklist-templates/${templateId}/tools`).then(r => r.items),
    enabled: !!templateId,
  });
}

export function useCreateChecklistTemplateTool(templateId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<ChecklistTemplateTool, "id" | "template_id" | "created_at" | "updated_at">) =>
      apiFetch<ChecklistTemplateTool>("/api/checklist-template-tools", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, template_id: templateId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: CHECKLIST_TEMPLATE_TOOLS_KEY(templateId) });
    },
  });
}

export function useDeleteChecklistTemplateTool(templateId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/checklist-template-tools/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: CHECKLIST_TEMPLATE_TOOLS_KEY(templateId) });
    },
  });
}

// EPPs
const CHECKLIST_TEMPLATE_EPPS_KEY = (templateId: string) => ["checklist-template-epps", templateId];

export function useChecklistTemplateEPPs(templateId: string) {
  return useQuery({
    queryKey: CHECKLIST_TEMPLATE_EPPS_KEY(templateId),
    queryFn: () => apiFetch<{ items: ChecklistTemplateEPP[]; total: number }>(`/api/checklist-templates/${templateId}/epps`).then(r => r.items),
    enabled: !!templateId,
  });
}

export function useCreateChecklistTemplateEPP(templateId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<ChecklistTemplateEPP, "id" | "template_id" | "created_at" | "updated_at">) =>
      apiFetch<ChecklistTemplateEPP>("/api/checklist-template-epps", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, template_id: templateId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: CHECKLIST_TEMPLATE_EPPS_KEY(templateId) });
    },
  });
}

export function useDeleteChecklistTemplateEPP(templateId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/checklist-template-epps/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: CHECKLIST_TEMPLATE_EPPS_KEY(templateId) });
    },
  });
}
