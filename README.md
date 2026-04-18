# K8s Service Auth Dashboard

A full-stack dashboard for managing Kubernetes Gateway API resources (HTTPRoutes, SecurityPolicies) and htpasswd authentication. Built with Next.js 16 (React 19), Go (Gin), and client-go.

## Features

- **User Management** — Admin-only CRUD with bcrypt passwords in K8s Secrets, JWT authentication
- **Service Overview** — Cross-namespace Service listing with HTTPRoute and SecurityPolicy status
- **Htpasswd Management** — Per-namespace htpasswd Secret CRUD for Envoy Gateway BasicAuth
- **No Database** — All data stored as Kubernetes resources

## Quick Start

### Local Development (Docker Compose)

```bash
# Requires: ~/.kube/config with K8s cluster access
docker compose up -d --build

# Access at http://localhost:3000
# Default login: admin / admin123
```

### Kubernetes Deployment (Helm)

```bash
# Import images into containerd (if using containerd runtime)
docker save ghcr.io/mohaijiang/k8s-service-auth-dashboard/backend:latest | ctr -n k8s.io images import -
docker save ghcr.io/mohaijiang/k8s-service-auth-dashboard/frontend:latest | ctr -n k8s.io images import -

# Install Helm chart
helm install auth-dashboard ./charts/k8s-service-auth-dashboard --namespace dashboard-auth-system --create-namespace

# Get credentials
export ADMIN_USER=$(kubectl get secret -n dashboard-auth-system auth-dashboard-k8s-service-auth-dashboard-admin-user -o jsonpath='{.data.username}' | base64 -d)
export ADMIN_PASS=$(kubectl get secret -n dashboard-auth-system auth-dashboard-k8s-service-auth-dashboard-admin-user -o jsonpath='{.data.password}' | base64 -d)

# Port-forward (or use NodePort/Ingress)
kubectl port-forward -n dashboard-auth-system svc/auth-dashboard-k8s-service-auth-dashboard-frontend 3000:3000
```

## Architecture

```
┌─────────────┐     /api/*      ┌─────────────┐     K8s API     ┌──────────────┐
│   Browser   │─────────────────▶│   Frontend  │─────────────────▶│   Kubernetes  │
│             │                 │  (Next.js)  │                 │              │
└─────────────┘                 └─────────────┘                 │              │
                                                                  │              │
                                                                  │              │
                                                                  │              │
                                                        ┌─────────┴──────────────┐
                                                        │                        │
                                                        ▼                        ▼
                                                  ┌─────────────┐        ┌──────────────┐
                                                  │   Backend   │        │   Secrets    │
                                                  │     (Go)     │        │              │
                                                  └─────────────┘        └──────────────┘
```

### API Proxying

The frontend uses **runtime environment variable** `BACKEND_URL` for API proxying:

| Environment | `BACKEND_URL` | Mechanism |
|-------------|---------------|-----------|
| Local `npm run dev` | (unset) | Next.js rewrites → `localhost:8080` |
| Docker Compose | `http://backend:8080` | Middleware proxy (Docker network) |
| Kubernetes | `http://<release>-backend:8080` | Middleware proxy (K8s DNS) |

## Project Structure

```
k8s-service-auth-dashboard/
├── backend/              # Go + Gin + client-go
│   ├── cmd/server/       # Entry point
│   ├── internal/         # handlers, k8s client, auth
│   └── Dockerfile
├── frontend/             # Next.js 16 + React 19
│   ├── src/
│   │   ├── app/          # App router pages
│   │   ├── components/   # UI components
│   │   ├── lib/          # API client, auth utilities
│   │   └── middleware.ts # Runtime API proxy
│   └── Dockerfile
├── charts/               # Helm chart
│   └── k8s-service-auth-dashboard/
├── docker-compose.yml
└── CLAUDE.md             # AI assistant guide
```

## Configuration

### Backend Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `ENV` | `development` | Environment |
| `NAMESPACE` | `dashboard-auth-system` | Dashboard K8s namespace |
| `INIT_ADMIN_USERNAME` | `admin` | Initial admin username |
| `INIT_ADMIN_PASSWORD` | — | Initial admin password (auto-gen if empty) |
| `JWT_EXPIRY` | `24h` | JWT token expiry |

### Frontend Environment Variables

| Variable | Description |
|----------|-------------|
| `BACKEND_URL` | Backend API URL (runtime). Example: `http://backend:8080` |
| `NODE_ENV` | Set to `production` for production builds |

### Helm Chart Values

Key values in `charts/k8s-service-auth-dashboard/values.yaml`:

```yaml
backend:
  enabled: true
  replicaCount: 1
  service:
    type: ClusterIP
    port: 8080

frontend:
  enabled: true
  replicaCount: 1
  service:
    type: NodePort       # Default for easy access
    port: 3000
  ingress:
    enabled: false       # Enable for production with Ingress controller
```

## Development

### Backend

```bash
cd backend
go mod tidy
go run ./cmd/server       # Start on :8080
go test ./...             # Run tests
```

### Frontend

```bash
cd frontend
npm run dev               # Start dev server on :3000
npm run build             # Production build
npm run lint              # ESLint check
npx vitest run            # Run tests
```

### Testing

- **Backend**: `go test ./...` — covers htpasswd, HTTPRoute, SecurityPolicy parsing
- **Frontend**: `npx vitest run` — covers API client, components, auth context

## Kubernetes Resources

| Resource | Purpose |
|----------|---------|
| `Secret` | User credentials (bcrypt), JWT key, htpasswd files |
| `HTTPRoute` | Gateway API routing rules |
| `SecurityPolicy` | Envoy Gateway auth policies |
| `Service` | Cluster services (read-only view) |

## Association Chain

```
Gateway
  └── HTTPRoute (spec.hostnames, spec.rules)
        ├── backendRefs[] → Service
        └── SecurityPolicy (targetRefs.name)
              └── basicAuth.users.name → Secret (.htpasswd)
```

## License

MIT
