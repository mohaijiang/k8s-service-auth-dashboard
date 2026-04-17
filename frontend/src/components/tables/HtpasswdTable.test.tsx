import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import HtpasswdTable from "./HtpasswdTable";
import type { HtpasswdSecretSummary } from "@/lib/api";

const mockSecrets: HtpasswdSecretSummary[] = [
  {
    name: "my-app-htpasswd",
    namespace: "production",
    userCount: 3,
    createdAt: "2026-04-17T00:00:00Z",
  },
  {
    name: "another-htpasswd",
    namespace: "production",
    userCount: 1,
    createdAt: "2026-04-18T00:00:00Z",
  },
];

describe("HtpasswdTable", () => {
  it("renders table headers", () => {
    render(<HtpasswdTable secrets={[]} loading={false} onSelect={vi.fn()} selectedName={undefined} />);

    expect(screen.getByText("Name")).toBeInTheDocument();
    expect(screen.getByText("User Count")).toBeInTheDocument();
    expect(screen.getByText("Created At")).toBeInTheDocument();
  });

  it("renders secret data rows", () => {
    render(<HtpasswdTable secrets={mockSecrets} loading={false} onSelect={vi.fn()} selectedName={undefined} />);

    expect(screen.getByText("my-app-htpasswd")).toBeInTheDocument();
    expect(screen.getByText("another-htpasswd")).toBeInTheDocument();
    expect(screen.getByText("3")).toBeInTheDocument();
    expect(screen.getByText("1")).toBeInTheDocument();
  });

  it("renders formatted created dates", () => {
    render(<HtpasswdTable secrets={mockSecrets} loading={false} onSelect={vi.fn()} selectedName={undefined} />);

    expect(screen.queryByText("2026-04-17T00:00:00Z")).not.toBeInTheDocument();
  });

  it("calls onSelect when row is clicked", () => {
    const onSelect = vi.fn();
    render(<HtpasswdTable secrets={mockSecrets} loading={false} onSelect={onSelect} selectedName={undefined} />);

    const row = screen.getByText("my-app-htpasswd").closest("tr");
    fireEvent.click(row!);
    expect(onSelect).toHaveBeenCalledWith("my-app-htpasswd");
  });

  it("highlights selected row", () => {
    render(<HtpasswdTable secrets={mockSecrets} loading={false} onSelect={vi.fn()} selectedName="my-app-htpasswd" />);

    const row = screen.getByText("my-app-htpasswd").closest("tr");
    expect(row?.className).toContain("brand");
  });

  it("shows loading state", () => {
    render(<HtpasswdTable secrets={[]} loading={true} onSelect={vi.fn()} selectedName={undefined} />);

    expect(screen.getByText("Loading...")).toBeInTheDocument();
  });

  it("shows empty state", () => {
    render(<HtpasswdTable secrets={[]} loading={false} onSelect={vi.fn()} selectedName={undefined} />);

    expect(screen.getByText("No htpasswd secrets found")).toBeInTheDocument();
  });
});
