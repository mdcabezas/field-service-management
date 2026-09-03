"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useRoutes } from "@/hooks/use-routes";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { TableSkeleton } from "@/components/ui/skeleton";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";
import { useEffect } from "react";
import type { Route } from "@/types/api";

const columns: ColumnDef<Route, unknown>[] = [
  {
    accessorKey: "type_name",
    accessorFn: (row) => row.type_name || row.type,
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
  },
];

export default function RoutesPage() {
  const router = useRouter();
  const { data: routes, isLoading, error } = useRoutes();
  const { user } = useAuthStore();

  useEffect(() => {
    if (error) toast.error("Error al cargar rutas");
  }, [error]);

  const handleRowClick = (row: Route) => {
    router.push(`/planning/routes/${row.id}`);
  };

  if (isLoading) return <TableSkeleton />;
  if (error) return <div className="font-mono text-red-600">Error al cargar datos</div>;

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
        searchPlaceholder="Buscar por tipo o estado..."
        searchColumn={["type", "type_name", "status"]}
        onRowClick={handleRowClick}
      />
    </div>
  );
}
