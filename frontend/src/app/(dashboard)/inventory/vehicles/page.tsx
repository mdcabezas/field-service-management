"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useVehicles } from "@/hooks/use-vehicles";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import type { Vehicle } from "@/types/api";

const columns: ColumnDef<Vehicle, unknown>[] = [
  { accessorKey: "license_plate", header: "Patente" },
  { accessorKey: "name", header: "Nombre" },
  { accessorKey: "brand", header: "Marca" },
  {
    accessorKey: "status",
    header: "Estado",
    cell: ({ row }) => <span className="uppercase">{row.original.status}</span>,
  },];

export default function VehiclesPage() {
  const router = useRouter();
  const { data, isLoading } = useVehicles();
  const { user } = useAuthStore();

  if (isLoading) return <div className="font-mono">Cargando...</div>;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Vehículos
        </h1>
        {canPerformAction(user?.role, "create") && (
          <Button onClick={() => router.push("/inventory/vehicles/new")}>
            + Crear
          </Button>
        )}
      </div>
      <DataTable
        columns={columns}
        data={data || []}
        searchPlaceholder="Buscar por patente..."
        searchColumn="license_plate"
        onRowClick={(row) => router.push(`/inventory/vehicles/${row.id}`)}
      />
    </div>
  );
}
