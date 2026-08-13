"use client";

import { useRouter, useParams } from "next/navigation";
import { useVehicle, useCreateVehicle, useUpdateVehicle, useDeleteVehicle } from "@/hooks/use-vehicles";
import { useVehicleTypes } from "@/hooks/use-vehicle-types";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useEffect } from "react";

const schema = z.object({
  license_plate: z.string().min(1, "Patente es requerida"),
  name: z.string().min(1, "Nombre es requerido"),
  brand: z.string().optional(),
  type: z.string().optional(),
  status: z.string().min(1, "Estado es requerido"),
});

type FormData = z.infer<typeof schema>;

export default function VehicleDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const isNew = id === "new";
  const { data, isLoading } = useVehicle(isNew ? "" : id);
  const createMutation = useCreateVehicle();
  const updateMutation = useUpdateVehicle();
  const deleteMutation = useDeleteVehicle();
  const { user } = useAuthStore();

  const { data: vehicleTypes } = useVehicleTypes();
  const { register, handleSubmit, formState: { errors }, reset } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: isNew
      ? { license_plate: "", name: "", brand: "", type: "", status: "" }
      : {
          license_plate: data?.license_plate || "",
          name: data?.name || "",
          brand: data?.brand || "",
          type: data?.type || "",
          status: data?.status || "",
        },
  });

  // Reset form when data loads
  useEffect(() => {
    if (!isNew && data) {
      reset({
        license_plate: data.license_plate || "",
        name: data.name || "",
        brand: data.brand || "",
        type: data.type || "",
        status: data.status || "",
      });
    }
  }, [data, isNew, reset]);

  const onSubmit = async (formData: FormData) => {
    try {
      const payload = {
        ...formData,
        type: formData.type || undefined,
      };
      if (isNew) {
        await createMutation.mutateAsync(payload);
        toast.success("Vehículo creado");
      } else {
        await updateMutation.mutateAsync({ id, data: payload });
        toast.success("Vehículo actualizado");
      }
      router.back();
    } catch {
      toast.error("Error al guardar");
    }
  };

  const handleDelete = async () => {
    if (confirm("¿Estás seguro de eliminar este vehículo?")) {
      try {
        await deleteMutation.mutateAsync(id);
        toast.success("Vehículo eliminado");
        router.back();
      } catch {
        toast.error("Error al eliminar");
      }
    }
  };

  if (!isNew && isLoading) return <div className="font-mono">Cargando...</div>;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          {isNew ? "Crear Vehículo" : "Editar Vehículo"}
        </h1>
        {!isNew && canPerformAction(user?.role, "delete") && (
          <Button variant="destructive" onClick={handleDelete}>
            Eliminar
          </Button>
        )}
      </div>
      <Card className="max-w-md">
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="license_plate">Patente</Label>
              <Input id="license_plate" {...register("license_plate")} />
              {errors.license_plate && <p className="text-sm text-red-500">{errors.license_plate.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="name">Nombre</Label>
              <Input id="name" {...register("name")} />
              {errors.name && <p className="text-sm text-red-500">{errors.name.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="brand">Marca</Label>
              <Input id="brand" {...register("brand")} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="type">Tipo</Label>
              <select
                id="type"
                {...register("type")}
                className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
              >
                <option value="">Seleccionar...</option>
                {vehicleTypes?.map((vt) => (
                  <option key={vt.id} value={vt.id}>{vt.name}</option>
                ))}
              </select>
              {errors.type && <p className="text-sm text-red-500">{errors.type.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="status">Estado</Label>
              <select
                id="status"
                {...register("status")}
                className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
              >
                <option value="">Seleccionar...</option>
                <option value="available">Disponible</option>
                <option value="in_use">En Uso</option>
                <option value="maintenance">Mantenimiento</option>
                <option value="retired">Retirado</option>
              </select>
              {errors.status && <p className="text-sm text-red-500">{errors.status.message}</p>}
            </div>
            <div className="flex gap-2">
              <Button type="submit">{isNew ? "Crear" : "Guardar"}</Button>
              <Button type="button" variant="outline" onClick={() => router.back()}>
                Cancelar
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
