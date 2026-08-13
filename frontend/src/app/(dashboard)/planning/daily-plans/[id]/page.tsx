"use client";

import { useRouter, useParams } from "next/navigation";
import { useDailyPlan, useCreateDailyPlan, useUpdateDailyPlan, useDeleteDailyPlan } from "@/hooks/use-daily-plans";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useEffect } from "react";

const dailyPlanSchema = z.object({
  name: z.string().min(1, "El nombre es requerido"),
  date: z.string().min(1, "La fecha es requerida"),
  notes: z.string().optional(),
});

type DailyPlanFormData = z.infer<typeof dailyPlanSchema>;

export default function DailyPlanDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const isNew = id === "new";
  const { data: dailyPlan, isLoading } = useDailyPlan(isNew ? "" : id);
  const createDailyPlan = useCreateDailyPlan();
  const updateDailyPlan = useUpdateDailyPlan();
  const deleteDailyPlan = useDeleteDailyPlan();
  const { user: currentUser } = useAuthStore();

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<DailyPlanFormData>({
    resolver: zodResolver(dailyPlanSchema),
    defaultValues: isNew
      ? { name: "", date: "", notes: "" }
      : {
          name: dailyPlan?.name || "",
          date: dailyPlan?.date || "",
          notes: dailyPlan?.notes || "",
        },
  });

  // Reset form when dailyPlan data loads
  useEffect(() => {
    if (!isNew && dailyPlan) {
      // Format date for HTML date input (YYYY-MM-DD)
      const formattedDate = dailyPlan.date ? dailyPlan.date.split('T')[0] : "";
      reset({
        name: dailyPlan.name || "",
        date: formattedDate,
        notes: dailyPlan.notes || "",
      });
    }
  }, [dailyPlan, isNew, reset]);

  const onSubmit = async (data: DailyPlanFormData) => {
    try {
      if (isNew) {
        await createDailyPlan.mutateAsync(data);
        toast.success("Plan diario creado");
      } else {
        await updateDailyPlan.mutateAsync({ id, data });
        toast.success("Plan diario actualizado");
      }
      router.push("/planning/daily-plans");
    } catch (error) {
      toast.error(isNew ? "Error al crear plan diario" : "Error al actualizar plan diario");
    }
  };

  const handleDelete = async () => {
    if (confirm("¿Estás seguro de eliminar este plan diario?")) {
      try {
        await deleteDailyPlan.mutateAsync(id);
        toast.success("Plan diario eliminado");
        router.push("/planning/daily-plans");
      } catch (error) {
        toast.error("Error al eliminar plan diario");
      }
    }
  };

  if (!isNew && isLoading) {
    return <div className="font-mono">Cargando...</div>;
  }

  if (!isNew && !dailyPlan) {
    return <div className="font-mono">Plan diario no encontrado</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          {isNew ? "Crear Plan Diario" : "Editar Plan Diario"}
        </h1>
        {!isNew && canPerformAction(currentUser?.role, "delete") && (
          <Button variant="destructive" onClick={handleDelete}>
            Eliminar
          </Button>
        )}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Información del Plan</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="name">Nombre</Label>
                <Input
                  id="name"
                  {...register("name")}
                />
                {errors.name && (
                  <p className="text-sm text-red-500">{errors.name.message}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="date">Fecha</Label>
                <Input
                  id="date"
                  type="date"
                  {...register("date")}
                />
                {errors.date && (
                  <p className="text-sm text-red-500">{errors.date.message}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="notes">Notas</Label>
                <Input
                  id="notes"
                  {...register("notes")}
                />
              </div>
            </div>

            <div className="flex justify-end gap-2">
              <Button
                type="button"
                variant="outline"
                onClick={() => router.push("/planning/daily-plans")}
              >
                Cancelar
              </Button>
              <Button type="submit">
                {isNew ? "Crear" : "Guardar"}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
