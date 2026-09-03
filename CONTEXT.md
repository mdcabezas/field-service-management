# GasConecta SpA — FSM Application Context

## 1. Original Request and Goals

Build a full-stack Field Service Management (FSM) application for **GasConecta SpA**, a Chilean gas services company. The system manages gas installations, maintenance, PH tests, and related field operations.

**Core Goals:**
- Single-tenant architecture (no multi-tenancy/RLS)
- Generic FSM core + Industry Pack pattern (gas-specific data in separate schema)
- Full CRUD for all entities (9 bounded contexts, ~49 core tables)
- JWT authentication with PostgreSQL (migrated from LDAP)
- React/Next.js frontend with Go backend API
- Docker-based development environment

---

## 2. What's Been Accomplished

### Database (PostgreSQL 16 + PostGIS 3.5)
- **9 core schemas**: `core`, `partners`, `customers`, `inventory`, `operations`, `planning`, `notifications`, `geocoding`, `shared`
- **27 native PostgreSQL enums** (status, priority, type, etc.)
- **~49 core tables** with proper FKs, triggers, and constraints
- **Gas Industry Pack** (`domain_gas` schema): 4 tables, 12 FKs to core
- **Init scripts**: 6 files in `docker/init/` (roles, schemas, DDL, seeds, comments, pre-request)
- **Naming**: All English (industry standard snake_case)

### Backend (Go + Gin + pgx)
- **Framework**: Gin v1.12.0 + pgx/v5 (direct PostgreSQL, no ORM)
- **Auth**: PostgreSQL bcrypt passwords → JWT (HS256, access + refresh tokens)
- **Architecture**: Handler → Service → Repository pattern
- **~50 handler files**, ~15 service files, ~50 repository files
- **Security**: Rate limiting (login 10/min, API 100/min), CORS allowlist, security headers, body limit 10MB
- **Logging**: slog (structured), Request ID propagation
- **Graceful shutdown** support

### Frontend (Next.js 15 + React 19 + TypeScript)
- **UI Framework**: shadcn/ui (Radix primitives) + Tailwind CSS v4
- **State Management**: Zustand (auth store with persistence)
- **Data Fetching**: TanStack React Query v5
- **Forms**: React Hook Form + Zod validation
- **Tables**: TanStack React Table v8
- **Maps**: MapLibre GL v6
- **Auth Flow**: JWT stored in Zustand + cookie for middleware

### Infrastructure
- **Docker Compose**: postgres, traefik, go-backend, frontend
- **Traefik v3**: API gateway, JWT validation, routing
- **Ports**: Traefik :8088, PostgreSQL :5432, Frontend :3000, Backend :8081

---

## 3. Current Work: Frontend Testing with Vitest

### Test Configuration
- **Framework**: Vitest v4.1.11 with jsdom environment
- **Setup**: `@testing-library/jest-dom/vitest` matchers
- **Coverage**: v8 provider, text/html/lcov/json-summary reporters
- **Thresholds**: statements: 10%, branches: 15%, functions: 6%, lines: 10%

### Test Files Created
```
frontend/src/__tests__/
├── setup.ts                    # jest-dom matchers
├── middleware.test.ts           # Route protection tests
├── components/
│   ├── customer-form.test.tsx   # Form validation, submission
│   ├── partner-form.test.tsx    # Partner form tests
│   ├── tech-role-form.test.tsx  # Tech role form tests
│   ├── user-form.test.tsx       # User form tests
│   └── visit-form.test.tsx      # Visit form tests
├── hooks/
│   ├── hooks-crud-factory.test.tsx  # Parameterized CRUD hook tests (23 entities)
│   ├── hooks-scoped-factory.test.tsx # Scoped hook tests
│   ├── use-auth.test.tsx        # Auth hook tests
│   ├── use-customers.test.tsx   # Customer hook tests
│   └── use-visit-checkpoints.test.tsx # Visit checkpoint tests
├── lib/
│   ├── api.test.ts              # API wrapper, auto-refresh, CRUD factory
│   ├── auth.test.ts             # Auth helpers
│   ├── rbac.test.ts             # Role-based access control
│   └── utils.test.ts            # Utility functions
└── stores/
    └── auth-store.test.ts       # Zustand auth store
```

### Testing Patterns Used
1. **Mocking Strategy**: `vi.mock()` for modules, `vi.fn()` for functions
2. **Wrapper Pattern**: QueryClientProvider wrapper for hooks requiring React Query
3. **Parameterized Tests**: `describe.each()` for testing 23 CRUD hooks with same assertions
4. **User Event Simulation**: `@testing-library/user-event` for form interactions
5. **Mock Store State**: `useAuthStore.setState()` for auth-dependent tests

---

## 4. Files Modified/Created (Key Paths)

