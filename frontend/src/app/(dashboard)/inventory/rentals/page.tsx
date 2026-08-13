"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useRentals } from "@/hooks/use-rentals";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import type { Rental } from "@/types/api";

const columns: ColumnDef<Rental, unknown>[] = [
  { accessorKey: "type", header: "Tipo" },
  { accessorKey: "supplier", header: "Proveedor" },
  { accessorKey: "item_description", header: "Descripción" },
  { accessorKey: "daily_cost", header: "Costo Diario" },];

export default function RentalsPage() {
  const router = useRouter();
  const { data, isLoading } = useRentals();
  const { user } = useAuthStore();

  if (isLoading) return <div className="font-mono">Cargando...</div>;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Arriendos
        </h1>
        {canPerformAction(user?.role, "create") && (
          <Button onClick={() => router.push("/inventory/rentals/new")}>
            + Crear
          </Button>
        )}
      </div>
      <DataTable
        columns={columns}
        data={data || []}
        searchPlaceholder="Buscar arriendo..."
        searchColumn="type"
        onRowClick={(row) => router.push(`/inventory/rentals/${row.id}`)}
      />
    </div>
  );
}
