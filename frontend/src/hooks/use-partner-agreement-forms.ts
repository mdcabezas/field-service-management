"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { PartnerAgreementForm } from "@/types/api";

const PARTNER_AGREEMENT_FORMS_KEY = (agreementId: string) => ["partner-agreement-forms", agreementId];

export function usePartnerAgreementForms(agreementId: string) {
  return useQuery({
    queryKey: PARTNER_AGREEMENT_FORMS_KEY(agreementId),
    queryFn: () => apiFetch<{ items: PartnerAgreementForm[]; total: number }>(`/api/partner-agreements/${agreementId}/forms`).then(r => r.items),
    enabled: !!agreementId,
  });
}

export function useCreatePartnerAgreementForm(agreementId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<PartnerAgreementForm, "id" | "agreement_id" | "created_at" | "updated_at">) =>
      apiFetch<PartnerAgreementForm>("/api/partner-agreement-forms", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, agreement_id: agreementId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PARTNER_AGREEMENT_FORMS_KEY(agreementId) });
    },
  });
}

export function useDeletePartnerAgreementForm(agreementId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/partner-agreement-forms/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PARTNER_AGREEMENT_FORMS_KEY(agreementId) });
    },
  });
}
