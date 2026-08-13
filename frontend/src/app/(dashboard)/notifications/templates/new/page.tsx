"use client";

import { NotificationTemplateForm } from "@/components/forms/notification-template-form";

export default function NewNotificationTemplatePage() {
  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Crear Plantilla de Notificación
      </h1>
      <NotificationTemplateForm />
    </div>
  );
}
