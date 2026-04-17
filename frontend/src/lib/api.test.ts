import { describe, it, expect, beforeEach, vi } from "vitest";
import {
  apiClient, login, createUser, listUsers, deleteUser, listServices, listNamespaces,
  listHtpasswdSecrets, getHtpasswdSecret, createHtpasswdSecret,
  addHtpasswdUser, removeHtpasswdUser, deleteHtpasswdSecret,
} from "./api";

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

describe("listServices", () => {
  beforeEach(() => {
    localStorage.clear();
    vi.restoreAllMocks();
  });

  it("calls GET /api/services without namespace", async () => {
    localStorage.setItem("auth_token", "test-jwt");
    const mockData = {
      success: true,
      data: [
        {
          name: "my-app",
          namespace: "default",
          clusterIP: "10.0.0.1",
          ports: [{ name: "http", port: 80, targetPort: "8080", protocol: "TCP" }],
          selector: { app: "my-app" },
          httpRoute: null,
          securityPolicy: null,
        },
      ],
    };
    vi.spyOn(globalThis, "fetch").mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockData),
    } as Response);

    const result = await listServices();
    expect(globalThis.fetch).toHaveBeenCalledWith(
      "/api/services",
      expect.any(Object)
    );
    expect(result.data).toHaveLength(1);
    expect(result.data[0].name).toBe("my-app");
    expect(result.data[0].httpRoute).toBeNull();
  });

  it("calls GET /api/services with namespace query param", async () => {
    localStorage.setItem("auth_token", "test-jwt");
    vi.spyOn(globalThis, "fetch").mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ success: true, data: [] }),
    } as Response);

    await listServices("production");
    expect(globalThis.fetch).toHaveBeenCalledWith(
      "/api/services?namespace=production",
      expect.any(Object)
    );
  });

  it("parses service with HTTPRoute and SecurityPolicy", async () => {
    localStorage.setItem("auth_token", "test-jwt");
    const mockData = {
      success: true,
      data: [
        {
          name: "my-app",
          namespace: "prod",
          clusterIP: "10.0.0.1",
          ports: [],
          selector: {},
          httpRoute: {
            name: "my-app-route",
            namespace: "prod",
            hostnames: ["app.example.com"],
            parentRefs: [{ name: "gateway", namespace: "gw" }],
          },
          securityPolicy: {
            name: "my-app-auth",
            namespace: "prod",
            hasBasicAuth: true,
            hasTLS: false,
          },
        },
      ],
    };
    vi.spyOn(globalThis, "fetch").mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockData),
    } as Response);

    const result = await listServices();
    expect(result.data[0].httpRoute?.name).toBe("my-app-route");
    expect(result.data[0].httpRoute?.hostnames).toContain("app.example.com");
    expect(result.data[0].securityPolicy?.hasBasicAuth).toBe(true);
    expect(result.data[0].securityPolicy?.hasTLS).toBe(false);
  });
});

describe("listNamespaces", () => {
  beforeEach(() => {
    localStorage.clear();
    vi.restoreAllMocks();
  });

  it("calls GET /api/namespaces and returns string array", async () => {
    localStorage.setItem("auth_token", "test-jwt");
    vi.spyOn(globalThis, "fetch").mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ success: true, data: ["default", "kube-system", "production"] }),
    } as Response);

    const result = await listNamespaces();
    expect(globalThis.fetch).toHaveBeenCalledWith(
      "/api/namespaces",
      expect.any(Object)
    );
    expect(result.data).toEqual(["default", "kube-system", "production"]);
  });
});

describe("listHtpasswdSecrets", () => {
  beforeEach(() => {
    localStorage.clear();
    vi.restoreAllMocks();
  });

  it("calls GET /api/namespaces/:namespace/htpasswd", async () => {
    localStorage.setItem("auth_token", "test-jwt");
    const mockData = {
      success: true,
      data: [
        { name: "my-app-htpasswd", namespace: "production", userCount: 3, createdAt: "2026-04-17T00:00:00Z" },
      ],
    };
    vi.spyOn(globalThis, "fetch").mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockData),
    } as Response);

    const result = await listHtpasswdSecrets("production");
    expect(globalThis.fetch).toHaveBeenCalledWith(
      "/api/namespaces/production/htpasswd",
      expect.any(Object)
    );
    expect(result.data).toHaveLength(1);
    expect(result.data[0].name).toBe("my-app-htpasswd");
    expect(result.data[0].userCount).toBe(3);
  });
});