### Frontend Structure
```
frontend/
├── src/
│   ├── app/
│   │   ├── (auth)/login/page.tsx          # Login page
│   │   ├── (dashboard)/
│   │   │   ├── layout.tsx                 # Dashboard layout (sidebar + topbar)
│   │   │   ├── page.tsx                   # Dashboard home
│   │   │   ├── core/                      # Users, tech roles
│   │   │   ├── customers/                 # Customer management
│   │   │   ├── inventory/                 # Materials, tools, EPP, vehicles, etc.
│   │   │   ├── operations/                # Visits, reports, assignments
│   │   │   ├── planning/                  # Daily plans, routes
│   │   │   ├── partners/                  # Partner management
│   │   │   ├── notifications/             # Notification templates
│   │   │   ├── geocoding/                 # Address management
│   │   │   └── shared/                    # Report templates
│   │   ├── layout.tsx                     # Root layout
│   │   └── globals.css                    # Global styles
│   ├── components/
│   │   ├── ui/                            # shadcn/ui components (15)
│   │   ├── forms/                         # Form components (13)
│   │   ├── layout/                        # Sidebar, topbar, breadcrumbs
│   │   ├── data-table/                    # Reusable data table
│   │   ├── customer/                      # Customer-specific components
│   │   ├── map/                           # Map components
│   │   ├── partner/                       # Partner components
│   │   └── providers.tsx                  # React Query provider
│   ├── hooks/                             # Custom hooks (54 files)
│   ├── lib/
│   │   ├── api.ts                         # API wrapper with auto-refresh
│   │   ├── auth.ts                        # Auth helpers
│   │   ├── rbac.ts                        # Role-based access control
│   │   ├── utils.ts                       # Utility functions (cn)
│   │   └── route-types-config.ts          # Route type configuration
│   ├── stores/
│   │   └── auth-store.ts                  # Zustand auth store
│   ├── types/
│   │   └── api.ts                         # TypeScript interfaces (674 lines)
│   └── middleware.ts                       # Route protection + photo proxy
├── vitest.config.ts                       # Test configuration
└── package.json                           # Dependencies
```

### Backend Structure
```
backend/
├── cmd/server/main.go                     # Entry point
├── internal/
│   ├── api/                               # Router setup
│   ├── auth/                              # JWT auth (handler, jwt, middleware, revocation)
│   ├── config/                            # Configuration
│   ├── handler/                           # HTTP handlers (12 domains)
│   │   ├── core/                          # Users, tech roles
│   │   ├── customers/                     # Customers, properties
│   │   ├── geocoding/                     # Addresses
│   │   ├── inventory/                     # Materials, tools, vehicles
│   │   ├── operations/                    # Visits, measurements, photos
│   │   ├── partners/                      # Partners, agreements, SLAs
│   │   ├── planning/                      # Daily plans, routes
│   │   ├── notifications/                 # Notifications
│   │   └── shared/                        # Report templates
│   ├── middleware/                         # CORS, rate limiting, security headers
│   ├── model/                             # Domain models (10 domains)
│   ├── repository/                        # Data access (interfaces + postgres impl)
│   ├── service/                           # Business logic (8 domains)
│   ├── testutil/                          # Test utilities
│   └── types/                             # Shared types
└── Dockerfile
```

---

## 5. Technical Details and Patterns

### Authentication Flow
```
1. Login: email + password → POST /auth/login
2. Backend: bcrypt.CompareHashAndPassword → Generate JWT (HS256)
3. Response: { access_token, refresh_token }
4. Frontend: Store in Zustand + cookie (for middleware)
5. API calls: Authorization: Bearer <access_token>
6. Auto-refresh: 401 → POST /auth/refresh → retry original request
```

### JWT Configuration
- **Algorithm**: HS256 (hardcoded)
- **Access token TTL**: 1 hour
- **Refresh token TTL**: 7 days
- **Claims**: `{ sub: user_id, role, token_type, jti, iss: "localis-backend" }`
- **Secret**: ≥32 bytes from `JWT_SECRET` env var

### RBAC Model
- **Roles**: admin, manager, supervisor, technician
- **Module Access**: Route-based (e.g., `/core` → admin only)
- **Action Permissions**: create/edit → admin+manager, delete → admin only

### API Pattern
```typescript
// Generic CRUD factory
export function createCrudApi<T extends { id: string }>(basePath: string) {
  return {
    list: (params?) => apiFetch<T[]>(url),
    get: (id) => apiFetch<T>(`${basePath}/${id}`),
    create: (data) => apiFetch<T>(basePath, { method: "POST", body: JSON.stringify(data) }),
    update: (id, data) => apiFetch<T>(`${basePath}/${id}`, { method: "PUT" }),
    delete: (id) => apiFetch<void>(`${basePath}/${id}`, { method: "DELETE" }),
  };
}
```

