"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { VisitChecklistMaterial, VisitChecklistTool, VisitChecklistEPP } from "@/types/api";

// Materials
const VISIT_CHECKLIST_MATERIALS_KEY = (visitId: string) => ["visit-checklist-materials", visitId];

export function useVisitChecklistMaterials(visitId: string) {
  return useQuery({
    queryKey: VISIT_CHECKLIST_MATERIALS_KEY(visitId),
    queryFn: () => apiFetch<VisitChecklistMaterial[]>(`/api/visits/${visitId}/checklist-materials`),
    enabled: !!visitId,
  });
}

export function useCreateVisitChecklistMaterial(visitId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<VisitChecklistMaterial, "id" | "visit_id" | "created_at" | "updated_at">) =>
      apiFetch<VisitChecklistMaterial>("/api/visit-checklist-materials", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, visit_id: visitId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISIT_CHECKLIST_MATERIALS_KEY(visitId) });
    },
  });
}

export function useDeleteVisitChecklistMaterial(visitId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/visit-checklist-materials/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISIT_CHECKLIST_MATERIALS_KEY(visitId) });
    },
  });
}

// Tools
const VISIT_CHECKLIST_TOOLS_KEY = (visitId: string) => ["visit-checklist-tools", visitId];

export function useVisitChecklistTools(visitId: string) {
  return useQuery({
    queryKey: VISIT_CHECKLIST_TOOLS_KEY(visitId),
    queryFn: () => apiFetch<VisitChecklistTool[]>(`/api/visits/${visitId}/checklist-tools`),
    enabled: !!visitId,
  });
}

export function useCreateVisitChecklistTool(visitId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<VisitChecklistTool, "id" | "visit_id" | "created_at" | "updated_at">) =>
      apiFetch<VisitChecklistTool>("/api/visit-checklist-tools", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, visit_id: visitId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISIT_CHECKLIST_TOOLS_KEY(visitId) });
    },
  });
}

export function useDeleteVisitChecklistTool(visitId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/visit-checklist-tools/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISIT_CHECKLIST_TOOLS_KEY(visitId) });
    },
  });
}

// EPPs
const VISIT_CHECKLIST_EPPS_KEY = (visitId: string) => ["visit-checklist-epps", visitId];

export function useVisitChecklistEPPs(visitId: string) {
  return useQuery({
    queryKey: VISIT_CHECKLIST_EPPS_KEY(visitId),
    queryFn: () => apiFetch<VisitChecklistEPP[]>(`/api/visits/${visitId}/checklist-epps`),
    enabled: !!visitId,
  });
}

export function useCreateVisitChecklistEPP(visitId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<VisitChecklistEPP, "id" | "visit_id" | "created_at" | "updated_at">) =>
      apiFetch<VisitChecklistEPP>("/api/visit-checklist-epps", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, visit_id: visitId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISIT_CHECKLIST_EPPS_KEY(visitId) });
    },
  });
}

export function useDeleteVisitChecklistEPP(visitId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/visit-checklist-epps/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VISIT_CHECKLIST_EPPS_KEY(visitId) });
    },
  });
}
