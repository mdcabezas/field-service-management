"use client";

import { useRouter, useParams } from "next/navigation";
import { useMaintenanceSchedule, useCreateMaintenanceSchedule, useUpdateMaintenanceSchedule, useDeleteMaintenanceSchedule } from "@/hooks/use-maintenance-schedules";
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
  reference_id: z.string().min(1, "ID de referencia es requerido"),
  type: z.string().min(1, "Tipo es requerido"),
  frequency_days: z.coerce.number().min(1, "Frecuencia en días es requerida").optional(),
  active: z.boolean(),
});

type FormData = z.infer<typeof schema>;

export default function MaintenanceScheduleDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const isNew = id === "new";
  const { data, isLoading } = useMaintenanceSchedule(isNew ? "" : id);
  const createMutation = useCreateMaintenanceSchedule();
  const updateMutation = useUpdateMaintenanceSchedule();
  const deleteMutation = useDeleteMaintenanceSchedule();
  const { user } = useAuthStore();

  const { register, handleSubmit, formState: { errors }, reset } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: isNew
      ? { reference_id: "", type: "", frequency_days: 0, active: true }
      : {
          reference_id: data?.reference_id || "",
          type: data?.type || "",
          frequency_days: data?.frequency_days || 0,
          active: data?.active || true,
        },
  });

  // Reset form when data loads
  useEffect(() => {
    if (!isNew && data) {
      reset({
        reference_id: data.reference_id || "",
        type: data.type || "",
        frequency_days: data.frequency_days || 0,
        active: data.active ?? true,
      });
    }
  }, [data, isNew, reset]);

  const onSubmit = async (formData: FormData) => {
    try {
      if (isNew) {
        await createMutation.mutateAsync(formData);
        toast.success("Programación creada");
      } else {
        await updateMutation.mutateAsync({ id, data: formData });
        toast.success("Programación actualizada");
      }
      router.back();
    } catch {
      toast.error("Error al guardar");
    }
  };

  const handleDelete = async () => {
    if (confirm("¿Estás seguro de eliminar esta programación?")) {
      try {
        await deleteMutation.mutateAsync(id);
        toast.success("Programación eliminada");
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
          {isNew ? "Crear Programación" : "Editar Programación"}
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
              <Label htmlFor="reference_id">ID de Referencia</Label>
              <Input id="reference_id" {...register("reference_id")} />
              {errors.reference_id && <p className="text-sm text-red-500">{errors.reference_id.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="type">Tipo</Label>
              <select
                id="type"
                {...register("type")}
                className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
              >
                <option value="">Seleccionar...</option>
                <option value="vehicle">Vehículo</option>
                <option value="tool">Herramienta</option>
                <option value="equipment">Equipo</option>
              </select>
              {errors.type && <p className="text-sm text-red-500">{errors.type.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="frequency_days">Frecuencia (días)</Label>
              <Input id="frequency_days" type="number" {...register("frequency_days")} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="active">Activo</Label>
              <input
                id="active"
                type="checkbox"
                {...register("active")}
                className="h-4 w-4"
              />
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
