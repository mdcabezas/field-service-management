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

// Hooks that return {items, total} and unwrap via .then(r => r.items)
const UNWRAPPING_SCOPED_HOOKS = [
  { label: "ReportEntries", module: "@/hooks/use-report-entries", parentId: "r1", listPath: "/api/visit-reports/r1/entries", createPath: "/api/report-entries", deletePath: (id: string) => `/api/report-entries/${id}`, fkField: "report_id", listHook: "useReportEntries", createHook: "useCreateReportEntry", deleteHook: "useDeleteReportEntry" },
  { label: "ReportImages", module: "@/hooks/use-report-images", parentId: "r1", listPath: "/api/visit-reports/r1/images", createPath: "/api/report-images", deletePath: (id: string) => `/api/report-images/${id}`, fkField: "report_id", listHook: "useReportImages", createHook: "useCreateReportImage", deleteHook: "useDeleteReportImage" },
  { label: "VisitRentals", module: "@/hooks/use-visit-rentals", parentId: "v1", listPath: "/api/visits/v1/rentals", createPath: "/api/visit-rentals", deletePath: (id: string) => `/api/visit-rentals/${id}`, fkField: "visit_id", listHook: "useVisitRentals", createHook: "useCreateVisitRental", deleteHook: "useDeleteVisitRental" },
  { label: "DailyLoadMaterials", module: "@/hooks/use-daily-load-materials", parentId: "dp1", listPath: "/api/daily-plans/dp1/materials", createPath: "/api/daily-load-materials", deletePath: (id: string) => `/api/daily-load-materials/${id}`, fkField: "daily_plan_id", listHook: "useDailyLoadMaterials", createHook: "useCreateDailyLoadMaterial", deleteHook: "useDeleteDailyLoadMaterial" },
  { label: "DailyLoadTools", module: "@/hooks/use-daily-load-tools", parentId: "dp1", listPath: "/api/daily-plans/dp1/tools", createPath: "/api/daily-load-tools", deletePath: (id: string) => `/api/daily-load-tools/${id}`, fkField: "daily_plan_id", listHook: "useDailyLoadTools", createHook: "useCreateDailyLoadTool", deleteHook: "useDeleteDailyLoadTool" },
  { label: "DailyLoadEPPs", module: "@/hooks/use-daily-load-epps", parentId: "dp1", listPath: "/api/daily-plans/dp1/epps", createPath: "/api/daily-load-epps", deletePath: (id: string) => `/api/daily-load-epps/${id}`, fkField: "daily_plan_id", listHook: "useDailyLoadEPPs", createHook: "useCreateDailyLoadEPP", deleteHook: "useDeleteDailyLoadEPP" },
  { label: "PartnerAgreements", module: "@/hooks/use-partner-agreements", parentId: "p1", listPath: "/api/partners/p1/agreements", createPath: "/api/partner-agreements", deletePath: (id: string) => `/api/partner-agreements/${id}`, fkField: "partner_id", listHook: "usePartnerAgreements", createHook: "useCreatePartnerAgreement", deleteHook: "useDeletePartnerAgreement" },
  { label: "SLAs", module: "@/hooks/use-slas", parentId: "p1", listPath: "/api/partners/p1/slas", createPath: "/api/slas", deletePath: (id: string) => `/api/slas/${id}`, fkField: "partner_id", listHook: "useSLAs", createHook: "useCreateSLA", deleteHook: "useDeleteSLA" },
];

