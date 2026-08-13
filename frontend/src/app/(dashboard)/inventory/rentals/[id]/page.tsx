"use client";

import { useRouter, useParams } from "next/navigation";
import { useRental, useCreateRental, useUpdateRental, useDeleteRental } from "@/hooks/use-rentals";
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
  supplier: z.string().optional(),
  item_description: z.string().optional(),
  daily_cost: z.coerce.number().min(0, "Costo diario debe ser ≥ 0").optional(),
  notes: z.string().optional(),
});

type FormData = z.infer<typeof schema>;

export default function RentalDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const isNew = id === "new";
  const { data, isLoading } = useRental(isNew ? "" : id);
  const createMutation = useCreateRental();
  const updateMutation = useUpdateRental();
  const deleteMutation = useDeleteRental();
  const { user } = useAuthStore();

  const { register, handleSubmit, formState: { errors }, reset } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: isNew
      ? { type: "", supplier: "", item_description: "", daily_cost: 0, notes: "" }
      : {
          type: data?.type || "",
          supplier: data?.supplier || "",
          item_description: data?.item_description || "",
          daily_cost: data?.daily_cost || 0,
          notes: data?.notes || "",
        },
  });

  // Reset form when data loads
  useEffect(() => {
    if (!isNew && data) {
      reset({
        type: data.type || "",
        supplier: data.supplier || "",
        item_description: data.item_description || "",
        daily_cost: data.daily_cost || 0,
        notes: data.notes || "",
      });
    }
  }, [data, isNew, reset]);

  const onSubmit = async (formData: FormData) => {
    try {
      if (isNew) {
        await createMutation.mutateAsync(formData);
        toast.success("Arriendo creado");
      } else {
        await updateMutation.mutateAsync({ id, data: formData });
        toast.success("Arriendo actualizado");
      }
      router.back();
    } catch {
      toast.error("Error al guardar");
    }
  };

  const handleDelete = async () => {
    if (confirm("¿Estás seguro de eliminar este arriendo?")) {
      try {
        await deleteMutation.mutateAsync(id);
        toast.success("Arriendo eliminado");
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
          {isNew ? "Crear Arriendo" : "Editar Arriendo"}
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
                <option value="vehicle">Vehículo</option>
                <option value="tool">Herramienta</option>
                <option value="equipment">Equipo</option>
              </select>
              {errors.type && <p className="text-sm text-red-500">{errors.type.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="supplier">Proveedor</Label>
              <Input id="supplier" {...register("supplier")} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="item_description">Descripción del Item</Label>
              <Input id="item_description" {...register("item_description")} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="daily_cost">Costo Diario</Label>
              <Input id="daily_cost" type="number" step="0.01" {...register("daily_cost")} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="notes">Notas</Label>
              <Input id="notes" {...register("notes")} />
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