### Hook Pattern (TanStack React Query)
```typescript
export function useMaterials() {
  return useQuery({
    queryKey: ["materials"],
    queryFn: () => apiFetch<{ items: Material[]; total: number }>("/api/materials").then(r => r.items),
  });
}

export function useCreateMaterial() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data) => apiFetch<Material>("/api/materials", { method: "POST", body: JSON.stringify(data) }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["materials"] }),
  });
}
```

### Form Pattern (React Hook Form + Zod)
```typescript
const schema = z.object({
  name: z.string().min(1, "Required"),
  email: z.string().email("Invalid email").optional(),
});

export function FormComponent() {
  const { register, handleSubmit, formState: { errors } } = useForm({
    resolver: zodResolver(schema),
  });
  // ...
}
```

### Database Schema Patterns
- **UUIDs**: Primary keys (`gen_random_uuid()`)
- **Timestamps**: `created_at`, `updated_at` with triggers
- **Soft deletes**: Not implemented (hard deletes)
- **Audit**: `plan_audit_logs` table with triggers
- **Geometry**: PostGIS `geography(POINT, 4326)` for locations

---

## 6. Decisions and Issues

### Architecture Decisions
1. **Single-tenant**: No `company_id`, no RLS (simplified for v1)
2. **Generic Core + Industry Packs**: Gas-specific data in `domain_gas` schema
3. **Direct pgx**: No ORM (hand-written repos for control)
4. **JWT in Zustand + cookie**: Cookie needed for Next.js middleware
5. **Photo proxy**: Next.js middleware rewrites `/api/photos/*` to backend

### Issues Resolved
1. **LDAP → PostgreSQL auth**: Migrated from GLAuth to direct PostgreSQL bcrypt
2. **License plates for vehicles**: Fallback to plate number when name is empty
3. **Tool codes**: Use `code` field instead of `name` for tools/equipment
4. **Photo URLs**: Constructed from photo ID at application layer (not stored)
5. **Dark mode**: Not implemented (using white theme with black accents)
6. **Mobile sync 404**: Fixed technician lookup to use `user_id` instead of `employee_number`
7. **Visit form defaults**: Changed from `"medium"`/`"pending"` to `"normal"`/`"scheduled"` to match select options

### Current Limitations
1. **No unit tests for backend**: Only frontend Vitest tests
2. **No E2E tests**: Playwright configured but not implemented
3. **No i18n**: Spanish hardcoded in UI
4. **No offline support**: No service worker/PWA
5. **No real-time**: No WebSocket/SSE implementation

---

## 7. What's Left To Do

### High Priority
- [ ] Backend unit tests (Go testing with testify)
- [ ] API integration tests
- [ ] Complete CRUD pages for all entities
- [ ] Form validation on all forms
- [ ] Error handling and loading states

### Medium Priority
- [ ] E2E tests with Playwright
- [ ] Photo upload functionality
- [ ] Map integration for geocoding
- [ ] Report generation (PDF)
- [ ] Notification system (email/SMS)

### Low Priority
- [ ] Offline support (PWA)
- [ ] Real-time updates (WebSocket)
- [ ] Internationalization (i18n)
- [ ] Dark mode support
- [ ] Performance optimization (virtual scrolling, lazy loading)

### Testing Goals
- [ ] Frontend coverage: 50% statements, 60% branches
- [ ] Backend coverage: 70% functions
- [ ] E2E critical paths: Login, Visit CRUD, Material management

---

## 8. Development Commands

```bash
# Start development stack
docker compose up -d

# With Gas Industry Pack
docker compose -f docker-compose.yml -f industry-packs/gas/docker-compose.override.yml up -d

# Frontend development
cd frontend
npm run dev          # Start Next.js dev server
npm run test        # Run Vitest tests
npm run test:coverage  # Run with coverage

# Backend development
cd backend
go run cmd/server/main.go  # Start Go server
go test ./...              # Run Go tests

# Database access
psql -h localhost -p 5432 -U fsm_admin -d fsm_fsm
```

---

## 9. Environment Variables

```bash
# .env
POSTGRES_PASSWORD=your_secure_password
JWT_SECRET=your_jwt_secret_at_least_32_bytes

# Frontend
NEXT_PUBLIC_API_URL=http://localhost:8081
```

---

## 10. Key Contacts

- **Company**: GasConecta SpA
- **Domain**: Gas installations, maintenance, PH tests (Chile)
- **Regulatory**: SEC (Superintendencia de Electricidad y Combustibles)
- **Certifications**: CL1, CL2, CL3, GLP, TC1, TC2, TC6, GREEN SEAL, DS66
