"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";
import { useAuthStore } from "@/stores/auth-store";
import { canAccessModule } from "@/lib/rbac";
import { useState } from "react";

interface SubModule {
  name: string;
  path: string;
}

interface Module {
  name: string;
  path: string;
  icon: string;
  children?: SubModule[];
}

const modules: Module[] = [
  {
    name: "Core",
    path: "/core",
    icon: "⚙",
    children: [
      { name: "Usuarios", path: "/core/users" },
      { name: "Roles Técnicos", path: "/core/tech-roles" },
    ],
  },
  {
    name: "Partners",
    path: "/partners",
    icon: "🤝",
  },
  {
    name: "Customers",
    path: "/customers",
    icon: "👥",
  },
  {
    name: "Inventory",
    path: "/inventory",
    icon: "📦",
    children: [
      { name: "Materiales", path: "/inventory/materials" },
      { name: "Herramientas", path: "/inventory/tools" },
      { name: "EPP", path: "/inventory/epp" },
      { name: "Vehículos", path: "/inventory/vehicles" },
      { name: "Arriendos", path: "/inventory/rentals" },
      { name: "Tarifas", path: "/inventory/cost-rates" },
      { name: "Mantención", path: "/inventory/maintenance-records" },
      { name: "Calendario", path: "/inventory/maintenance-schedules" },
      { name: "Checklists", path: "/inventory/checklist-templates" },
    ],
  },
  {
    name: "Planning",
    path: "/planning",
    icon: "📅",
    children: [
      { name: "Planes Diarios", path: "/planning/daily-plans" },
      { name: "Asignaciones", path: "/planning/assignments" },
      { name: "Rutas", path: "/planning/routes" },
    ],
  },
  {
    name: "Operations",
    path: "/operations",
    icon: "🔧",
    children: [
      { name: "Visitas", path: "/operations/visits" },
      { name: "Reportes", path: "/operations/reports" },
      { name: "Asig. Vehículos", path: "/operations/vehicle-assignments" },
      { name: "SLAs", path: "/operations/sla-trackings" },
    ],
  },
  {
    name: "Notifications",
    path: "/notifications",
    icon: "🔔",
    children: [
      { name: "Plantillas", path: "/notifications/templates" },
      { name: "Enviados", path: "/notifications/list" },
    ],
  },
  {
    name: "Geocoding",
    path: "/geocoding",
    icon: "📍",
  },
  {
    name: "Shared",
    path: "/shared",
    icon: "📋",
    children: [
      { name: "Plantillas Informes", path: "/shared/report-templates" },
    ],
  },
];

export function Sidebar() {
  const pathname = usePathname();
  const { user } = useAuthStore();
  const [expanded, setExpanded] = useState<Record<string, boolean>>({});

  const toggleExpand = (path: string) => {
    setExpanded((prev) => ({ ...prev, [path]: !prev[path] }));
  };

  return (
    <aside className="w-64 h-screen border-r-2 border-black bg-black text-white flex flex-col">
      <div className="p-4 border-b-2 border-white">
        <h1 className="font-mono text-lg font-bold uppercase tracking-wider">
          LOCALIS FSM
        </h1>
      </div>

      <nav className="flex-1 overflow-y-auto">
        <ul className="p-2">
          {modules.map((module) => {
            const hasAccess = canAccessModule(user?.role, module.path);
            if (!hasAccess) return null;

            const isActive = pathname.startsWith(module.path);
            const isExpanded = expanded[module.path] ?? isActive;
            const hasChildren = module.children && module.children.length > 0;

            return (
              <li key={module.path}>
                <div className="flex items-center">
                  <Link
                    href={module.path}
                    className={cn(
                      "flex-1 flex items-center gap-3 px-3 py-2 font-mono text-sm uppercase tracking-wider transition-colors",
                      isActive
                        ? "bg-white text-black"
                        : "text-white/70 hover:bg-white/10 hover:text-white"
                    )}
                  >
                    <span>{module.icon}</span>
                    <span>{module.name}</span>
                  </Link>
                  {hasChildren && (
                    <button
                      onClick={() => toggleExpand(module.path)}
                      className="px-2 py-2 text-white/50 hover:text-white"
                    >
                      {isExpanded ? "▼" : "▶"}
                    </button>
                  )}
                </div>

                {hasChildren && isExpanded && (
                  <ul className="ml-6 border-l border-white/20">
                    {module.children!.map((child) => {
                      const childActive = pathname === child.path;
                      return (
                        <li key={child.path}>
                          <Link
                            href={child.path}
                            className={cn(
                              "block px-3 py-1.5 font-mono text-xs uppercase tracking-wider transition-colors",
                              childActive
                                ? "bg-white/20 text-white"
                                : "text-white/50 hover:bg-white/10 hover:text-white"
                            )}
                          >
                            {child.name}
                          </Link>
                        </li>
                      );
                    })}
                  </ul>
                )}
              </li>
            );
          })}
        </ul>
      </nav>

      <div className="p-4 border-t-2 border-white">
        <p className="font-mono text-xs text-white/50">
          v0.1.0
        </p>
      </div>
    </aside>
  );
}
