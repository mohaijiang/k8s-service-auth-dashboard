# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**K8s Service Auth Dashboard** — a full-stack application for managing Kubernetes cluster networking and authentication. The project uses a frontend/backend separation architecture.

### Core Purpose

This dashboard provides a management interface for:

1. **User Management** — Admin-only user CRUD with bcrypt passwords stored in K8s Secrets, stateless JWT authentication.
2. **Service Overview** — Cross-namespace K8s Service listing with HTTPRoute and SecurityPolicy status indicators.
3. **Htpasswd Management** — Per-namespace .htpasswd Secret CRUD for Envoy Gateway SecurityPolicy BasicAuth.

### Architecture

```
k8s-service-auth-dashboard/
├── frontend/          # Next.js 16 + React 19 + TypeScript + Tailwind CSS v4
├── backend/           # Go (Gin) + client-go + JWT
├── CLAUDE.md
└── .gitignore
```

### Request Flow

```
Browser → Frontend (Next.js :3001) → Backend API (Go :8080) → Kubernetes API
                                                            ├── Secret (user accounts, JWT key, .htpasswd)
                                                            ├── HTTPRoute CRD (Gateway API)
                                                            ├── SecurityPolicy CRD (Envoy Gateway)
                                                            └── Service (read-only)
```

Frontend proxies `/api/*` requests to backend via Next.js rewrites (no CORS issues).

### Key Design Decisions

1. **No database** — All data stored as K8s resources (Secrets, CRDs).
2. **Backend**: Go + Gin + client-go, not Next.js API Routes.
3. **Authentication**: Stateless JWT signed by key stored in K8s Secret.
4. **Authorization**: All dashboard users share admin permissions.
5. **Concurrency**: Optimistic locking via K8s `resourceVersion` with retry on 409 Conflict.
6. **Namespace**: Dashboard-owned resources live in `dashboard-auth-system` namespace.
7. **Dual-mode K8s credentials**: Local dev (`~/.kube/config`) or in-cluster (ServiceAccount token).

### Kubernetes Resource Model

| Resource | Purpose | Key Fields |
|----------|---------|------------|
| Secret | User credentials, JWT key, htpasswd | `type: Opaque`, labeled by purpose |
| HTTPRoute | Gateway API routing rule | hostnames, rules, backendRefs, parentRefs |
| SecurityPolicy | Envoy Gateway auth policy | targetRefs, basicAuth.users.name |
| Service | Cluster-internal endpoint (read-only) | clusterIP, ports, selector |

### Association Chain

```
Gateway (parentRefs.name)
  └── HTTPRoute (spec.hostnames, spec.rules)
        ├── backendRefs[].name → Service
        └── SecurityPolicy (targetRefs.name matches HTTPRoute name)
              └── basicAuth.users.name → Secret (.htpasswd)
```

## Requirements Status

### REQ-1: User Management (Done)

- Admin user init via env vars (`INIT_ADMIN_USERNAME`, `INIT_ADMIN_PASSWORD`)
- Web UI for user CRUD (create, list, delete)
- Passwords stored as bcrypt hash in K8s Secret (`dashboard-user-<username>`)
- Stateless JWT auth, rate-limited login

### REQ-2: Service Overview (Done)

- Cross-namespace Service list with HTTPRoute and SecurityPolicy status badges
- Association matching: Service → HTTPRoute → SecurityPolicy
- Namespace filtering

### REQ-3: Htpasswd Management (Done)

- Per-namespace .htpasswd Secret CRUD (SHA1 format for Envoy Gateway)
- User add/remove within Secrets
- Association view: which SecurityPolicies reference each Secret
- Optimistic concurrency retry (3 attempts) on 409 Conflict

## Backend (`backend/`)

### Commands

```bash
cd backend
go mod tidy
go run ./cmd/server       # Start on :8080
go test ./...             # Run all tests
go build ./...            # Build check
```

### Structure

```
backend/
├── cmd/server/main.go          # Entry point, Gin router setup
├── internal/
│   ├── auth/                   # JWT, bcrypt, middleware, rate limiting
│   ├── bootstrap/              # Admin user initialization
│   ├── config/                 # Env-based configuration
│   ├── handler/                # HTTP handlers (auth, user, service, htpasswd)
│   ├── k8s/                    # K8s client factory + resource operations
│   ├── model/                  # Request/response data models
│   └── validator/              # Input validation
├── Dockerfile
└── go.mod
```

### API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/health` | No | Health check |
| POST | `/api/auth/login` | No | Login (rate-limited) |
| GET | `/api/users` | JWT | List users |
| POST | `/api/users` | JWT | Create user |
| DELETE | `/api/users/:username` | JWT | Delete user |
| GET | `/api/services` | JWT | List services (optional `?namespace=`) |
| GET | `/api/namespaces` | JWT | List namespaces |
| GET | `/api/namespaces/:ns/htpasswd` | JWT | List htpasswd secrets |
| GET | `/api/namespaces/:ns/htpasswd/:name` | JWT | Get htpasswd detail |
| POST | `/api/namespaces/:ns/htpasswd` | JWT | Create htpasswd secret |
| POST | `/api/namespaces/:ns/htpasswd/:name/users` | JWT | Add user |
| DELETE | `/api/namespaces/:ns/htpasswd/:name/users/:username` | JWT | Remove user |
| DELETE | `/api/namespaces/:ns/htpasswd/:name` | JWT | Delete secret |

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `ENV` | `development` | Environment |
| `NAMESPACE` | `dashboard-auth-system` | Dashboard K8s namespace |
| `INIT_ADMIN_USERNAME` | — | Initial admin username |
| `INIT_ADMIN_PASSWORD` | — | Initial admin password |
| `JWT_EXPIRY` | `24h` | JWT token expiry |
| `CORS_ALLOW_ORIGIN` | `*` | CORS allowed origin |

## Frontend (`frontend/`)

### Commands

```bash
cd frontend
npm run dev       # Start dev server
npm run build     # Production build
npm run lint      # ESLint check
npm run start     # Start production server
npx vitest run    # Run tests
```

### Structure

```
frontend/src/
├── app/
│   ├── (admin)/                 # Admin layout (sidebar + header)
│   │   ├── page.tsx             # Dashboard home
│   │   ├── services/page.tsx    # Service overview (REQ-2)
│   │   └── htpasswd/page.tsx    # Htpasswd management (REQ-3)
│   └── (full-width-pages)/      # Auth pages (signin, signup, 404)
├── components/
│   ├── tables/                  # ServiceTable, HtpasswdTable
│   ├── common/                  # PageBreadcrumb, ComponentCard
│   └── ui/                      # Button, Badge, Modal, Table, etc.
├── context/                     # ThemeContext, SidebarContext
├── hooks/                       # useModal, useGoBack
├── layout/                      # AppSidebar, AppHeader, Backdrop
└── lib/
    ├── api.ts                   # API client + all types/functions
    └── auth.ts                  # Token storage utilities
```

### Path Alias

`@/*` maps to `./src/*` (configured in `tsconfig.json`).

### Key Libraries

- **Framework**: Next.js 16, React 19
- **Styling**: Tailwind CSS v4
- **Testing**: Vitest + jsdom + @testing-library/react
- **Charts**: ApexCharts (template-only, removable)

## Testing

### Backend

```bash
cd backend && go test ./...
```

Covers: htpasswd generation/parsing, HTTPRoute/SecurityPolicy parsing, handler logic.

### Frontend

```bash
cd frontend && npx vitest run
```

Covers: apiClient, all API functions, ServiceTable, HtpasswdTable, auth context.
