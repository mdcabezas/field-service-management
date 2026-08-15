"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useCreateVisitSLATrackingGlobal } from "@/hooks/use-all-visit-sla-trackings";
import { useVisits } from "@/hooks/use-visits";
import { usePartners } from "@/hooks/use-partners";
import { useSLAs } from "@/hooks/use-slas";
import type { VisitSLATracking } from "@/types/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent } from "@/components/ui/card";
import { toast } from "sonner";

const slaTrackingSchema = z.object({
  visit_id: z.string().min(1, "La visita es requerida"),
  partner_id: z.string().min(1, "El partner es requerido"),
  sla_id: z.string().min(1, "El SLA es requerido"),
  requested_at: z.string().optional(),
  responded_at: z.string().optional(),
  resolved_at: z.string().optional(),
});

type SLATrackingFormData = z.infer<typeof slaTrackingSchema>;

export function SLATrackingForm() {
  const router = useRouter();
  const createSLATracking = useCreateVisitSLATrackingGlobal();
  const { data: visits, isLoading: loadingVisits } = useVisits();
  const { data: partners, isLoading: loadingPartners } = usePartners();

  const { register, handleSubmit, watch, setValue, formState: { errors, isSubmitting } } = useForm<SLATrackingFormData>({
    resolver: zodResolver(slaTrackingSchema),
    defaultValues: { visit_id: "", partner_id: "", sla_id: "", requested_at: "", responded_at: "", resolved_at: "" },
  });

  const selectedPartnerId = watch("partner_id");
  const { data: slas, isLoading: loadingSLAs } = useSLAs(selectedPartnerId);

  // Reset sla_id when partner changes
  useEffect(() => {
    setValue("sla_id", "");
  }, [selectedPartnerId, setValue]);

  const onSubmit = async (data: SLATrackingFormData) => {
    try {
      const payload: Omit<VisitSLATracking, "id" | "created_at" | "updated_at"> = {
        visit_id: data.visit_id,
        sla_id: data.sla_id,
      };
      if (data.requested_at) payload.requested_at = `${data.requested_at}:00Z`;
      if (data.responded_at) payload.responded_at = `${data.responded_at}:00Z`;
      if (data.resolved_at) payload.resolved_at = `${data.resolved_at}:00Z`;
      await createSLATracking.mutateAsync(payload);
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
            <select
              id="visit_id"
              {...register("visit_id")}
              disabled={loadingVisits}
              className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
            >
              <option value="">{loadingVisits ? "Cargando..." : "Seleccionar visita"}</option>
              {visits?.map((visit) => (
                <option key={visit.id} value={visit.id}>
                  {visit.type_name || visit.type} - {new Date(visit.scheduled_at).toLocaleDateString()}
                </option>
              ))}
            </select>
            {errors.visit_id && <p className="text-sm text-red-500">{errors.visit_id.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="partner_id">Partner</Label>
            <select
              id="partner_id"
              {...register("partner_id")}
              disabled={loadingPartners}
              className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
            >
              <option value="">{loadingPartners ? "Cargando..." : "Seleccionar partner"}</option>
              {partners?.map((partner) => (
                <option key={partner.id} value={partner.id}>
                  {partner.name}
                </option>
              ))}
            </select>
            {errors.partner_id && <p className="text-sm text-red-500">{errors.partner_id.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="sla_id">SLA</Label>
            <select
              id="sla_id"
              {...register("sla_id")}
              disabled={loadingSLAs || !selectedPartnerId}
              className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
            >
              <option value="">
                {!selectedPartnerId
                  ? "Primero selecciona partner"
                  : loadingSLAs
                  ? "Cargando SLAs..."
                  : "Seleccionar SLA"}
              </option>
              {slas?.map((sla) => (
                <option key={sla.id} value={sla.id}>
                  {sla.name}
                </option>
              ))}
            </select>
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