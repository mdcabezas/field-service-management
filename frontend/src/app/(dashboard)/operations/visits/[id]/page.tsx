"use client";

import { useState } from "react";
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
import { useVisitMaterialUsages } from "@/hooks/use-visit-material-usages";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import { toast } from "sonner";
import type {
  VisitAssignment,
  VisitChecklistMaterial,
  VisitMaterialUsage,
  VisitMeasurement,
  VisitCheckpoint,
  VisitPhoto,
  VisitReport,
  VisitSLATracking,
} from "@/types/api";

const assignmentColumns: ColumnDef<VisitAssignment, unknown>[] = [
  {
    accessorKey: "tech_name",
    header: "Técnico",
    cell: ({ row }) => (
      <span className="font-mono text-sm">{row.original.tech_name || row.original.tech_id}</span>
    ),
  },
  {
    accessorKey: "role_name",
    header: "Rol",
    cell: ({ row }) => (
      <span className="font-mono text-sm">{row.original.role_name || row.original.role_id}</span>
    ),
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

const materialUsageColumns: ColumnDef<VisitMaterialUsage, unknown>[] = [
  {
    accessorKey: "material_name",
    header: "Material",
    cell: ({ row }) => row.original.material_name || row.original.material_id,
  },
  {
    accessorKey: "quantity",
    header: "Cantidad",
  },
  {
    accessorKey: "notes",
    header: "Notas",
  },
  {
    accessorKey: "created_at",
    header: "Fecha",
    cell: ({ row }) => {
      if (!row.original.created_at) return "-";
      return new Date(row.original.created_at).toLocaleDateString("es-CL");
    },
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
  const [selectedPhoto, setSelectedPhoto] = useState<VisitPhoto | null>(null);
  const [failedPhotos, setFailedPhotos] = useState<Set<string>>(new Set());

  const { data: assignments } = useVisitAssignments(id);
  const { data: checklistMaterials } = useVisitChecklistMaterials(id);
  const { data: measurements } = useVisitMeasurements(id);
  const { data: checkpoints } = useVisitCheckpoints(id);
  const { data: photos, isLoading: photosLoading, error: photosError } = useVisitPhotos(id);
  const { data: reports } = useVisitReports(id);
  const { data: slaTrackings } = useVisitSLATrackings(id);
  const { data: materialUsages } = useVisitMaterialUsages(id);

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
          <TabsTrigger value="material-usages">Materiales</TabsTrigger>
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

        <TabsContent value="material-usages">
          <Card>
            <CardHeader>
              <CardTitle className="font-mono text-lg uppercase">
                Materiales Utilizados
              </CardTitle>
            </CardHeader>
            <CardContent>
              <DataTable
                columns={materialUsageColumns}
                data={materialUsages || []}
                searchPlaceholder="Buscar por material..."
                searchColumn="material_name"
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
              {photosLoading && <div className="font-mono text-center py-4">Cargando fotos...</div>}
              {photosError && (
                <div className="font-mono text-center py-4 text-red-600">
                  Error al cargar fotos: {photosError.message || "Error desconocido"}
                </div>
              )}
              {!photosLoading && !photosError && photos && photos.length === 0 && (
                <div className="font-mono text-center py-8 text-gray-500">No hay fotos registradas</div>
              )}
              {!photosLoading && !photosError && photos && photos.length > 0 && (
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                  {photos.map((photo) => (
                    <div
                      key={photo.id}
                      className={`border-2 border-black relative group ${photo.thumbnail_url && !failedPhotos.has(photo.id) ? "cursor-pointer" : ""}`}
                      onClick={photo.thumbnail_url && !failedPhotos.has(photo.id) ? () => setSelectedPhoto(photo) : undefined}
                    >
                      {photo.thumbnail_url && !failedPhotos.has(photo.id) ? (
                        <img
                          src={photo.thumbnail_url}
                          alt={`Foto ${photo.id.slice(0, 8)}`}
                          className="w-full h-48 object-cover"
                          loading="lazy"
                          onError={() => setFailedPhotos(prev => new Set(prev).add(photo.id))}
                        />
                      ) : (
                        <div className="w-full h-48 flex items-center justify-center bg-gray-100 border-2 border-dashed border-gray-400">
                          <svg
                            className="w-12 h-12 text-gray-400"
                            viewBox="0 0 16 16"
                            xmlns="http://www.w3.org/2000/svg"
                          >
                            <path
                              d="m 6.5 0 c -0.265625 0 -0.519531 0.105469 -0.707031 0.292969 l -1.707031 1.707031 h -1.023438 l -1.53125 -1.53125 l -1.0625 1.0625 l 14 14 l 1.0625 -1.0625 l -0.386719 -0.386719 c 0.527344 -0.539062 0.855469 -1.277343 0.855469 -2.082031 v -7 c 0 -1.644531 -1.355469 -3 -3 -3 h -1.085938 l -1.707031 -1.707031 c -0.1875 -0.1875 -0.441406 -0.292969 -0.707031 -0.292969 z m 0.414062 2 h 2.171876 l 1.707031 1.707031 c 0.1875 0.1875 0.441406 0.292969 0.707031 0.292969 h 1.5 c 0.570312 0 1 0.429688 1 1 v 7 c 0 0.269531 -0.097656 0.503906 -0.257812 0.679688 l -2.4375 -2.4375 c 0.4375 -0.640626 0.695312 -1.414063 0.695312 -2.242188 c 0 -2.199219 -1.800781 -4 -4 -4 c -0.828125 0 -1.601562 0.257812 -2.242188 0.695312 l -0.808593 -0.808593 c 0.09375 -0.046875 0.183593 -0.105469 0.257812 -0.179688 z m -6.492187 1.484375 c -0.265625 0.445313 -0.421875 0.964844 -0.421875 1.515625 v 7 c 0 1.644531 1.355469 3 3 3 h 8.9375 l -2 -2 h -6.9375 c -0.570312 0 -1 -0.429688 -1 -1 v -6.9375 z m 7.578125 2.515625 c 1.117188 0 2 0.882812 2 2 c 0 0.277344 -0.058594 0.539062 -0.15625 0.78125 l -2.625 -2.625 c 0.242188 -0.097656 0.503906 -0.15625 0.78125 -0.15625 z m -3.90625 1.15625 c -0.058594 0.273438 -0.09375 0.554688 -0.09375 0.84375 c 0 2.199219 1.800781 4 4 4 c 0.289062 0 0.570312 -0.035156 0.84375 -0.09375 z m 0 0"
                              fill="currentColor"
                            />
                          </svg>
                        </div>
                      )}
                      {photo.thumbnail_url && !failedPhotos.has(photo.id) && (
                      <div className="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-end p-2">
                        <div className="font-mono text-xs text-white w-full">
                          {photo.stage && <span>Etapa: {photo.stage}</span>}
                          {photo.finding_type && (
                            <span className="ml-2">Hallazgo: {photo.finding_type}</span>
                          )}
                        </div>
                      </div>
                      )}
                    </div>
                  ))}
                </div>
              )}
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

      <Dialog open={!!selectedPhoto} onOpenChange={(open) => !open && setSelectedPhoto(null)}>
        <DialogContent className="max-w-4xl max-h-[90vh] p-0 sm:max-w-4xl">
          {selectedPhoto && (
            <>
              <DialogHeader className="p-4 border-b-2 border-black">
                <DialogTitle className="font-mono text-sm">
                  {selectedPhoto.stage || "Foto"} — {selectedPhoto.id.slice(0, 8)}
                </DialogTitle>
              </DialogHeader>
              <div className="p-4 flex justify-center overflow-auto max-h-[70vh]">
                <img
                  src={selectedPhoto.url}
                  alt={`Foto ${selectedPhoto.id.slice(0, 8)}`}
                  className="max-w-full h-auto border-2 border-black"
                />
              </div>
              <DialogFooter className="border-t-2 border-black px-4 py-2">
                <div className="font-mono text-xs text-white flex gap-4">
                  {selectedPhoto.stage && <span>Etapa: {selectedPhoto.stage}</span>}
                  {selectedPhoto.finding_type && (
                    <span>Hallazgo: {selectedPhoto.finding_type}</span>
                  )}
                  {selectedPhoto.timestamp && (
                    <span>{new Date(selectedPhoto.timestamp).toLocaleString("es-CL")}</span>
                  )}
                </div>
              </DialogFooter>
            </>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}
