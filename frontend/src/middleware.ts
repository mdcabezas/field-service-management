// ============================================================================
// Middleware — Route protection for authenticated routes + API proxy
// ============================================================================

import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

// Protected routes that require authentication
const protectedPrefixes = [
  "/core",
  "/partners",
  "/customers",
  "/inventory",
  "/planning",
  "/operations",
  "/notifications",
  "/geocoding",
  "/shared",
];

// Auth routes (redirect to / if already logged in)
const authRoutes = ["/login"];

// Backend URL for internal Docker network
const BACKEND_URL = "http://go-backend:8080";

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  // Proxy photo requests to backend
  if (pathname.startsWith("/api/photos/")) {
    const backendUrl = new URL(pathname, BACKEND_URL);
    return NextResponse.rewrite(backendUrl);
  }

  // Check for access token in cookies
  const accessToken = request.cookies.get("access_token")?.value;

  // If accessing protected route without token, redirect to login
  const isProtected =
    pathname === "/" || protectedPrefixes.some((prefix) => pathname.startsWith(prefix));

  if (isProtected) {
    if (!accessToken) {
      const loginUrl = new URL("/login", request.url);
      loginUrl.searchParams.set("redirect", pathname);
      return NextResponse.redirect(loginUrl);
    }
  }

  // If accessing auth route with token, redirect to /
  if (authRoutes.some((route) => pathname.startsWith(route))) {
    if (accessToken) {
      return NextResponse.redirect(new URL("/", request.url));
    }
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    "/",
    "/core/:path*",
    "/partners/:path*",
    "/customers/:path*",
    "/inventory/:path*",
    "/planning/:path*",
    "/operations/:path*",
    "/notifications/:path*",
    "/geocoding/:path*",
    "/shared/:path*",
    "/login",
    "/api/photos/:path*",
  ],
};
