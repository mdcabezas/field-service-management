"use client";

import { useRouter, useParams } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useVisit, useDeleteVisit } from "@/hooks/use-visits";
import { useVisitAssignments } from "@/hooks/use-visit-assignments";
import { useVisitChecklistMaterials } from "@/hooks/use-visit-checklist-items";
import { useVisitMeasurements } from "@/hooks/use-visit-measurements";
import { useVisitCheckpoints } from "@/hooks/use-visit-checkpoints";
import { useVisitPhotos } from "@/hooks/use-visit-photos";
import { useVisitReports } from "@/hooks/use-visit-reports";
import { useVisitSLATrackings } from "@/hooks/use-visit-sla-trackings";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";
import type {
  VisitAssignment,
  VisitChecklistMaterial,
  VisitMeasurement,
  VisitCheckpoint,
  VisitPhoto,
  VisitReport,
  VisitSLATracking,
} from "@/types/api";

const assignmentColumns: ColumnDef<VisitAssignment, unknown>[] = [
  {
    accessorKey: "tech_id",
    header: "Técnico",
  },
  {
    accessorKey: "role_id",
    header: "Rol",
  },
];

const checklistColumns: ColumnDef<VisitChecklistMaterial, unknown>[] = [
  {
    accessorKey: "material_id",
    header: "Material",
  },
  {
    accessorKey: "planned_quantity",
    header: "Cantidad",
  },
  {
    accessorKey: "confirmed",
    header: "Confirmado",
    cell: ({ row }) => <span className="uppercase">{row.original.confirmed ? "Sí" : "No"}</span>,
  },
];

const measurementColumns: ColumnDef<VisitMeasurement, unknown>[] = [
  {
    accessorKey: "type",
    header: "Tipo",
  },
  {
    accessorKey: "value",
    header: "Valor",
  },
  {
    accessorKey: "unit",
    header: "Unidad",
  },
  {
    accessorKey: "notes",
    header: "Notas",
  },
];

const checkpointColumns: ColumnDef<VisitCheckpoint, unknown>[] = [
  {
    accessorKey: "type",
    header: "Tipo",
  },
  {
    accessorKey: "timestamp",
    header: "Fecha",
    cell: ({ row }) => {
      if (!row.original.timestamp) return "-";
      return new Date(row.original.timestamp).toLocaleDateString("es-CL");
    },
  },
];

const photoColumns: ColumnDef<VisitPhoto, unknown>[] = [
  {
    accessorKey: "url",
    header: "URL",
    cell: ({ row }) => (
      <a
        href={row.original.url}
        target="_blank"
        rel="noopener noreferrer"
        className="text-blue-600 underline"
      >
        Ver foto
      </a>
    ),
  },
  {
    accessorKey: "stage",
    header: "Etapa",
  },
  {
    accessorKey: "finding_type",
    header: "Tipo de hallazgo",
  },
];

const reportColumns: ColumnDef<VisitReport, unknown>[] = [
  {
    accessorKey: "report_template_id",
    header: "Plantilla",
  },
  {
    accessorKey: "source",
    header: "Fuente",
  },
  {
    accessorKey: "recorded_at",
    header: "Fecha",
    cell: ({ row }) => {
      if (!row.original.recorded_at) return "-";
      return new Date(row.original.recorded_at).toLocaleDateString("es-CL");
    },
  },
];

const slaColumns: ColumnDef<VisitSLATracking, unknown>[] = [
  {
    accessorKey: "sla_id",
    header: "SLA",
  },
  {
    accessorKey: "requested_at",
    header: "Solicitado",
    cell: ({ row }) => {
      if (!row.original.requested_at) return "-";
      return new Date(row.original.requested_at).toLocaleDateString("es-CL");
    },
  },
  {
    accessorKey: "resolved_at",
    header: "Resuelto",
    cell: ({ row }) => {
      if (!row.original.resolved_at) return "-";
      return new Date(row.original.resolved_at).toLocaleDateString("es-CL");
    },
  },
];

