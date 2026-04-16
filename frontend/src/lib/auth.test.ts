import { describe, it, expect, beforeEach } from "vitest";
import { getToken, setToken, removeToken } from "./auth";

describe("auth token storage", () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it("returns null when no token is stored", () => {
    expect(getToken()).toBeNull();
  });

  it("stores and retrieves a token", () => {
    setToken("test-jwt-token");
    expect(getToken()).toBe("test-jwt-token");
  });

  it("removes a stored token", () => {
    setToken("test-jwt-token");
    removeToken();
    expect(getToken()).toBeNull();
  });

  it("overwrites an existing token", () => {
    setToken("old-token");
    setToken("new-token");
    expect(getToken()).toBe("new-token");
  });
});
