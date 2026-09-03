"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useDailyPlans } from "@/hooks/use-daily-plans";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { TableSkeleton } from "@/components/ui/skeleton";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";
import { useEffect } from "react";
import type { DailyPlan } from "@/types/api";

const columns: ColumnDef<DailyPlan, unknown>[] = [
  {
    accessorKey: "date",
    header: "Fecha",
  },
  {
    accessorKey: "name",
    header: "Nombre",
  },
  {
    accessorKey: "notes",
    header: "Notas",
    cell: ({ row }) => row.original.notes || "-",
  },
];

export default function DailyPlansPage() {
  const router = useRouter();
  const { data: dailyPlans, isLoading, error } = useDailyPlans();
  const { user } = useAuthStore();

  useEffect(() => {
    if (error) toast.error("Error al cargar planes diarios");
  }, [error]);

  const handleRowClick = (row: DailyPlan) => {
    router.push(`/planning/daily-plans/${row.id}`);
  };

  if (isLoading) return <TableSkeleton />;
  if (error) return <div className="font-mono text-red-600">Error al cargar datos</div>;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Planes Diarios
        </h1>
        {canPerformAction(user?.role, "create") && (
          <Button onClick={() => router.push("/planning/daily-plans/new")}>
            + Crear Plan Diario
          </Button>
        )}
      </div>

      <DataTable
        columns={columns}
        data={dailyPlans || []}
        searchPlaceholder="Buscar por fecha o notas..."
        searchColumn={["date", "notes", "name"]}
        onRowClick={handleRowClick}
      />
    </div>
  );
}
