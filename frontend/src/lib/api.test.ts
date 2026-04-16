import { describe, it, expect, beforeEach, vi } from "vitest";
import { apiClient, login, createUser, listUsers, deleteUser } from "./api";

describe("apiClient", () => {
  beforeEach(() => {
    localStorage.clear();
    vi.restoreAllMocks();
  });

  it("sends request without auth header when no token", async () => {
    const mockResponse = { status: "ok" };
    vi.spyOn(globalThis, "fetch").mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockResponse),
    } as Response);

    const result = await apiClient("/health");
    expect(result).toEqual(mockResponse);
    expect(globalThis.fetch).toHaveBeenCalledWith(
      expect.any(String),
      expect.objectContaining({
        headers: expect.not.objectContaining({
          Authorization: expect.anything(),
        }),
      })
    );
  });

  it("includes Authorization header when token exists", async () => {
    localStorage.setItem("auth_token", "test-jwt");
    vi.spyOn(globalThis, "fetch").mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({}),
    } as Response);

    await apiClient("/api/users");
    expect(globalThis.fetch).toHaveBeenCalledWith(
      expect.any(String),
      expect.objectContaining({
        headers: { Authorization: "Bearer test-jwt", "Content-Type": "application/json" },
      })
    );
  });

  it("throws on non-ok response", async () => {
    vi.spyOn(globalThis, "fetch").mockResolvedValueOnce({
      ok: false,
      status: 401,
      json: () => Promise.resolve({ error: "unauthorized" }),
    } as Response);

    await expect(apiClient("/api/users")).rejects.toThrow("unauthorized");
  });
});

describe("login", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it("calls POST /api/auth/login with credentials", async () => {
    const mockResponse = {
      token: "jwt-token",
      user: { username: "admin", createdAt: "2026-01-01T00:00:00Z" },
    };
    vi.spyOn(globalThis, "fetch").mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockResponse),
    } as Response);

    const result = await login({ username: "admin", password: "pass123" });
    expect(result).toEqual(mockResponse);
    expect(globalThis.fetch).toHaveBeenCalledWith(
      expect.stringContaining("/api/auth/login"),
      expect.objectContaining({ method: "POST" })
    );
  });
});

describe("listUsers", () => {
  beforeEach(() => {
    localStorage.clear();
    vi.restoreAllMocks();
  });

  it("calls GET /api/users", async () => {
    const mockResponse = {
      success: true,
      data: [{ username: "admin", createdAt: "2026-01-01T00:00:00Z" }],
    };
    localStorage.setItem("auth_token", "test-jwt");
    vi.spyOn(globalThis, "fetch").mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockResponse),
    } as Response);

    const result = await listUsers();
    expect(result).toEqual(mockResponse);
  });
});

describe("createUser", () => {
  beforeEach(() => {
    localStorage.clear();
    vi.restoreAllMocks();
  });

  it("calls POST /api/users with user data", async () => {
    localStorage.setItem("auth_token", "test-jwt");
    vi.spyOn(globalThis, "fetch").mockResolvedValueOnce({
      ok: true,
      json: () =>
        Promise.resolve({
          success: true,
          data: { username: "newuser", createdAt: "2026-01-01T00:00:00Z" },
        }),
    } as Response);

    const result = await createUser({ username: "newuser", password: "password123" });
    expect(result.success).toBe(true);
  });
});

describe("deleteUser", () => {
  beforeEach(() => {
    localStorage.clear();
    vi.restoreAllMocks();
  });

  it("calls DELETE /api/users/:username", async () => {
    localStorage.setItem("auth_token", "test-jwt");
    vi.spyOn(globalThis, "fetch").mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ success: true, message: "user deleted" }),
    } as Response);

    const result = await deleteUser("testuser");
    expect(globalThis.fetch).toHaveBeenCalledWith(
      expect.stringContaining("/api/users/testuser"),
      expect.objectContaining({ method: "DELETE" })
    );
    expect(result.success).toBe(true);
  });
});
