"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useCustomers } from "@/hooks/use-customers";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import type { Customer } from "@/types/api";

const columns: ColumnDef<Customer, unknown>[] = [
  {
    accessorKey: "name",
    header: "Nombre",
  },
  {
    accessorKey: "tax_id",
    header: "RUT",
  },
  {
    accessorKey: "phone",
    header: "Teléfono",
  },
  {
    accessorKey: "email",
    header: "Email",
  },];

export default function CustomersPage() {
  const router = useRouter();
  const { data: customers, isLoading } = useCustomers();
  const { user } = useAuthStore();

  const handleRowClick = (row: Customer) => {
    router.push(`/customers/${row.id}`);
  };

  if (isLoading) {
    return <div className="font-mono">Cargando...</div>;
  }

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
        searchPlaceholder="Buscar por nombre o RUT..."
        searchColumn="name"
        onRowClick={handleRowClick}
      />
    </div>
  );
}
