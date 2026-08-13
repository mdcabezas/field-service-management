"use client";

import { useQuery } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { RouteType } from "@/types/api";

export type { RouteType };

export interface RouteTypesResult {
  types: RouteType[];
  isDropdown: boolean;
  isLoading: boolean;
}

export function useRouteTypes(): RouteTypesResult {
  const { data, isLoading } = useQuery({
    queryKey: ["route-types"],
    queryFn: () => apiFetch<{ items: RouteType[]; total: number }>("/api/route-types").then(r => r.items),
  });

  return {
    types: data || [],
    isDropdown: true,
    isLoading,
  };
}
