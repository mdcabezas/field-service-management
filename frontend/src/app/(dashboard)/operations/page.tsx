"use client";

import { useRouter } from "next/navigation";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

const modules = [
  {
    title: "Visitas",
    description: "Gestión de visitas técnicas, asignaciones, checklists y más",
    href: "/operations/visits",
  },
  {
    title: "Reportes de Visita",
    description: "Informes generados a partir de las visitas realizadas",
    href: "/operations/reports",
  },
  {
    title: "Asignaciones de Vehículos",
    description: "Control de vehículos asignados a cada visita",
    href: "/operations/vehicle-assignments",
  },
  {
    title: "Seguimiento SLA",
    description: "Monitoreo del cumplimiento de acuerdos de nivel de servicio",
    href: "/operations/sla-trackings",
  },
];

export default function OperationsPage() {
  const router = useRouter();

  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Operaciones
      </h1>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        {modules.map((mod) => (
          <Card
            key={mod.href}
            className="cursor-pointer border-2 border-black transition-colors hover:bg-gray-50"
            onClick={() => router.push(mod.href)}
          >
            <CardHeader>
              <CardTitle className="font-mono text-lg uppercase">
                {mod.title}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="font-mono text-sm text-gray-600">
                {mod.description}
              </p>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
