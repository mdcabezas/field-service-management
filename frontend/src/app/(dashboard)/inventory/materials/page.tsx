"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useMaterials } from "@/hooks/use-materials";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { TableSkeleton } from "@/components/ui/skeleton";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";
import { useEffect } from "react";
import type { Material } from "@/types/api";

const columns: ColumnDef<Material, unknown>[] = [
  { accessorKey: "sku", header: "SKU" },
  { accessorKey: "name", header: "Nombre" },
  { accessorKey: "unit", header: "Unidad" },
  {
    accessorKey: "unit_cost",
    header: "Costo Unit.",
    cell: ({ row }) =>
      row.original.unit_cost != null
        ? `$${row.original.unit_cost.toLocaleString("es-CL")}`
        : "-",
  },
];

export default function MaterialsPage() {
  const router = useRouter();
  const { data, isLoading, error } = useMaterials();
  const { user } = useAuthStore();

  useEffect(() => {
    if (error) toast.error("Error al cargar materiales");
  }, [error]);

  if (isLoading) return <TableSkeleton />;
  if (error) return <div className="font-mono text-red-600">Error al cargar datos</div>;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Materiales
        </h1>
        {canPerformAction(user?.role, "create") && (
          <Button onClick={() => router.push("/inventory/materials/new")}>
            + Crear
          </Button>
        )}
      </div>
      <DataTable
        columns={columns}
        data={data || []}
        searchPlaceholder="Buscar por nombre o SKU..."
        searchColumn={["name", "sku"]}
        onRowClick={(row) => router.push(`/inventory/materials/${row.id}`)}
        pageSize={5}
      />
    </div>
  );
}
