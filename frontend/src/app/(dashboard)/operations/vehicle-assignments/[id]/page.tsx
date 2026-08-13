"use client";

import { useRouter, useParams } from "next/navigation";
import {
  useVehicleAssignment,
  useDeleteVehicleAssignmentGlobal,
} from "@/hooks/use-all-vehicle-assignments";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";

export default function VehicleAssignmentDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const { data: assignment, isLoading } = useVehicleAssignment(id);
  const deleteAssignment = useDeleteVehicleAssignmentGlobal();
  const { user: currentUser } = useAuthStore();

  const handleDelete = async () => {
    if (confirm("¿Estás seguro de eliminar esta asignación?")) {
      try {
        await deleteAssignment.mutateAsync(id);
        toast.success("Asignación eliminada");
        router.push("/operations/vehicle-assignments");
      } catch (error) {
        toast.error("Error al eliminar asignación");
      }
    }
  };

  if (isLoading) {
    return <div className="font-mono">Cargando...</div>;
  }

  if (!assignment) {
    return <div className="font-mono">Asignación no encontrada</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Asignación de Vehículo
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
            Información de la Asignación
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="font-mono text-sm font-medium text-gray-500">
                Vehículo
              </label>
              <p className="font-mono">{assignment.vehicle_id}</p>
            </div>
            <div>
              <label className="font-mono text-sm font-medium text-gray-500">
                Salida
              </label>
              <p className="font-mono">
                {assignment.departure_time
                  ? new Date(assignment.departure_time).toLocaleDateString("es-CL")
                  : "-"}
              </p>
            </div>
            <div>
              <label className="font-mono text-sm font-medium text-gray-500">
                Retorno
              </label>
              <p className="font-mono">
                {assignment.return_time
                  ? new Date(assignment.return_time).toLocaleDateString("es-CL")
                  : "-"}
              </p>
            </div>
            <div>
              <label className="font-mono text-sm font-medium text-gray-500">
                Notas
              </label>
              <p className="font-mono">{assignment.notes || "-"}</p>
            </div>
            <div>
              <label className="font-mono text-sm font-medium text-gray-500">
                Creado
              </label>
              <p className="font-mono">
                {new Date(assignment.created_at).toLocaleDateString("es-CL")}
              </p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
