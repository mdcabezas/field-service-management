"use client";

import { useRouter, useParams } from "next/navigation";
import { useUser, useDeleteUser } from "@/hooks/use-users";
import { UserForm } from "@/components/forms/user-form";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";

export default function UserDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const { data: user, isLoading } = useUser(id);
  const deleteUser = useDeleteUser();
  const { user: currentUser } = useAuthStore();

  const handleDelete = async () => {
    if (confirm("¿Estás seguro de eliminar este usuario?")) {
      try {
        await deleteUser.mutateAsync(id);
        toast.success("Usuario eliminado");
        router.push("/core/users");
      } catch (error) {
        toast.error("Error al eliminar usuario");
      }
    }
  };

  if (isLoading) {
    return <div className="font-mono">Cargando...</div>;
  }

  if (!user) {
    return <div className="font-mono">Usuario no encontrado</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Editar Usuario
        </h1>
        {canPerformAction(currentUser?.role, "delete") && (
          <Button variant="destructive" onClick={handleDelete}>
            Eliminar
          </Button>
        )}
      </div>

      <UserForm initialData={user} isEdit />
    </div>
  );
}