describe("getHtpasswdSecret", () => {
  beforeEach(() => {
    localStorage.clear();
    vi.restoreAllMocks();
  });

  it("calls GET /api/namespaces/:namespace/htpasswd/:name with linked policies", async () => {
    localStorage.setItem("auth_token", "test-jwt");
    const mockData = {
      success: true,
      data: {
        name: "my-app-htpasswd",
        namespace: "production",
        users: [{ username: "admin" }, { username: "user1" }],
        userCount: 2,
        createdAt: "2026-04-17T00:00:00Z",
        linkedSecurityPolicies: [
          { name: "my-app-auth", namespace: "production", targetRef: { name: "my-app-route", kind: "HTTPRoute" } },
        ],
      },
    };
    vi.spyOn(globalThis, "fetch").mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockData),
    } as Response);

    const result = await getHtpasswdSecret("production", "my-app-htpasswd");
    expect(globalThis.fetch).toHaveBeenCalledWith(
      "/api/namespaces/production/htpasswd/my-app-htpasswd",
      expect.any(Object)
    );
    expect(result.data.users).toHaveLength(2);
    expect(result.data.linkedSecurityPolicies).toHaveLength(1);
    expect(result.data.linkedSecurityPolicies[0].targetRef.kind).toBe("HTTPRoute");
  });
});

describe("createHtpasswdSecret", () => {
  beforeEach(() => {
    localStorage.clear();
    vi.restoreAllMocks();
  });

  it("calls POST /api/namespaces/:namespace/htpasswd with users", async () => {
    localStorage.setItem("auth_token", "test-jwt");
    vi.spyOn(globalThis, "fetch").mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({
        success: true,
        data: { name: "new-htpasswd", namespace: "production", userCount: 1 },
      }),
    } as Response);

    const result = await createHtpasswdSecret("production", {
      name: "new-htpasswd",
      users: [{ username: "admin", password: "pass12345678" }],
    });
    expect(globalThis.fetch).toHaveBeenCalledWith(
      "/api/namespaces/production/htpasswd",
      expect.objectContaining({ method: "POST" })
    );
    expect(result.data.name).toBe("new-htpasswd");
    expect(result.data.userCount).toBe(1);
  });
});

describe("addHtpasswdUser", () => {
  beforeEach(() => {
    localStorage.clear();
    vi.restoreAllMocks();
  });

  it("calls POST /api/namespaces/:namespace/htpasswd/:name/users", async () => {
    localStorage.setItem("auth_token", "test-jwt");
    vi.spyOn(globalThis, "fetch").mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ success: true, message: "user added" }),
    } as Response);

    const result = await addHtpasswdUser("production", "my-app-htpasswd", {
      username: "user2",
      password: "pass99999999",
    });
    expect(globalThis.fetch).toHaveBeenCalledWith(
      "/api/namespaces/production/htpasswd/my-app-htpasswd/users",
      expect.objectContaining({ method: "POST" })
    );
    expect(result.success).toBe(true);
  });
});

describe("removeHtpasswdUser", () => {
  beforeEach(() => {
    localStorage.clear();
    vi.restoreAllMocks();
  });

  it("calls DELETE /api/namespaces/:namespace/htpasswd/:name/users/:username", async () => {
    localStorage.setItem("auth_token", "test-jwt");
    vi.spyOn(globalThis, "fetch").mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ success: true, message: "user removed" }),
    } as Response);

    const result = await removeHtpasswdUser("production", "my-app-htpasswd", "user2");
    expect(globalThis.fetch).toHaveBeenCalledWith(
      "/api/namespaces/production/htpasswd/my-app-htpasswd/users/user2",
      expect.objectContaining({ method: "DELETE" })
    );
    expect(result.success).toBe(true);
  });
});

describe("deleteHtpasswdSecret", () => {
  beforeEach(() => {
    localStorage.clear();
    vi.restoreAllMocks();
  });

  it("calls DELETE /api/namespaces/:namespace/htpasswd/:name", async () => {
    localStorage.setItem("auth_token", "test-jwt");
    vi.spyOn(globalThis, "fetch").mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ success: true, message: "deleted" }),
    } as Response);

    const result = await deleteHtpasswdSecret("production", "my-app-htpasswd");
    expect(globalThis.fetch).toHaveBeenCalledWith(
      "/api/namespaces/production/htpasswd/my-app-htpasswd",
      expect.objectContaining({ method: "DELETE" })
    );
    expect(result.success).toBe(true);
  });
});
