import { describe, it, expect } from "vitest";
import { canAccessModule, canPerformAction, getVisibleModules } from "@/lib/rbac";
import type { UserRole } from "@/types/api";

describe("canAccessModule", () => {
  it("returns true for admin on any module", () => {
    expect(canAccessModule("admin", "/core")).toBe(true);
    expect(canAccessModule("admin", "/partners")).toBe(true);
    expect(canAccessModule("admin", "/operations")).toBe(true);
    expect(canAccessModule("admin", "/nonexistent")).toBe(true);
  });

  it("returns false when role is undefined", () => {
    expect(canAccessModule(undefined, "/core")).toBe(false);
  });

  it("restricts manager from /core", () => {
    expect(canAccessModule("manager", "/core")).toBe(false);
  });

  it("allows manager on /partners", () => {
    expect(canAccessModule("manager", "/partners")).toBe(true);
  });

  it("allows manager on /customers", () => {
    expect(canAccessModule("manager", "/customers")).toBe(true);
  });

  it("allows manager on /inventory", () => {
    expect(canAccessModule("manager", "/inventory")).toBe(true);
  });

  it("allows manager on /planning", () => {
    expect(canAccessModule("manager", "/planning")).toBe(true);
  });

  it("allows manager on /operations", () => {
    expect(canAccessModule("manager", "/operations")).toBe(true);
  });

  it("allows supervisor only on /operations", () => {
    expect(canAccessModule("supervisor", "/operations")).toBe(true);
    expect(canAccessModule("supervisor", "/core")).toBe(false);
    expect(canAccessModule("supervisor", "/partners")).toBe(false);
    expect(canAccessModule("supervisor", "/customers")).toBe(false);
    expect(canAccessModule("supervisor", "/inventory")).toBe(false);
    expect(canAccessModule("supervisor", "/planning")).toBe(false);
  });

  it("technician cannot access any module", () => {
    expect(canAccessModule("technician", "/core")).toBe(false);
    expect(canAccessModule("technician", "/operations")).toBe(false);
    expect(canAccessModule("technician", "/partners")).toBe(false);
  });

  it("returns false for unknown module path", () => {
    expect(canAccessModule("manager", "/no-existe")).toBe(false);
  });

  it("matches prefix paths", () => {
    expect(canAccessModule("manager", "/partners/123")).toBe(true);
    expect(canAccessModule("manager", "/customers/456/edit")).toBe(true);
  });
});

describe("canPerformAction", () => {
  it("admin can perform all actions", () => {
    expect(canPerformAction("admin", "create")).toBe(true);
    expect(canPerformAction("admin", "edit")).toBe(true);
    expect(canPerformAction("admin", "delete")).toBe(true);
  });

  it("returns false when role is undefined", () => {
    expect(canPerformAction(undefined, "create")).toBe(false);
  });

  it("manager can create and edit but not delete", () => {
    expect(canPerformAction("manager", "create")).toBe(true);
    expect(canPerformAction("manager", "edit")).toBe(true);
    expect(canPerformAction("manager", "delete")).toBe(false);
  });

  it("supervisor cannot create, edit, or delete", () => {
    expect(canPerformAction("supervisor", "create")).toBe(false);
    expect(canPerformAction("supervisor", "edit")).toBe(false);
    expect(canPerformAction("supervisor", "delete")).toBe(false);
  });

  it("technician cannot create, edit, or delete", () => {
    expect(canPerformAction("technician", "create")).toBe(false);
    expect(canPerformAction("technician", "edit")).toBe(false);
    expect(canPerformAction("technician", "delete")).toBe(false);
  });
});

describe("getVisibleModules", () => {
  it("admin sees all modules", () => {
    const modules = getVisibleModules("admin");
    expect(modules).toContain("/core");
    expect(modules).toContain("/partners");
    expect(modules).toContain("/customers");
    expect(modules).toContain("/inventory");
    expect(modules).toContain("/planning");
    expect(modules).toContain("/operations");
    expect(modules).toContain("/notifications");
    expect(modules).toContain("/geocoding");
    expect(modules).toContain("/shared");
  });

  it("returns empty array for undefined role", () => {
    expect(getVisibleModules(undefined)).toEqual([]);
  });

  it("manager sees non-core modules", () => {
    const modules = getVisibleModules("manager");
    expect(modules).toContain("/partners");
    expect(modules).toContain("/customers");
    expect(modules).toContain("/inventory");
    expect(modules).toContain("/planning");
    expect(modules).toContain("/operations");
    expect(modules).toContain("/notifications");
    expect(modules).toContain("/geocoding");
    expect(modules).not.toContain("/core");
    expect(modules).not.toContain("/shared");
  });

  it("supervisor sees only /operations", () => {
    const modules = getVisibleModules("supervisor");
    expect(modules).toEqual(["/operations"]);
  });

  it("technician sees no modules", () => {
    const modules = getVisibleModules("technician");
    expect(modules).toEqual([]);
  });
});
