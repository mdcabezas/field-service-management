"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { VehicleAssignment } from "@/types/api";

const VEHICLE_ASSIGNMENTS_KEY = ["vehicle-assignments"];

export function useVehicleAssignments() {
  return useQuery({
    queryKey: VEHICLE_ASSIGNMENTS_KEY,
    queryFn: () => apiFetch<{ items: VehicleAssignment[]; total: number }>("/api/vehicle-assignments").then(r => r.items),
  });
}

export function useCreateVehicleAssignment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Omit<VehicleAssignment, "id" | "created_at" | "updated_at">) =>
      apiFetch<VehicleAssignment>("/api/vehicle-assignments", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VEHICLE_ASSIGNMENTS_KEY });
    },
  });
}

export function useDeleteVehicleAssignment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiFetch<void>(`/api/vehicle-assignments/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VEHICLE_ASSIGNMENTS_KEY });
    },
  });
}
