"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useChecklistTemplates } from "@/hooks/use-checklist-templates";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import type { ChecklistTemplate } from "@/types/api";

const columns: ColumnDef<ChecklistTemplate, unknown>[] = [
  { accessorKey: "name", header: "Nombre" },
  { accessorKey: "description", header: "Descripción" },
  { accessorKey: "visit_type", header: "Tipo de Visita" },];

export default function ChecklistTemplatesPage() {
  const router = useRouter();
  const { data, isLoading } = useChecklistTemplates();
  const { user } = useAuthStore();

  if (isLoading) return <div className="font-mono">Cargando...</div>;

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
        searchPlaceholder="Buscar plantilla..."
        searchColumn="name"
        onRowClick={(row) => router.push(`/inventory/checklist-templates/${row.id}`)}
      />
    </div>
  );
}
