"use client";

import { useRouter, useParams } from "next/navigation";
import { useCostRate, useCreateCostRate, useUpdateCostRate, useDeleteCostRate } from "@/hooks/use-cost-rates";
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
  type: z.string().min(1, "Tipo es requerido"),
  value: z.coerce.number().min(0, "Valor debe ser ≥ 0"),
  unit: z.string().min(1, "Unidad es requerida"),
  reference_type: z.string().optional(),
  valid_from: z.string().optional(),
  valid_until: z.string().optional(),
});

type FormData = z.infer<typeof schema>;

export default function CostRateDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const isNew = id === "new";
  const { data, isLoading } = useCostRate(isNew ? "" : id);
  const createMutation = useCreateCostRate();
  const updateMutation = useUpdateCostRate();
  const deleteMutation = useDeleteCostRate();
  const { user } = useAuthStore();

  const { register, handleSubmit, formState: { errors }, reset } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: isNew
      ? { type: "", value: 0, unit: "", reference_type: "", valid_from: "", valid_until: "" }
      : {
          type: data?.type || "",
          value: data?.value || 0,
          unit: data?.unit || "",
          reference_type: data?.reference_type || "",
          valid_from: data?.valid_from || "",
          valid_until: data?.valid_until || "",
        },
  });

  // Reset form when data loads
  useEffect(() => {
    if (!isNew && data) {
      reset({
        type: data.type || "",
        value: data.value || 0,
        unit: data.unit || "",
        reference_type: data.reference_type || "",
        valid_from: data.valid_from || "",
        valid_until: data.valid_until || "",
      });
    }
  }, [data, isNew, reset]);

  const onSubmit = async (formData: FormData) => {
    try {
      if (isNew) {
        await createMutation.mutateAsync(formData);
        toast.success("Tarifa creada");
      } else {
        await updateMutation.mutateAsync({ id, data: formData });
        toast.success("Tarifa actualizada");
      }
      router.back();
    } catch {
      toast.error("Error al guardar");
    }
  };

  const handleDelete = async () => {
    if (confirm("¿Estás seguro de eliminar esta tarifa?")) {
      try {
        await deleteMutation.mutateAsync(id);
        toast.success("Tarifa eliminada");
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
          {isNew ? "Crear Tarifa" : "Editar Tarifa"}
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
              <Label htmlFor="type">Tipo</Label>
              <select
                id="type"
                {...register("type")}
                className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
              >
                <option value="">Seleccionar...</option>
                <option value="labor">Mano de Obra</option>
                <option value="vehicle">Vehículo</option>
                <option value="tool_depreciation">Depreciación Herramienta</option>
                <option value="epp">EPP</option>
                <option value="overhead">Gastos Generales</option>
              </select>
              {errors.type && <p className="text-sm text-red-500">{errors.type.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="value">Valor</Label>
              <Input id="value" type="number" step="0.01" {...register("value")} />
              {errors.value && <p className="text-sm text-red-500">{errors.value.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="unit">Unidad</Label>
              <select
                id="unit"
                {...register("unit")}
                className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
              >
                <option value="">Seleccionar...</option>
                <option value="hour">Hora</option>
                <option value="km">Kilómetro</option>
                <option value="day">Día</option>
                <option value="visit">Visita</option>
                <option value="unit">Unidad</option>
              </select>
              {errors.unit && <p className="text-sm text-red-500">{errors.unit.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="reference_type">Tipo de Referencia</Label>
              <select
                id="reference_type"
                {...register("reference_type")}
                className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
              >
                <option value="">Ninguno</option>
                <option value="vehicle">Vehículo</option>
                <option value="tool">Herramienta</option>
                <option value="equipment">Equipo</option>
              </select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="valid_from">Vigente desde</Label>
              <Input id="valid_from" type="date" {...register("valid_from")} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="valid_until">Vigente hasta</Label>
              <Input id="valid_until" type="date" {...register("valid_until")} />
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
