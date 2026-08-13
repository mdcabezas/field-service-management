"use client";

import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useCreateReportTemplate } from "@/hooks/use-report-templates";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent } from "@/components/ui/card";
import { toast } from "sonner";

const reportTemplateSchema = z.object({
  name: z.string().min(1, "Nombre requerido"),
  fields_json: z.string().min(1, "Campos requeridos"),
});

type ReportTemplateFormData = z.infer<typeof reportTemplateSchema>;

export function ReportTemplateForm() {
  const router = useRouter();
  const createTemplate = useCreateReportTemplate();

  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<ReportTemplateFormData>({
    resolver: zodResolver(reportTemplateSchema),
    defaultValues: { name: "", fields_json: "{}" },
  });

  const onSubmit = async (data: ReportTemplateFormData) => {
    try {
      let parsedFields: Record<string, unknown>;
      try {
        parsedFields = JSON.parse(data.fields_json);
      } catch {
        toast.error("Los campos deben ser JSON válido");
        return;
      }
      await createTemplate.mutateAsync({ name: data.name, fields_json: parsedFields });
      toast.success("Plantilla creada");
      router.push("/shared/report-templates");
    } catch {
      toast.error("Error al crear plantilla");
    }
  };

  return (
    <Card className="max-w-md">
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">Nombre</Label>
            <Input id="name" {...register("name")} />
            {errors.name && <p className="text-sm text-red-500">{errors.name.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="fields_json">Campos (JSON)</Label>
            <textarea
              id="fields_json"
              {...register("fields_json")}
              rows={10}
              className="w-full px-3 py-2 font-mono text-sm border-2 border-black bg-white rounded-none"
            />
            {errors.fields_json && <p className="text-sm text-red-500">{errors.fields_json.message}</p>}
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
