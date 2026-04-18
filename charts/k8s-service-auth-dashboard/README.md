# k8s-service-auth-dashboard Helm Chart

A Helm chart for deploying the Kubernetes Service Authentication Dashboard.

## Description

This chart deploys a full-stack application for managing Kubernetes Gateway API resources:
- **Frontend**: Next.js dashboard UI for Services, HTTPRoutes, SecurityPolicies, and Htpasswd management
- **Backend**: Go API server with JWT authentication and Kubernetes client integration

## Features

- Service overview with HTTPRoute and SecurityPolicy status
- HTTPRoute CRUD operations with optional SecurityPolicy creation
- Per-namespace Htpasswd secret management for Envoy Gateway BasicAuth
- User authentication with stateless JWT tokens
- Automatic K8s service account with admin token for backend API access

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- Gateway API CRDs installed (gateway.networking.k8s.io/v1)
- Envoy Gateway installed (gateway.envoyproxy.io/v1alpha1)

## Installation

### Default installation

```bash
helm install my-dashboard ./k8s-service-auth-dashboard
```

### With custom values

```bash
helm install my-dashboard ./k8s-service-auth-dashboard -f custom-values.yaml
```

### With admin password override

```bash
helm install my-dashboard ./k8s-service-auth-dashboard \
  --set backend.adminUser.password=your-secure-password
```

## Configuration

See `values.yaml` for all configurable options.

### Key Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `backend.enabled` | Enable backend deployment | `true` |
| `backend.image.repository` | Backend image repository | `ghcr.io/mohaijiang/k8s-service-auth-dashboard/backend` |
| `backend.adminUser.create` | Create initial admin user | `true` |
| `backend.adminUser.username` | Admin username | `admin` |
| `backend.adminUser.password` | Admin password (random if not set) | `""` |
| `frontend.enabled` | Enable frontend deployment | `true` |
| `frontend.ingress.enabled` | Enable ingress for frontend | `false` |
| `clusterRole.create` | Create admin ClusterRole | `true` |

## RBAC

The chart creates the following RBAC resources:

1. **ServiceAccount**: `backend` service account with mounted K8s token
2. **ClusterRole**: Permissions for Gateway API, SecurityPolicy, Secrets, Services, Namespaces
3. **ClusterRoleBinding**: Binds ClusterRole to backend ServiceAccount
4. **Secret**: Admin K8s token secret mounted to backend pod

The backend uses the mounted service account token to communicate with the Kubernetes API.

## Uninstallation

```bash
helm uninstall my-dashboard
```

## Notes

- The admin password is stored in a Secret. Retrieve it using the notes provided after installation.
- For production use, provide a strong admin password during installation.
- The backend requires access to Gateway API CRDs and Envoy Gateway SecurityPolicy CRD.
