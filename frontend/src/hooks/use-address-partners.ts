"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { CustomerAddressPartner } from "@/types/api";

const ADDRESS_PARTNERS_KEY = (addressId: string) => ["address-partners", addressId];

export function useAddressPartners(addressId: string) {
  return useQuery({
    queryKey: ADDRESS_PARTNERS_KEY(addressId),
    queryFn: () => apiFetch<{ items: CustomerAddressPartner[]; total: number }>(`/api/customer-addresses/${addressId}/partners`).then(r => r.items),
    enabled: !!addressId,
  });
}

export function useCreateAddressPartner(addressId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<CustomerAddressPartner, "id" | "address_id" | "created_at" | "updated_at">) =>
      apiFetch<CustomerAddressPartner>("/api/customer-address-partners", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, address_id: addressId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ADDRESS_PARTNERS_KEY(addressId) });
    },
  });
}

export function useDeleteAddressPartner(addressId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/customer-address-partners/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ADDRESS_PARTNERS_KEY(addressId) });
    },
  });
}
