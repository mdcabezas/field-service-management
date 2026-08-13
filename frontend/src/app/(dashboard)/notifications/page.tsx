"use client";

import { useRouter } from "next/navigation";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";

export default function NotificationsPage() {
  const router = useRouter();
  const { user } = useAuthStore();

  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Notificaciones
      </h1>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Card
          className="cursor-pointer hover:border-black transition-colors"
          onClick={() => router.push("/notifications/templates")}
        >
          <CardHeader>
            <CardTitle>Plantillas</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="font-mono text-sm text-gray-500">
              Gestionar plantillas de notificaciones
            </p>
          </CardContent>
        </Card>

        <Card
          className="cursor-pointer hover:border-black transition-colors"
          onClick={() => router.push("/notifications/list")}
        >
          <CardHeader>
            <CardTitle>Notificaciones Enviadas</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="font-mono text-sm text-gray-500">
              Historial de notificaciones enviadas
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
