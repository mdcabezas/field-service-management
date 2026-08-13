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
import { useCreatePartner, useUpdatePartner } from "@/hooks/use-partners";
import { toast } from "sonner";
import type { Partner } from "@/types/api";

const partnerSchema = z.object({
  name: z.string().min(1, "Nombre requerido"),
  tax_id: z.string().optional(),
  status: z.string().optional(),
});

type PartnerFormData = z.infer<typeof partnerSchema>;

interface PartnerFormProps {
  initialData?: Partner;
  isEdit?: boolean;
}

export function PartnerForm({ initialData, isEdit = false }: PartnerFormProps) {
  const router = useRouter();
  const createPartner = useCreatePartner();
  const updatePartner = useUpdatePartner();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<PartnerFormData>({
    resolver: zodResolver(partnerSchema),
    defaultValues: initialData
      ? {
          name: initialData.name,
          tax_id: initialData.tax_id,
          status: initialData.status,
        }
      : undefined,
  });

  // Reset form when initialData changes (e.g., async data loads)
  useEffect(() => {
    if (isEdit && initialData) {
      reset({
        name: initialData.name,
        tax_id: initialData.tax_id,
        status: initialData.status,
      });
    }
  }, [initialData, isEdit, reset]);

  const onSubmit = async (data: PartnerFormData) => {
    try {
      if (isEdit && initialData) {
        await updatePartner.mutateAsync({
          id: initialData.id,
          data,
        });
        toast.success("Socio actualizado");
      } else {
        await createPartner.mutateAsync({
          name: data.name,
          tax_id: data.tax_id,
          status: data.status || "active",
        });
        toast.success("Socio creado");
      }
      router.push("/partners");
    } catch (error) {
      toast.error(isEdit ? "Error al actualizar" : "Error al crear");
    }
  };

  return (
    <Card className="max-w-md">
      <CardHeader>
        <CardTitle>{isEdit ? "Editar Socio" : "Crear Socio"}</CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">Nombre</Label>
            <Input id="name" {...register("name")} />
            {errors.name && (
              <p className="text-sm text-red-600">{errors.name.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="tax_id">RUT</Label>
            <Input id="tax_id" {...register("tax_id")} />
          </div>

          <div className="space-y-2">
            <Label htmlFor="status">Estado</Label>
            <select
              id="status"
              {...register("status")}
              className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
            >
              <option value="active">Activo</option>
              <option value="inactive">Inactivo</option>
            </select>
          </div>

          <div className="flex gap-2">
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? "Guardando..." : isEdit ? "Actualizar" : "Crear"}
            </Button>
            <Button
              type="button"
              variant="outline"
              onClick={() => router.push("/partners")}
            >
              Cancelar
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}