"use client";

import { UserForm } from "@/components/forms/user-form";

export default function NewUserPage() {
  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Crear Usuario
      </h1>
      <UserForm />
    </div>
  );
}
