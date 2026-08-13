"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useUsers, useDeleteUser } from "@/hooks/use-users";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import type { User } from "@/types/api";
import { toast } from "sonner";

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
  const { data: users, isLoading } = useUsers();
  const deleteUser = useDeleteUser();
  const { user } = useAuthStore();

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

  if (isLoading) {
    return <div className="font-mono">Cargando...</div>;
  }

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
        searchPlaceholder="Buscar por usuario..."
        searchColumn="name"
        onRowClick={handleRowClick}
        pageSize={5}
      />
    </div>
  );
}
