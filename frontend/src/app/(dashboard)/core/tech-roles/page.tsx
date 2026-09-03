"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useTechRoles } from "@/hooks/use-tech-roles";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { TableSkeleton } from "@/components/ui/skeleton";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";
import { useEffect } from "react";
import type { TechRole } from "@/types/api";

const columns: ColumnDef<TechRole, unknown>[] = [
  {
    accessorKey: "name",
    header: "Nombre",
  },
  {
    accessorKey: "code",
    header: "Código",
  },
  {
    accessorKey: "active",
    header: "Activo",
    cell: ({ row }) => (
      <span className="uppercase">{row.original.active ? "Sí" : "No"}</span>
    ),
  },
];

export default function TechRolesPage() {
  const router = useRouter();
  const { data: techRoles, isLoading, error } = useTechRoles();
  const { user } = useAuthStore();

  useEffect(() => {
    if (error) toast.error("Error al cargar roles de técnico");
  }, [error]);

  const handleRowClick = (row: TechRole) => {
    router.push(`/core/tech-roles/${row.id}`);
  };

  if (isLoading) return <TableSkeleton />;
  if (error) return <div className="font-mono text-red-600">Error al cargar datos</div>;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Roles de Técnico
        </h1>
        {canPerformAction(user?.role, "create") && (
          <Button onClick={() => router.push("/core/tech-roles/new")}>
            + Crear Rol
          </Button>
        )}
      </div>

      <DataTable
        columns={columns}
        data={techRoles || []}
        searchPlaceholder="Buscar por nombre o código..."
        searchColumn={["name", "code"]}
        onRowClick={handleRowClick}
      />
    </div>
  );
}
