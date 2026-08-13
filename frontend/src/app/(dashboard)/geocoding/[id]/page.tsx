"use client";

import { useRouter, useParams } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useGeocodedAddress, useUpdateGeocodedAddress, useDeleteGeocodedAddress } from "@/hooks/use-geocoded-addresses";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";
import { useEffect } from "react";

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

export default function GeocodedAddressDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const { data: address, isLoading } = useGeocodedAddress(id);
  const updateAddress = useUpdateGeocodedAddress();
  const deleteAddress = useDeleteGeocodedAddress();
  const { user: currentUser } = useAuthStore();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<AddressFormData>({
    resolver: zodResolver(addressSchema),
    defaultValues: address
      ? {
          street: address.street,
          number: address.number,
          city: address.city,
          region: address.region,
          neighborhood: address.neighborhood,
          apartment: address.apartment,
          location_references: address.location_references,
          postal_code: address.postal_code,
        }
      : undefined,
  });

  // Reset form when address data loads
  useEffect(() => {
    if (address) {
      reset({
        street: address.street,
        number: address.number,
        city: address.city,
        region: address.region,
        neighborhood: address.neighborhood || "",
        apartment: address.apartment || "",
        location_references: address.location_references || "",
        postal_code: address.postal_code || "",
      });
    }
  }, [address, reset]);

  const onSubmit = async (data: AddressFormData) => {
    try {
      await updateAddress.mutateAsync({
        id,
        data,
      });
      toast.success("Dirección actualizada");
    } catch (error) {
      toast.error("Error al actualizar dirección");
    }
  };

  const handleDelete = async () => {
    if (confirm("¿Estás seguro de eliminar esta dirección?")) {
      try {
        await deleteAddress.mutateAsync(id);
        toast.success("Dirección eliminada");
        router.push("/geocoding");
      } catch (error) {
        toast.error("Error al eliminar dirección");
      }
    }
  };

  if (isLoading) {
    return <div className="font-mono">Cargando...</div>;
  }

  if (!address) {
    return <div className="font-mono">Dirección no encontrada</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          {address.street} {address.number}
        </h1>
        {canPerformAction(currentUser?.role, "delete") && (
          <Button variant="destructive" onClick={handleDelete}>
            Eliminar
          </Button>
        )}
      </div>

      <Card className="max-w-md">
        <CardHeader>
          <CardTitle>Editar Dirección</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="street">Calle</Label>
              <Input id="street" {...register("street")} />
              {errors.street && (
                <p className="text-sm text-red-600">{errors.street.message}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="number">Número</Label>
              <Input id="number" {...register("number")} />
              {errors.number && (
                <p className="text-sm text-red-600">{errors.number.message}</p>
              )}
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
                {isSubmitting ? "Guardando..." : "Actualizar"}
              </Button>
              <Button
                type="button"
                variant="outline"
                onClick={() => router.push("/geocoding")}
              >
                Cancelar
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
