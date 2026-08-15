"use client";

import { useState } from "react";
import { DataTable } from "@/components/data-table/data-table";
import { type ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { useSLAs, useCreateSLA, useDeleteSLA } from "@/hooks/use-slas";
import type { SLA } from "@/types/api";

interface PartnerSLAsTableProps {
  partnerId: string;
}

interface SLADraft {
  name: string;
  description: string;
  response_hours: string;
  resolution_hours: string;
  compliance_target: string;
  active: boolean;
  valid_from: string;
  valid_until: string;
}

const emptyDraft: SLADraft = {
  name: "",
  description: "",
  response_hours: "",
  resolution_hours: "",
  compliance_target: "",
  active: true,
  valid_from: "",
  valid_until: "",
};

export function PartnerSLAsTable({ partnerId }: PartnerSLAsTableProps) {
  const { data: slas = [], isLoading, refetch } = useSLAs(partnerId);
  const createSLA = useCreateSLA(partnerId);
  const deleteSLA = useDeleteSLA(partnerId);

  const [showDialog, setShowDialog] = useState(false);
  const [draft, setDraft] = useState<SLADraft>(emptyDraft);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleOpenDialog = () => {
    setDraft(emptyDraft);
    setError(null);
    setShowDialog(true);
  };

  const handleDelete = async (id: string) => {
    if (!confirm("¿Estás seguro de eliminar este SLA?")) return;
    try {
      await deleteSLA.mutateAsync(id);
      refetch();
    } catch (err) {
      console.error("Error deleting SLA:", err);
    }
  };

  const handleSubmit = async () => {
    if (!draft.name.trim()) {
      setError("El nombre es requerido");
      return;
    }
    if (!draft.valid_from) {
      setError("La fecha de inicio es requerida");
      return;
    }
    setIsSubmitting(true);
    setError(null);
    try {
      await createSLA.mutateAsync({
        name: draft.name.trim(),
        description: draft.description || undefined,
        response_hours: draft.response_hours ? Number(draft.response_hours) : undefined,
        resolution_hours: draft.resolution_hours ? Number(draft.resolution_hours) : undefined,
        compliance_target: draft.compliance_target ? Number(draft.compliance_target) : undefined,
        active: draft.active,
        valid_from: `${draft.valid_from}:00Z`,
        valid_until: draft.valid_until ? `${draft.valid_until}:00Z` : undefined,
      });
      setShowDialog(false);
      refetch();
    } catch (err) {
      setError("Error al crear el SLA");
      console.error(err);
    } finally {
      setIsSubmitting(false);
    }
  };

  const set = (key: keyof SLADraft, value: string | boolean) =>
    setDraft((d) => ({ ...d, [key]: value }));

  const columns: ColumnDef<SLA, unknown>[] = [
    {
      accessorKey: "name",
      header: "Nombre",
      cell: ({ row }) => (
        <div className="font-mono text-sm">
          <div className="font-medium">{row.original.name}</div>
          {row.original.description && (
            <div className="text-gray-500">{row.original.description}</div>
          )}
        </div>
      ),
    },
    {
      accessorKey: "response_hours",
      header: "Respuesta (h)",
      cell: ({ row }) => <span className="font-mono text-sm">{row.original.response_hours ?? "—"}</span>,
    },
    {
      accessorKey: "resolution_hours",
      header: "Resolución (h)",
      cell: ({ row }) => <span className="font-mono text-sm">{row.original.resolution_hours ?? "—"}</span>,
    },
    {
      accessorKey: "compliance_target",
      header: "Cumplimiento",
      cell: ({ row }) => (
        <span className="font-mono text-sm">
          {row.original.compliance_target != null
            ? `${(row.original.compliance_target * 100).toFixed(0)}%`
            : "—"}
        </span>
      ),
    },
    {
      accessorKey: "active",
      header: "Estado",
      cell: ({ row }) => (
        <Badge variant={row.original.active ? "default" : "outline"} className="font-mono text-xs">
          {row.original.active ? "Activo" : "Inactivo"}
        </Badge>
      ),
    },
    {
      accessorKey: "valid_from",
      header: "Vigencia",
      cell: ({ row }) => (
        <span className="font-mono text-sm">
          {row.original.valid_from
            ? new Date(row.original.valid_from).toLocaleDateString("es-CL")
            : "—"}
          {row.original.valid_until
            ? ` → ${new Date(row.original.valid_until).toLocaleDateString("es-CL")}`
            : ""}
        </span>
      ),
    },
    {
      id: "actions",
      header: "Acciones",
      cell: ({ row }) => (
        <Button
          variant="destructive"
          size="sm"
          onClick={() => handleDelete(row.original.id)}
          className="h-8 px-3 font-mono text-xs"
        >
          Eliminar
        </Button>
      ),
    },
  ];

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64 font-mono text-gray-500">
        Cargando SLAs...
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="font-mono text-lg font-semibold">SLAs del Socio</h3>
        <Button onClick={handleOpenDialog}>+ Agregar SLA</Button>
      </div>

      <DataTable
        columns={columns}
        data={slas}
        searchPlaceholder="Buscar por nombre..."
        searchColumn="name"
      />

      <Dialog open={showDialog} onOpenChange={setShowDialog}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>Agregar SLA</DialogTitle>
          </DialogHeader>

          <div className="space-y-4 p-4">
            <div className="space-y-2">
              <Label>Nombre *</Label>
              <Input
                value={draft.name}
                onChange={(e) => set("name", e.target.value)}
                placeholder="SLA Prioridad Alta"
              />
            </div>

            <div className="space-y-2">
              <Label>Descripción</Label>
              <Input
                value={draft.description}
                onChange={(e) => set("description", e.target.value)}
                placeholder="Respuesta en 4h, resolución en 24h"
              />
            </div>

            <div className="grid grid-cols-3 gap-4">
              <div className="space-y-2">
                <Label>Respuesta (h)</Label>
                <Input
                  type="number"
                  value={draft.response_hours}
                  onChange={(e) => set("response_hours", e.target.value)}
                  placeholder="24"
                />
              </div>
              <div className="space-y-2">
                <Label>Resolución (h)</Label>
                <Input
                  type="number"
                  value={draft.resolution_hours}
                  onChange={(e) => set("resolution_hours", e.target.value)}
                  placeholder="72"
                />
              </div>
              <div className="space-y-2">
                <Label>Cumplimiento</Label>
                <Input
                  type="number"
                  step="0.01"
                  value={draft.compliance_target}
                  onChange={(e) => set("compliance_target", e.target.value)}
                  placeholder="0.95"
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Vigencia desde *</Label>
                <Input
                  type="datetime-local"
                  value={draft.valid_from}
                  onChange={(e) => set("valid_from", e.target.value)}
                />
              </div>
              <div className="space-y-2">
                <Label>Vigencia hasta</Label>
                <Input
                  type="datetime-local"
                  value={draft.valid_until}
                  onChange={(e) => set("valid_until", e.target.value)}
                />
              </div>
            </div>

            <label className="flex items-center gap-2 font-mono text-sm">
              <input
                type="checkbox"
                checked={draft.active}
                onChange={(e) => set("active", e.target.checked)}
              />
              Activo
            </label>

            {error && (
              <div className="p-3 bg-red-50 border-2 border-red-500 text-red-700 font-mono text-sm">
                {error}
              </div>
            )}
          </div>

          <DialogFooter className="border-t border-black bg-white p-4">
            <Button variant="outline" onClick={() => setShowDialog(false)} disabled={isSubmitting}>
              Cancelar
            </Button>
            <Button onClick={handleSubmit} disabled={isSubmitting}>
              {isSubmitting ? "Guardando..." : "Agregar"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
