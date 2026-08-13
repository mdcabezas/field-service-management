"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

export function Breadcrumbs() {
  const pathname = usePathname();
  const segments = pathname.split("/").filter(Boolean);

  const breadcrumbs = segments.map((segment, index) => {
    const path = `/${segments.slice(0, index + 1).join("/")}`;
    const label = segment.charAt(0).toUpperCase() + segment.slice(1);
    const isLast = index === segments.length - 1;

    return { path, label, isLast };
  });

  return (
    <nav className="font-mono text-sm uppercase tracking-wider">
      <ol className="flex items-center gap-2">
        <li>
          <Link href="/dashboard" className="text-gray-500 hover:text-black">
            Dashboard
          </Link>
        </li>
        {breadcrumbs.map((breadcrumb) => (
          <li key={breadcrumb.path} className="flex items-center gap-2">
            <span className="text-gray-400">/</span>
            {breadcrumb.isLast ? (
              <span className="text-black font-bold">{breadcrumb.label}</span>
            ) : (
              <Link href={breadcrumb.path} className="text-gray-500 hover:text-black">
                {breadcrumb.label}
              </Link>
            )}
          </li>
        ))}
      </ol>
    </nav>
  );
}
