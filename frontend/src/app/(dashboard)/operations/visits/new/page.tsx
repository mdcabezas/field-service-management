"use client";

import { VisitForm } from "@/components/forms/visit-form";

export default function NewVisitPage() {
  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Crear Visita
      </h1>
      <VisitForm />
    </div>
  );
}
