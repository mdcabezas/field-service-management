"use client";

import { useRouter } from "next/navigation";
import { useCreateVehicle } from "@/hooks/use-vehicles";
import { useVehicleTypes } from "@/hooks/use-vehicle-types";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

const schema = z.object({
  license_plate: z.string().min(1, "Patente es requerida"),
  name: z.string().min(1, "Nombre es requerido"),
  brand: z.string().optional(),
  type: z.string().optional(),
  status: z.string().min(1, "Estado es requerido"),
});

type FormData = z.infer<typeof schema>;

export default function NewVehiclePage() {
  const router = useRouter();
  const createVehicle = useCreateVehicle();
  const { data: vehicleTypes } = useVehicleTypes();

  const { register, handleSubmit, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: { license_plate: "", name: "", brand: "", type: "", status: "" },
  });

  const onSubmit = async (data: FormData) => {
    try {
      const payload = {
        ...data,
        brand: data.brand || undefined,
        type: data.type || undefined,
      };
      await createVehicle.mutateAsync(payload);
      toast.success("Vehículo creado");
      router.back();
    } catch {
      toast.error("Error al guardar");
    }
  };

  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Crear Vehículo
      </h1>
      <Card className="max-w-md">
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="license_plate">Patente</Label>
              <Input id="license_plate" {...register("license_plate")} />
              {errors.license_plate && <p className="text-sm text-red-500">{errors.license_plate.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="name">Nombre</Label>
              <Input id="name" {...register("name")} />
              {errors.name && <p className="text-sm text-red-500">{errors.name.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="brand">Marca</Label>
              <Input id="brand" {...register("brand")} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="type">Tipo</Label>
              <select
                id="type"
                {...register("type")}
                className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
              >
                <option value="">Seleccionar...</option>
                {vehicleTypes?.map((vt) => (
                  <option key={vt.id} value={vt.id}>{vt.name}</option>
                ))}
              </select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="status">Estado</Label>
              <select
                id="status"
                {...register("status")}
                className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
              >
                <option value="">Seleccionar...</option>
                <option value="available">Disponible</option>
                <option value="in_use">En Uso</option>
                <option value="maintenance">Mantenimiento</option>
                <option value="retired">Retirado</option>
              </select>
              {errors.status && <p className="text-sm text-red-500">{errors.status.message}</p>}
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
