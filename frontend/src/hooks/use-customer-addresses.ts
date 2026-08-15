"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { CustomerAddress } from "@/types/api";

const CUSTOMER_ADDRESSES_KEY = (customerId: string) => ["customer-addresses", customerId];

export function useCustomerAddresses(customerId: string) {
  return useQuery({
    queryKey: CUSTOMER_ADDRESSES_KEY(customerId),
    queryFn: () => apiFetch<{ items: CustomerAddress[]; total: number }>(`/api/customers/${customerId}/addresses`).then(r => r.items),
    enabled: !!customerId,
  });
}

export function useCreateCustomerAddress(customerId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<CustomerAddress, "id" | "customer_id" | "created_at" | "updated_at">) =>
      apiFetch<CustomerAddress>("/api/customer-addresses", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, customer_id: customerId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: CUSTOMER_ADDRESSES_KEY(customerId) });
    },
  });
}

export function useDeleteCustomerAddress(customerId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/customer-addresses/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: CUSTOMER_ADDRESSES_KEY(customerId) });
    },
  });
}

export function useUpdateCustomerAddress(customerId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<CustomerAddress> }) =>
      apiFetch<CustomerAddress>(`/api/customer-addresses/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: CUSTOMER_ADDRESSES_KEY(customerId) });
    },
  });
}
