"use client";

import { useRouter, useParams } from "next/navigation";
import { useTool, useCreateTool, useUpdateTool, useDeleteTool } from "@/hooks/use-tools";
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
  code: z.string().min(1, "Código es requerido"),
  status: z.string().min(1, "Estado es requerido"),
});

type FormData = z.infer<typeof schema>;

export default function ToolDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const isNew = id === "new";
  const { data, isLoading } = useTool(isNew ? "" : id);
  const createMutation = useCreateTool();
  const updateMutation = useUpdateTool();
  const deleteMutation = useDeleteTool();
  const { user } = useAuthStore();

  const { register, handleSubmit, formState: { errors }, reset } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: isNew
      ? { name: "", code: "", status: "" }
      : {
          name: data?.name || "",
          code: data?.code || "",
          status: data?.status || "",
        },
  });

  // Reset form when data loads
  useEffect(() => {
    if (!isNew && data) {
      reset({
        name: data.name || "",
        code: data.code || "",
        status: data.status || "",
      });
    }
  }, [data, isNew, reset]);

  const onSubmit = async (formData: FormData) => {
    try {
      if (isNew) {
        await createMutation.mutateAsync(formData);
        toast.success("Herramienta creada");
      } else {
        await updateMutation.mutateAsync({ id, data: formData });
        toast.success("Herramienta actualizada");
      }
      router.back();
    } catch {
      toast.error("Error al guardar");
    }
  };

  const handleDelete = async () => {
    if (confirm("¿Estás seguro de eliminar esta herramienta?")) {
      try {
        await deleteMutation.mutateAsync(id);
        toast.success("Herramienta eliminada");
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
          {isNew ? "Crear Herramienta" : "Editar Herramienta"}
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
              <Label htmlFor="code">Código</Label>
              <Input id="code" {...register("code")} />
              {errors.code && <p className="text-sm text-red-500">{errors.code.message}</p>}
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
