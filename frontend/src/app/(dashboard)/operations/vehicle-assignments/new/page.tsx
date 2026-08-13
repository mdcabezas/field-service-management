"use client";

import { VehicleAssignmentForm } from "@/components/forms/vehicle-assignment-form";

export default function NewVehicleAssignmentPage() {
  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Crear Asignación de Vehículo
      </h1>
      <VehicleAssignmentForm />
    </div>
  );
}
