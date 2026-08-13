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
import { useCreateCustomer, useUpdateCustomer } from "@/hooks/use-customers";
import { toast } from "sonner";
import type { Customer } from "@/types/api";

const customerSchema = z.object({
  name: z.string().min(1, "Nombre requerido"),
  tax_id: z.string().optional(),
  phone: z.string().optional(),
  email: z.string().email("Email inválido").optional(),
});

type CustomerFormData = z.infer<typeof customerSchema>;

interface CustomerFormProps {
  initialData?: Customer;
  isEdit?: boolean;
}

export function CustomerForm({ initialData, isEdit = false }: CustomerFormProps) {
  const router = useRouter();
  const createCustomer = useCreateCustomer();
  const updateCustomer = useUpdateCustomer();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<CustomerFormData>({
    resolver: zodResolver(customerSchema),
    defaultValues: initialData
      ? {
          name: initialData.name,
          tax_id: initialData.tax_id,
          phone: initialData.phone,
          email: initialData.email,
        }
      : undefined,
  });

  // Reset form when initialData changes (e.g., async data loads)
  useEffect(() => {
    if (isEdit && initialData) {
      reset({
        name: initialData.name,
        tax_id: initialData.tax_id,
        phone: initialData.phone,
        email: initialData.email,
      });
    }
  }, [initialData, isEdit, reset]);

  const onSubmit = async (data: CustomerFormData) => {
    try {
      if (isEdit && initialData) {
        await updateCustomer.mutateAsync({
          id: initialData.id,
          data,
        });
        toast.success("Cliente actualizado");
      } else {
        await createCustomer.mutateAsync({
          name: data.name,
          tax_id: data.tax_id,
          phone: data.phone || "",
          email: data.email || "",
        });
        toast.success("Cliente creado");
      }
      router.push("/customers");
    } catch (error) {
      toast.error(isEdit ? "Error al actualizar" : "Error al crear");
    }
  };

  return (
    <Card className="max-w-md">
      <CardHeader>
        <CardTitle>{isEdit ? "Editar Cliente" : "Crear Cliente"}</CardTitle>
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
            <Label htmlFor="phone">Teléfono</Label>
            <Input id="phone" {...register("phone")} />
          </div>

          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <Input id="email" type="email" {...register("email")} />
            {errors.email && (
              <p className="text-sm text-red-600">{errors.email.message}</p>
            )}
          </div>

          <div className="flex gap-2">
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? "Guardando..." : isEdit ? "Actualizar" : "Crear"}
            </Button>
            <Button
              type="button"
              variant="outline"
              onClick={() => router.push("/customers")}
            >
              Cancelar
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}
