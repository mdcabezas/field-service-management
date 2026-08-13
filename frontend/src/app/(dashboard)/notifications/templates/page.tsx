"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useNotificationTemplates } from "@/hooks/use-notification-templates";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import type { NotificationTemplate } from "@/types/api";

const columns: ColumnDef<NotificationTemplate, unknown>[] = [
  {
    accessorKey: "type",
    header: "Tipo",
  },
  {
    accessorKey: "channel",
    header: "Canal",
  },
  {
    accessorKey: "subject",
    header: "Asunto",
  },
  {
    accessorKey: "active",
    header: "Activo",
    cell: ({ row }) => (
      <span className="uppercase">{row.original.active ? "Sí" : "No"}</span>
    ),
  },];

export default function NotificationTemplatesPage() {
  const router = useRouter();
  const { data: templates, isLoading } = useNotificationTemplates();
  const { user } = useAuthStore();

  const handleRowClick = (row: NotificationTemplate) => {
    router.push(`/notifications/templates/${row.id}`);
  };

  if (isLoading) {
    return <div className="font-mono">Cargando...</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Plantillas de Notificaciones
        </h1>
        {canPerformAction(user?.role, "create") && (
          <Button onClick={() => router.push("/notifications/templates/new")}>
            + Crear Plantilla
          </Button>
        )}
      </div>

      <DataTable
        columns={columns}
        data={templates || []}
        searchPlaceholder="Buscar por asunto..."
        searchColumn="subject"
        onRowClick={handleRowClick}
      />
    </div>
  );
}
