"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { VehicleAssignment } from "@/types/api";

const ALL_VEHICLE_ASSIGNMENTS_KEY = ["vehicle-assignments"];

export function useAllVehicleAssignments() {
  return useQuery({
    queryKey: ALL_VEHICLE_ASSIGNMENTS_KEY,
    queryFn: () => apiFetch<{ items: VehicleAssignment[]; total: number }>("/api/vehicle-assignments").then(r => r.items),
  });
}

export function useVehicleAssignment(id: string) {
  return useQuery({
    queryKey: [...ALL_VEHICLE_ASSIGNMENTS_KEY, id],
    queryFn: () => apiFetch<VehicleAssignment>(`/api/vehicle-assignments/${id}`),
    enabled: !!id,
  });
}

export function useCreateVehicleAssignmentGlobal() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<VehicleAssignment, "id" | "created_at" | "updated_at">) =>
      apiFetch<VehicleAssignment>("/api/vehicle-assignments", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ALL_VEHICLE_ASSIGNMENTS_KEY });
    },
  });
}

export function useUpdateVehicleAssignment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<VehicleAssignment> }) =>
      apiFetch<VehicleAssignment>(`/api/vehicle-assignments/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ALL_VEHICLE_ASSIGNMENTS_KEY });
    },
  });
}

export function useDeleteVehicleAssignmentGlobal() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/vehicle-assignments/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ALL_VEHICLE_ASSIGNMENTS_KEY });
    },
  });
}
