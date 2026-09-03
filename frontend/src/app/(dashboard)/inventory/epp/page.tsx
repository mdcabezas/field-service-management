"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useEPPItems } from "@/hooks/use-epp-items";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { TableSkeleton } from "@/components/ui/skeleton";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";
import { useEffect } from "react";
import type { EPPItem } from "@/types/api";

const columns: ColumnDef<EPPItem, unknown>[] = [
  { accessorKey: "name", header: "Nombre" },
  { accessorKey: "type", header: "Tipo" },
  { accessorKey: "lifecycle", header: "Ciclo de vida" },
];

export default function EPPPage() {
  const router = useRouter();
  const { data, isLoading, error } = useEPPItems();
  const { user } = useAuthStore();

  useEffect(() => {
    if (error) toast.error("Error al cargar elementos EPP");
  }, [error]);

  if (isLoading) return <TableSkeleton />;
  if (error) return <div className="font-mono text-red-600">Error al cargar datos</div>;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Elementos de Protección Personal
        </h1>
        {canPerformAction(user?.role, "create") && (
          <Button onClick={() => router.push("/inventory/epp/new")}>
            + Crear
          </Button>
        )}
      </div>
      <DataTable
        columns={columns}
        data={data || []}
        searchPlaceholder="Buscar por nombre o tipo..."
        searchColumn={["name", "type"]}
        onRowClick={(row) => router.push(`/inventory/epp/${row.id}`)}
      />
    </div>
  );
}
