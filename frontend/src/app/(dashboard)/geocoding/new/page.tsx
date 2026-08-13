"use client";

import { GeocodingForm } from "@/components/forms/geocoding-form";

export default function NewGeocodingPage() {
  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Crear Dirección
      </h1>
      <GeocodingForm />
    </div>
  );
}
