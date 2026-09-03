import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import React from "react";

vi.mock("@/lib/api", () => ({ apiFetch: vi.fn() }));
import { apiFetch } from "@/lib/api";
const mockApiFetch = vi.mocked(apiFetch);

function createWrapper() {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  return function Wrapper({ children }: { children: React.ReactNode }) {
    return React.createElement(QueryClientProvider, { client: qc }, children);
  };
}

// Hooks that follow the standard CRUD pattern: no-arg list/get/create/update/delete
// All unwrap list via .then(r => r.items) so query.data is T[]
const CRUD_HOOKS = [
  { label: "Tools", module: "@/hooks/use-tools", listPath: "/api/tools", getPath: (id: string) => `/api/tools/${id}`, createPath: "/api/tools", updatePath: (id: string) => `/api/tools/${id}`, deletePath: (id: string) => `/api/tools/${id}`, listResult: { items: [{ id: "1", name: "Hammer" }], total: 1 }, singleResult: { id: "1", name: "Hammer" }, listHook: "useTools", getHook: "useTool", createHook: "useCreateTool", updateHook: "useUpdateTool", deleteHook: "useDeleteTool" },
  { label: "Materials", module: "@/hooks/use-materials", listPath: "/api/materials", getPath: (id: string) => `/api/materials/${id}`, createPath: "/api/materials", updatePath: (id: string) => `/api/materials/${id}`, deletePath: (id: string) => `/api/materials/${id}`, listResult: { items: [{ id: "1", name: "Pipe" }], total: 1 }, singleResult: { id: "1", name: "Pipe" }, listHook: "useMaterials", getHook: "useMaterial", createHook: "useCreateMaterial", updateHook: "useUpdateMaterial", deleteHook: "useDeleteMaterial" },
  { label: "Vehicles", module: "@/hooks/use-vehicles", listPath: "/api/vehicles", getPath: (id: string) => `/api/vehicles/${id}`, createPath: "/api/vehicles", updatePath: (id: string) => `/api/vehicles/${id}`, deletePath: (id: string) => `/api/vehicles/${id}`, listResult: { items: [{ id: "1", license_plate: "AB-12" }], total: 1 }, singleResult: { id: "1", license_plate: "AB-12" }, listHook: "useVehicles", getHook: "useVehicle", createHook: "useCreateVehicle", updateHook: "useUpdateVehicle", deleteHook: "useDeleteVehicle" },
  { label: "Visits", module: "@/hooks/use-visits", listPath: "/api/visits", getPath: (id: string) => `/api/visits/${id}`, createPath: "/api/visits", updatePath: (id: string) => `/api/visits/${id}`, deletePath: (id: string) => `/api/visits/${id}`, listResult: { items: [{ id: "1", status: "pending" }], total: 1 }, singleResult: { id: "1", status: "pending" }, listHook: "useVisits", getHook: "useVisit", createHook: "useCreateVisit", updateHook: "useUpdateVisit", deleteHook: "useDeleteVisit" },
  { label: "Users", module: "@/hooks/use-users", listPath: "/api/users", getPath: (id: string) => `/api/users/${id}`, createPath: "/api/users", updatePath: (id: string) => `/api/users/${id}`, deletePath: (id: string) => `/api/users/${id}`, listResult: { items: [{ id: "1", email: "a@b.com" }], total: 1 }, singleResult: { id: "1", email: "a@b.com" }, listHook: "useUsers", getHook: "useUser", createHook: "useCreateUser", updateHook: "useUpdateUser", deleteHook: "useDeleteUser" },
  { label: "Routes", module: "@/hooks/use-routes", listPath: "/api/routes", getPath: (id: string) => `/api/routes/${id}`, createPath: "/api/routes", updatePath: (id: string) => `/api/routes/${id}`, deletePath: (id: string) => `/api/routes/${id}`, listResult: { items: [{ id: "1", status: "scheduled" }], total: 1 }, singleResult: { id: "1", status: "scheduled" }, listHook: "useRoutes", getHook: "useRoute", createHook: "useCreateRoute", updateHook: "useUpdateRoute", deleteHook: "useDeleteRoute" },
  { label: "Notifications", module: "@/hooks/use-notifications", listPath: "/api/notifications", getPath: (id: string) => `/api/notifications/${id}`, createPath: "/api/notifications", updatePath: (id: string) => `/api/notifications/${id}`, deletePath: (id: string) => `/api/notifications/${id}`, listResult: { items: [{ id: "1", status: "sent" }], total: 1 }, singleResult: { id: "1", status: "sent" }, listHook: "useNotifications", getHook: "useNotification", createHook: "useCreateNotification", updateHook: "useUpdateNotification", deleteHook: "useDeleteNotification" },
  { label: "GeocodedAddresses", module: "@/hooks/use-geocoded-addresses", listPath: "/api/addresses", getPath: (id: string) => `/api/addresses/${id}`, createPath: "/api/addresses", updatePath: (id: string) => `/api/addresses/${id}`, deletePath: (id: string) => `/api/addresses/${id}`, listResult: { items: [{ id: "1", street: "Main" }], total: 1 }, singleResult: { id: "1", street: "Main" }, listHook: "useGeocodedAddresses", getHook: "useGeocodedAddress", createHook: "useCreateGeocodedAddress", updateHook: "useUpdateGeocodedAddress", deleteHook: "useDeleteGeocodedAddress" },
  { label: "CostRates", module: "@/hooks/use-cost-rates", listPath: "/api/cost-rates", getPath: (id: string) => `/api/cost-rates/${id}`, createPath: "/api/cost-rates", updatePath: (id: string) => `/api/cost-rates/${id}`, deletePath: (id: string) => `/api/cost-rates/${id}`, listResult: { items: [{ id: "1", value: 100 }], total: 1 }, singleResult: { id: "1", value: 100 }, listHook: "useCostRates", getHook: "useCostRate", createHook: "useCreateCostRate", updateHook: "useUpdateCostRate", deleteHook: "useDeleteCostRate" },
  { label: "ChecklistTemplates", module: "@/hooks/use-checklist-templates", listPath: "/api/checklist-templates", getPath: (id: string) => `/api/checklist-templates/${id}`, createPath: "/api/checklist-templates", updatePath: (id: string) => `/api/checklist-templates/${id}`, deletePath: (id: string) => `/api/checklist-templates/${id}`, listResult: { items: [{ id: "1", name: "T1" }], total: 1 }, singleResult: { id: "1", name: "T1" }, listHook: "useChecklistTemplates", getHook: "useChecklistTemplate", createHook: "useCreateChecklistTemplate", updateHook: "useUpdateChecklistTemplate", deleteHook: "useDeleteChecklistTemplate" },
  { label: "MaintenanceSchedules", module: "@/hooks/use-maintenance-schedules", listPath: "/api/maintenance-schedules", getPath: (id: string) => `/api/maintenance-schedules/${id}`, createPath: "/api/maintenance-schedules", updatePath: (id: string) => `/api/maintenance-schedules/${id}`, deletePath: (id: string) => `/api/maintenance-schedules/${id}`, listResult: { items: [{ id: "1", frequency_days: 30 }], total: 1 }, singleResult: { id: "1", frequency_days: 30 }, listHook: "useMaintenanceSchedules", getHook: "useMaintenanceSchedule", createHook: "useCreateMaintenanceSchedule", updateHook: "useUpdateMaintenanceSchedule", deleteHook: "useDeleteMaintenanceSchedule" },
  { label: "TechRoles", module: "@/hooks/use-tech-roles", listPath: "/api/tech-roles", getPath: (id: string) => `/api/tech-roles/${id}`, createPath: "/api/tech-roles", updatePath: (id: string) => `/api/tech-roles/${id}`, deletePath: (id: string) => `/api/tech-roles/${id}`, listResult: { items: [{ id: "1", code: "TECH" }], total: 1 }, singleResult: { id: "1", code: "TECH" }, listHook: "useTechRoles", getHook: "useTechRole", createHook: "useCreateTechRole", updateHook: "useUpdateTechRole", deleteHook: "useDeleteTechRole" },
  { label: "NotificationTemplates", module: "@/hooks/use-notification-templates", listPath: "/api/notification-templates", getPath: (id: string) => `/api/notification-templates/${id}`, createPath: "/api/notification-templates", updatePath: (id: string) => `/api/notification-templates/${id}`, deletePath: (id: string) => `/api/notification-templates/${id}`, listResult: { items: [{ id: "1", type: "email" }], total: 1 }, singleResult: { id: "1", type: "email" }, listHook: "useNotificationTemplates", getHook: "useNotificationTemplate", createHook: "useCreateNotificationTemplate", updateHook: "useUpdateNotificationTemplate", deleteHook: "useDeleteNotificationTemplate" },
  { label: "ReportTemplates", module: "@/hooks/use-report-templates", listPath: "/api/report-templates", getPath: (id: string) => `/api/report-templates/${id}`, createPath: "/api/report-templates", updatePath: (id: string) => `/api/report-templates/${id}`, deletePath: (id: string) => `/api/report-templates/${id}`, listResult: { items: [{ id: "1", name: "R1" }], total: 1 }, singleResult: { id: "1", name: "R1" }, listHook: "useReportTemplates", getHook: "useReportTemplate", createHook: "useCreateReportTemplate", updateHook: "useUpdateReportTemplate", deleteHook: "useDeleteReportTemplate" },
  { label: "EPPItems", module: "@/hooks/use-epp-items", listPath: "/api/epp-items", getPath: (id: string) => `/api/epp-items/${id}`, createPath: "/api/epp-items", updatePath: (id: string) => `/api/epp-items/${id}`, deletePath: (id: string) => `/api/epp-items/${id}`, listResult: { items: [{ id: "1", name: "Helmet" }], total: 1 }, singleResult: { id: "1", name: "Helmet" }, listHook: "useEPPItems", getHook: "useEPPItem", createHook: "useCreateEPPItem", updateHook: "useUpdateEPPItem", deleteHook: "useDeleteEPPItem" },
  { label: "Rentals", module: "@/hooks/use-rentals", listPath: "/api/rentals", getPath: (id: string) => `/api/rentals/${id}`, createPath: "/api/rentals", updatePath: (id: string) => `/api/rentals/${id}`, deletePath: (id: string) => `/api/rentals/${id}`, listResult: { items: [{ id: "1", type: "crane" }], total: 1 }, singleResult: { id: "1", type: "crane" }, listHook: "useRentals", getHook: "useRental", createHook: "useCreateRental", updateHook: "useUpdateRental", deleteHook: "useDeleteRental" },
  { label: "MaintenanceRecords", module: "@/hooks/use-maintenance-records", listPath: "/api/maintenance-records", getPath: (id: string) => `/api/maintenance-records/${id}`, createPath: "/api/maintenance-records", updatePath: (id: string) => `/api/maintenance-records/${id}`, deletePath: (id: string) => `/api/maintenance-records/${id}`, listResult: { items: [{ id: "1", date: "2024-01-01" }], total: 1 }, singleResult: { id: "1", date: "2024-01-01" }, listHook: "useMaintenanceRecords", getHook: "useMaintenanceRecord", createHook: "useCreateMaintenanceRecord", updateHook: "useUpdateMaintenanceRecord", deleteHook: "useDeleteMaintenanceRecord" },
  { label: "DailyPlans", module: "@/hooks/use-daily-plans", listPath: "/api/daily-plans", getPath: (id: string) => `/api/daily-plans/${id}`, createPath: "/api/daily-plans", updatePath: (id: string) => `/api/daily-plans/${id}`, deletePath: (id: string) => `/api/daily-plans/${id}`, listResult: { items: [{ id: "1", name: "Plan1" }], total: 1 }, singleResult: { id: "1", name: "Plan1" }, listHook: "useDailyPlans", getHook: "useDailyPlan", createHook: "useCreateDailyPlan", updateHook: "useUpdateDailyPlan", deleteHook: "useDeleteDailyPlan" },
  { label: "Partners", module: "@/hooks/use-partners", listPath: "/api/partners", getPath: (id: string) => `/api/partners/${id}`, createPath: "/api/partners", updatePath: (id: string) => `/api/partners/${id}`, deletePath: (id: string) => `/api/partners/${id}`, listResult: { items: [{ id: "1", name: "P1" }], total: 1 }, singleResult: { id: "1", name: "P1" }, listHook: "usePartners", getHook: "usePartner", createHook: "useCreatePartner", updateHook: "useUpdatePartner", deleteHook: "useDeletePartner" },
  { label: "VisitTypes", module: "@/hooks/use-visit-types", listPath: "/api/visit-types", getPath: (id: string) => `/api/visit-types/${id}`, createPath: "/api/visit-types", updatePath: (id: string) => `/api/visit-types/${id}`, deletePath: (id: string) => `/api/visit-types/${id}`, listResult: { items: [{ id: "1", code: "install" }], total: 1 }, singleResult: { id: "1", code: "install" }, listHook: "useVisitTypes", getHook: "useVisitType", createHook: "useCreateVisitType", updateHook: "useUpdateVisitType", deleteHook: "useDeleteVisitType" },
  { label: "VehicleAssignments", module: "@/hooks/use-vehicle-assignments", listPath: "/api/vehicle-assignments", getPath: (id: string) => `/api/vehicle-assignments/${id}`, createPath: "/api/vehicle-assignments", updatePath: (id: string) => `/api/vehicle-assignments/${id}`, deletePath: (id: string) => `/api/vehicle-assignments/${id}`, listResult: { items: [{ id: "1", vehicle_id: "v1" }], total: 1 }, singleResult: { id: "1", vehicle_id: "v1" }, listHook: "useVehicleAssignments", getHook: null, createHook: "useCreateVehicleAssignment", updateHook: null, deleteHook: "useDeleteVehicleAssignment" },
];

