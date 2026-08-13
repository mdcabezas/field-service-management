"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useMaterials } from "@/hooks/use-materials";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import type { Material } from "@/types/api";

const columns: ColumnDef<Material, unknown>[] = [
  { accessorKey: "sku", header: "SKU" },
  { accessorKey: "name", header: "Nombre" },
  { accessorKey: "unit", header: "Unidad" },
  { accessorKey: "unit_cost", header: "Costo Unit." },];

export default function MaterialsPage() {
  const router = useRouter();
  const { data, isLoading } = useMaterials();
  const { user } = useAuthStore();

  if (isLoading) return <div className="font-mono">Cargando...</div>;

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
        searchPlaceholder="Buscar material..."
        searchColumn="name"
        onRowClick={(row) => router.push(`/inventory/materials/${row.id}`)}
        pageSize={5}
      />
    </div>
  );
}
