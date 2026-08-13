"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { PartnerAgreementDoc } from "@/types/api";

const PARTNER_AGREEMENT_DOCS_KEY = (agreementId: string) => ["partner-agreement-docs", agreementId];

export function usePartnerAgreementDocs(agreementId: string) {
  return useQuery({
    queryKey: PARTNER_AGREEMENT_DOCS_KEY(agreementId),
    queryFn: () => apiFetch<{ items: PartnerAgreementDoc[]; total: number }>(`/api/partner-agreements/${agreementId}/docs`).then(r => r.items),
    enabled: !!agreementId,
  });
}

export function useCreatePartnerAgreementDoc(agreementId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<PartnerAgreementDoc, "id" | "agreement_id" | "created_at" | "updated_at">) =>
      apiFetch<PartnerAgreementDoc>("/api/partner-agreement-docs", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, agreement_id: agreementId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PARTNER_AGREEMENT_DOCS_KEY(agreementId) });
    },
  });
}

export function useDeletePartnerAgreementDoc(agreementId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/partner-agreement-docs/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PARTNER_AGREEMENT_DOCS_KEY(agreementId) });
    },
  });
}
