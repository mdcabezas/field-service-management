"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useChecklistTemplates } from "@/hooks/use-checklist-templates";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { TableSkeleton } from "@/components/ui/skeleton";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";
import { useEffect } from "react";
import type { ChecklistTemplate } from "@/types/api";

const columns: ColumnDef<ChecklistTemplate, unknown>[] = [
  { accessorKey: "name", header: "Nombre" },
  { accessorKey: "description", header: "Descripción" },
  { accessorKey: "visit_type_display", header: "Tipo de Visita" },
  { accessorKey: "work_type_display", header: "Tipo de Trabajo" },
];

export default function ChecklistTemplatesPage() {
  const router = useRouter();
  const { data, isLoading, error } = useChecklistTemplates();
  const { user } = useAuthStore();

  useEffect(() => {
    if (error) toast.error("Error al cargar plantillas de checklist");
  }, [error]);

  if (isLoading) return <TableSkeleton />;
  if (error) return <div className="font-mono text-red-600">Error al cargar datos</div>;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Plantillas de Checklist
        </h1>
        {canPerformAction(user?.role, "create") && (
          <Button onClick={() => router.push("/inventory/checklist-templates/new")}>
            + Crear
          </Button>
        )}
      </div>
      <DataTable
        columns={columns}
        data={data || []}
        searchPlaceholder="Buscar por nombre o tipo..."
        searchColumn={["name", "visit_type", "work_type"]}
        onRowClick={(row) => router.push(`/inventory/checklist-templates/${row.id}`)}
      />
    </div>
  );
}
