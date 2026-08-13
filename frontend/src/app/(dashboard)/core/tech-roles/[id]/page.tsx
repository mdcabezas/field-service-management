"use client";

import { useRouter, useParams } from "next/navigation";
import { useTechRole, useDeleteTechRole } from "@/hooks/use-tech-roles";
import { TechRoleForm } from "@/components/forms/tech-role-form";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";

export default function TechRoleDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const { data: techRole, isLoading } = useTechRole(id);
  const deleteTechRole = useDeleteTechRole();
  const { user: currentUser } = useAuthStore();

  const handleDelete = async () => {
    if (confirm("¿Estás seguro de eliminar este rol?")) {
      try {
        await deleteTechRole.mutateAsync(id);
        toast.success("Rol eliminado");
        router.push("/core/tech-roles");
      } catch (error) {
        toast.error("Error al eliminar rol");
      }
    }
  };

  if (isLoading) {
    return <div className="font-mono">Cargando...</div>;
  }

  if (!techRole) {
    return <div className="font-mono">Rol no encontrado</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Editar Rol
        </h1>
        {canPerformAction(currentUser?.role, "delete") && (
          <Button variant="destructive" onClick={handleDelete}>
            Eliminar
          </Button>
        )}
      </div>

      <TechRoleForm initialData={techRole} isEdit />
    </div>
  );
}
