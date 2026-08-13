"use client";

import { useRouter, useParams } from "next/navigation";
import {
  useVisitSLATracking,
  useDeleteVisitSLATrackingGlobal,
} from "@/hooks/use-all-visit-sla-trackings";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";

export default function SLATrackingDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const { data: slaTracking, isLoading } = useVisitSLATracking(id);
  const deleteSLATracking = useDeleteVisitSLATrackingGlobal();
  const { user: currentUser } = useAuthStore();

  const handleDelete = async () => {
    if (confirm("¿Estás seguro de eliminar este seguimiento SLA?")) {
      try {
        await deleteSLATracking.mutateAsync(id);
        toast.success("Seguimiento SLA eliminado");
        router.push("/operations/sla-trackings");
      } catch (error) {
        toast.error("Error al eliminar seguimiento SLA");
      }
    }
  };

  if (isLoading) {
    return <div className="font-mono">Cargando...</div>;
  }

  if (!slaTracking) {
    return <div className="font-mono">Seguimiento SLA no encontrado</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Seguimiento SLA
        </h1>
        {canPerformAction(currentUser?.role, "delete") && (
          <Button variant="destructive" onClick={handleDelete}>
            Eliminar
          </Button>
        )}
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="font-mono text-lg uppercase">
            Información del Seguimiento
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="font-mono text-sm font-medium text-gray-500">
                Visita
              </label>
              <p className="font-mono">{slaTracking.visit_id}</p>
            </div>
            <div>
              <label className="font-mono text-sm font-medium text-gray-500">
                SLA
              </label>
              <p className="font-mono">{slaTracking.sla_id}</p>
            </div>
            <div>
              <label className="font-mono text-sm font-medium text-gray-500">
                Solicitado
              </label>
              <p className="font-mono">
                {slaTracking.requested_at
                  ? new Date(slaTracking.requested_at).toLocaleDateString("es-CL")
                  : "-"}
              </p>
            </div>
            <div>
              <label className="font-mono text-sm font-medium text-gray-500">
                Respondido
              </label>
              <p className="font-mono">
                {slaTracking.responded_at
                  ? new Date(slaTracking.responded_at).toLocaleDateString("es-CL")
                  : "-"}
              </p>
            </div>
            <div>
              <label className="font-mono text-sm font-medium text-gray-500">
                Resuelto
              </label>
              <p className="font-mono">
                {slaTracking.resolved_at
                  ? new Date(slaTracking.resolved_at).toLocaleDateString("es-CL")
                  : "-"}
              </p>
            </div>
            <div>
              <label className="font-mono text-sm font-medium text-gray-500">
                Creado
              </label>
              <p className="font-mono">
                {new Date(slaTracking.created_at).toLocaleDateString("es-CL")}
              </p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
