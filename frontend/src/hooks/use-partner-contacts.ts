"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { PartnerContact } from "@/types/api";

const PARTNER_CONTACTS_KEY = (partnerId: string) => ["partner-contacts", partnerId];

export function usePartnerContacts(partnerId: string) {
  return useQuery({
    queryKey: PARTNER_CONTACTS_KEY(partnerId),
    queryFn: () => apiFetch<{ items: PartnerContact[]; total: number }>(`/api/partners/${partnerId}/contacts`).then(r => r.items),
    enabled: !!partnerId,
  });
}

export function useCreatePartnerContact(partnerId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<PartnerContact, "id" | "partner_id" | "created_at" | "updated_at">) =>
      apiFetch<PartnerContact>("/api/partner-contacts", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, partner_id: partnerId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PARTNER_CONTACTS_KEY(partnerId) });
    },
  });
}

export function useDeletePartnerContact(partnerId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/partner-contacts/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PARTNER_CONTACTS_KEY(partnerId) });
    },
  });
}
