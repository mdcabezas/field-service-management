"use client";

import { useRouter } from "next/navigation";
import { useCreateRoute } from "@/hooks/use-routes";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useRouteTypes } from "@/lib/route-types-config";

const routeSchema = z.object({
  type: z.string().min(1, "El tipo es requerido"),
  date: z.string().min(1, "La fecha es requerida"),
  daily_plan_id: z.string().min(1, "El plan diario es requerido"),
  status: z.string().min(1, "El estado es requerido"),
});

type RouteFormData = z.infer<typeof routeSchema>;

export default function NewRoutePage() {
  const router = useRouter();
  const createRoute = useCreateRoute();
  const { types: routeTypes, isDropdown } = useRouteTypes();
  const { register, handleSubmit, formState: { errors } } = useForm<RouteFormData>({
    resolver: zodResolver(routeSchema),
    defaultValues: { type: "", date: "", daily_plan_id: "", status: "" },
  });

  const onSubmit = async (data: RouteFormData) => {
    try {
      await createRoute.mutateAsync(data);
      toast.success("Ruta creada");
      router.push("/planning/routes");
    } catch {
      toast.error("Error al crear ruta");
    }
  };

  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Crear Ruta
      </h1>
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
                  <Input id="type" {...register("type")} />
                )}
                {errors.type && <p className="text-sm text-red-500">{errors.type.message}</p>}
              </div>
              <div className="space-y-2">
                <Label htmlFor="date">Fecha</Label>
                <Input id="date" type="date" {...register("date")} />
                {errors.date && <p className="text-sm text-red-500">{errors.date.message}</p>}
              </div>
              <div className="space-y-2">
                <Label htmlFor="daily_plan_id">ID Plan Diario</Label>
                <Input id="daily_plan_id" {...register("daily_plan_id")} />
                {errors.daily_plan_id && <p className="text-sm text-red-500">{errors.daily_plan_id.message}</p>}
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
                {errors.status && <p className="text-sm text-red-500">{errors.status.message}</p>}
              </div>
            </div>
            <div className="flex justify-end gap-2">
              <Button type="button" variant="outline" onClick={() => router.push("/planning/routes")}>
                Cancelar
              </Button>
              <Button type="submit">Crear</Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
