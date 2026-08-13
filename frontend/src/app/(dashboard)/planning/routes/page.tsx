"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useRoutes } from "@/hooks/use-routes";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import type { Route } from "@/types/api";

const columns: ColumnDef<Route, unknown>[] = [
  {
    accessorKey: "type",
    header: "Tipo",
  },
  {
    accessorKey: "date",
    header: "Fecha",
  },
  {
    accessorKey: "status",
    header: "Estado",
    cell: ({ row }) => (
      <span className="uppercase">{row.original.status}</span>
    ),
  },];

export default function RoutesPage() {
  const router = useRouter();
  const { data: routes, isLoading } = useRoutes();
  const { user } = useAuthStore();

  const handleRowClick = (row: Route) => {
    router.push(`/planning/routes/${row.id}`);
  };

  if (isLoading) {
    return <div className="font-mono">Cargando...</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Rutas
        </h1>
        {canPerformAction(user?.role, "create") && (
          <Button onClick={() => router.push("/planning/routes/new")}>
            + Crear Ruta
          </Button>
        )}
      </div>

      <DataTable
        columns={columns}
        data={routes || []}
        searchPlaceholder="Buscar por tipo..."
        searchColumn="type"
        onRowClick={handleRowClick}
      />
    </div>
  );
}
