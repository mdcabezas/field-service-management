"use client";

import { useRouter, useParams } from "next/navigation";
import { useCustomer, useDeleteCustomer } from "@/hooks/use-customers";
import { CustomerForm } from "@/components/forms/customer-form";
import { CustomerAddressesTable } from "@/components/customer/customer-addresses-table";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";

export default function CustomerDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const { data: customer, isLoading } = useCustomer(id);
  const deleteCustomer = useDeleteCustomer();
  const { user: currentUser } = useAuthStore();

  const handleDelete = async () => {
    if (confirm("¿Estás seguro de eliminar este cliente?")) {
      try {
        await deleteCustomer.mutateAsync(id);
        toast.success("Cliente eliminado");
        router.push("/customers");
      } catch (error) {
        toast.error("Error al eliminar cliente");
      }
    }
  };

  if (isLoading) {
    return <div className="font-mono">Cargando...</div>;
  }

  if (!customer) {
    return <div className="font-mono">Cliente no encontrado</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          {customer.name}
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
          <TabsTrigger value="addresses">Direcciones</TabsTrigger>
        </TabsList>

        <TabsContent value="info">
          <CustomerForm initialData={customer} isEdit />
        </TabsContent>

        <TabsContent value="addresses">
          <CustomerAddressesTable customerId={id} />
        </TabsContent>
      </Tabs>
    </div>
  );
}
