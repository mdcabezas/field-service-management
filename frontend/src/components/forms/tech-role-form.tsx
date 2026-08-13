"use client";

import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useEffect } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useCreateTechRole, useUpdateTechRole } from "@/hooks/use-tech-roles";
import { toast } from "sonner";
import type { TechRole } from "@/types/api";

const techRoleSchema = z.object({
  code: z.string().min(1, "Código requerido"),
  name: z.string().min(1, "Nombre requerido"),
  active: z.boolean(),
});

type TechRoleFormData = z.infer<typeof techRoleSchema>;

interface TechRoleFormProps {
  initialData?: TechRole;
  isEdit?: boolean;
}

export function TechRoleForm({ initialData, isEdit = false }: TechRoleFormProps) {
  const router = useRouter();
  const createTechRole = useCreateTechRole();
  const updateTechRole = useUpdateTechRole();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<TechRoleFormData>({
    resolver: zodResolver(techRoleSchema),
    defaultValues: initialData
      ? {
          code: initialData.code,
          name: initialData.name,
          active: initialData.active,
        }
      : undefined,
  });

  // Reset form when initialData changes (e.g., async data loads)
  useEffect(() => {
    if (isEdit && initialData) {
      reset({
        code: initialData.code,
        name: initialData.name,
        active: initialData.active,
      });
    }
  }, [initialData, isEdit, reset]);

  const onSubmit = async (data: TechRoleFormData) => {
    try {
      if (isEdit && initialData) {
        await updateTechRole.mutateAsync({
          id: initialData.id,
          data: {
            code: data.code,
            name: data.name,
            active: data.active,
          },
        });
        toast.success("Rol actualizado");
      } else {
        await createTechRole.mutateAsync({
          code: data.code,
          name: data.name,
          active: data.active,
        });
        toast.success("Rol creado");
      }
      router.push("/core/tech-roles");
    } catch (error) {
      toast.error(isEdit ? "Error al actualizar" : "Error al crear");
    }
  };

  return (
    <Card className="max-w-md">
      <CardHeader>
        <CardTitle>{isEdit ? "Editar Rol" : "Crear Rol"}</CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="code">Código</Label>
            <Input id="code" {...register("code")} />
            {errors.code && (
              <p className="text-sm text-red-600">{errors.code.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="name">Nombre</Label>
            <Input id="name" {...register("name")} />
            {errors.name && (
              <p className="text-sm text-red-600">{errors.name.message}</p>
            )}
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
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? "Guardando..." : isEdit ? "Actualizar" : "Crear"}
            </Button>
            <Button
              type="button"
              variant="outline"
              onClick={() => router.push("/core/tech-roles")}
            >
              Cancelar
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}