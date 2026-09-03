"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useUsers, useDeleteUser } from "@/hooks/use-users";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { TableSkeleton } from "@/components/ui/skeleton";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";
import { useEffect } from "react";
import type { User } from "@/types/api";

const columns: ColumnDef<User, unknown>[] = [
  {
    accessorKey: "name",
    header: "Nombre",
  },
  {
    accessorKey: "email",
    header: "Email",
  },
  {
    accessorKey: "role",
    header: "Rol",
    cell: ({ row }) => (
      <span className="uppercase">{row.original.role}</span>
    ),
  },
];

export default function UsersPage() {
  const router = useRouter();
  const { data: users, isLoading, error } = useUsers();
  const deleteUser = useDeleteUser();
  const { user } = useAuthStore();

  useEffect(() => {
    if (error) toast.error("Error al cargar usuarios");
  }, [error]);

  const handleRowClick = (row: User) => {
    router.push(`/core/users/${row.id}`);
  };

  const handleDelete = async (id: string) => {
    if (confirm("¿Estás seguro de eliminar este usuario?")) {
      try {
        await deleteUser.mutateAsync(id);
        toast.success("Usuario eliminado");
      } catch (error) {
        toast.error("Error al eliminar usuario");
      }
    }
  };

  if (isLoading) return <TableSkeleton />;
  if (error) return <div className="font-mono text-red-600">Error al cargar datos</div>;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Usuarios
        </h1>
        {canPerformAction(user?.role, "create") && (
          <Button onClick={() => router.push("/core/users/new")}>
            + Crear Usuario
          </Button>
        )}
      </div>

      <DataTable
        columns={columns}
        data={users || []}
        searchPlaceholder="Buscar por nombre o email..."
        searchColumn={["name", "email"]}
        onRowClick={handleRowClick}
        pageSize={5}
      />
    </div>
  );
}
