"use client";

import { useRouter, useParams } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useReportTemplate, useUpdateReportTemplate, useDeleteReportTemplate } from "@/hooks/use-report-templates";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";
import { useEffect } from "react";

const reportTemplateSchema = z.object({
  name: z.string().min(1, "Nombre requerido"),
  fields_json: z.string().min(1, "Campos requeridos"),
});

type ReportTemplateFormData = z.infer<typeof reportTemplateSchema>;

export default function ReportTemplateDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const { data: template, isLoading } = useReportTemplate(id);
  const updateTemplate = useUpdateReportTemplate();
  const deleteTemplate = useDeleteReportTemplate();
  const { user: currentUser } = useAuthStore();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<ReportTemplateFormData>({
    resolver: zodResolver(reportTemplateSchema),
    defaultValues: template
      ? {
          name: template.name,
          fields_json: JSON.stringify(template.fields_json, null, 2),
        }
      : undefined,
  });

  // Reset form when template data loads
  useEffect(() => {
    if (template) {
      reset({
        name: template.name,
        fields_json: JSON.stringify(template.fields_json, null, 2),
      });
    }
  }, [template, reset]);

  const onSubmit = async (data: ReportTemplateFormData) => {
    try {
      let parsedFields: Record<string, unknown>;
      try {
        parsedFields = JSON.parse(data.fields_json);
      } catch {
        toast.error("Los campos deben ser JSON válido");
        return;
      }

      await updateTemplate.mutateAsync({
        id,
        data: {
          name: data.name,
          fields_json: parsedFields,
        },
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
        router.push("/shared/report-templates");
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
          {template.name}
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
              <Label htmlFor="name">Nombre</Label>
              <Input id="name" {...register("name")} />
              {errors.name && (
                <p className="text-sm text-red-600">{errors.name.message}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="fields_json">Campos (JSON)</Label>
              <textarea
                id="fields_json"
                {...register("fields_json")}
                rows={10}
                className="w-full px-3 py-2 font-mono text-sm border-2 border-black bg-white rounded-none"
              />
              {errors.fields_json && (
                <p className="text-sm text-red-600">{errors.fields_json.message}</p>
              )}
            </div>

            <div className="flex gap-2">
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? "Guardando..." : "Actualizar"}
              </Button>
              <Button
                type="button"
                variant="outline"
                onClick={() => router.push("/shared/report-templates")}
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
