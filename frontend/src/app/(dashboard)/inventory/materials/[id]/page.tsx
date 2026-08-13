"use client";

import { useRouter, useParams } from "next/navigation";
import { useMaterial, useCreateMaterial, useUpdateMaterial, useDeleteMaterial } from "@/hooks/use-materials";
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
  name: z.string().min(1, "Nombre es requerido"),
  sku: z.string().min(1, "SKU es requerido"),
  unit: z.string().min(1, "Unidad es requerida"),
  unit_cost: z.coerce.number().min(0, "Costo unitario debe ser ≥ 0"),
});

type FormData = z.infer<typeof schema>;

export default function MaterialDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const isNew = id === "new";
  const { data, isLoading } = useMaterial(isNew ? "" : id);
  const createMutation = useCreateMaterial();
  const updateMutation = useUpdateMaterial();
  const deleteMutation = useDeleteMaterial();
  const { user } = useAuthStore();

  const { register, handleSubmit, formState: { errors }, reset } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: isNew
      ? { name: "", sku: "", unit: "", unit_cost: 0 }
      : {
          name: data?.name || "",
          sku: data?.sku || "",
          unit: data?.unit || "",
          unit_cost: data?.unit_cost || 0,
        },
  });

  // Reset form when data loads
  useEffect(() => {
    if (!isNew && data) {
      reset({
        name: data.name || "",
        sku: data.sku || "",
        unit: data.unit || "",
        unit_cost: data.unit_cost || 0,
      });
    }
  }, [data, isNew, reset]);

  const onSubmit = async (formData: FormData) => {
    try {
      if (isNew) {
        await createMutation.mutateAsync(formData);
        toast.success("Material creado");
      } else {
        await updateMutation.mutateAsync({ id, data: formData });
        toast.success("Material actualizado");
      }
      router.back();
    } catch {
      toast.error("Error al guardar");
    }
  };

  const handleDelete = async () => {
    if (confirm("¿Estás seguro de eliminar este material?")) {
      try {
        await deleteMutation.mutateAsync(id);
        toast.success("Material eliminado");
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
          {isNew ? "Crear Material" : "Editar Material"}
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
              <Label htmlFor="name">Nombre</Label>
              <Input id="name" {...register("name")} />
              {errors.name && <p className="text-sm text-red-500">{errors.name.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="sku">SKU</Label>
              <Input id="sku" {...register("sku")} />
              {errors.sku && <p className="text-sm text-red-500">{errors.sku.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="unit">Unidad</Label>
              <Input id="unit" {...register("unit")} />
              {errors.unit && <p className="text-sm text-red-500">{errors.unit.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="unit_cost">Costo Unitario</Label>
              <Input id="unit_cost" type="number" step="0.01" {...register("unit_cost")} />
              {errors.unit_cost && <p className="text-sm text-red-500">{errors.unit_cost.message}</p>}
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
