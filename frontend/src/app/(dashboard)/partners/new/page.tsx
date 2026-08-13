"use client";

import { PartnerForm } from "@/components/forms/partner-form";

export default function NewPartnerPage() {
  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Crear Socio
      </h1>
      <PartnerForm />
    </div>
  );
}
