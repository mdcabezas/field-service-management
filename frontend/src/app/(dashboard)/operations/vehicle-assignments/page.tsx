"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useAllVehicleAssignments } from "@/hooks/use-all-vehicle-assignments";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { TableSkeleton } from "@/components/ui/skeleton";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";
import { useEffect } from "react";
import type { VehicleAssignment } from "@/types/api";

const columns: ColumnDef<VehicleAssignment, unknown>[] = [
  {
    accessorKey: "vehicle_name",
    header: "Vehículo",
    cell: ({ row }) => (
      <span className="font-mono text-sm">{row.original.vehicle_name || row.original.vehicle_id}</span>
    ),
  },
  {
    accessorKey: "daily_plan_name",
    header: "Plan Diario",
    cell: ({ row }) => (
      <span className="font-mono text-sm">{row.original.daily_plan_name || "-"}</span>
    ),
  },
  {
    accessorKey: "visit_label",
    header: "Visita",
    cell: ({ row }) => (
      <span className="font-mono text-sm">{row.original.visit_label || "-"}</span>
    ),
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
  },
];

export default function VehicleAssignmentsPage() {
  const router = useRouter();
  const { data: assignments, isLoading, error } = useAllVehicleAssignments();
  const { user } = useAuthStore();

  useEffect(() => {
    if (error) toast.error("Error al cargar asignaciones de vehículos");
  }, [error]);

  const handleRowClick = (row: VehicleAssignment) => {
    router.push(`/operations/vehicle-assignments/${row.id}`);
  };

  if (isLoading) return <TableSkeleton />;
  if (error) return <div className="font-mono text-red-600">Error al cargar datos</div>;

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
        searchPlaceholder="Buscar por vehículo, plan o visita..."
        searchColumn={["vehicle_name", "daily_plan_name", "visit_label"]}
        onRowClick={handleRowClick}
      />
    </div>
  );
}
