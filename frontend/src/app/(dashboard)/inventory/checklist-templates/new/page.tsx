"use client";

import { useRouter } from "next/navigation";
import { useCreateChecklistTemplate } from "@/hooks/use-checklist-templates";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

const schema = z.object({
  name: z.string().min(1, "Nombre es requerido"),
  description: z.string().optional(),
  visit_type: z.string().optional(),
  work_type: z.string().min(1, "Tipo de trabajo es requerido"),
});

type FormData = z.infer<typeof schema>;

export default function NewChecklistTemplatePage() {
  const router = useRouter();
  const createTemplate = useCreateChecklistTemplate();

  const { register, handleSubmit, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: { name: "", description: "", visit_type: "", work_type: "" },
  });

  const onSubmit = async (data: FormData) => {
    try {
      await createTemplate.mutateAsync(data);
      toast.success("Plantilla creada");
      router.back();
    } catch {
      toast.error("Error al guardar");
    }
  };

  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Crear Plantilla
      </h1>
      <Card className="max-w-md">
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">Nombre</Label>
              <Input id="name" {...register("name")} />
              {errors.name && <p className="text-sm text-red-500">{errors.name.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="description">Descripción</Label>
              <Input id="description" {...register("description")} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="visit_type">Tipo de Visita</Label>
              <Input id="visit_type" {...register("visit_type")} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="work_type">Tipo de Trabajo</Label>
              <Input id="work_type" {...register("work_type")} />
              {errors.work_type && <p className="text-sm text-red-500">{errors.work_type.message}</p>}
            </div>
            <div className="flex gap-2">
              <Button type="submit">Crear</Button>
              <Button type="button" variant="outline" onClick={() => router.back()}>Cancelar</Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
