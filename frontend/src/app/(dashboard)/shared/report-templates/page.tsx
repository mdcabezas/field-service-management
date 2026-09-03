"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useReportTemplates } from "@/hooks/use-report-templates";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { TableSkeleton } from "@/components/ui/skeleton";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";
import { useEffect } from "react";
import type { ReportTemplate } from "@/types/api";

const columns: ColumnDef<ReportTemplate, unknown>[] = [
  {
    accessorKey: "name",
    header: "Nombre",
  },
  {
    accessorKey: "fields_json",
    header: "Campos",
    cell: ({ row }) => {
      const fields = row.original.fields_json;
      const keys = fields ? Object.keys(fields) : [];
      if (keys.length === 0) return <span className="text-gray-400">—</span>;
      return (
        <span className="font-mono text-xs" title={JSON.stringify(fields)}>
          {keys.length} campo(s): {keys.join(", ")}
        </span>
      );
    },
  },
];

export default function ReportTemplatesPage() {
  const router = useRouter();
  const { data: templates, isLoading, error } = useReportTemplates();
  const { user } = useAuthStore();

  useEffect(() => {
    if (error) toast.error("Error al cargar plantillas de reportes");
  }, [error]);

  const handleRowClick = (row: ReportTemplate) => {
    router.push(`/shared/report-templates/${row.id}`);
  };

  if (isLoading) return <TableSkeleton />;
  if (error) return <div className="font-mono text-red-600">Error al cargar datos</div>;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Plantillas de Reportes
        </h1>
        {canPerformAction(user?.role, "create") && (
          <Button onClick={() => router.push("/shared/report-templates/new")}>
            + Crear Plantilla
          </Button>
        )}
      </div>

      <DataTable
        columns={columns}
        data={templates || []}
        searchPlaceholder="Buscar por nombre..."
        searchColumn="name"
        onRowClick={handleRowClick}
      />
    </div>
  );
}
