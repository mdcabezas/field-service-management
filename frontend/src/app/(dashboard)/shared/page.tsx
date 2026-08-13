"use client";

import { useRouter } from "next/navigation";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export default function SharedPage() {
  const router = useRouter();

  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Compartido
      </h1>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Card
          className="cursor-pointer hover:border-black transition-colors"
          onClick={() => router.push("/shared/report-templates")}
        >
          <CardHeader>
            <CardTitle>Plantillas de Reportes</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="font-mono text-sm text-gray-500">
              Gestionar plantillas de reportes
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
