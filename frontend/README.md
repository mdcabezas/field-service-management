# Frontend - Localis FSM

Frontend web para el sistema Field Service Management.

## Stack

- **Framework:** Next.js 15 (App Router)
- **Language:** TypeScript
- **Styling:** Tailwind CSS + shadcn/ui
- **Tables:** TanStack Table
- **Forms:** React Hook Form + Zod
- **State:** TanStack Query (server) + Zustand (client)
- **Design:** Brutalismo

## Desarrollo

```bash
# Instalar dependencias
npm install

# Ejecutar en desarrollo
npm run dev

# Build de producción
npm run build

# Ejecutar producción
npm start
```

## Variables de Entorno

Copiar `.env.example` a `.env.local`:

```bash
NEXT_PUBLIC_API_URL=http://localhost:8081
```

## Estructura

```
frontend/
├── app/
│   ├── (auth)/login/page.tsx    # Login page
│   └── (dashboard)/             # Dashboard layout
│       ├── page.tsx             # Dashboard principal
│       ├── core/                # Users, Tech Roles
│       ├── partners/            # Partners, Contacts, Agreements, SLAs
│       ├── customers/           # Customers, Addresses, Technicians
│       ├── inventory/           # Materials, Tools, EPP, Vehicles
│       ├── planning/            # Daily Plans, Routes
│       ├── operations/          # Visits, Reports
│       ├── notifications/       # Templates, Notifications
│       ├── geocoding/           # Addresses
│       └── shared/              # Report Templates
├── components/
│   ├── ui/                      # shadcn components (brutalist)
│   ├── layout/                  # Sidebar, Topbar, Breadcrumbs
│   ├── data-table/              # TanStack Table wrapper
│   └── forms/                   # React Hook Form components
├── hooks/                       # TanStack Query hooks
├── lib/
│   ├── api.ts                   # API wrapper con JWT
│   ├── auth.ts                  # Cookie helpers
│   └── rbac.ts                  # Role-based access control
├── stores/                      # Zustand stores
├── types/                       # TypeScript types
└── middleware.ts                 # Route protection
```

## Docker

```bash
# Desarrollo
docker compose --profile dev up frontend

# Producción
docker compose -f docker-compose.prod.yml up frontend
```

## Autenticación

El frontend usa JWT con httpOnly cookies:
- Login via LDAP (employee_number + password)
- Access token: 1 hora
- Refresh token: 7 días
- Refresh automático antes de expiración

## RBAC

| Rol | Módulos |
|-----|---------|
| admin | Todos |
| manager | Partners, Customers, Inventory, Planning, Operations, Notifications |
| supervisor | Operations (limitado) |
| technician | (pendiente) |
