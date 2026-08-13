"use client";

import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useCreateVisit } from "@/hooks/use-visits";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent } from "@/components/ui/card";
import { toast } from "sonner";

const visitSchema = z.object({
  scheduled_at: z.string().min(1, "La fecha es requerida"),
  priority: z.string().min(1, "La prioridad es requerida"),
  status: z.string().min(1, "El estado es requerido"),
  type: z.string().min(1, "El tipo es requerido"),
  billing_to: z.string().min(1, "El destino de facturación es requerido"),
  notes: z.string().optional(),
});

type VisitFormData = z.infer<typeof visitSchema>;

export function VisitForm() {
  const router = useRouter();
  const createVisit = useCreateVisit();

  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<VisitFormData>({
    resolver: zodResolver(visitSchema),
    defaultValues: { scheduled_at: "", priority: "medium", status: "pending", type: "", billing_to: "", notes: "" },
  });

  const onSubmit = async (data: VisitFormData) => {
    try {
      await createVisit.mutateAsync(data);
      toast.success("Visita creada");
      router.push("/operations/visits");
    } catch {
      toast.error("Error al crear visita");
    }
  };

  return (
    <Card className="max-w-md">
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="scheduled_at">Fecha Programada</Label>
            <Input id="scheduled_at" type="datetime-local" {...register("scheduled_at")} />
            {errors.scheduled_at && <p className="text-sm text-red-500">{errors.scheduled_at.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="priority">Prioridad</Label>
            <Input id="priority" {...register("priority")} />
            {errors.priority && <p className="text-sm text-red-500">{errors.priority.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="type">Tipo</Label>
            <Input id="type" {...register("type")} />
            {errors.type && <p className="text-sm text-red-500">{errors.type.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="status">Estado</Label>
            <Input id="status" {...register("status")} />
            {errors.status && <p className="text-sm text-red-500">{errors.status.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="billing_to">Destino de Facturación</Label>
            <Input id="billing_to" {...register("billing_to")} />
            {errors.billing_to && <p className="text-sm text-red-500">{errors.billing_to.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="notes">Notas</Label>
            <Input id="notes" {...register("notes")} />
          </div>
          <div className="flex gap-2">
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? "Creando..." : "Crear"}
            </Button>
            <Button type="button" variant="outline" onClick={() => router.back()}>
              Cancelar
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}
