"use client";

import type React from "react";
import { createContext, useState, useContext, useEffect } from "react";
import { login as apiLogin } from "@/lib/api";
import { getToken, setToken, removeToken } from "@/lib/auth";

interface User {
  username: string;
  createdAt: string;
}

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  login: (credentials: { username: string; password: string }) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    const token = getToken();
    const stored = localStorage.getItem("auth_user");
    if (token && stored) {
      try {
        setUser(JSON.parse(stored));
      } catch {
        removeToken();
        localStorage.removeItem("auth_user");
      }
    }
  }, []);

  const login = async (credentials: { username: string; password: string }) => {
    const response = await apiLogin(credentials);
    setToken(response.token);
    const userData = { username: response.user.username, createdAt: response.user.createdAt };
    localStorage.setItem("auth_user", JSON.stringify(userData));
    setUser(userData);
  };

  const logout = () => {
    removeToken();
    localStorage.removeItem("auth_user");
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, isAuthenticated: !!user, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
