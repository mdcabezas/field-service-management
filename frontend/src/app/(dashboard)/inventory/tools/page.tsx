"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useTools } from "@/hooks/use-tools";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { TableSkeleton } from "@/components/ui/skeleton";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";
import { useEffect } from "react";
import type { Tool } from "@/types/api";

const columns: ColumnDef<Tool, unknown>[] = [
  { accessorKey: "code", header: "Código" },
  { accessorKey: "name", header: "Nombre" },
  {
    accessorKey: "status",
    header: "Estado",
    cell: ({ row }) => <span className="uppercase">{row.original.status}</span>,
  },
];

export default function ToolsPage() {
  const router = useRouter();
  const { data, isLoading, error } = useTools();
  const { user } = useAuthStore();

  useEffect(() => {
    if (error) toast.error("Error al cargar herramientas");
  }, [error]);

  if (isLoading) return <TableSkeleton />;
  if (error) return <div className="font-mono text-red-600">Error al cargar datos</div>;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Herramientas
        </h1>
        {canPerformAction(user?.role, "create") && (
          <Button onClick={() => router.push("/inventory/tools/new")}>
            + Crear
          </Button>
        )}
      </div>
      <DataTable
        columns={columns}
        data={data || []}
        searchPlaceholder="Buscar por código o nombre..."
        searchColumn={["code", "name"]}
        onRowClick={(row) => router.push(`/inventory/tools/${row.id}`)}
      />
    </div>
  );
}
