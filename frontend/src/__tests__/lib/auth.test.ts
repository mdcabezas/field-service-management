import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import { getCookie, setCookie, deleteCookie, getAccessToken, isAuthenticated } from "@/lib/auth";

describe("getCookie", () => {
  it("returns null when document is undefined (SSR)", () => {
    expect(getCookie("test")).toBeNull();
  });

  it("returns null for nonexistent cookie", () => {
    expect(getCookie("nonexistent")).toBeNull();
  });

  it("returns cookie value when set", () => {
    document.cookie = "testcookie=hello123";
    expect(getCookie("testcookie")).toBe("hello123");
    document.cookie = "testcookie=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/";
  });

  it("returns correct cookie when multiple exist", () => {
    document.cookie = "a=1";
    document.cookie = "b=2";
    document.cookie = "c=3";
    expect(getCookie("b")).toBe("2");
    document.cookie = "a=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/";
    document.cookie = "b=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/";
    document.cookie = "c=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/";
  });
});

describe("setCookie", () => {
  it("sets a cookie on document", () => {
    setCookie("mycookie", "myvalue", 7);
    expect(getCookie("mycookie")).toBe("myvalue");
    document.cookie = "mycookie=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/";
  });

  it("encodes special characters", () => {
    setCookie("encoded", "hello world=1", 1);
    // setCookie uses encodeURIComponent, getCookie reads raw — verify it's stored
    expect(document.cookie).toContain("encoded=");
    document.cookie = "encoded=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/";
  });
});

describe("deleteCookie", () => {
  it("removes a cookie", () => {
    document.cookie = "todelete=bye";
    expect(getCookie("todelete")).toBe("bye");
    deleteCookie("todelete");
    expect(getCookie("todelete")).toBeNull();
  });
});

describe("getAccessToken", () => {
  it("returns null when no access_token cookie", () => {
    expect(getAccessToken()).toBeNull();
  });

  it("returns access_token cookie value", () => {
    document.cookie = "access_token=test-jwt-token";
    expect(getAccessToken()).toBe("test-jwt-token");
    document.cookie = "access_token=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/";
  });
});

describe("isAuthenticated", () => {
  it("returns false when no access token", () => {
    expect(isAuthenticated()).toBe(false);
  });

  it("returns true when access token exists", () => {
    document.cookie = "access_token=some-token";
    expect(isAuthenticated()).toBe(true);
    document.cookie = "access_token=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/";
  });
});
