"use client";

import { useRouter, useParams } from "next/navigation";
import { usePartner, useDeletePartner } from "@/hooks/use-partners";
import { PartnerForm } from "@/components/forms/partner-form";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";

export default function PartnerDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const { data: partner, isLoading } = usePartner(id);
  const deletePartner = useDeletePartner();
  const { user: currentUser } = useAuthStore();

  const handleDelete = async () => {
    if (confirm("¿Estás seguro de eliminar este socio?")) {
      try {
        await deletePartner.mutateAsync(id);
        toast.success("Socio eliminado");
        router.push("/partners");
      } catch (error) {
        toast.error("Error al eliminar socio");
      }
    }
  };

  if (isLoading) {
    return <div className="font-mono">Cargando...</div>;
  }

  if (!partner) {
    return <div className="font-mono">Socio no encontrado</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          {partner.name}
        </h1>
        {canPerformAction(currentUser?.role, "delete") && (
          <Button variant="destructive" onClick={handleDelete}>
            Eliminar
          </Button>
        )}
      </div>

      <Tabs defaultValue="info">
        <TabsList>
          <TabsTrigger value="info">Información</TabsTrigger>
          <TabsTrigger value="contacts">Contactos</TabsTrigger>
          <TabsTrigger value="agreements">Acuerdos</TabsTrigger>
          <TabsTrigger value="slas">SLAs</TabsTrigger>
        </TabsList>

        <TabsContent value="info">
          <PartnerForm initialData={partner} isEdit />
        </TabsContent>

        <TabsContent value="contacts">
          <div className="p-4 border-2 border-black">
            <p className="font-mono text-gray-500">Contactos - Próximamente</p>
          </div>
        </TabsContent>

        <TabsContent value="agreements">
          <div className="p-4 border-2 border-black">
            <p className="font-mono text-gray-500">Acuerdos - Próximamente</p>
          </div>
        </TabsContent>

        <TabsContent value="slas">
          <div className="p-4 border-2 border-black">
            <p className="font-mono text-gray-500">SLAs - Próximamente</p>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}
