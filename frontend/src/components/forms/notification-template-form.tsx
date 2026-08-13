"use client";

import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useCreateNotificationTemplate } from "@/hooks/use-notification-templates";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent } from "@/components/ui/card";
import { toast } from "sonner";

const templateSchema = z.object({
  type: z.string().min(1, "Tipo requerido"),
  channel: z.string().min(1, "Canal requerido"),
  subject: z.string().optional(),
  body: z.string().optional(),
  active: z.boolean(),
});

type TemplateFormData = z.infer<typeof templateSchema>;

export function NotificationTemplateForm() {
  const router = useRouter();
  const createTemplate = useCreateNotificationTemplate();

  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<TemplateFormData>({
    resolver: zodResolver(templateSchema),
    defaultValues: { type: "email", channel: "", subject: "", body: "", active: true },
  });

  const onSubmit = async (data: TemplateFormData) => {
    try {
      await createTemplate.mutateAsync(data);
      toast.success("Plantilla creada");
      router.push("/notifications/templates");
    } catch {
      toast.error("Error al crear plantilla");
    }
  };

  return (
    <Card className="max-w-md">
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="type">Tipo</Label>
            <select
              id="type"
              {...register("type")}
              className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
            >
              <option value="visit_assigned">Visita Asignada</option>
              <option value="visit_reminder">Recordatorio de Visita</option>
              <option value="visit_cancelled">Visita Cancelada</option>
              <option value="visit_rescheduled">Visita Reprogramada</option>
              <option value="checklist_pending">Checklist Pendiente</option>
              <option value="report_pending">Reporte Pendiente</option>
              <option value="sla_warning">Alerta SLA</option>
              <option value="sla_expired">SLA Vencido</option>
              <option value="maintenance_scheduled">Mantenimiento Programado</option>
              <option value="certification_expired">Certificación Vencida</option>
              <option value="other">Otro</option>
            </select>
            {errors.type && <p className="text-sm text-red-500">{errors.type.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="channel">Canal</Label>
            <select
              id="channel"
              {...register("channel")}
              className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
            >
              <option value="email">Email</option>
              <option value="sms">SMS</option>
              <option value="whatsapp">WhatsApp</option>
              <option value="push">Push</option>
              <option value="in_app">In App</option>
            </select>
            {errors.channel && <p className="text-sm text-red-500">{errors.channel.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="subject">Asunto</Label>
            <Input id="subject" {...register("subject")} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="body">Cuerpo</Label>
            <textarea
              id="body"
              {...register("body")}
              rows={6}
              className="w-full px-3 py-2 font-mono text-sm border-2 border-black bg-white rounded-none"
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="active">Activo</Label>
            <input id="active" type="checkbox" {...register("active")} className="h-4 w-4" />
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
