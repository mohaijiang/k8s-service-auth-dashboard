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
