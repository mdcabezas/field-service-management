"use client";

import { useRouter } from "next/navigation";
import { useCreateRental } from "@/hooks/use-rentals";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

const schema = z.object({
  type: z.string().min(1, "Tipo es requerido"),
  supplier: z.string().optional(),
  item_description: z.string().optional(),
  daily_cost: z.coerce.number().min(0, "Costo diario debe ser ≥ 0").optional(),
  notes: z.string().optional(),
});

type FormData = z.infer<typeof schema>;

export default function NewRentalPage() {
  const router = useRouter();
  const createRental = useCreateRental();

  const { register, handleSubmit, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: { type: "", supplier: "", item_description: "", daily_cost: 0, notes: "" },
  });

  const onSubmit = async (data: FormData) => {
    try {
      await createRental.mutateAsync(data);
      toast.success("Arriendo creado");
      router.back();
    } catch {
      toast.error("Error al guardar");
    }
  };

  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Crear Arriendo
      </h1>
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
                <option value="">Seleccionar...</option>
                <option value="vehicle">Vehículo</option>
                <option value="tool">Herramienta</option>
                <option value="equipment">Equipo</option>
              </select>
              {errors.type && <p className="text-sm text-red-500">{errors.type.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="supplier">Proveedor</Label>
              <Input id="supplier" {...register("supplier")} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="item_description">Descripción del Item</Label>
              <Input id="item_description" {...register("item_description")} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="daily_cost">Costo Diario</Label>
              <Input id="daily_cost" type="number" step="0.01" {...register("daily_cost")} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="notes">Notas</Label>
              <Input id="notes" {...register("notes")} />
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
