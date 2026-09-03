"use client";

import { useAuthStore } from "@/stores/auth-store";
import { useAuth } from "@/hooks/use-auth";
import { Button } from "@/components/ui/button";

export function Topbar() {
  const { user } = useAuthStore();
  const { logout } = useAuth();

  return (
    <header className="h-16 border-b-2 border-black bg-white flex items-center justify-between px-6">
      <div className="flex items-center gap-4">
        <h2 className="font-mono text-sm font-bold uppercase tracking-wider">
          Panel de Control
        </h2>
      </div>

      <div className="flex items-center gap-4">
        <div className="text-right">
          <p className="font-mono text-sm font-bold">
            {user?.email}
          </p>
          <p className="font-mono text-xs text-gray-500 uppercase">
            {user?.role}
          </p>
        </div>

        <Button
          variant="outline"
          size="sm"
          onClick={logout}
        >
          Cerrar Sesión
        </Button>
      </div>
    </header>
  );
}
