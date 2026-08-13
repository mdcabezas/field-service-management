"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useMaintenanceRecords } from "@/hooks/use-maintenance-records";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import type { MaintenanceRecord } from "@/types/api";

const columns: ColumnDef<MaintenanceRecord, unknown>[] = [
  { accessorKey: "reference_id", header: "Referencia" },
  { accessorKey: "type", header: "Tipo" },
  { accessorKey: "date", header: "Fecha" },
  { accessorKey: "cost", header: "Costo" },
  { accessorKey: "supplier", header: "Proveedor" },];

export default function MaintenanceRecordsPage() {
  const router = useRouter();
  const { data, isLoading } = useMaintenanceRecords();
  const { user } = useAuthStore();

  if (isLoading) return <div className="font-mono">Cargando...</div>;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Registros de Mantención
        </h1>
        {canPerformAction(user?.role, "create") && (
          <Button onClick={() => router.push("/inventory/maintenance-records/new")}>
            + Crear
          </Button>
        )}
      </div>
      <DataTable
        columns={columns}
        data={data || []}
        searchPlaceholder="Buscar registro..."
        searchColumn="reference_id"
        onRowClick={(row) => router.push(`/inventory/maintenance-records/${row.id}`)}
      />
    </div>
  );
}
