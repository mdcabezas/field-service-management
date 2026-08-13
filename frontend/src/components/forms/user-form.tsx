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
import { useCreateUser, useUpdateUser } from "@/hooks/use-users";
import { toast } from "sonner";
import type { User } from "@/types/api";

const userSchema = z.object({
  email: z.string().email("Email inválido").min(1, "Email requerido"),
  name: z.string().min(1, "Nombre requerido"),
  password: z.string().min(6, "Mínimo 6 caracteres").optional(),
  role: z.enum(["admin", "manager", "supervisor", "technician"]),
});

type UserFormData = z.infer<typeof userSchema>;

interface UserFormProps {
  initialData?: User;
  isEdit?: boolean;
}

export function UserForm({ initialData, isEdit = false }: UserFormProps) {
  const router = useRouter();
  const createUser = useCreateUser();
  const updateUser = useUpdateUser();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<UserFormData>({
    resolver: zodResolver(userSchema),
    defaultValues: initialData
      ? {
          email: initialData.email,
          name: initialData.name,
          role: initialData.role,
        }
      : undefined,
  });

  // Reset form when initialData changes (e.g., async data loads)
  useEffect(() => {
    if (isEdit && initialData) {
      reset({
        email: initialData.email,
        name: initialData.name,
        role: initialData.role,
      });
    }
  }, [initialData, isEdit, reset]);

  const onSubmit = async (data: UserFormData) => {
    try {
      if (isEdit && initialData) {
        await updateUser.mutateAsync({
          id: initialData.id,
          data: {
            email: data.email,
            name: data.name,
            role: data.role,
          },
        });
        toast.success("Usuario actualizado");
      } else {
        await createUser.mutateAsync({
          email: data.email,
          name: data.name,
          role: data.role,
          password: data.password || "",
        });
        toast.success("Usuario creado");
      }
      router.push("/core/users");
    } catch (error) {
      toast.error(isEdit ? "Error al actualizar" : "Error al crear");
    }
  };

  return (
    <Card className="max-w-md">
      <CardHeader>
        <CardTitle>{isEdit ? "Editar Usuario" : "Crear Usuario"}</CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              type="email"
              {...register("email")}
              disabled={isEdit}
            />
            {errors.email && (
              <p className="text-sm text-red-600">{errors.email.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="name">Nombre</Label>
            <Input id="name" {...register("name")} />
            {errors.name && (
              <p className="text-sm text-red-600">{errors.name.message}</p>
            )}
          </div>

          {!isEdit && (
            <div className="space-y-2">
              <Label htmlFor="password">Contraseña</Label>
              <Input id="password" type="password" {...register("password")} />
              {errors.password && (
                <p className="text-sm text-red-600">{errors.password.message}</p>
              )}
            </div>
          )}

          <div className="space-y-2">
            <Label htmlFor="role">Rol</Label>
            <select
              id="role"
              {...register("role")}
              className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
            >
              <option value="admin">Admin</option>
              <option value="manager">Manager</option>
              <option value="supervisor">Supervisor</option>
              <option value="technician">Technician</option>
            </select>
            {errors.role && (
              <p className="text-sm text-red-600">{errors.role.message}</p>
            )}
          </div>

          <div className="flex gap-2">
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? "Guardando..." : isEdit ? "Actualizar" : "Crear"}
            </Button>
            <Button
              type="button"
              variant="outline"
              onClick={() => router.push("/core/users")}
            >
              Cancelar
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}