// Hooks that return plain arrays directly (no .items unwrapping)
const ARRAY_SCOPED_HOOKS = [
  { label: "VisitAssignments", module: "@/hooks/use-visit-assignments", parentId: "v1", listPath: "/api/visits/v1/assignments", createPath: "/api/visit-assignments", deletePath: (id: string) => `/api/visit-assignments/${id}`, fkField: "visit_id", listHook: "useVisitAssignments", createHook: "useCreateVisitAssignment", deleteHook: "useDeleteVisitAssignment" },
  { label: "VisitReports", module: "@/hooks/use-visit-reports", parentId: "v1", listPath: "/api/visits/v1/reports", createPath: "/api/visit-reports", deletePath: (id: string) => `/api/visit-reports/${id}`, fkField: "visit_id", listHook: "useVisitReports", createHook: "useCreateVisitReport", deleteHook: "useDeleteVisitReport" },
  { label: "VisitPhotos", module: "@/hooks/use-visit-photos", parentId: "v1", listPath: "/api/visits/v1/photos", createPath: "/api/visit-photos", deletePath: (id: string) => `/api/visit-photos/${id}`, fkField: "visit_id", listHook: "useVisitPhotos", createHook: "useCreateVisitPhoto", deleteHook: "useDeleteVisitPhoto" },
  { label: "VisitMeasurements", module: "@/hooks/use-visit-measurements", parentId: "v1", listPath: "/api/visits/v1/measurements", createPath: "/api/visit-measurements", deletePath: (id: string) => `/api/visit-measurements/${id}`, fkField: "visit_id", listHook: "useVisitMeasurements", createHook: "useCreateVisitMeasurement", deleteHook: "useDeleteVisitMeasurement" },
  { label: "VisitChecklistMaterials", module: "@/hooks/use-visit-checklist-items", parentId: "v1", listPath: "/api/visits/v1/checklist-materials", createPath: "/api/visit-checklist-materials", deletePath: (id: string) => `/api/visit-checklist-materials/${id}`, fkField: "visit_id", listHook: "useVisitChecklistMaterials", createHook: "useCreateVisitChecklistMaterial", deleteHook: "useDeleteVisitChecklistMaterial" },
];

// DailyPlanAssignments: create hook does NOT inject FK (consumer must provide it)
const DAILY_PLAN_ASSIGNMENTS_DEF = { label: "DailyPlanAssignments", module: "@/hooks/use-daily-plan-assignments", parentId: "dp1", listPath: "/api/daily-plans/dp1/assignments", createPath: "/api/daily-plan-assignments", deletePath: (id: string) => `/api/daily-plan-assignments/${id}`, fkField: "daily_plan_id", listHook: "useDailyPlanAssignments", createHook: "useCreateDailyPlanAssignment", deleteHook: "useDeleteDailyPlanAssignment" };

describe.each(UNWRAPPING_SCOPED_HOOKS)("$label scoped hooks (unwrapping)", (def) => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.resetModules();
  });

  it("list returns scoped items", async () => {
    const mod = await import(def.module);
    mockApiFetch.mockResolvedValueOnce({ items: [{ id: "1", [def.fkField]: def.parentId }], total: 1 });
    const wrapper = createWrapper();
    const { result } = renderHook(() => mod[def.listHook](def.parentId), { wrapper });
    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toHaveLength(1);
    expect(mockApiFetch).toHaveBeenCalledWith(def.listPath);
  });

  it("list does not fetch when parent id is empty", async () => {
    const mod = await import(def.module);
    const wrapper = createWrapper();
    renderHook(() => mod[def.listHook](""), { wrapper });
    expect(mockApiFetch).not.toHaveBeenCalled();
  });

  it("create injects parent FK", async () => {
    const mod = await import(def.module);
    mockApiFetch.mockResolvedValueOnce({ id: "new", [def.fkField]: def.parentId });
    const wrapper = createWrapper();
    const { result } = renderHook(() => mod[def.createHook](def.parentId), { wrapper });
    await result.current.mutateAsync({ name: "test" } as any);
    const body = mockApiFetch.mock.calls[0]?.[1]?.body;
    const callBody = JSON.parse(body as string);
    expect(callBody[def.fkField]).toBe(def.parentId);
  });

  it("delete calls DELETE", async () => {
    const mod = await import(def.module);
    mockApiFetch.mockResolvedValueOnce(undefined as any);
    const wrapper = createWrapper();
    const { result } = renderHook(() => mod[def.deleteHook](def.parentId), { wrapper });
    await result.current.mutateAsync("item-id");
    expect(mockApiFetch).toHaveBeenCalledWith(def.deletePath("item-id"), expect.objectContaining({ method: "DELETE" }));
  });
});

