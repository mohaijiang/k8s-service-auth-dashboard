# Backend — K8s Service Auth Dashboard

Go backend service that manages Kubernetes networking and authentication resources via the Kubernetes API.

## Tech Stack

- **Go 1.25** with [Gin](https://gin-gonic.com/) web framework
- **client-go** for Kubernetes API access
- **JWT** (golang-jwt/v5) for stateless authentication
- **bcrypt** for password hashing

## Prerequisites

- Go 1.24+
- Access to a Kubernetes cluster (local `~/.kube/config` or in-cluster ServiceAccount)

## Quick Start

```bash
# Install dependencies
go mod tidy

# Set required environment variables
export INIT_ADMIN_USERNAME=admin
export INIT_ADMIN_PASSWORD=admin12345678

# Run (uses ~/.kube/config for K8s access)
go run ./cmd/server
```

Server starts on `http://localhost:8080`.

## Running Tests

```bash
go test ./...
```

## Docker

```bash
docker build -t k8s-auth-dashboard-backend .
docker run -p 8080:8080 \
  -e INIT_ADMIN_USERNAME=admin \
  -e INIT_ADMIN_PASSWORD=admin12345678 \
  -v ~/.kube/config:/root/.kube/config \
  k8s-auth-dashboard-backend
```

When running in-cluster, K8s credentials are loaded from the ServiceAccount automatically.

## Project Structure

```
backend/
├── cmd/server/main.go              # Entry point, Gin router, route registration
├── internal/
│   ├── auth/                       # Authentication layer
│   │   ├── bcrypt.go               # Password hashing utilities
│   │   ├── jwt.go                  # JWT token generation and validation
│   │   ├── middleware.go           # JWT authentication middleware
│   │   └── ratelimit.go            # Login rate limiter (1 req/5 sec)
│   ├── bootstrap/
│   │   └── admin.go                # Initial admin user creation
│   ├── config/
│   │   └── config.go               # Environment-based configuration
│   ├── handler/                    # HTTP request handlers
│   │   ├── auth_handler.go         # POST /api/auth/login
│   │   ├── user_handler.go         # User CRUD endpoints
│   │   ├── service_handler.go      # GET /api/services, /api/namespaces
│   │   ├── service_association.go  # Service↔HTTPRoute↔SecurityPolicy matching
│   │   └── htpasswd_handler.go     # Htpasswd secret CRUD + user management
│   ├── k8s/                        # Kubernetes API interactions
│   │   ├── client.go               # K8s client factory (in-cluster or kubeconfig)
│   │   ├── user_secret.go          # User credential secret operations
│   │   ├── jwt_secret.go           # JWT signing key management
│   │   ├── htpasswd.go             # Htpasswd secret CRUD with retry
│   │   ├── service.go              # Service listing
│   │   ├── httproute.go            # HTTPRoute CRD parsing
│   │   └── securitypolicy.go       # SecurityPolicy CRD parsing
│   ├── model/                      # Request/response data models
│   │   ├── user.go                 # User-related types
│   │   ├── service.go              # Service overview types
│   │   └── htpasswd.go             # Htpasswd types
│   └── validator/
│       └── username.go             # Username format validation
├── Dockerfile                      # Multi-stage build (Go 1.25 Alpine)
└── go.mod
```

## API Endpoints

### Public

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Health check |
| `POST` | `/api/auth/login` | Login with username/password (rate-limited) |

### Protected (JWT required)

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/users` | List all users |
| `POST` | `/api/users` | Create user (`username`, `password`) |
| `DELETE` | `/api/users/:username` | Delete user (cannot delete self) |
| `GET` | `/api/services` | List services with associations (`?namespace=`) |
| `GET` | `/api/namespaces` | List all namespaces |
| `GET` | `/api/namespaces/:ns/htpasswd` | List htpasswd secrets |
| `GET` | `/api/namespaces/:ns/htpasswd/:name` | Get secret detail with linked policies |
| `POST` | `/api/namespaces/:ns/htpasswd` | Create htpasswd secret |
| `POST` | `/api/namespaces/:ns/htpasswd/:name/users` | Add user to secret |
| `DELETE` | `/api/namespaces/:ns/htpasswd/:name/users/:username` | Remove user from secret |
| `DELETE` | `/api/namespaces/:ns/htpasswd/:name` | Delete htpasswd secret |

### Response Format

All endpoints return JSON with a consistent envelope:

```json
{
  "success": true,
  "data": { ... },
  "error": "optional error message"
}
```

## API Usage Examples

```bash
# Login
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin12345678"}' | jq -r '.token')

# List users
curl -s http://localhost:8080/api/users \
  -H "Authorization: Bearer $TOKEN" | jq

# List services
curl -s http://localhost:8080/api/services \
  -H "Authorization: Bearer $TOKEN" | jq

# List namespaces
curl -s http://localhost:8080/api/namespaces \
  -H "Authorization: Bearer $TOKEN" | jq

# Create htpasswd secret
curl -s -X POST http://localhost:8080/api/namespaces/default/htpasswd \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"my-app-htpasswd","users":[{"username":"admin","password":"pass12345678"}]}' | jq

# List htpasswd secrets
curl -s http://localhost:8080/api/namespaces/default/htpasswd \
  -H "Authorization: Bearer $TOKEN" | jq

# Add user to htpasswd
curl -s -X POST http://localhost:8080/api/namespaces/default/htpasswd/my-app-htpasswd/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","password":"pass99999999"}' | jq

# Delete htpasswd secret
curl -s -X DELETE http://localhost:8080/api/namespaces/default/htpasswd/my-app-htpasswd \
  -H "Authorization: Bearer $TOKEN" | jq
```

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `ENV` | `development` | `development` or `production` |
| `NAMESPACE` | `dashboard-auth-system` | K8s namespace for dashboard resources |
| `INIT_ADMIN_USERNAME` | — | Initial admin username (required on first run) |
| `INIT_ADMIN_PASSWORD` | — | Initial admin password (min 8 chars, required on first run) |
| `JWT_EXPIRY` | `24h` | JWT token expiration duration |
| `CORS_ALLOW_ORIGIN` | `*` | Allowed CORS origin |

## Kubernetes Resources

The backend manages these K8s resources:

- **User Secrets**: `dashboard-user-<username>` in `dashboard-auth-system` namespace, labeled `app.kubernetes.io/type: dashboard-user`
- **JWT Key Secret**: `dashboard-jwt-key` in `dashboard-auth-system` namespace
- **Htpasswd Secrets**: User-created, labeled `app.kubernetes.io/type: htpasswd` + `app.kubernetes.io/part-of: k8s-service-auth-dashboard`
- **HTTPRoutes**: Read-only access to `gateway.networking.k8s.io/v1` HTTPRoute CRDs
- **SecurityPolicies**: Read-only access to `gateway.envoyproxy.io/v1alpha1` SecurityPolicy CRDs
- **Services**: Read-only access to core K8s Services
