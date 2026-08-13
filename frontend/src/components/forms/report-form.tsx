"use client";

import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useCreateVisitReportGlobal } from "@/hooks/use-all-visit-reports";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent } from "@/components/ui/card";
import { toast } from "sonner";

const reportSchema = z.object({
  visit_id: z.string().min(1, "La visita es requerida"),
  report_template_id: z.string().min(1, "La plantilla es requerida"),
  source: z.string().min(1, "La fuente es requerida"),
  recorded_at: z.string().optional(),
});

type ReportFormData = z.infer<typeof reportSchema>;

export function ReportForm() {
  const router = useRouter();
  const createReport = useCreateVisitReportGlobal();

  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<ReportFormData>({
    resolver: zodResolver(reportSchema),
    defaultValues: { visit_id: "", report_template_id: "", source: "", recorded_at: "" },
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
            <Input id="visit_id" {...register("visit_id")} />
            {errors.visit_id && <p className="text-sm text-red-500">{errors.visit_id.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="report_template_id">Plantilla</Label>
            <Input id="report_template_id" {...register("report_template_id")} />
            {errors.report_template_id && <p className="text-sm text-red-500">{errors.report_template_id.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="source">Fuente</Label>
            <Input id="source" {...register("source")} />
            {errors.source && <p className="text-sm text-red-500">{errors.source.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="recorded_at">Registrado</Label>
            <Input id="recorded_at" type="datetime-local" {...register("recorded_at")} />
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
