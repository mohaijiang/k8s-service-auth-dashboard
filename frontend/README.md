# Frontend — K8s Service Auth Dashboard

Next.js admin dashboard for managing Kubernetes service networking and authentication.

## Tech Stack

- **Next.js 16** (App Router) + **React 19** + **TypeScript**
- **Tailwind CSS v4** for styling
- **Vitest** + **@testing-library/react** for testing

## Prerequisites

- Node.js 18.x or later (recommended: 20.x+)

## Quick Start

```bash
# Install dependencies
npm install

# Start dev server (port 3001)
npm run dev

# Production build
npm run build
npm run start
```

The frontend proxies `/api/*` requests to the backend at `http://localhost:8080` via Next.js rewrites.

## Running Tests

```bash
npx vitest run           # Run all tests
npx vitest run --watch   # Watch mode
```

## Docker

```bash
docker build -t k8s-auth-dashboard-frontend .
docker run -p 3000:3000 k8s-auth-dashboard-frontend
```

## Project Structure

```
src/
├── app/
│   ├── layout.tsx                     # Root layout (ThemeProvider + SidebarProvider)
│   ├── globals.css                    # Tailwind v4 theme tokens
│   ├── (admin)/                       # Admin pages (sidebar + header layout)
│   │   ├── layout.tsx                 # Admin layout wrapper
│   │   ├── page.tsx                   # Dashboard home
│   │   ├── services/page.tsx          # Service overview (REQ-2)
│   │   └── htpasswd/page.tsx          # Htpasswd management (REQ-3)
│   └── (full-width-pages)/            # Auth pages (no sidebar)
│       ├── signin/page.tsx
│       ├── signup/page.tsx
│       └── error-404/page.tsx
├── components/
│   ├── auth/                          # SignUpForm, auth-related components
│   ├── common/                        # PageBreadcrumb, ComponentCard, ThemeToggler
│   ├── tables/
│   │   ├── ServiceTable.tsx           # K8s services table (REQ-2)
│   │   ├── ServiceTable.test.tsx
│   │   ├── HtpasswdTable.tsx          # Htpasswd secrets table (REQ-3)
│   │   └── HtpasswdTable.test.tsx
│   └── ui/                            # Button, Badge, Modal, Table, Alert, etc.
├── context/
│   ├── ThemeContext.tsx                # Light/dark theme toggle
│   └── SidebarContext.tsx             # Sidebar state management
├── hooks/
│   ├── useModal.ts                    # Modal open/close state
│   └── useGoBack.ts                   # Navigation helper
├── layout/
│   ├── AppSidebar.tsx                 # Main sidebar navigation
│   ├── AppHeader.tsx                  # Top header bar
│   └── Backdrop.tsx                   # Mobile sidebar backdrop
└── lib/
    ├── api.ts                         # API client + all types and functions
    ├── api.test.ts                    # API function tests
    └── auth.ts                        # JWT token storage (localStorage)
```

## Pages

| Path | Description |
|------|-------------|
| `/` | Dashboard home |
| `/services` | K8s service overview with namespace filter and HTTPRoute/SecurityPolicy status badges |
| `/htpasswd` | Htpasswd secret management per namespace with user CRUD and linked SecurityPolicy view |
| `/signin` | Login page |
| `/signup` | Registration page |

## API Layer

All backend communication goes through `src/lib/api.ts`, which provides:

### Auth
- `login({ username, password })` — Authenticate and receive JWT

### User Management
- `listUsers()` — Get all dashboard users
- `createUser({ username, password })` — Create new user
- `deleteUser(username)` — Delete user

### Services
- `listServices(namespace?)` — List K8s services with HTTPRoute/SecurityPolicy status
- `listNamespaces()` — List all K8s namespaces

### Htpasswd Management
- `listHtpasswdSecrets(namespace)` — List htpasswd secrets in namespace
- `getHtpasswdSecret(namespace, name)` — Get secret detail with linked policies
- `createHtpasswdSecret(namespace, { name, users })` — Create secret
- `addHtpasswdUser(namespace, name, { username, password })` — Add user
- `removeHtpasswdUser(namespace, name, username)` — Remove user
- `deleteHtpasswdSecret(namespace, name)` — Delete secret

### API Client

All functions use the shared `apiClient<T>(path, options)` helper which:
- Adds `Authorization: Bearer <token>` header from localStorage
- Sets `Content-Type: application/json`
- Throws on non-ok responses with the error message from the backend

## Path Alias

`@/*` maps to `./src/*` — configured in `tsconfig.json`.

## Styling

- Tailwind CSS v4 with custom theme tokens defined in `globals.css` via `@theme` directive
- Dark mode via `.dark` class (toggled by ThemeContext)
- Custom UI components in `src/components/ui/` (Button, Badge, Modal, Table)

## Key Patterns

- **Page components**: Client components (`"use client"`) using `useState` + `useEffect` for data fetching
- **Namespace filtering**: Dropdown populated by `listNamespaces()`, triggers re-fetch on change
- **Modals**: Using `<Modal>` component + `useModal` hook for create/delete dialogs
- **Tables**: Custom `HtpasswdTable`/`ServiceTable` components using shared `Table` UI primitives
- **Error display**: Red alert boxes for both fetch errors and action errors
- **Null safety**: API response arrays use `?? []` fallback for nullable fields
