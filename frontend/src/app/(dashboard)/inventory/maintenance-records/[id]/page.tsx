"use client";

import { useRouter, useParams } from "next/navigation";
import { useMaintenanceRecord, useCreateMaintenanceRecord, useUpdateMaintenanceRecord, useDeleteMaintenanceRecord } from "@/hooks/use-maintenance-records";
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
  date: z.string().min(1, "Fecha es requerida"),
  cost: z.coerce.number().min(0, "Costo debe ser ≥ 0").optional(),
  description: z.string().optional(),
});

type FormData = z.infer<typeof schema>;

export default function MaintenanceRecordDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const isNew = id === "new";
  const { data, isLoading } = useMaintenanceRecord(isNew ? "" : id);
  const createMutation = useCreateMaintenanceRecord();
  const updateMutation = useUpdateMaintenanceRecord();
  const deleteMutation = useDeleteMaintenanceRecord();
  const { user } = useAuthStore();

  const { register, handleSubmit, formState: { errors }, reset } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: isNew
      ? { reference_id: "", type: "", date: "", cost: 0, description: "" }
      : {
          reference_id: data?.reference_id || "",
          type: data?.type || "",
          date: data?.date || "",
          cost: data?.cost || 0,
          description: data?.description || "",
        },
  });

  // Reset form when data loads
  useEffect(() => {
    if (!isNew && data) {
      const formattedDate = data.date ? data.date.split('T')[0] : "";
      reset({
        reference_id: data.reference_id || "",
        type: data.type || "",
        date: formattedDate,
        cost: data.cost || 0,
        description: data.description || "",
      });
    }
  }, [data, isNew, reset]);

  const onSubmit = async (formData: FormData) => {
    try {
      const payload = {
        ...formData,
        date: formData.date ? new Date(formData.date + 'T00:00:00Z').toISOString() : formData.date,
      };
      if (isNew) {
        await createMutation.mutateAsync(payload);
        toast.success("Registro creado");
      } else {
        await updateMutation.mutateAsync({ id, data: payload });
        toast.success("Registro actualizado");
      }
      router.back();
    } catch {
      toast.error("Error al guardar");
    }
  };

  const handleDelete = async () => {
    if (confirm("¿Estás seguro de eliminar este registro?")) {
      try {
        await deleteMutation.mutateAsync(id);
        toast.success("Registro eliminado");
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
          {isNew ? "Crear Registro" : "Editar Registro"}
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
              <Label htmlFor="date">Fecha</Label>
              <Input id="date" type="date" {...register("date")} />
              {errors.date && <p className="text-sm text-red-500">{errors.date.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="cost">Costo</Label>
              <Input id="cost" type="number" step="0.01" {...register("cost")} />
              {errors.cost && <p className="text-sm text-red-500">{errors.cost.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="description">Descripción</Label>
              <Input id="description" {...register("description")} />
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
