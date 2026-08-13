"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useMaintenanceSchedules } from "@/hooks/use-maintenance-schedules";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import type { MaintenanceSchedule } from "@/types/api";

const columns: ColumnDef<MaintenanceSchedule, unknown>[] = [
  { accessorKey: "reference_id", header: "Referencia" },
  { accessorKey: "type", header: "Tipo" },
  { accessorKey: "frequency_days", header: "Frecuencia (días)" },
  {
    accessorKey: "active",
    header: "Activo",
    cell: ({ row }) => <span className="uppercase">{row.original.active ? "Sí" : "No"}</span>,
  },];

export default function MaintenanceSchedulesPage() {
  const router = useRouter();
  const { data, isLoading } = useMaintenanceSchedules();
  const { user } = useAuthStore();

  if (isLoading) return <div className="font-mono">Cargando...</div>;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Programación de Mantención
        </h1>
        {canPerformAction(user?.role, "create") && (
          <Button onClick={() => router.push("/inventory/maintenance-schedules/new")}>
            + Crear
          </Button>
        )}
      </div>
      <DataTable
        columns={columns}
        data={data || []}
        searchPlaceholder="Buscar programación..."
        searchColumn="type"
        onRowClick={(row) => router.push(`/inventory/maintenance-schedules/${row.id}`)}
      />
    </div>
  );
}
