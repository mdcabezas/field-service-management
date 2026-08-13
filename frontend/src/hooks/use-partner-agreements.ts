"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { PartnerAgreement } from "@/types/api";

const PARTNER_AGREEMENTS_KEY = (partnerId: string) => ["partner-agreements", partnerId];

export function usePartnerAgreements(partnerId: string) {
  return useQuery({
    queryKey: PARTNER_AGREEMENTS_KEY(partnerId),
    queryFn: () => apiFetch<{ items: PartnerAgreement[]; total: number }>(`/api/partners/${partnerId}/agreements`).then(r => r.items),
    enabled: !!partnerId,
  });
}

export function useCreatePartnerAgreement(partnerId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<PartnerAgreement, "id" | "partner_id" | "created_at" | "updated_at">) =>
      apiFetch<PartnerAgreement>("/api/partner-agreements", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, partner_id: partnerId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PARTNER_AGREEMENTS_KEY(partnerId) });
    },
  });
}

export function useDeletePartnerAgreement(partnerId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/partner-agreements/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PARTNER_AGREEMENTS_KEY(partnerId) });
    },
  });
}
