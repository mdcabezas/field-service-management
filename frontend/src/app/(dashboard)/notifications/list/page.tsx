"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useNotifications } from "@/hooks/use-notifications";
import { DataTable } from "@/components/data-table/data-table";
import { TableSkeleton } from "@/components/ui/skeleton";
import { toast } from "sonner";
import { useEffect } from "react";
import type { Notification } from "@/types/api";

const columns: ColumnDef<Notification, unknown>[] = [
  {
    accessorKey: "recipient",
    header: "Destinatario",
  },
  {
    accessorKey: "type",
    header: "Plantilla",
  },
  {
    accessorKey: "channel",
    header: "Canal",
  },
  {
    accessorKey: "status",
    header: "Estado",
    cell: ({ row }) => (
      <span className="uppercase">{row.original.status}</span>
    ),
  },
  {
    accessorKey: "sent_at",
    header: "Enviado",
    cell: ({ row }) =>
      row.original.sent_at
        ? new Date(row.original.sent_at).toLocaleString()
        : "-",
  },
];

export default function NotificationsListPage() {
  const router = useRouter();
  const { data: notifications, isLoading, error } = useNotifications();

  useEffect(() => {
    if (error) toast.error("Error al cargar notificaciones");
  }, [error]);

  const handleRowClick = (row: Notification) => {
    router.push(`/notifications/list/${row.id}`);
  };

  if (isLoading) return <TableSkeleton />;
  if (error) return <div className="font-mono text-red-600">Error al cargar datos</div>;

  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Notificaciones Enviadas
      </h1>

      <DataTable
        columns={columns}
        data={notifications || []}
        searchPlaceholder="Buscar por destinatario, tipo o canal..."
        searchColumn={["recipient", "type", "channel"]}
        onRowClick={handleRowClick}
      />
    </div>
  );
}
