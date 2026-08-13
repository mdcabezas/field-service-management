"use client";

import { ReportTemplateForm } from "@/components/forms/report-template-form";

export default function NewReportTemplatePage() {
  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Crear Plantilla
      </h1>
      <ReportTemplateForm />
    </div>
  );
}