export default function VisitDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const { data: visit, isLoading } = useVisit(id);
  const deleteVisit = useDeleteVisit();
  const { user: currentUser } = useAuthStore();

  const { data: assignments } = useVisitAssignments(id);
  const { data: checklistMaterials } = useVisitChecklistMaterials(id);
  const { data: measurements } = useVisitMeasurements(id);
  const { data: checkpoints } = useVisitCheckpoints(id);
  const { data: photos } = useVisitPhotos(id);
  const { data: reports } = useVisitReports(id);
  const { data: slaTrackings } = useVisitSLATrackings(id);

  const handleDelete = async () => {
    if (confirm("¿Estás seguro de eliminar esta visita?")) {
      try {
        await deleteVisit.mutateAsync(id);
        toast.success("Visita eliminada");
        router.push("/operations/visits");
      } catch (error) {
        toast.error("Error al eliminar visita");
      }
    }
  };

  if (isLoading) {
    return <div className="font-mono">Cargando...</div>;
  }

  if (!visit) {
    return <div className="font-mono">Visita no encontrada</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Visita - {new Date(visit.scheduled_at).toLocaleDateString("es-CL")}
        </h1>
        {canPerformAction(currentUser?.role, "delete") && (
          <Button variant="destructive" onClick={handleDelete}>
            Eliminar
          </Button>
        )}
      </div>

      <Tabs defaultValue="info">
        <TabsList>
          <TabsTrigger value="info">Información</TabsTrigger>
          <TabsTrigger value="assignments">Asignaciones</TabsTrigger>
          <TabsTrigger value="checklist">Checklist</TabsTrigger>
          <TabsTrigger value="measurements">Mediciones</TabsTrigger>
          <TabsTrigger value="checkpoints">Checkpoints</TabsTrigger>
          <TabsTrigger value="photos">Fotos</TabsTrigger>
          <TabsTrigger value="reports">Reportes</TabsTrigger>
          <TabsTrigger value="vehicle-assignments">Vehículos</TabsTrigger>
          <TabsTrigger value="slas">SLAs</TabsTrigger>
        </TabsList>

        <TabsContent value="info">
          <Card>
            <CardHeader>
              <CardTitle className="font-mono text-lg uppercase">
                Información de la Visita
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="font-mono text-sm font-medium text-gray-500">
                    Fecha
                  </label>
                  <p className="font-mono">
                    {new Date(visit.scheduled_at).toLocaleDateString("es-CL")}
                  </p>
                </div>
                <div>
                  <label className="font-mono text-sm font-medium text-gray-500">
                    Prioridad
                  </label>
                  <p className="font-mono uppercase">{visit.priority}</p>
                </div>
                <div>
                  <label className="font-mono text-sm font-medium text-gray-500">
                    Estado
                  </label>
                  <p className="font-mono uppercase">{visit.status}</p>
                </div>
                <div>
                  <label className="font-mono text-sm font-medium text-gray-500">
                    Notas
                  </label>
                  <p className="font-mono">{visit.notes || "-"}</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="assignments">
          <Card>
            <CardHeader>
              <CardTitle className="font-mono text-lg uppercase">
                Asignaciones
              </CardTitle>
            </CardHeader>
            <CardContent>
              <DataTable
                columns={assignmentColumns}
                data={assignments || []}
                searchPlaceholder="Buscar por técnico..."
                searchColumn="tech_id"
              />
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="checklist">
          <Card>
            <CardHeader>
              <CardTitle className="font-mono text-lg uppercase">
                Checklist de Materiales
              </CardTitle>
            </CardHeader>
            <CardContent>
              <DataTable
                columns={checklistColumns}
                data={checklistMaterials || []}
                searchPlaceholder="Buscar por material..."
                searchColumn="material_id"
              />
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="measurements">
          <Card>
            <CardHeader>
              <CardTitle className="font-mono text-lg uppercase">
                Mediciones
              </CardTitle>
            </CardHeader>
            <CardContent>
              <DataTable
                columns={measurementColumns}
                data={measurements || []}
                searchPlaceholder="Buscar por tipo..."
                searchColumn="type"
              />
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="checkpoints">
          <Card>
            <CardHeader>
              <CardTitle className="font-mono text-lg uppercase">
                Checkpoints
              </CardTitle>
            </CardHeader>
            <CardContent>
              <DataTable
                columns={checkpointColumns}
                data={checkpoints || []}
                searchPlaceholder="Buscar por nombre..."
                searchColumn="type"
              />
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="photos">
          <Card>
            <CardHeader>
              <CardTitle className="font-mono text-lg uppercase">
                Fotos
              </CardTitle>
            </CardHeader>
            <CardContent>
              <DataTable
                columns={photoColumns}
                data={photos || []}
                searchPlaceholder="Buscar por descripción..."
                searchColumn="stage"
              />
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="reports">
          <Card>
            <CardHeader>
              <CardTitle className="font-mono text-lg uppercase">
                Reportes
              </CardTitle>
            </CardHeader>
            <CardContent>
              <DataTable
                columns={reportColumns}
                data={reports || []}
                searchPlaceholder="Buscar por plantilla..."
                searchColumn="report_template_id"
              />
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="vehicle-assignments">
          <Card>
            <CardHeader>
              <CardTitle className="font-mono text-lg uppercase">
                Asignaciones de Vehículos
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="font-mono text-gray-500">
                Asignaciones de vehículos - Próximamente
              </p>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="slas">
          <Card>
            <CardHeader>
              <CardTitle className="font-mono text-lg uppercase">
                Seguimiento SLA
              </CardTitle>
            </CardHeader>
            <CardContent>
              <DataTable
                columns={slaColumns}
                data={slaTrackings || []}
                searchPlaceholder="Buscar por SLA..."
                searchColumn="sla_id"
              />
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}
