"use client";

import { useRouter } from "next/navigation";
import { useCreateDailyPlanAssignment } from "@/hooks/use-daily-plan-assignments";
import { useDailyPlans } from "@/hooks/use-daily-plans";
import { useAllTechnicians } from "@/hooks/use-technicians";
import { useTechRoles } from "@/hooks/use-tech-roles";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

const assignmentSchema = z.object({
  daily_plan_id: z.string().min(1, "El plan diario es requerido"),
  tech_id: z.string().min(1, "El técnico es requerido"),
  role_id: z.string().min(1, "El rol es requerido"),
});

type AssignmentFormData = z.infer<typeof assignmentSchema>;

export default function NewAssignmentPage() {
  const router = useRouter();
  const createAssignment = useCreateDailyPlanAssignment();
  const { data: dailyPlans } = useDailyPlans();
  const { data: technicians } = useAllTechnicians();
  const { data: techRoles } = useTechRoles();

  const activeRoles = techRoles?.filter((r) => r.active) || [];
  const allTechnicians = technicians || [];

  const { register, handleSubmit, formState: { errors } } = useForm<AssignmentFormData>({
    resolver: zodResolver(assignmentSchema),
    defaultValues: { daily_plan_id: "", tech_id: "", role_id: "" },
  });

  const onSubmit = async (data: AssignmentFormData) => {
    try {
      await createAssignment.mutateAsync(data);
      toast.success("Asignación creada");
      router.push("/planning/assignments");
    } catch {
      toast.error("Error al crear asignación");
    }
  };

  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Crear Asignación
      </h1>
      <Card>
        <CardHeader>
          <CardTitle>Información de la Asignación</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="daily_plan_id">Plan Diario</Label>
                <select id="daily_plan_id" {...register("daily_plan_id")} className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50">
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
                <select id="tech_id" {...register("tech_id")} className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50">
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
                <select id="role_id" {...register("role_id")} className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50">
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
              <Button type="button" variant="outline" onClick={() => router.push("/planning/assignments")}>
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
