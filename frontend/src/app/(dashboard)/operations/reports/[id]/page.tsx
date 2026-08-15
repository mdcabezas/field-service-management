"use client";

import { useRouter, useParams } from "next/navigation";
import { useVisitReport, useDeleteVisitReportGlobal } from "@/hooks/use-all-visit-reports";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";

export default function ReportDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const { data: report, isLoading } = useVisitReport(id);
  const deleteReport = useDeleteVisitReportGlobal();
  const { user: currentUser } = useAuthStore();

  const handleDelete = async () => {
    if (confirm("¿Estás seguro de eliminar este reporte?")) {
      try {
        await deleteReport.mutateAsync(id);
        toast.success("Reporte eliminado");
        router.push("/operations/reports");
      } catch (error) {
        toast.error("Error al eliminar reporte");
      }
    }
  };

  if (isLoading) {
    return <div className="font-mono">Cargando...</div>;
  }

  if (!report) {
    return <div className="font-mono">Reporte no encontrado</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Reporte
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
            Información del Reporte
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="font-mono text-sm font-medium text-gray-500">
                Visita
              </label>
              <p className="font-mono">{report.visit_label || report.visit_id}</p>
            </div>
            <div>
              <label className="font-mono text-sm font-medium text-gray-500">
                Plantilla
              </label>
              <p className="font-mono">{report.report_template_name || report.report_template_id}</p>
            </div>
            <div>
              <label className="font-mono text-sm font-medium text-gray-500">
                Fuente
              </label>
              <p className="font-mono uppercase">{report.source || "-"}</p>
            </div>
            <div>
              <label className="font-mono text-sm font-medium text-gray-500">
                Registrado
              </label>
              <p className="font-mono">
                {report.recorded_at
                  ? new Date(report.recorded_at).toLocaleDateString("es-CL")
                  : "-"}
              </p>
            </div>
            <div>
              <label className="font-mono text-sm font-medium text-gray-500">
                Creado
              </label>
              <p className="font-mono">
                {new Date(report.created_at).toLocaleDateString("es-CL")}
              </p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
