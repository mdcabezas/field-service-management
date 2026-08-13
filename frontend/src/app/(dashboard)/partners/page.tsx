"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { usePartners } from "@/hooks/use-partners";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import type { Partner } from "@/types/api";

const columns: ColumnDef<Partner, unknown>[] = [
  {
    accessorKey: "name",
    header: "Nombre",
  },
  {
    accessorKey: "tax_id",
    header: "RUT",
  },
  {
    accessorKey: "status",
    header: "Estado",
    cell: ({ row }) => (
      <span className="uppercase">{row.original.status}</span>
    ),
  },];

export default function PartnersPage() {
  const router = useRouter();
  const { data: partners, isLoading } = usePartners();
  const { user } = useAuthStore();

  const handleRowClick = (row: Partner) => {
    router.push(`/partners/${row.id}`);
  };

  if (isLoading) {
    return <div className="font-mono">Cargando...</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Socios
        </h1>
        {canPerformAction(user?.role, "create") && (
          <Button onClick={() => router.push("/partners/new")}>
            + Crear Socio
          </Button>
        )}
      </div>

      <DataTable
        columns={columns}
        data={partners || []}
        searchPlaceholder="Buscar por nombre o RUT..."
        searchColumn="name"
        onRowClick={handleRowClick}
      />
    </div>
  );
}
