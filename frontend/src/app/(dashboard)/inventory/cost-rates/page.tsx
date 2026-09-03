"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useCostRates } from "@/hooks/use-cost-rates";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { TableSkeleton } from "@/components/ui/skeleton";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";
import { useEffect } from "react";
import type { CostRate } from "@/types/api";

const columns: ColumnDef<CostRate, unknown>[] = [
  { accessorKey: "type", header: "Tipo" },
  {
    accessorKey: "value",
    header: "Valor",
    cell: ({ row }) =>
      row.original.value != null
        ? `$${row.original.value.toLocaleString("es-CL")}`
        : "-",
  },
  { accessorKey: "unit", header: "Unidad" },
  { accessorKey: "valid_from", header: "Vigente desde" },
];

export default function CostRatesPage() {
  const router = useRouter();
  const { data, isLoading, error } = useCostRates();
  const { user } = useAuthStore();

  useEffect(() => {
    if (error) toast.error("Error al cargar tarifas");
  }, [error]);

  if (isLoading) return <TableSkeleton />;
  if (error) return <div className="font-mono text-red-600">Error al cargar datos</div>;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Tarifas de Costo
        </h1>
        {canPerformAction(user?.role, "create") && (
          <Button onClick={() => router.push("/inventory/cost-rates/new")}>
            + Crear
          </Button>
        )}
      </div>
      <DataTable
        columns={columns}
        data={data || []}
        searchPlaceholder="Buscar por tipo o unidad..."
        searchColumn={["type", "unit"]}
        onRowClick={(row) => router.push(`/inventory/cost-rates/${row.id}`)}
      />
    </div>
  );
}
