"use client";

import { useRouter } from "next/navigation";
import { useCreateDailyPlan } from "@/hooks/use-daily-plans";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

const dailyPlanSchema = z.object({
  name: z.string().min(1, "El nombre es requerido"),
  date: z.string().min(1, "La fecha es requerida"),
  notes: z.string().optional(),
});

type DailyPlanFormData = z.infer<typeof dailyPlanSchema>;

export default function NewDailyPlanPage() {
  const router = useRouter();
  const createDailyPlan = useCreateDailyPlan();

  const { register, handleSubmit, formState: { errors } } = useForm<DailyPlanFormData>({
    resolver: zodResolver(dailyPlanSchema),
    defaultValues: { name: "", date: "", notes: "" },
  });

  const onSubmit = async (data: DailyPlanFormData) => {
    try {
      const payload = {
        ...data,
        date: data.date ? `${data.date}T00:00:00Z` : data.date,
      };
      await createDailyPlan.mutateAsync(payload);
      toast.success("Plan diario creado");
      router.push("/planning/daily-plans");
    } catch {
      toast.error("Error al crear plan diario");
    }
  };

  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Crear Plan Diario
      </h1>
      <Card>
        <CardHeader>
          <CardTitle>Información del Plan</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="name">Nombre</Label>
                <Input id="name" {...register("name")} />
                {errors.name && <p className="text-sm text-red-500">{errors.name.message}</p>}
              </div>
              <div className="space-y-2">
                <Label htmlFor="date">Fecha</Label>
                <Input id="date" type="date" {...register("date")} />
                {errors.date && <p className="text-sm text-red-500">{errors.date.message}</p>}
              </div>
            </div>
            <div className="grid gap-4">
              <div className="space-y-2">
                <Label htmlFor="notes">Notas</Label>
                <Input id="notes" {...register("notes")} />
              </div>
            </div>
            <div className="flex justify-end gap-2">
              <Button type="button" variant="outline" onClick={() => router.push("/planning/daily-plans")}>
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
