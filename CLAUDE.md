# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**K8s Service Auth Dashboard** — a web frontend for managing Kubernetes cluster networking and authentication. Built on Next.js 16, React 19, TypeScript, and Tailwind CSS v4, using the TailAdmin template as the UI shell.

### Core Purpose

This dashboard provides a management interface for:

1. **Kubernetes Services** — List, inspect, and manage Service objects across namespaces.
2. **HTTPRoute Management** — Create and edit HTTPRoute CRDs (Kubernetes Gateway API) that route external traffic to Services, including host/path rules and backend references.
3. **Security Policies** — Manage security policy CRDs attached to Services for domain-level authentication and encryption (TLS/mTLS).
4. **Domain & Auth Configuration** — Configure custom domains, TLS certificates, and encrypted login (e.g., basic-auth, OAuth proxy) for each Service's HTTPRoute.

### Architecture Pattern

Frontend (this repo) communicates with a backend API that interacts with the Kubernetes cluster. The backend is responsible for CRUD operations on Service, HTTPRoute, and SecurityPolicy CRDs via the Kubernetes API server.

- **Frontend**: Next.js App Router, React components, Tailwind CSS
- **Backend API**: TBD — will be a Go operator, Python FastAPI, or Next.js API routes with `@kubernetes/client-node`
- **Kubernetes Resources**: Service, HTTPRoute (Gateway API), SecurityPolicy (custom or standard CRD)

### Kubernetes Resource Model

| Resource | Purpose | Key Fields |
|----------|---------|------------|
| Service | Cluster-internal service endpoint | name, namespace, clusterIP, ports, selector |
| HTTPRoute | Gateway API routing rule | hostnames, rules (matches, backends), parentRefs |
| SecurityPolicy | Auth/encryption policy for a route | TLS config, auth type (basic/OIDC/mTLS), credentials ref |

## Commands

```bash
npm run dev       # Start dev server
npm run build     # Production build
npm run lint      # ESLint check
npm run start     # Start production server
```

No test framework is configured.

## Architecture

### Path Alias

`@/*` maps to `./src/*` (configured in tsconfig.json).

### Route Groups

The app uses Next.js App Router with two route groups:

- `src/app/(admin)/` — Pages wrapped in the admin layout (sidebar + header). The dashboard home page is at `src/app/(admin)/page.tsx`.
- `src/app/(full-width-pages)/` — Full-width pages without sidebar: auth (signin, signup) and error pages (404).

### Layouts

- **Root layout** (`src/app/layout.tsx`): Wraps everything in `ThemeProvider` then `SidebarProvider`. Uses the "Outfit" font.
- **Admin layout** (`src/app/(admin)/layout.tsx`): Client component that renders `AppSidebar`, `Backdrop`, and `AppHeader` around page content. Sidebar state controls main content margin.
- **Full-width layout** (`src/app/(full-width-pages)/layout.tsx`): Minimal passthrough wrapper.

### State Management

Two React Context providers (both client components):

- `src/context/ThemeContext.tsx` — Light/dark theme toggle, persisted to localStorage, toggles `.dark` class on `<html>`.
- `src/context/SidebarContext.tsx` — Sidebar expand/collapse, mobile drawer, hover state, active nav item, and submenu toggling.

### Component Organization

`src/components/` groups by feature: `auth`, `calendar`, `charts`, `ecommerce`, `header`, `tables`, `ui`, `form`, `user-profile`, `videos`, `common`, `example`.

Template components should be replaced or adapted as the K8s dashboard features are built. Key feature areas to develop:

- **Service list & detail** — Table and detail views for Kubernetes Services
- **HTTPRoute editor** — Form/table for creating and editing HTTPRoute rules
- **SecurityPolicy editor** — Form for TLS, auth type, and credential configuration
- **Domain management** — Domain assignment and certificate status per Service

`src/layout/` contains the three layout-level components: `AppHeader.tsx`, `AppSidebar.tsx`, `Backdrop.tsx`. The sidebar nav items in `AppSidebar.tsx` should be updated to reflect K8s dashboard routes.

`src/hooks/` has `useGoBack` and `useModal`. Add cluster data fetching hooks here.

### Styling

Tailwind CSS v4 with custom theme tokens defined in `src/app/globals.css` using `@theme` directive. Color palette uses semantic names: `brand-*`, `gray-*`, `success-*`, `error-*`, `warning-*`, `blue-light-*`. Dark mode uses the `dark` variant (`@custom-variant dark (&:is(.dark *))`).

SVGs are loaded via `@svgr/webpack` and exported from `src/icons/index.tsx` as React components.

### Key Libraries

- **Charts**: ApexCharts via `react-apexcharts` — useful for service health/status dashboards
- **Tables**: Existing table components should be adapted for listing K8s resources
- **Forms**: Existing form components for HTTPRoute and SecurityPolicy editing
- **Date picker**: flatpickr (direct, not react-flatpickr)

Template-only libraries (jVectormap, FullCalendar, Swiper) can be removed when no longer needed.

### ESLint

Configuration in `eslint.config.mjs` uses `eslint-config-next` with core-web-vitals and TypeScript presets. Ignores `.next/`, `out/`, `build/`, `next-env.d.ts`.
