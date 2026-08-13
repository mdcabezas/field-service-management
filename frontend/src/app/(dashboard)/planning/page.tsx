"use client";

import { useRouter } from "next/navigation";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

const modules = [
  {
    title: "Planes Diarios",
    description: "Gestión de planes de trabajo diarios para técnicos",
    path: "/planning/daily-plans",
  },
  {
    title: "Asignaciones",
    description: "Asignación de visitas a planes diarios con horarios y prioridades",
    path: "/planning/assignments",
  },
  {
    title: "Rutas",
    description: "Definición y gestión de rutas de trabajo",
    path: "/planning/routes",
  },
];

export default function PlanningPage() {
  const router = useRouter();

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Planificación
        </h1>
      </div>

      <div className="grid gap-4 md:grid-cols-3">
        {modules.map((module) => (
          <Card key={module.path}>
            <CardHeader>
              <CardTitle>{module.title}</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="font-mono text-sm text-gray-600 mb-4">
                {module.description}
              </p>
              <Button
                className="w-full"
                onClick={() => router.push(module.path)}
              >
                Ver {module.title}
              </Button>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
