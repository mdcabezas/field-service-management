"use client";

import { useRouter, useParams } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useNotificationTemplate, useUpdateNotificationTemplate, useDeleteNotificationTemplate } from "@/hooks/use-notification-templates";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";
import { useEffect } from "react";

const templateSchema = z.object({
  type: z.string().min(1, "Tipo requerido"),
  channel: z.string().min(1, "Canal requerido"),
  subject: z.string().optional(),
  body: z.string().optional(),
  active: z.boolean(),
});

type TemplateFormData = z.infer<typeof templateSchema>;

export default function NotificationTemplateDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const { data: template, isLoading } = useNotificationTemplate(id);
  const updateTemplate = useUpdateNotificationTemplate();
  const deleteTemplate = useDeleteNotificationTemplate();
  const { user: currentUser } = useAuthStore();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<TemplateFormData>({
    resolver: zodResolver(templateSchema),
    defaultValues: template
      ? {
          type: template.type,
          channel: template.channel,
          subject: template.subject || "",
          body: template.body || "",
          active: template.active,
        }
      : undefined,
  });

  // Reset form when template data loads
  useEffect(() => {
    if (template) {
      reset({
        type: template.type,
        channel: template.channel,
        subject: template.subject || "",
        body: template.body || "",
        active: template.active,
      });
    }
  }, [template, reset]);

  const onSubmit = async (data: TemplateFormData) => {
    try {
      await updateTemplate.mutateAsync({
        id,
        data,
      });
      toast.success("Plantilla actualizada");
    } catch (error) {
      toast.error("Error al actualizar plantilla");
    }
  };

  const handleDelete = async () => {
    if (confirm("¿Estás seguro de eliminar esta plantilla?")) {
      try {
        await deleteTemplate.mutateAsync(id);
        toast.success("Plantilla eliminada");
        router.push("/notifications/templates");
      } catch (error) {
        toast.error("Error al eliminar plantilla");
      }
    }
  };

  if (isLoading) {
    return <div className="font-mono">Cargando...</div>;
  }

  if (!template) {
    return <div className="font-mono">Plantilla no encontrada</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Plantilla de Notificación
        </h1>
        {canPerformAction(currentUser?.role, "delete") && (
          <Button variant="destructive" onClick={handleDelete}>
            Eliminar
          </Button>
        )}
      </div>

      <Card className="max-w-md">
        <CardHeader>
          <CardTitle>Editar Plantilla</CardTitle>
        </CardHeader>
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
              {errors.type && (
                <p className="text-sm text-red-600">{errors.type.message}</p>
              )}
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
              {errors.channel && (
                <p className="text-sm text-red-600">{errors.channel.message}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="subject">Asunto</Label>
              <Input id="subject" {...register("subject")} />
              {errors.subject && (
                <p className="text-sm text-red-600">{errors.subject.message}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="body">Cuerpo</Label>
              <textarea
                id="body"
                {...register("body")}
                rows={6}
                className="w-full px-3 py-2 font-mono text-sm border-2 border-black bg-white rounded-none"
              />
              {errors.body && (
                <p className="text-sm text-red-600">{errors.body.message}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="active">Activo</Label>
              <input
                id="active"
                type="checkbox"
                {...register("active")}
                className="h-4 w-4"
              />
            </div>

            <div className="flex gap-2">
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? "Guardando..." : "Actualizar"}
              </Button>
              <Button
                type="button"
                variant="outline"
                onClick={() => router.push("/notifications/templates")}
              >
                Cancelar
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
