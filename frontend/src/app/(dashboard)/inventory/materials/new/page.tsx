"use client";

import { useRouter } from "next/navigation";
import { useCreateMaterial } from "@/hooks/use-materials";
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
  sku: z.string().min(1, "SKU es requerido"),
  unit: z.string().min(1, "Unidad es requerida"),
  unit_cost: z.coerce.number().min(0, "Costo unitario debe ser ≥ 0"),
});

type FormData = z.infer<typeof schema>;

export default function NewMaterialPage() {
  const router = useRouter();
  const createMaterial = useCreateMaterial();

  const { register, handleSubmit, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: { name: "", sku: "", unit: "", unit_cost: 0 },
  });

  const onSubmit = async (data: FormData) => {
    try {
      await createMaterial.mutateAsync(data);
      toast.success("Material creado");
      router.back();
    } catch {
      toast.error("Error al guardar");
    }
  };

  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Crear Material
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
              <Label htmlFor="sku">SKU</Label>
              <Input id="sku" {...register("sku")} />
              {errors.sku && <p className="text-sm text-red-500">{errors.sku.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="unit">Unidad</Label>
              <Input id="unit" {...register("unit")} />
              {errors.unit && <p className="text-sm text-red-500">{errors.unit.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="unit_cost">Costo Unitario</Label>
              <Input id="unit_cost" type="number" step="0.01" {...register("unit_cost")} />
              {errors.unit_cost && <p className="text-sm text-red-500">{errors.unit_cost.message}</p>}
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
