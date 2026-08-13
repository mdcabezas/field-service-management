"use client";

import { ReportForm } from "@/components/forms/report-form";

export default function NewReportPage() {
  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Crear Reporte
      </h1>
      <ReportForm />
    </div>
  );
}
