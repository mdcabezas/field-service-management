"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useCustomers } from "@/hooks/use-customers";
import { usePartners } from "@/hooks/use-partners";
import { useVisits } from "@/hooks/use-visits";
import { useMaterials } from "@/hooks/use-materials";

export default function DashboardPage() {
  const { data: customers } = useCustomers();
  const { data: partners } = usePartners();
  const { data: visits } = useVisits();
  const { data: materials } = useMaterials();

  const metrics = [
    {
      title: "Clientes",
      value: customers?.length || 0,
      icon: "👥",
    },
    {
      title: "Socios",
      value: partners?.length || 0,
      icon: "🤝",
    },
    {
      title: "Visitas",
      value: visits?.length || 0,
      icon: "🔧",
    },
    {
      title: "Materiales",
      value: materials?.length || 0,
      icon: "📦",
    },
  ];

  return (
    <div className="space-y-6">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Dashboard
      </h1>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {metrics.map((metric) => (
          <Card key={metric.title}>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <span>{metric.icon}</span>
                <span>{metric.title}</span>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="font-mono text-3xl font-bold">{metric.value}</p>
            </CardContent>
          </Card>
        ))}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Accesos Rápidos</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <a
              href="/customers/new"
              className="border-2 border-black p-4 font-mono text-sm uppercase tracking-wider hover:bg-black hover:text-white transition-colors"
            >
              + Nuevo Cliente
            </a>
            <a
              href="/partners/new"
              className="border-2 border-black p-4 font-mono text-sm uppercase tracking-wider hover:bg-black hover:text-white transition-colors"
            >
              + Nuevo Socio
            </a>
            <a
              href="/operations/visits/new"
              className="border-2 border-black p-4 font-mono text-sm uppercase tracking-wider hover:bg-black hover:text-white transition-colors"
            >
              + Nueva Visita
            </a>
            <a
              href="/planning/daily-plans/new"
              className="border-2 border-black p-4 font-mono text-sm uppercase tracking-wider hover:bg-black hover:text-white transition-colors"
            >
              + Plan Diario
            </a>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
