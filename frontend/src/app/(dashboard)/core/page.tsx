"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

export default function CorePage() {
  const router = useRouter();

  useEffect(() => {
    router.replace("/core/users");
  }, [router]);

  return <div className="font-mono">Redirigiendo...</div>;
}
