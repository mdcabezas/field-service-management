"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { Property } from "@/types/api";

const PROPERTIES_KEY = (addressId: string) => ["properties", addressId];

export function useProperties(addressId: string) {
  return useQuery({
    queryKey: PROPERTIES_KEY(addressId),
    queryFn: () => apiFetch<{ items: Property[]; total: number }>(`/api/customer-addresses/${addressId}/properties`).then(r => r.items),
    enabled: !!addressId,
  });
}

export function useCreateProperty(addressId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<Property, "id" | "address_id" | "created_at" | "updated_at">) =>
      apiFetch<Property>("/api/properties", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, address_id: addressId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PROPERTIES_KEY(addressId) });
    },
  });
}

export function useDeleteProperty(addressId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/properties/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PROPERTIES_KEY(addressId) });
    },
  });
}
