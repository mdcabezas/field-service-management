"use client";

import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useCreateVisit } from "@/hooks/use-visits";
import { useVisitTypes } from "@/hooks/use-visit-types";
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
  source: z.string().min(1, "La fuente es requerida"),
  notes: z.string().optional(),
});

type VisitFormData = z.infer<typeof visitSchema>;

export function VisitForm() {
  const router = useRouter();
  const createVisit = useCreateVisit();
  const { data: visitTypes, isLoading: loadingVisitTypes } = useVisitTypes();

  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<VisitFormData>({
    resolver: zodResolver(visitSchema),
    defaultValues: {
      scheduled_at: "",
      priority: "medium",
      status: "pending",
      type: "",
      billing_to: "partner",
      source: "portal",
      notes: "",
    },
  });

  const onSubmit = async (data: VisitFormData) => {
    try {
      // Convert datetime-local to RFC3339 format
      const payload = {
        ...data,
        scheduled_at: data.scheduled_at ? `${data.scheduled_at}:00Z` : "",
      };
      await createVisit.mutateAsync(payload);
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
            <select
              id="priority"
              {...register("priority")}
              className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
            >
              <option value="low">Baja</option>
              <option value="normal">Normal</option>
              <option value="urgent">Urgente</option>
              <option value="emergency">Emergencia</option>
            </select>
            {errors.priority && <p className="text-sm text-red-500">{errors.priority.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="status">Estado</Label>
            <select
              id="status"
              {...register("status")}
              className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
            >
              <option value="scheduled">Programada</option>
              <option value="en_route">En Ruta</option>
              <option value="in_progress">En Progreso</option>
              <option value="completed">Completada</option>
              <option value="cancelled">Cancelada</option>
              <option value="rescheduled">Reprogramada</option>
            </select>
            {errors.status && <p className="text-sm text-red-500">{errors.status.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="type">Tipo</Label>
            <select
              id="type"
              {...register("type")}
              disabled={loadingVisitTypes}
              className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
            >
              <option value="">{loadingVisitTypes ? "Cargando..." : "Seleccionar tipo"}</option>
              {visitTypes?.map((vt) => (
                <option key={vt.id} value={vt.id}>
                  {vt.name} ({vt.code})
                </option>
              ))}
            </select>
            {errors.type && <p className="text-sm text-red-500">{errors.type.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="billing_to">Destino de Facturación</Label>
            <select
              id="billing_to"
              {...register("billing_to")}
              className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
            >
              <option value="partner">Partner</option>
              <option value="customer">Cliente</option>
            </select>
            {errors.billing_to && <p className="text-sm text-red-500">{errors.billing_to.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="source">Fuente</Label>
            <select
              id="source"
              {...register("source")}
              className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
            >
              <option value="email">Email</option>
              <option value="phone">Teléfono</option>
              <option value="portal">Portal</option>
              <option value="whatsapp">WhatsApp</option>
              <option value="other">Otro</option>
            </select>
            {errors.source && <p className="text-sm text-red-500">{errors.source.message}</p>}
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