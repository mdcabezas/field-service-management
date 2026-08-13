// ============================================================================
// RBAC — Role-based access control utilities
// ============================================================================

import type { UserRole } from "@/types/api";

// Module access by role
const moduleAccess: Record<string, UserRole[]> = {
  "/core": ["admin"],
  "/partners": ["admin", "manager"],
  "/customers": ["admin", "manager"],
  "/inventory": ["admin", "manager"],
  "/planning": ["admin", "manager"],
  "/operations": ["admin", "manager", "supervisor"],
  "/notifications": ["admin", "manager"],
  "/geocoding": ["admin", "manager"],
  "/shared": ["admin"],
};

// Action permissions by role
const actionPermissions: Record<string, UserRole[]> = {
  create: ["admin", "manager"],
  edit: ["admin", "manager"],
  delete: ["admin"],
};

export function canAccessModule(role: UserRole | undefined, modulePath: string): boolean {
  if (!role) return false;

  // Admin can access everything
  if (role === "admin") return true;

  // Check module-specific access
  for (const [path, allowedRoles] of Object.entries(moduleAccess)) {
    if (modulePath.startsWith(path)) {
      return allowedRoles.includes(role);
    }
  }

  return false;
}

export function canPerformAction(role: UserRole | undefined, action: "create" | "edit" | "delete"): boolean {
  if (!role) return false;

  // Admin can do everything
  if (role === "admin") return true;

  const allowedRoles = actionPermissions[action];
  return allowedRoles?.includes(role) ?? false;
}

export function getVisibleModules(role: UserRole | undefined): string[] {
  if (!role) return [];

  // Admin sees all modules
  if (role === "admin") {
    return Object.keys(moduleAccess);
  }

  // Other roles see only their permitted modules
  return Object.entries(moduleAccess)
    .filter(([_, allowedRoles]) => allowedRoles.includes(role))
    .map(([path]) => path);
}
