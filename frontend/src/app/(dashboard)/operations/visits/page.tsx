"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useVisits } from "@/hooks/use-visits";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import type { Visit } from "@/types/api";

const columns: ColumnDef<Visit, unknown>[] = [
  {
    accessorKey: "scheduled_at",
    header: "Fecha",
    cell: ({ row }) => {
      const date = new Date(row.original.scheduled_at);
      return date.toLocaleDateString("es-CL");
    },
  },
  {
    accessorKey: "priority",
    header: "Prioridad",
    cell: ({ row }) => (
      <span className="uppercase">{row.original.priority}</span>
    ),
  },
  {
    accessorKey: "status",
    header: "Estado",
    cell: ({ row }) => (
      <span className="uppercase">{row.original.status}</span>
    ),
  },];

export default function VisitsPage() {
  const router = useRouter();
  const { data: visits, isLoading } = useVisits();
  const { user } = useAuthStore();

  const handleRowClick = (row: Visit) => {
    router.push(`/operations/visits/${row.id}`);
  };

  if (isLoading) {
    return <div className="font-mono">Cargando...</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Visitas
        </h1>
        {canPerformAction(user?.role, "create") && (
          <Button onClick={() => router.push("/operations/visits/new")}>
            + Crear Visita
          </Button>
        )}
      </div>

      <DataTable
        columns={columns}
        data={visits || []}
        searchPlaceholder="Buscar por estado o prioridad..."
        searchColumn="status"
        onRowClick={handleRowClick}
      />
    </div>
  );
}
