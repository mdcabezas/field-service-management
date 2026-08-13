"use client";

import { useRouter, useParams } from "next/navigation";
import { useAllDailyPlanAssignments, useCreateDailyPlanAssignment, useUpdateDailyPlanAssignment, useDeleteDailyPlanAssignment } from "@/hooks/use-daily-plan-assignments";
import { useDailyPlans } from "@/hooks/use-daily-plans";
import { useAllTechnicians } from "@/hooks/use-technicians";
import { useTechRoles } from "@/hooks/use-tech-roles";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useEffect } from "react";

const assignmentSchema = z.object({
  daily_plan_id: z.string().min(1, "El plan diario es requerido"),
  tech_id: z.string().min(1, "El técnico es requerido"),
  role_id: z.string().min(1, "El rol es requerido"),
});

type AssignmentFormData = z.infer<typeof assignmentSchema>;

export default function AssignmentDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const isNew = id === "new";
  const { data: assignments } = useAllDailyPlanAssignments();
  const { data: dailyPlans } = useDailyPlans();
  const { data: technicians } = useAllTechnicians();
  const { data: techRoles } = useTechRoles();
  const createAssignment = useCreateDailyPlanAssignment();
  const updateAssignment = useUpdateDailyPlanAssignment();
  const deleteAssignment = useDeleteDailyPlanAssignment();
  const { user: currentUser } = useAuthStore();

  const activeRoles = techRoles?.filter((r) => r.active) || [];
  const allTechnicians = technicians || [];

  const assignment = assignments?.find((a) => a.id === id);

  const { register, handleSubmit, formState: { errors }, reset } = useForm<AssignmentFormData>({
    resolver: zodResolver(assignmentSchema),
    defaultValues: isNew
      ? { daily_plan_id: "", tech_id: "", role_id: "" }
      : {
          daily_plan_id: assignment?.daily_plan_id || "",
          tech_id: assignment?.tech_id || "",
          role_id: assignment?.role_id || "",
        },
  });

  // Reset form when assignment data loads
  useEffect(() => {
    if (!isNew && assignment) {
      reset({
        daily_plan_id: assignment.daily_plan_id || "",
        tech_id: assignment.tech_id || "",
        role_id: assignment.role_id || "",
      });
    }
  }, [assignment, isNew, reset]);

  const onSubmit = async (data: AssignmentFormData) => {
    try {
      if (isNew) {
        await createAssignment.mutateAsync(data);
        toast.success("Asignación creada");
      } else {
        await updateAssignment.mutateAsync({ id, data });
        toast.success("Asignación actualizada");
      }
      router.push("/planning/assignments");
    } catch (error) {
      toast.error(isNew ? "Error al crear asignación" : "Error al actualizar asignación");
    }
  };

  const handleDelete = async () => {
    if (confirm("¿Estás seguro de eliminar esta asignación?")) {
      try {
        await deleteAssignment.mutateAsync(id);
        toast.success("Asignación eliminada");
        router.push("/planning/assignments");
      } catch (error) {
        toast.error("Error al eliminar asignación");
      }
    }
  };

  if (!isNew && !assignment) {
    return <div className="font-mono">Asignación no encontrada</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          {isNew ? "Crear Asignación" : "Detalle de Asignación"}
        </h1>
        {!isNew && canPerformAction(currentUser?.role, "delete") && (
          <Button variant="destructive" onClick={handleDelete}>
            Eliminar
          </Button>
        )}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Información de la Asignación</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="daily_plan_id">Plan Diario</Label>
                <select
                  id="daily_plan_id"
                  {...register("daily_plan_id")}
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  disabled={!isNew}
                >
                  <option value="">Seleccionar plan diario</option>
                  {dailyPlans?.map((plan) => (
                    <option key={plan.id} value={plan.id}>
                      {plan.date} - {plan.notes || "Sin notas"}
                    </option>
                  ))}
                </select>
                {errors.daily_plan_id && <p className="text-sm text-red-500">{errors.daily_plan_id.message}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="tech_id">Técnico</Label>
                <select
                  id="tech_id"
                  {...register("tech_id")}
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                >
                  <option value="">Seleccionar técnico</option>
                  {allTechnicians.map((tech) => (
                    <option key={tech.id} value={tech.id}>
                      {tech.name}
                    </option>
                  ))}
                </select>
                {errors.tech_id && <p className="text-sm text-red-500">{errors.tech_id.message}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="role_id">Rol</Label>
                <select
                  id="role_id"
                  {...register("role_id")}
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                >
                  <option value="">Seleccionar rol</option>
                  {activeRoles.map((role) => (
                    <option key={role.id} value={role.id}>
                      {role.name}
                    </option>
                  ))}
                </select>
                {errors.role_id && <p className="text-sm text-red-500">{errors.role_id.message}</p>}
              </div>
            </div>

            <div className="flex justify-end gap-2">
              <Button
                type="button"
                variant="outline"
                onClick={() => router.push("/planning/assignments")}
              >
                Cancelar
              </Button>
              {isNew && (
                <Button type="submit">
                  Crear
                </Button>
              )}
              {!isNew && canPerformAction(currentUser?.role, "edit") && (
                <Button type="submit">
                  Guardar
                </Button>
              )}
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}