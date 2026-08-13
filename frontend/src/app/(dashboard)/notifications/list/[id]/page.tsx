"use client";

import { useRouter, useParams } from "next/navigation";
import { useNotification } from "@/hooks/use-notifications";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export default function NotificationDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const { data: notification, isLoading } = useNotification(id);

  if (isLoading) {
    return <div className="font-mono">Cargando...</div>;
  }

  if (!notification) {
    return <div className="font-mono">Notificación no encontrada</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Notificación
        </h1>
        <Button
          variant="outline"
          onClick={() => router.push("/notifications/list")}
        >
          Volver
        </Button>
      </div>

      <Card className="max-w-md">
        <CardHeader>
          <CardTitle>Detalles de la Notificación</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <div>
            <span className="font-mono text-sm text-gray-500">ID:</span>
            <p className="font-mono">{notification.id}</p>
          </div>
          <div>
            <span className="font-mono text-sm text-gray-500">Destinatario:</span>
            <p className="font-mono">{notification.recipient}</p>
          </div>
          <div>
            <span className="font-mono text-sm text-gray-500">Tipo:</span>
            <p className="font-mono">{notification.type}</p>
          </div>
          <div>
            <span className="font-mono text-sm text-gray-500">Canal:</span>
            <p className="font-mono">{notification.channel}</p>
          </div>
          <div>
            <span className="font-mono text-sm text-gray-500">Asunto:</span>
            <p className="font-mono">{notification.subject || "-"}</p>
          </div>
          <div>
            <span className="font-mono text-sm text-gray-500">Cuerpo:</span>
            <p className="font-mono">{notification.body || "-"}</p>
          </div>
          <div>
            <span className="font-mono text-sm text-gray-500">Estado:</span>
            <p className="font-mono uppercase">{notification.status}</p>
          </div>
          <div>
            <span className="font-mono text-sm text-gray-500">Enviado:</span>
            <p className="font-mono">
              {notification.sent_at
                ? new Date(notification.sent_at).toLocaleString()
                : "-"}
            </p>
          </div>
          <div>
            <span className="font-mono text-sm text-gray-500">Creado:</span>
            <p className="font-mono">
              {new Date(notification.created_at).toLocaleString()}
            </p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
