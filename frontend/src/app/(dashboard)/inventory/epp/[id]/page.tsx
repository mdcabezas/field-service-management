"use client";

import { useRouter, useParams } from "next/navigation";
import { useEPPItem, useCreateEPPItem, useUpdateEPPItem, useDeleteEPPItem } from "@/hooks/use-epp-items";
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
  type: z.string().min(1, "Tipo es requerido"),
  lifecycle: z.string().min(1, "Ciclo de vida es requerido"),
});

type FormData = z.infer<typeof schema>;

export default function EPPDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const isNew = id === "new";
  const { data, isLoading } = useEPPItem(isNew ? "" : id);
  const createMutation = useCreateEPPItem();
  const updateMutation = useUpdateEPPItem();
  const deleteMutation = useDeleteEPPItem();
  const { user } = useAuthStore();

  const { register, handleSubmit, formState: { errors }, reset } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: isNew
      ? { name: "", type: "", lifecycle: "" }
      : {
          name: data?.name || "",
          type: data?.type || "",
          lifecycle: data?.lifecycle || "",
        },
  });

  // Reset form when data loads
  useEffect(() => {
    if (!isNew && data) {
      reset({
        name: data.name || "",
        type: data.type || "",
        lifecycle: data.lifecycle || "",
      });
    }
  }, [data, isNew, reset]);

  const onSubmit = async (formData: FormData) => {
    try {
      if (isNew) {
        await createMutation.mutateAsync(formData);
        toast.success("EPP creado");
      } else {
        await updateMutation.mutateAsync({ id, data: formData });
        toast.success("EPP actualizado");
      }
      router.back();
    } catch {
      toast.error("Error al guardar");
    }
  };

  const handleDelete = async () => {
    if (confirm("¿Estás seguro de eliminar este EPP?")) {
      try {
        await deleteMutation.mutateAsync(id);
        toast.success("EPP eliminado");
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
          {isNew ? "Crear EPP" : "Editar EPP"}
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
              <Label htmlFor="type">Tipo</Label>
              <select
                id="type"
                {...register("type")}
                className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
              >
                <option value="">Seleccionar...</option>
                <option value="head">Cabeza</option>
                <option value="hands">Manos</option>
                <option value="body">Cuerpo</option>
                <option value="eyes">Ojos</option>
                <option value="ears">Oídos</option>
                <option value="feet">Pies</option>
                <option value="respiratory">Respiratorio</option>
              </select>
              {errors.type && <p className="text-sm text-red-500">{errors.type.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="lifecycle">Ciclo de Vida</Label>
              <select
                id="lifecycle"
                {...register("lifecycle")}
                className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
              >
                <option value="">Seleccionar...</option>
                <option value="disposable">Desechable</option>
                <option value="reusable">Reutilizable</option>
              </select>
              {errors.lifecycle && <p className="text-sm text-red-500">{errors.lifecycle.message}</p>}
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
