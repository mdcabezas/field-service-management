import { describe, it, expect } from "vitest";
import { NextRequest, NextResponse } from "next/server";
import { middleware } from "@/middleware";

function makeRequest(path: string, cookies?: Record<string, string>) {
  const url = `http://localhost:3000${path}`;
  const req = new NextRequest(url);
  if (cookies) {
    for (const [key, value] of Object.entries(cookies)) {
      req.cookies.set(key, value);
    }
  }
  return req;
}

describe("middleware", () => {
  it("redirects to /login when accessing protected route without token", () => {
    const req = makeRequest("/core/users");
    const res = middleware(req);
    expect(res).toBeInstanceOf(NextResponse);
    const location = res.headers.get("location") || "";
    expect(location).toContain("/login");
    expect(location).toContain("redirect=%2Fcore%2Fusers");
  });

  it("redirects to /login for root path without token", () => {
    const req = makeRequest("/");
    const res = middleware(req);
    const location = res.headers.get("location") || "";
    expect(location).toContain("/login");
  });

  it("allows protected route when token exists", () => {
    const req = makeRequest("/core/users", { access_token: "valid-token" });
    const res = middleware(req);
    // Should not redirect — either NextResponse.next() or rewrite
    const location = res.headers.get("location");
    expect(location).toBeNull();
  });

  it("redirects authenticated user away from /login", () => {
    const req = makeRequest("/login", { access_token: "valid-token" });
    const res = middleware(req);
    const location = res.headers.get("location") || "";
    expect(location).toContain("http://localhost:3000/");
  });

  it("allows /login when not authenticated", () => {
    const req = makeRequest("/login");
    const res = middleware(req);
    const location = res.headers.get("location");
    expect(location).toBeNull();
  });

  it("rewrites /api/photos/* to backend", () => {
    const req = makeRequest("/api/photos/123/thumbnail");
    const res = middleware(req);
    // Rewrite returns a response with rewritten URL
    expect(res).toBeDefined();
  });

  it("passes through non-protected routes without token", () => {
    const req = makeRequest("/some-public-page");
    const res = middleware(req);
    const location = res.headers.get("location");
    expect(location).toBeNull();
  });
});
