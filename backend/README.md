# K8s Service Auth Dashboard - Backend

Go backend service for the K8s Service Auth Dashboard.

## Prerequisites

- Go 1.24+
- Access to a Kubernetes cluster (local `~/.kube/config` or in-cluster ServiceAccount)

## Quick Start

```bash
go mod tidy
INIT_ADMIN_PASSWORD=testpass123 go run ./cmd/server
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `ENV` | Environment (`development` or `production`) | `development` |
| `NAMESPACE` | K8s namespace for dashboard resources | `dashboard-auth-system` |
| `INIT_ADMIN_USERNAME` | Initial admin username | `admin` |
| `INIT_ADMIN_PASSWORD` | Initial admin password (only used if no users exist) | - |
| `JWT_EXPIRY` | JWT token expiry duration | `24h` |
| `CORS_ALLOW_ORIGIN` | CORS allowed origin | `http://localhost:3000` |

## API Endpoints

### Public

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/auth/login` | Authenticate and receive JWT token |
| GET | `/health` | Server health check |

### Protected (require JWT)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/users` | List all users |
| POST | `/api/users` | Create a new user |
| DELETE | `/api/users/:username` | Delete a user |

## Project Structure

```
backend/
├── cmd/server/main.go          # Entry point
├── internal/
│   ├── config/                 # Environment configuration
│   ├── k8s/                    # Kubernetes client and Secret operations
│   ├── auth/                   # JWT, bcrypt, middleware
│   ├── handler/                # HTTP request handlers
│   ├── model/                  # Domain models
│   ├── validator/              # Input validation
│   └── bootstrap/              # Admin initialization
├── go.mod
└── go.sum
```

## Testing

```bash
go test -race ./...
```

## API Usage Examples

```bash
# Login
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"testpass123"}' | jq -r '.token')

# List users
curl -s http://localhost:8080/api/users \
  -H "Authorization: Bearer $TOKEN" | jq

# Create user
curl -s -X POST http://localhost:8080/api/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}' | jq

# Delete user
curl -s -X DELETE http://localhost:8080/api/users/testuser \
  -H "Authorization: Bearer $TOKEN" | jq
```
