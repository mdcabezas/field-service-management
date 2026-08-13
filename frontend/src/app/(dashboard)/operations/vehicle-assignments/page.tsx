"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useAllVehicleAssignments } from "@/hooks/use-all-vehicle-assignments";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import type { VehicleAssignment } from "@/types/api";

const columns: ColumnDef<VehicleAssignment, unknown>[] = [
  {
    accessorKey: "vehicle_id",
    header: "Vehículo",
  },
  {
    accessorKey: "departure_time",
    header: "Salida",
    cell: ({ row }) => {
      if (!row.original.departure_time) return "-";
      const date = new Date(row.original.departure_time);
      return date.toLocaleDateString("es-CL");
    },
  },
  {
    accessorKey: "return_time",
    header: "Retorno",
    cell: ({ row }) => {
      if (!row.original.return_time) return "-";
      const date = new Date(row.original.return_time);
      return date.toLocaleDateString("es-CL");
    },
  },];

export default function VehicleAssignmentsPage() {
  const router = useRouter();
  const { data: assignments, isLoading } = useAllVehicleAssignments();
  const { user } = useAuthStore();

  const handleRowClick = (row: VehicleAssignment) => {
    router.push(`/operations/vehicle-assignments/${row.id}`);
  };

  if (isLoading) {
    return <div className="font-mono">Cargando...</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Asignaciones de Vehículos
        </h1>
        {canPerformAction(user?.role, "create") && (
          <Button onClick={() => router.push("/operations/vehicle-assignments/new")}>
            + Crear Asignación
          </Button>
        )}
      </div>

      <DataTable
        columns={columns}
        data={assignments || []}
        searchPlaceholder="Buscar por vehículo o visita..."
        searchColumn="vehicle_id"
        onRowClick={handleRowClick}
      />
    </div>
  );
}
