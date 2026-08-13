"use client";

import { useRouter } from "next/navigation";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

const sections = [
  { title: "Materiales", description: "Gestión de materiales e insumos", href: "/inventory/materials" },
  { title: "Herramientas", description: "Control de herramientas asignadas", href: "/inventory/tools" },
  { title: "Elementos de Protección", description: "Gestión de EPPs", href: "/inventory/epp" },
  { title: "Vehículos", description: "Flota de vehículos", href: "/inventory/vehicles" },
  { title: "Arriendos", description: "Registro de arriendos de equipos", href: "/inventory/rentals" },
  { title: "Tarifas de Costo", description: "Tarifas vigentes por unidad", href: "/inventory/cost-rates" },
  { title: "Registros de Mantención", description: "Historial de mantención vehicular", href: "/inventory/maintenance-records" },
  { title: "Programación de Mantención", description: "Calendario de mantención", href: "/inventory/maintenance-schedules" },
  { title: "Plantillas de Checklist", description: "Checklists de visita por tipo", href: "/inventory/checklist-templates" },
];

export default function InventoryPage() {
  const router = useRouter();

  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Inventario
      </h1>
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {sections.map((section) => (
          <Card
            key={section.href}
            className="cursor-pointer hover:border-black transition-colors"
            onClick={() => router.push(section.href)}
          >
            <CardHeader>
              <CardTitle className="font-mono text-sm uppercase tracking-wider">
                {section.title}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="font-mono text-sm text-gray-500">
                {section.description}
              </p>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
