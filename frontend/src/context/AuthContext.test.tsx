import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { ReactNode } from "react";
import { AuthProvider, useAuth } from "./AuthContext";

// Mock the api module
vi.mock("@/lib/api", () => ({
  login: vi.fn(),
}));

import { login } from "@/lib/api";
const mockLogin = vi.mocked(login);

function wrapper({ children }: { children: ReactNode }) {
  return <AuthProvider>{children}</AuthProvider>;
}

describe("AuthContext", () => {
  beforeEach(() => {
    localStorage.clear();
    vi.restoreAllMocks();
  });

  it("starts unauthenticated when no token in storage", () => {
    const { result } = renderHook(() => useAuth(), { wrapper });
    expect(result.current.isAuthenticated).toBe(false);
    expect(result.current.user).toBeNull();
  });

  it("starts authenticated when token exists in storage", () => {
    localStorage.setItem("auth_token", "existing-jwt");
    localStorage.setItem("auth_user", JSON.stringify({ username: "admin", createdAt: "2026-01-01" }));

    const { result } = renderHook(() => useAuth(), { wrapper });
    expect(result.current.isAuthenticated).toBe(true);
    expect(result.current.user?.username).toBe("admin");
  });

  it("login stores token and user", async () => {
    mockLogin.mockResolvedValueOnce({
      token: "new-jwt",
      user: { username: "admin", createdAt: "2026-01-01T00:00:00Z" },
    });

    const { result } = renderHook(() => useAuth(), { wrapper });

    await act(async () => {
      await result.current.login({ username: "admin", password: "pass123" });
    });

    expect(result.current.isAuthenticated).toBe(true);
    expect(result.current.user?.username).toBe("admin");
    expect(localStorage.getItem("auth_token")).toBe("new-jwt");
  });

  it("logout clears token and user", async () => {
    localStorage.setItem("auth_token", "existing-jwt");
    localStorage.setItem("auth_user", JSON.stringify({ username: "admin", createdAt: "2026-01-01" }));

    const { result } = renderHook(() => useAuth(), { wrapper });

    await act(async () => {
      result.current.logout();
    });

    expect(result.current.isAuthenticated).toBe(false);
    expect(result.current.user).toBeNull();
    expect(localStorage.getItem("auth_token")).toBeNull();
  });

  it("login error does not change auth state", async () => {
    mockLogin.mockRejectedValueOnce(new Error("invalid credentials"));

    const { result } = renderHook(() => useAuth(), { wrapper });

    await act(async () => {
      try {
        await result.current.login({ username: "admin", password: "wrong" });
      } catch {
        // expected
      }
    });

    expect(result.current.isAuthenticated).toBe(false);
    expect(result.current.user).toBeNull();
  });

  it("useAuth throws when used outside AuthProvider", () => {
    expect(() => {
      renderHook(() => useAuth());
    }).toThrow("useAuth must be used within an AuthProvider");
  });
});
