"use client";

import { useQuery } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { VehicleType } from "@/types/api";

const VEHICLE_TYPES_KEY = ["vehicle-types"];

export function useVehicleTypes() {
  return useQuery({
    queryKey: VEHICLE_TYPES_KEY,
    queryFn: () =>
      apiFetch<{ items: VehicleType[]; total: number }>("/api/vehicle-types").then(
        (r) => r.items
      ),
  });
}
