"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useCustomers } from "@/hooks/use-customers";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { TableSkeleton } from "@/components/ui/skeleton";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";
import { useEffect } from "react";
import type { Customer } from "@/types/api";

const columns: ColumnDef<Customer, unknown>[] = [
  {
    accessorKey: "name",
    header: "Nombre",
  },
  {
    accessorKey: "tax_id",
    header: "RUT",
    cell: ({ row }) => row.original.tax_id || "-",
  },
  {
    accessorKey: "phone",
    header: "Teléfono",
    cell: ({ row }) => row.original.phone || "-",
  },
  {
    accessorKey: "email",
    header: "Email",
    cell: ({ row }) => row.original.email || "-",
  },
];

export default function CustomersPage() {
  const router = useRouter();
  const { data: customers, isLoading, error } = useCustomers();
  const { user } = useAuthStore();

  useEffect(() => {
    if (error) toast.error("Error al cargar clientes");
  }, [error]);

  const handleRowClick = (row: Customer) => {
    router.push(`/customers/${row.id}`);
  };

  if (isLoading) return <TableSkeleton />;
  if (error) return <div className="font-mono text-red-600">Error al cargar datos</div>;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Clientes
        </h1>
        {canPerformAction(user?.role, "create") && (
          <Button onClick={() => router.push("/customers/new")}>
            + Crear Cliente
          </Button>
        )}
      </div>

      <DataTable
        columns={columns}
        data={customers || []}
        searchPlaceholder="Buscar por nombre, RUT o email..."
        searchColumn={["name", "tax_id", "email"]}
        onRowClick={handleRowClick}
      />
    </div>
  );
}