describe.each(ARRAY_SCOPED_HOOKS)("$label scoped hooks (array)", (def) => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.resetModules();
  });

  it("list returns scoped items", async () => {
    const mod = await import(def.module);
    mockApiFetch.mockResolvedValueOnce([{ id: "1", [def.fkField]: def.parentId }]);
    const wrapper = createWrapper();
    const { result } = renderHook(() => mod[def.listHook](def.parentId), { wrapper });
    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toHaveLength(1);
    expect(mockApiFetch).toHaveBeenCalledWith(def.listPath);
  });

  it("list does not fetch when parent id is empty", async () => {
    const mod = await import(def.module);
    const wrapper = createWrapper();
    renderHook(() => mod[def.listHook](""), { wrapper });
    expect(mockApiFetch).not.toHaveBeenCalled();
  });

  it("create injects parent FK", async () => {
    const mod = await import(def.module);
    mockApiFetch.mockResolvedValueOnce({ id: "new", [def.fkField]: def.parentId });
    const wrapper = createWrapper();
    const { result } = renderHook(() => mod[def.createHook](def.parentId), { wrapper });
    await result.current.mutateAsync({ name: "test" } as any);
    const body = mockApiFetch.mock.calls[0]?.[1]?.body;
    const callBody = JSON.parse(body as string);
    expect(callBody[def.fkField]).toBe(def.parentId);
  });

  it("delete calls DELETE", async () => {
    const mod = await import(def.module);
    mockApiFetch.mockResolvedValueOnce(undefined as any);
    const wrapper = createWrapper();
    const { result } = renderHook(() => mod[def.deleteHook](def.parentId), { wrapper });
    await result.current.mutateAsync("item-id");
    expect(mockApiFetch).toHaveBeenCalledWith(def.deletePath("item-id"), expect.objectContaining({ method: "DELETE" }));
  });
});

describe("DailyPlanAssignments scoped hooks", () => {
  const def = DAILY_PLAN_ASSIGNMENTS_DEF;

  beforeEach(() => {
    vi.clearAllMocks();
    vi.resetModules();
  });

  it("list returns scoped items", async () => {
    const mod = await import(def.module);
    mockApiFetch.mockResolvedValueOnce({ items: [{ id: "1", daily_plan_id: def.parentId }], total: 1 });
    const wrapper = createWrapper();
    const { result } = renderHook(() => mod[def.listHook](def.parentId), { wrapper });
    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toHaveLength(1);
    expect(mockApiFetch).toHaveBeenCalledWith(def.listPath);
  });

  it("list does not fetch when parent id is empty", async () => {
    const mod = await import(def.module);
    const wrapper = createWrapper();
    renderHook(() => mod[def.listHook](""), { wrapper });
    expect(mockApiFetch).not.toHaveBeenCalled();
  });

  it("create sends data with FK in body (consumer provides FK)", async () => {
    const mod = await import(def.module);
    mockApiFetch.mockResolvedValueOnce({ id: "new", daily_plan_id: def.parentId });
    const wrapper = createWrapper();
    const { result } = renderHook(() => mod[def.createHook](), { wrapper });
    await result.current.mutateAsync({ daily_plan_id: def.parentId, technician_id: "t1" } as any);
    const body = mockApiFetch.mock.calls[0]?.[1]?.body;
    const callBody = JSON.parse(body as string);
    expect(callBody.daily_plan_id).toBe(def.parentId);
  });

  it("delete calls DELETE", async () => {
    const mod = await import(def.module);
    mockApiFetch.mockResolvedValueOnce(undefined as any);
    const wrapper = createWrapper();
    const { result } = renderHook(() => mod[def.deleteHook](), { wrapper });
    await result.current.mutateAsync("item-id");
    expect(mockApiFetch).toHaveBeenCalledWith(def.deletePath("item-id"), expect.objectContaining({ method: "DELETE" }));
  });
});
