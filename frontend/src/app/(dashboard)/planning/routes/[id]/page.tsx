"use client";

import { useRouter, useParams } from "next/navigation";
import { useRoute, useCreateRoute, useUpdateRoute, useDeleteRoute } from "@/hooks/use-routes";
import { useDailyPlans } from "@/hooks/use-daily-plans";
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
import { useRouteTypes } from "@/lib/route-types-config";

const routeSchema = z.object({
  type: z.string().min(1, "El tipo es requerido"),
  date: z.string().min(1, "La fecha es requerida"),
  daily_plan_id: z.string().min(1, "El plan diario es requerido"),
  status: z.string().min(1, "El estado es requerido"),
});

type RouteFormData = z.infer<typeof routeSchema>;

export default function RouteDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const isNew = id === "new";
  const { data: route, isLoading } = useRoute(isNew ? "" : id);
  const createRoute = useCreateRoute();
  const updateRoute = useUpdateRoute();
  const deleteRoute = useDeleteRoute();
  const { user: currentUser } = useAuthStore();
  const { types: routeTypes, isDropdown } = useRouteTypes();
  const { data: dailyPlans } = useDailyPlans();

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<RouteFormData>({
    resolver: zodResolver(routeSchema),
    defaultValues: isNew
      ? { type: "", date: "", daily_plan_id: "", status: "" }
      : {
          type: route?.type || "",
          date: route?.date || "",
          daily_plan_id: route?.daily_plan_id || "",
          status: route?.status || "",
        },
  });

  // Reset form when route data loads
  useEffect(() => {
    if (!isNew && route) {
      const formattedDate = route.date ? route.date.split('T')[0] : "";
      reset({
        type: route.type || "",
        date: formattedDate,
        daily_plan_id: route.daily_plan_id || "",
        status: route.status || "",
      });
    }
  }, [route, isNew, reset]);

  const onSubmit = async (data: RouteFormData) => {
    try {
      if (isNew) {
        await createRoute.mutateAsync(data);
        toast.success("Ruta creada");
      } else {
        await updateRoute.mutateAsync({ id, data });
        toast.success("Ruta actualizada");
      }
      router.push("/planning/routes");
    } catch (error) {
      toast.error(isNew ? "Error al crear ruta" : "Error al actualizar ruta");
    }
  };

  const handleDelete = async () => {
    if (confirm("¿Estás seguro de eliminar esta ruta?")) {
      try {
        await deleteRoute.mutateAsync(id);
        toast.success("Ruta eliminada");
        router.push("/planning/routes");
      } catch (error) {
        toast.error("Error al eliminar ruta");
      }
    }
  };

  if (!isNew && isLoading) {
    return <div className="font-mono">Cargando...</div>;
  }

  if (!isNew && !route) {
    return <div className="font-mono">Ruta no encontrada</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          {isNew ? "Crear Ruta" : "Editar Ruta"}
        </h1>
        {!isNew && canPerformAction(currentUser?.role, "delete") && (
          <Button variant="destructive" onClick={handleDelete}>
            Eliminar
          </Button>
        )}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Información de la Ruta</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="type">Tipo</Label>
                {isDropdown ? (
                  <select
                    id="type"
                    {...register("type")}
                    className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
                  >
                    {routeTypes.map(t => (
                      <option key={t.id} value={t.id}>{t.name}</option>
                    ))}
                  </select>
                ) : (
                  <Input
                    id="type"
                    {...register("type")}
                  />
                )}
                {errors.type && (
                  <p className="text-sm text-red-500">{errors.type.message}</p>
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
                <Label htmlFor="daily_plan_id">Plan Diario</Label>
                <select
                  id="daily_plan_id"
                  {...register("daily_plan_id")}
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                >
                  <option value="">Seleccionar plan diario</option>
                  {dailyPlans?.map((plan) => (
                    <option key={plan.id} value={plan.id}>
                      {plan.name} ({plan.date.split('T')[0]})
                    </option>
                  ))}
                </select>
                {errors.daily_plan_id && (
                  <p className="text-sm text-red-500">{errors.daily_plan_id.message}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="status">Estado</Label>
                <select
                  id="status"
                  {...register("status")}
                  className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
                >
                  <option value="scheduled">Programada</option>
                  <option value="in_progress">En progreso</option>
                  <option value="completed">Completada</option>
                  <option value="cancelled">Cancelada</option>
                </select>
                {errors.status && (
                  <p className="text-sm text-red-500">{errors.status.message}</p>
                )}
              </div>
            </div>

            <div className="flex justify-end gap-2">
              <Button
                type="button"
                variant="outline"
                onClick={() => router.push("/planning/routes")}
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