describe.each(CRUD_HOOKS)("$label CRUD hooks", (def) => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.resetModules();
  });

  it("list returns items", async () => {
    const mod = await import(def.module);
    mockApiFetch.mockResolvedValueOnce(def.listResult);
    const wrapper = createWrapper();
    const { result } = renderHook(() => mod[def.listHook](), { wrapper });
    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toBeDefined();
    expect(mockApiFetch).toHaveBeenCalledWith(def.listPath);
  });

  if (def.getHook) {
    it("get fetches single item", async () => {
      const mod = await import(def.module);
      mockApiFetch.mockResolvedValueOnce(def.singleResult);
      const wrapper = createWrapper();
      const { result } = renderHook(() => mod[def.getHook]("test-id"), { wrapper });
      await waitFor(() => expect(result.current.isSuccess).toBe(true));
      expect(mockApiFetch).toHaveBeenCalledWith(def.getPath("test-id"));
    });

    it("get does not fetch when id is empty", async () => {
      const mod = await import(def.module);
      const wrapper = createWrapper();
      renderHook(() => mod[def.getHook](""), { wrapper });
      expect(mockApiFetch).not.toHaveBeenCalled();
    });
  }

  it("create calls POST", async () => {
    const mod = await import(def.module);
    mockApiFetch.mockResolvedValueOnce(def.singleResult);
    const wrapper = createWrapper();
    const { result } = renderHook(() => mod[def.createHook](), { wrapper });
    await result.current.mutateAsync({ name: "test" } as any);
    expect(mockApiFetch).toHaveBeenCalledWith(def.createPath, expect.objectContaining({ method: "POST" }));
  });

  if (def.updateHook) {
    it("update calls PUT", async () => {
      const mod = await import(def.module);
      mockApiFetch.mockResolvedValueOnce(def.singleResult);
      const wrapper = createWrapper();
      const { result } = renderHook(() => mod[def.updateHook](), { wrapper });
      await result.current.mutateAsync({ id: "test-id", data: { name: "updated" } });
      expect(mockApiFetch).toHaveBeenCalledWith(def.updatePath("test-id"), expect.objectContaining({ method: "PUT" }));
    });
  }

  it("delete calls DELETE", async () => {
    const mod = await import(def.module);
    mockApiFetch.mockResolvedValueOnce(undefined as any);
    const wrapper = createWrapper();
    const { result } = renderHook(() => mod[def.deleteHook](), { wrapper });
    await result.current.mutateAsync("test-id");
    expect(mockApiFetch).toHaveBeenCalledWith(def.deletePath("test-id"), expect.objectContaining({ method: "DELETE" }));
  });
});
