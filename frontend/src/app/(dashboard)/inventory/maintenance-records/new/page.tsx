"use client";

import { useRouter } from "next/navigation";
import { useCreateMaintenanceRecord } from "@/hooks/use-maintenance-records";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

const schema = z.object({
  reference_id: z.string().min(1, "ID de referencia es requerido"),
  type: z.string().min(1, "Tipo es requerido"),
  date: z.string().min(1, "Fecha es requerida"),
  cost: z.coerce.number().min(0, "Costo debe ser ≥ 0").optional(),
  description: z.string().optional(),
});

type FormData = z.infer<typeof schema>;

export default function NewMaintenanceRecordPage() {
  const router = useRouter();
  const createRecord = useCreateMaintenanceRecord();

  const { register, handleSubmit, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: { reference_id: "", type: "", date: "", cost: 0, description: "" },
  });

  const onSubmit = async (data: FormData) => {
    try {
      const payload = {
        ...data,
        date: data.date ? `${data.date}T00:00:00Z` : data.date,
      };
      await createRecord.mutateAsync(payload);
      toast.success("Registro creado");
      router.back();
    } catch {
      toast.error("Error al guardar");
    }
  };

  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Crear Registro
      </h1>
      <Card className="max-w-md">
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="reference_id">ID de Referencia</Label>
              <Input id="reference_id" {...register("reference_id")} placeholder="UUID del vehículo/herramienta" />
              {errors.reference_id && <p className="text-sm text-red-500">{errors.reference_id.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="type">Tipo</Label>
              <select
                id="type"
                {...register("type")}
                className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
              >
                <option value="">Seleccionar...</option>
                <option value="vehicle">Vehículo</option>
                <option value="tool">Herramienta</option>
                <option value="equipment">Equipo</option>
              </select>
              {errors.type && <p className="text-sm text-red-500">{errors.type.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="date">Fecha</Label>
              <Input id="date" type="date" {...register("date")} />
              {errors.date && <p className="text-sm text-red-500">{errors.date.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="cost">Costo</Label>
              <Input id="cost" type="number" step="0.01" {...register("cost")} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="description">Descripción</Label>
              <Input id="description" {...register("description")} />
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
