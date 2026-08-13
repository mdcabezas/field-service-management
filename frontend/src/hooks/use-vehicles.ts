"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { Vehicle } from "@/types/api";

const VEHICLES_KEY = ["vehicles"];

export function useVehicles() {
  return useQuery({
    queryKey: VEHICLES_KEY,
    queryFn: () => apiFetch<{ items: Vehicle[]; total: number }>("/api/vehicles").then(r => r.items),
  });
}

export function useVehicle(id: string) {
  return useQuery({
    queryKey: [...VEHICLES_KEY, id],
    queryFn: () => apiFetch<Vehicle>(`/api/vehicles/${id}`),
    enabled: !!id,
  });
}

export function useCreateVehicle() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<Vehicle, "id" | "created_at" | "updated_at">) =>
      apiFetch<Vehicle>("/api/vehicles", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VEHICLES_KEY });
    },
  });
}

export function useUpdateVehicle() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Vehicle> }) =>
      apiFetch<Vehicle>(`/api/vehicles/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VEHICLES_KEY });
    },
  });
}

export function useDeleteVehicle() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/vehicles/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VEHICLES_KEY });
    },
  });
}
