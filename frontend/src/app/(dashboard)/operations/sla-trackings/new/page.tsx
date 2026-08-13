"use client";

import { SLATrackingForm } from "@/components/forms/sla-tracking-form";

export default function NewSLATrackingPage() {
  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Crear Seguimiento SLA
      </h1>
      <SLATrackingForm />
    </div>
  );
}
