"use client";

import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useCreateVisitSLATrackingGlobal } from "@/hooks/use-all-visit-sla-trackings";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent } from "@/components/ui/card";
import { toast } from "sonner";

const slaTrackingSchema = z.object({
  visit_id: z.string().min(1, "La visita es requerida"),
  sla_id: z.string().min(1, "El SLA es requerido"),
  requested_at: z.string().optional(),
  responded_at: z.string().optional(),
  resolved_at: z.string().optional(),
});

type SLATrackingFormData = z.infer<typeof slaTrackingSchema>;

export function SLATrackingForm() {
  const router = useRouter();
  const createSLATracking = useCreateVisitSLATrackingGlobal();

  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<SLATrackingFormData>({
    resolver: zodResolver(slaTrackingSchema),
    defaultValues: { visit_id: "", sla_id: "", requested_at: "", responded_at: "", resolved_at: "" },
  });

  const onSubmit = async (data: SLATrackingFormData) => {
    try {
      await createSLATracking.mutateAsync(data);
      toast.success("Seguimiento SLA creado");
      router.push("/operations/sla-trackings");
    } catch {
      toast.error("Error al crear seguimiento SLA");
    }
  };

  return (
    <Card className="max-w-md">
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="visit_id">Visita</Label>
            <Input id="visit_id" {...register("visit_id")} />
            {errors.visit_id && <p className="text-sm text-red-500">{errors.visit_id.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="sla_id">SLA</Label>
            <Input id="sla_id" {...register("sla_id")} />
            {errors.sla_id && <p className="text-sm text-red-500">{errors.sla_id.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="requested_at">Solicitado</Label>
            <Input id="requested_at" type="datetime-local" {...register("requested_at")} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="responded_at">Respondido</Label>
            <Input id="responded_at" type="datetime-local" {...register("responded_at")} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="resolved_at">Resuelto</Label>
            <Input id="resolved_at" type="datetime-local" {...register("resolved_at")} />
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
