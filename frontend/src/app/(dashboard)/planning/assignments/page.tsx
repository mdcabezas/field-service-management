"use client";

import { useRouter } from "next/navigation";
import { type ColumnDef } from "@tanstack/react-table";
import { useAllDailyPlanAssignments } from "@/hooks/use-daily-plan-assignments";
import { useDailyPlans } from "@/hooks/use-daily-plans";
import { useAllTechnicians } from "@/hooks/use-technicians";
import { useTechRoles } from "@/hooks/use-tech-roles";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/stores/auth-store";
import { canPerformAction } from "@/lib/rbac";
import type { DailyPlanAssignment } from "@/types/api";

export default function AssignmentsPage() {
  const router = useRouter();
  const { data: assignments, isLoading } = useAllDailyPlanAssignments();
  const { data: dailyPlans } = useDailyPlans();
  const { data: technicians } = useAllTechnicians();
  const { data: techRoles } = useTechRoles();
  const { user } = useAuthStore();

  const columns: ColumnDef<DailyPlanAssignment, unknown>[] = [
    {
      accessorKey: "daily_plan_id",
      header: "Plan Diario",
      cell: ({ row }) => {
        const plan = dailyPlans?.find((p) => p.id === row.original.daily_plan_id);
        return plan ? plan.date : row.original.daily_plan_id;
      },
    },
    {
      accessorKey: "tech_id",
      header: "Técnico",
      cell: ({ row }) => {
        const tech = technicians?.find((t) => t.id === row.original.tech_id);
        return tech ? tech.name : row.original.tech_id;
      },
    },
    {
      accessorKey: "role_id",
      header: "Rol",
      cell: ({ row }) => {
        const role = techRoles?.find((r) => r.id === row.original.role_id);
        return role ? role.name : row.original.role_id;
      },
    },];

  const handleRowClick = (row: DailyPlanAssignment) => {
    router.push(`/planning/assignments/${row.id}`);
  };

  if (isLoading) {
    return <div className="font-mono">Cargando...</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
          Asignaciones
        </h1>
        {canPerformAction(user?.role, "create") && (
          <Button onClick={() => router.push("/planning/assignments/new")}>
            + Crear Asignación
          </Button>
        )}
      </div>

      <DataTable
        columns={columns}
        data={assignments || []}
        searchPlaceholder="Buscar por plan o visita..."
        searchColumn="daily_plan_id"
        onRowClick={handleRowClick}
      />
    </div>
  );
}