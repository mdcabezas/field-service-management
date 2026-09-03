"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useAllVisitReports } from "@/hooks/use-all-visit-reports";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { TableSkeleton } from "@/components/ui/skeleton";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";
import { useEffect } from "react";
import type { VisitReport } from "@/types/api";

const columns: ColumnDef<VisitReport, unknown>[] = [
  {
    accessorKey: "report_template_name",
    header: "Plantilla",
    cell: ({ row }) => row.original.report_template_name || "-",
  },
  {
    accessorKey: "visit_label",
    header: "Visita",
    cell: ({ row }) => row.original.visit_label || "-",
  },
  {
    accessorKey: "source",
    header: "Fuente",
  },
  {
    accessorKey: "recorded_at",
    header: "Fecha",
    cell: ({ row }) => {
      if (!row.original.recorded_at) return "-";
      const date = new Date(row.original.recorded_at);
      return date.toLocaleDateString("es-CL");
    },
  },
];

export default function ReportsPage() {
  const router = useRouter();
  const { data: reports, isLoading, error } = useAllVisitReports();
  const { user } = useAuthStore();

  useEffect(() => {
    if (error) toast.error("Error al cargar reportes de visita");
  }, [error]);

  const handleRowClick = (row: VisitReport) => {
    router.push(`/operations/reports/${row.id}`);
  };

  if (isLoading) return <TableSkeleton />;
  if (error) return <div className="font-mono text-red-600">Error al cargar datos</div>;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Reportes de Visita
        </h1>
        {canPerformAction(user?.role, "create") && (
          <Button onClick={() => router.push("/operations/reports/new")}>
            + Crear Reporte
          </Button>
        )}
      </div>

      <DataTable
        columns={columns}
        data={reports || []}
        searchPlaceholder="Buscar por plantilla, visita o fuente..."
        searchColumn={["report_template_name", "visit_label", "source"]}
        onRowClick={handleRowClick}
      />
    </div>
  );
}
