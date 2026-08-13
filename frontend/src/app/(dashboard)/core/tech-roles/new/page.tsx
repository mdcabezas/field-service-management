"use client";

import { TechRoleForm } from "@/components/forms/tech-role-form";

export default function NewTechRolePage() {
  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Crear Rol Técnico
      </h1>
      <TechRoleForm />
    </div>
  );
}
