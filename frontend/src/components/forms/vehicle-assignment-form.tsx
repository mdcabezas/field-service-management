"use client";

import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useCreateVehicleAssignmentGlobal } from "@/hooks/use-all-vehicle-assignments";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent } from "@/components/ui/card";
import { toast } from "sonner";

const vehicleAssignmentSchema = z.object({
  vehicle_id: z.string().min(1, "El vehículo es requerido"),
  departure_time: z.string().optional(),
  return_time: z.string().optional(),
  notes: z.string().optional(),
});

type VehicleAssignmentFormData = z.infer<typeof vehicleAssignmentSchema>;

export function VehicleAssignmentForm() {
  const router = useRouter();
  const createAssignment = useCreateVehicleAssignmentGlobal();

  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<VehicleAssignmentFormData>({
    resolver: zodResolver(vehicleAssignmentSchema),
    defaultValues: { vehicle_id: "", departure_time: "", return_time: "", notes: "" },
  });

  const onSubmit = async (data: VehicleAssignmentFormData) => {
    try {
      await createAssignment.mutateAsync(data);
      toast.success("Asignación creada");
      router.push("/operations/vehicle-assignments");
    } catch {
      toast.error("Error al crear asignación");
    }
  };

  return (
    <Card className="max-w-md">
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="vehicle_id">Vehículo</Label>
            <Input id="vehicle_id" {...register("vehicle_id")} />
            {errors.vehicle_id && <p className="text-sm text-red-500">{errors.vehicle_id.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="departure_time">Salida</Label>
            <Input id="departure_time" type="datetime-local" {...register("departure_time")} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="return_time">Retorno</Label>
            <Input id="return_time" type="datetime-local" {...register("return_time")} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="notes">Notas</Label>
            <Input id="notes" {...register("notes")} />
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
