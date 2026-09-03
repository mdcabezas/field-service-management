"use client";

import { useRouter, useParams } from "next/navigation";
import { useChecklistTemplate, useCreateChecklistTemplate, useUpdateChecklistTemplate, useDeleteChecklistTemplate } from "@/hooks/use-checklist-templates";
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
  description: z.string().optional(),
  visit_type: z.string().optional(),
  work_type: z.string().uuid().optional().or(z.literal("")),
});

type FormData = z.infer<typeof schema>;

export default function ChecklistTemplateDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const isNew = id === "new";
  const { data, isLoading } = useChecklistTemplate(isNew ? "" : id);
  const createMutation = useCreateChecklistTemplate();
  const updateMutation = useUpdateChecklistTemplate();
  const deleteMutation = useDeleteChecklistTemplate();
  const { user } = useAuthStore();

  const { register, handleSubmit, formState: { errors }, reset } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: isNew
      ? { name: "", description: "", visit_type: "", work_type: "" }
      : {
          name: data?.name || "",
          description: data?.description || "",
          visit_type: data?.visit_type || "",
          work_type: data?.work_type || "",
        },
  });

  // Reset form when data loads
  useEffect(() => {
    if (!isNew && data) {
      reset({
        name: data.name || "",
        description: data.description || "",
        visit_type: data.visit_type || "",
        work_type: data.work_type || "",
      });
    }
  }, [data, isNew, reset]);

  const onSubmit = async (formData: FormData) => {
    try {
      // Convert empty strings to null for UUID fields
      const submitData = {
        ...formData,
        work_type: formData.work_type === "" ? null : formData.work_type,
        visit_type: formData.visit_type === "" ? null : formData.visit_type,
      };
      if (isNew) {
        await createMutation.mutateAsync(submitData);
        toast.success("Plantilla creada");
      } else {
        await updateMutation.mutateAsync({ id, data: submitData });
        toast.success("Plantilla actualizada");
      }
      router.push("/inventory/checklist-templates");
    } catch {
      toast.error("Error al guardar");
    }
  };

  const handleDelete = async () => {
    if (confirm("¿Estás seguro de eliminar esta plantilla?")) {
      try {
        await deleteMutation.mutateAsync(id);
        toast.success("Plantilla eliminada");
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
          {isNew ? "Crear Plantilla" : "Editar Plantilla"}
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
              <Label htmlFor="description">Descripción</Label>
              <Input id="description" {...register("description")} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="visit_type">Tipo de Visita</Label>
              <Input id="visit_type" {...register("visit_type")} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="work_type">Tipo de Trabajo</Label>
              <Input id="work_type" {...register("work_type")} />
              {errors.work_type && <p className="text-sm text-red-500">{errors.work_type.message}</p>}
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
