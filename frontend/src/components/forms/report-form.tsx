"use client";

import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useCreateVisitReportGlobal } from "@/hooks/use-all-visit-reports";
import { useVisits } from "@/hooks/use-visits";
import { useReportTemplates } from "@/hooks/use-report-templates";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Card, CardContent } from "@/components/ui/card";
import { toast } from "sonner";

const reportSchema = z.object({
  visit_id: z.string().min(1, "La visita es requerida"),
  report_template_id: z.string().min(1, "La plantilla es requerida"),
  source: z.string().min(1, "La fuente es requerida"),
});

type ReportFormData = z.infer<typeof reportSchema>;

export function ReportForm() {
  const router = useRouter();
  const createReport = useCreateVisitReportGlobal();
  const { data: visits, isLoading: loadingVisits } = useVisits();
  const { data: templates, isLoading: loadingTemplates } = useReportTemplates();

  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<ReportFormData>({
    resolver: zodResolver(reportSchema),
    defaultValues: { visit_id: "", report_template_id: "", source: "" },
  });

  const onSubmit = async (data: ReportFormData) => {
    try {
      await createReport.mutateAsync(data);
      toast.success("Reporte creado");
      router.push("/operations/reports");
    } catch {
      toast.error("Error al crear reporte");
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
            <Label htmlFor="report_template_id">Plantilla</Label>
            <select
              id="report_template_id"
              {...register("report_template_id")}
              disabled={loadingTemplates}
              className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
            >
              <option value="">{loadingTemplates ? "Cargando..." : "Seleccionar plantilla"}</option>
              {templates?.map((template) => (
                <option key={template.id} value={template.id}>
                  {template.name}
                </option>
              ))}
            </select>
            {errors.report_template_id && <p className="text-sm text-red-500">{errors.report_template_id.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="source">Fuente</Label>
            <select
              id="source"
              {...register("source")}
              className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
            >
              <option value="">Seleccionar fuente</option>
              <option value="system">Sistema</option>
              <option value="paper">Papel</option>
            </select>
            {errors.source && <p className="text-sm text-red-500">{errors.source.message}</p>}
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