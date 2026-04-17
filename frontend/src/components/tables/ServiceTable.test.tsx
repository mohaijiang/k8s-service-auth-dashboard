import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import ServiceTable from "./ServiceTable";
import type { ServiceOverview } from "@/lib/api";

const mockServices: ServiceOverview[] = [
  {
    name: "my-app",
    namespace: "production",
    clusterIP: "10.96.0.42",
    ports: [{ name: "http", port: 80, targetPort: "8080", protocol: "TCP" }],
    selector: { app: "my-app" },
    httpRoute: {
      name: "my-app-route",
      namespace: "production",
      hostnames: ["app.example.com"],
      parentRefs: [{ name: "gateway", namespace: "gw" }],
    },
    securityPolicy: {
      name: "my-app-auth",
      namespace: "production",
      hasBasicAuth: true,
      hasTLS: false,
    },
  },
  {
    name: "simple-svc",
    namespace: "default",
    clusterIP: "10.96.0.10",
    ports: [{ name: "web", port: 8080, targetPort: "3000", protocol: "TCP" }],
    selector: { app: "simple" },
    httpRoute: null,
    securityPolicy: null,
  },
];

describe("ServiceTable", () => {
  it("renders table headers", () => {
    render(<ServiceTable services={[]} loading={false} />);

    expect(screen.getByText("Name")).toBeInTheDocument();
    expect(screen.getByText("Namespace")).toBeInTheDocument();
    expect(screen.getByText("Cluster IP")).toBeInTheDocument();
    expect(screen.getByText("Ports")).toBeInTheDocument();
    expect(screen.getByText("HTTPRoute")).toBeInTheDocument();
    expect(screen.getByText("SecurityPolicy")).toBeInTheDocument();
  });

  it("renders service rows with data", () => {
    render(<ServiceTable services={mockServices} loading={false} />);

    expect(screen.getByText("my-app")).toBeInTheDocument();
    expect(screen.getByText("simple-svc")).toBeInTheDocument();
    expect(screen.getByText("10.96.0.42")).toBeInTheDocument();
    expect(screen.getByText("10.96.0.10")).toBeInTheDocument();
  });

  it("shows port info as port/targetPort", () => {
    render(<ServiceTable services={mockServices} loading={false} />);

    expect(screen.getByText("80/8080")).toBeInTheDocument();
    expect(screen.getByText("8080/3000")).toBeInTheDocument();
  });

  it("shows hostname badge when HTTPRoute exists", () => {
    render(<ServiceTable services={mockServices} loading={false} />);

    expect(screen.getByText("app.example.com")).toBeInTheDocument();
  });

  it("shows 'None' badge when no HTTPRoute", () => {
    render(<ServiceTable services={mockServices} loading={false} />);

    // simple-svc has no HTTPRoute
    const noneBadges = screen.getAllByText("None");
    expect(noneBadges.length).toBeGreaterThanOrEqual(2); // HTTPRoute + SecurityPolicy both None
  });

  it("shows BasicAuth status when SecurityPolicy has basicAuth", () => {
    render(<ServiceTable services={mockServices} loading={false} />);

    expect(screen.getByText("BasicAuth")).toBeInTheDocument();
  });

  it("shows loading state", () => {
    render(<ServiceTable services={[]} loading={true} />);

    expect(screen.getByText("Loading...")).toBeInTheDocument();
  });

  it("shows empty state when no services", () => {
    render(<ServiceTable services={[]} loading={false} />);

    expect(screen.getByText("No services found")).toBeInTheDocument();
  });
});
