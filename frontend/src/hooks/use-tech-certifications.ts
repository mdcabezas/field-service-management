"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { TechCertification } from "@/types/api";

const TECH_CERTIFICATIONS_KEY = (technicianId: string) => ["tech-certifications", technicianId];

export function useTechCertifications(technicianId: string) {
  return useQuery({
    queryKey: TECH_CERTIFICATIONS_KEY(technicianId),
    queryFn: () => apiFetch<{ items: TechCertification[]; total: number }>(`/api/technicians/${technicianId}/certifications`).then(r => r.items),
    enabled: !!technicianId,
  });
}

export function useCreateTechCertification(technicianId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<TechCertification, "id" | "technician_id" | "created_at" | "updated_at">) =>
      apiFetch<TechCertification>("/api/tech-certifications", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, technician_id: technicianId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: TECH_CERTIFICATIONS_KEY(technicianId) });
    },
  });
}

export function useDeleteTechCertification(technicianId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/tech-certifications/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: TECH_CERTIFICATIONS_KEY(technicianId) });
    },
  });
}
