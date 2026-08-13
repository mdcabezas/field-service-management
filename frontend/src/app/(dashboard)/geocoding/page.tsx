"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { type ColumnDef } from "@tanstack/react-table";
import { useGeocodedAddresses, useNearbyAddresses } from "@/hooks/use-geocoded-addresses";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import type { GeocodedAddress } from "@/types/api";

const nearbySchema = z.object({
  lat: z.coerce.number().min(-90).max(90),
  lng: z.coerce.number().min(-180).max(180),
  radius: z.coerce.number().positive("Radio debe ser positivo"),
});

type NearbyFormData = z.infer<typeof nearbySchema>;

const columns: ColumnDef<GeocodedAddress, unknown>[] = [
  {
    accessorKey: "street",
    header: "Calle",
  },
  {
    accessorKey: "number",
    header: "Número",
  },
  {
    accessorKey: "commune",
    header: "Comuna",
  },
  {
    accessorKey: "city",
    header: "Ciudad",
  },
  {
    accessorKey: "region",
    header: "Región",
  },
  {
    accessorKey: "lat",
    header: "Lat",
  },
  {
    accessorKey: "lng",
    header: "Lng",
  },];

export default function GeocodingPage() {
  const router = useRouter();
  const { data: addresses, isLoading } = useGeocodedAddresses();
  const { user } = useAuthStore();

  const [nearbyParams, setNearbyParams] = useState<NearbyFormData | null>(null);
  const { data: nearbyAddresses, isLoading: isLoadingNearby } = useNearbyAddresses(
    nearbyParams?.lat ?? 0,
    nearbyParams?.lng ?? 0,
    nearbyParams?.radius ?? 0
  );

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<NearbyFormData>({
    resolver: zodResolver(nearbySchema),
    defaultValues: {
      lat: -33.4489,
      lng: -70.6693,
      radius: 1000,
    },
  });

  const onNearbySubmit = (data: NearbyFormData) => {
    setNearbyParams(data);
  };

  const handleRowClick = (row: GeocodedAddress) => {
    router.push(`/geocoding/${row.id}`);
  };

  if (isLoading) {
    return <div className="font-mono">Cargando...</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Geocodificación
        </h1>
        {canPerformAction(user?.role, "create") && (
          <Button onClick={() => router.push("/geocoding/new")}>
            + Crear Dirección
          </Button>
        )}
      </div>

      <DataTable
        columns={columns}
        data={addresses || []}
        searchPlaceholder="Buscar por calle o comuna..."
        searchColumn="street"
        onRowClick={handleRowClick}
      />

      <Card>
        <CardHeader>
          <CardTitle>Buscar Cercanos</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onNearbySubmit)} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="space-y-2">
                <Label htmlFor="lat">Latitud</Label>
                <Input id="lat" type="number" step="any" {...register("lat")} />
                {errors.lat && (
                  <p className="text-sm text-red-600">{errors.lat.message}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="lng">Longitud</Label>
                <Input id="lng" type="number" step="any" {...register("lng")} />
                {errors.lng && (
                  <p className="text-sm text-red-600">{errors.lng.message}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="radius">Radio (m)</Label>
                <Input id="radius" type="number" step="any" {...register("radius")} />
                {errors.radius && (
                  <p className="text-sm text-red-600">{errors.radius.message}</p>
                )}
              </div>
            </div>

            <Button type="submit">Buscar</Button>
          </form>

          {nearbyParams && (
            <div className="mt-4">
              {isLoadingNearby ? (
                <p className="font-mono text-sm text-gray-500">Buscando...</p>
              ) : nearbyAddresses && nearbyAddresses.length > 0 ? (
                <DataTable
                  columns={columns}
                  data={nearbyAddresses}
                  searchPlaceholder=""
                  searchColumn="street"
                  onRowClick={handleRowClick}
                />
              ) : (
                <p className="font-mono text-sm text-gray-500">
                  No se encontraron direcciones cercanas
                </p>
              )}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
