# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**K8s Service Auth Dashboard** — a full-stack application for managing Kubernetes cluster networking and authentication. The project uses a frontend/backend separation architecture.

### Core Purpose

This dashboard provides a management interface for:

1. **Kubernetes Services** — List, inspect, and manage Service objects across namespaces.
2. **HTTPRoute Management** — Create and edit HTTPRoute CRDs (Kubernetes Gateway API) that route external traffic to Services, including host/path rules and backend references.
3. **Security Policies** — Manage security policy CRDs attached to Services for domain-level authentication and encryption (TLS/mTLS).
4. **Domain & Auth Configuration** — Configure custom domains, TLS certificates, and encrypted login (e.g., basic-auth, OAuth proxy) for each Service's HTTPRoute.
5. **Basic-Auth (.htpasswd) Management** — Per-namespace .htpasswd credential management for Envoy Gateway SecurityPolicies.

### Architecture Pattern

Monorepo with separated frontend and backend:

```
k8s-service-auth-dashboard/
├── frontend/          # Next.js 16 + React 19 + TypeScript + Tailwind CSS v4
├── backend/           # Go (Golang) backend service
├── CLAUDE.md
└── .gitignore
```

- **Frontend** (`frontend/`): Next.js App Router, React components, Tailwind CSS. Communicates with backend API via HTTP.
- **Backend** (`backend/`): Go service using `client-go` to interact with the Kubernetes API server. Responsible for all CRUD operations on K8s resources (Secrets, ConfigMaps, HTTPRoutes, SecurityPolicies, Services).
- **Kubernetes Resources**: Service, HTTPRoute (Gateway API), SecurityPolicy (Envoy Gateway CRD)

### Architecture Flow

```
Browser → Frontend (Next.js) → Backend API (Go) → Kubernetes API
                                                     ├── Secret (user accounts, JWT key, .htpasswd)
                                                     ├── ConfigMap (dashboard config)
                                                     ├── HTTPRoute CRD (Gateway API)
                                                     ├── SecurityPolicy CRD (Envoy Gateway)
                                                     └── Service (read-only)
```

### Key Design Decisions

1. **No database** — All data stored as K8s resources (Secrets, ConfigMaps, CRDs), following Helm's pattern.
2. **Backend**: Go service with `client-go`, not Next.js API Routes. Chosen for better K8s ecosystem integration and type-safe K8s API access.
3. **Authentication**: Stateless JWT signed by key stored in K8s Secret. No session storage.
4. **Authorization**: All dashboard users share admin permissions. No per-user RBAC.
5. **Concurrency**: Optimistic locking via K8s `resourceVersion`. On 409 Conflict, return current state to frontend.
6. **Namespace**: Dashboard-owned resources live in `dashboard-auth-system` namespace.
7. **Dual-mode K8s credential loading**: Backend supports both local dev (`~/.kube/config`) and in-cluster (ServiceAccount token).

### Kubernetes Resource Model

| Resource | Purpose | Key Fields |
|----------|---------|------------|
| Service | Cluster-internal service endpoint | name, namespace, clusterIP, ports, selector |
| HTTPRoute | Gateway API routing rule (`gateway.networking.k8s.io/v1`) | hostnames, rules (matches, backends), parentRefs |
| SecurityPolicy | Envoy Gateway auth policy (`gateway.envoyproxy.io/v1alpha1`) | targetRefs, basicAuth, TLS config |

### Association Chain

```
Gateway (parentRefs.name)
  └── HTTPRoute (spec.hostnames, spec.rules)
        ├── backendRefs[].name → Service
        └── SecurityPolicy (targetRefs.name matches HTTPRoute name)
              └── basicAuth.users.name → Secret (.htpasswd)
```

## Requirements

### REQ-1: User Management with Admin-Only Roles

- Environment variable initialization (`INIT_ADMIN_USERNAME`, `INIT_ADMIN_PASSWORD`)
- Web UI for user CRUD (create, list, delete)
- Passwords stored as bcrypt hash in K8s Secret (`dashboard-user-<username>`)
- Stateless JWT auth (signed by K8s Secret key)
- All users share admin permissions

### REQ-2: Service Overview with HTTPRoute and SecurityPolicy Status

- Cross-namespace Service list with HTTPRoute and SecurityPolicy status indicators
- Association matching: Service → HTTPRoute → SecurityPolicy
- Namespace filtering

### REQ-3: Basic-Auth (.htpasswd) Management per Namespace

- Per-namespace .htpasswd Secret management (SHA1 format for Envoy Gateway)
- CRUD for Secrets and individual user entries
- Association view showing which SecurityPolicies reference each Secret

## Frontend (`frontend/`)

### Commands

```bash
cd frontend
npm run dev       # Start dev server
npm run build     # Production build
npm run lint      # ESLint check
npm run start     # Start production server
```

### Path Alias

`@/*` maps to `./src/*` (configured in `frontend/tsconfig.json`).

### Route Groups

The app uses Next.js App Router with two route groups:

- `src/app/(admin)/` — Pages wrapped in the admin layout (sidebar + header). The dashboard home page is at `src/app/(admin)/page.tsx`.
- `src/app/(full-width-pages)/` — Full-width pages without sidebar: auth (signin, signup) and error pages (404).

### Layouts

- **Root layout** (`src/app/layout.tsx`): Wraps everything in `ThemeProvider` then `SidebarProvider`. Uses the "Outfit" font.
- **Admin layout** (`src/app/(admin)/layout.tsx`): Client component that renders `AppSidebar`, `Backdrop`, and `AppHeader` around page content.
- **Full-width layout** (`src/app/(full-width-pages)/layout.tsx`): Minimal passthrough wrapper.

### State Management

Two React Context providers (both client components):

- `src/context/ThemeContext.tsx` — Light/dark theme toggle, persisted to localStorage.
- `src/context/SidebarContext.tsx` — Sidebar expand/collapse, mobile drawer, hover state, active nav item.

### Component Organization

`src/components/` groups by feature: `auth`, `calendar`, `charts`, `ecommerce`, `header`, `tables`, `ui`, `form`, `user-profile`, `videos`, `common`, `example`.

`src/layout/` contains: `AppHeader.tsx`, `AppSidebar.tsx`, `Backdrop.tsx`.

`src/hooks/` has `useGoBack` and `useModal`.

### Styling

Tailwind CSS v4 with custom theme tokens in `src/app/globals.css` using `@theme` directive. Dark mode via `.dark` class.

### Key Libraries

- **Charts**: ApexCharts via `react-apexcharts`
- **Date picker**: flatpickr

Template-only libraries (jVectormap, FullCalendar, Swiper) can be removed when no longer needed.

## Backend (`backend/`)

Go backend service (to be implemented). Will use:
- `client-go` for Kubernetes API access
- Standard Go net/http or Gin/Echo framework (TBD)
- JWT middleware for authentication

### Commands (planned)

```bash
cd backend
go mod tidy
go run ./cmd/server
go test ./...
```
