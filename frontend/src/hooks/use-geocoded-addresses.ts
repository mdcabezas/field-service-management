"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { GeocodedAddress } from "@/types/api";

const GEOCODED_ADDRESSES_KEY = ["geocoded-addresses"];

export function useGeocodedAddresses() {
  return useQuery({
    queryKey: GEOCODED_ADDRESSES_KEY,
    queryFn: () => apiFetch<{ items: GeocodedAddress[]; total: number }>("/api/addresses").then(r => r.items),
  });
}

export function useGeocodedAddress(id: string) {
  return useQuery({
    queryKey: [...GEOCODED_ADDRESSES_KEY, id],
    queryFn: () => apiFetch<GeocodedAddress>(`/api/addresses/${id}`),
    enabled: !!id,
  });
}

export function useNearbyAddresses(lat: number, lng: number, radius: number) {
  return useQuery({
    queryKey: [...GEOCODED_ADDRESSES_KEY, "nearby", lat, lng, radius],
    queryFn: () =>
      apiFetch<GeocodedAddress[]>(
        `/api/addresses/nearby?lat=${lat}&lng=${lng}&radius=${radius}`
      ),
    enabled: !!lat && !!lng,
  });
}

export function useCreateGeocodedAddress() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<GeocodedAddress, "id" | "created_at" | "updated_at">) =>
      apiFetch<GeocodedAddress>("/api/addresses", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: GEOCODED_ADDRESSES_KEY });
    },
  });
}

export function useUpdateGeocodedAddress() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<GeocodedAddress> }) =>
      apiFetch<GeocodedAddress>(`/api/addresses/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: GEOCODED_ADDRESSES_KEY });
    },
  });
}

export function useDeleteGeocodedAddress() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/addresses/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: GEOCODED_ADDRESSES_KEY });
    },
  });
}
