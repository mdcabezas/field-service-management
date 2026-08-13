"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useAllVisitSLATrackings } from "@/hooks/use-all-visit-sla-trackings";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import type { VisitSLATracking } from "@/types/api";

const columns: ColumnDef<VisitSLATracking, unknown>[] = [
  {
    accessorKey: "visit_id",
    header: "Visita",
  },
  {
    accessorKey: "sla_id",
    header: "SLA",
  },
  {
    accessorKey: "requested_at",
    header: "Solicitado",
    cell: ({ row }) => {
      if (!row.original.requested_at) return "-";
      const date = new Date(row.original.requested_at);
      return date.toLocaleDateString("es-CL");
    },
  },
  {
    accessorKey: "resolved_at",
    header: "Resuelto",
    cell: ({ row }) => {
      if (!row.original.resolved_at) return "-";
      const date = new Date(row.original.resolved_at);
      return date.toLocaleDateString("es-CL");
    },
  },];

export default function SLATrackingsPage() {
  const router = useRouter();
  const { data: slaTrackings, isLoading } = useAllVisitSLATrackings();
  const { user } = useAuthStore();

  const handleRowClick = (row: VisitSLATracking) => {
    router.push(`/operations/sla-trackings/${row.id}`);
  };

  if (isLoading) {
    return <div className="font-mono">Cargando...</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Seguimiento SLA
        </h1>
        {canPerformAction(user?.role, "create") && (
          <Button onClick={() => router.push("/operations/sla-trackings/new")}>
            + Crear Seguimiento
          </Button>
        )}
      </div>

      <DataTable
        columns={columns}
        data={slaTrackings || []}
        searchPlaceholder="Buscar por SLA o visita..."
        searchColumn="sla_id"
        onRowClick={handleRowClick}
      />
    </div>
  );
}
