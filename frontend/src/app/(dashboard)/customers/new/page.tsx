"use client";

import { CustomerForm } from "@/components/forms/customer-form";

export default function NewCustomerPage() {
  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Crear Cliente
      </h1>
      <CustomerForm />
    </div>
  );
}
