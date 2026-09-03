"use client";

import { useQuery } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { VisitMaterialUsage } from "@/types/api";

export function useVisitMaterialUsages(visitId: string) {
  return useQuery({
    queryKey: ["visit-material-usages", visitId],
    queryFn: () => apiFetch<VisitMaterialUsage[]>(`/api/visits/${visitId}/material-usages`),
    enabled: !!visitId,
  });
}
