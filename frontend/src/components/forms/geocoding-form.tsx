"use client";

import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useCreateGeocodedAddress } from "@/hooks/use-geocoded-addresses";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent } from "@/components/ui/card";
import { toast } from "sonner";

const addressSchema = z.object({
  street: z.string().min(1, "Calle requerida"),
  number: z.string().min(1, "Número requerido"),
  city: z.string().min(1, "Ciudad requerida"),
  region: z.string().min(1, "Región requerida"),
  neighborhood: z.string().optional(),
  apartment: z.string().optional(),
  location_references: z.string().optional(),
  postal_code: z.string().optional(),
});

type AddressFormData = z.infer<typeof addressSchema>;

export function GeocodingForm() {
  const router = useRouter();
  const createAddress = useCreateGeocodedAddress();

  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<AddressFormData>({
    resolver: zodResolver(addressSchema),
    defaultValues: { street: "", number: "", city: "", region: "", neighborhood: "", apartment: "", location_references: "", postal_code: "" },
  });

  const onSubmit = async (data: AddressFormData) => {
    try {
      await createAddress.mutateAsync(data);
      toast.success("Dirección creada");
      router.push("/geocoding");
    } catch {
      toast.error("Error al crear dirección");
    }
  };

  return (
    <Card className="max-w-md">
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="street">Calle</Label>
            <Input id="street" {...register("street")} />
            {errors.street && <p className="text-sm text-red-500">{errors.street.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="number">Número</Label>
            <Input id="number" {...register("number")} />
            {errors.number && <p className="text-sm text-red-500">{errors.number.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="city">Ciudad</Label>
            <Input id="city" {...register("city")} />
            {errors.city && <p className="text-sm text-red-500">{errors.city.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="region">Región</Label>
            <Input id="region" {...register("region")} />
            {errors.region && <p className="text-sm text-red-500">{errors.region.message}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="neighborhood">Barrio</Label>
            <Input id="neighborhood" {...register("neighborhood")} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="apartment">Departamento</Label>
            <Input id="apartment" {...register("apartment")} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="location_references">Referencias</Label>
            <Input id="location_references" {...register("location_references")} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="postal_code">Código Postal</Label>
            <Input id="postal_code" {...register("postal_code")} />
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
