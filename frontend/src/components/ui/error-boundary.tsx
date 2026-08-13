"use client";

import { useEffect } from "react";
import { Button } from "@/components/ui/button";

interface ErrorBoundaryProps {
  error: Error & { digest?: string };
  reset: () => void;
}

export function ErrorBoundary({ error, reset }: ErrorBoundaryProps) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <div className="min-h-[400px] flex flex-col items-center justify-center border-2 border-black p-8">
      <h2 className="font-mono text-xl font-bold uppercase tracking-wider mb-4">
        Error
      </h2>
      <p className="font-mono text-sm text-gray-600 mb-6">
        {error.message || "Algo salió mal"}
      </p>
      <Button onClick={reset} variant="outline">
        Intentar de nuevo
      </Button>
    </div>
  );
}
