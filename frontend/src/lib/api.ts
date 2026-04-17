import { getToken } from "./auth";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "";

export interface User {
  username: string;
  createdAt: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface ListUsersResponse {
  success: boolean;
  data: User[];
}

export interface CreateUserResponse {
  success: boolean;
  data: User;
}

export interface DeleteUserResponse {
  success: boolean;
  message: string;
}

export interface ServicePort {
  name: string;
  port: number;
  targetPort: string;
  protocol: string;
}

export interface ParentRef {
  name: string;
  namespace: string;
}

export interface HTTPRouteInfo {
  name: string;
  namespace: string;
  hostnames: string[];
  parentRefs: ParentRef[];
}

export interface SecurityPolicyInfo {
  name: string;
  namespace: string;
  hasBasicAuth: boolean;
  hasTLS: boolean;
}

export interface ServiceOverview {
  name: string;
  namespace: string;
  clusterIP: string;
  ports: ServicePort[];
  selector: Record<string, string>;
  httpRoute: HTTPRouteInfo | null;
  securityPolicy: SecurityPolicyInfo | null;
}

export interface ListServicesResponse {
  success: boolean;
  data: ServiceOverview[];
}

export interface ListNamespacesResponse {
  success: boolean;
  data: string[];
}

export async function apiClient<T = unknown>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const token = getToken();

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string>),
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers,
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error || `Request failed with status ${response.status}`);
  }

  return data as T;
}

export function login(credentials: { username: string; password: string }): Promise<LoginResponse> {
  return apiClient<LoginResponse>("/api/auth/login", {
    method: "POST",
    body: JSON.stringify(credentials),
  });
}

export function listUsers(): Promise<ListUsersResponse> {
  return apiClient<ListUsersResponse>("/api/users");
}

export function createUser(userData: { username: string; password: string }): Promise<CreateUserResponse> {
  return apiClient<CreateUserResponse>("/api/users", {
    method: "POST",
    body: JSON.stringify(userData),
  });
}

export function deleteUser(username: string): Promise<DeleteUserResponse> {
  return apiClient<DeleteUserResponse>(`/api/users/${username}`, {
    method: "DELETE",
  });
}

export function listServices(namespace?: string): Promise<ListServicesResponse> {
  const path = namespace ? `/api/services?namespace=${namespace}` : "/api/services";
  return apiClient<ListServicesResponse>(path);
}

export function listNamespaces(): Promise<ListNamespacesResponse> {
  return apiClient<ListNamespacesResponse>("/api/namespaces");
}
