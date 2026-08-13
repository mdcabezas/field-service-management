## Why

El backend Go expone 163 endpoints REST para el sistema Field Service Management, pero no existe interfaz gráfica. Los operadores y técnicos necesitan una UI para gestionar clientes, inventario, operaciones y planificación diaria. Sin frontend, el sistema solo es utilizable vía API directa o herramientas externas, lo cual limita su adopción y productividad.

## What Changes

- **Nuevo frontend Next.js 15** con TypeScript, Tailwind CSS y shadcn/ui
- **Login LDAP** con autenticación JWT (access + refresh tokens)
- **Layout autenticado** con sidebar brutalista, topbar y navegación por módulos
- **RBAC client-side**: admin ve todo, manager y supervisor ven módulos según permisos
- **CRUD completo** para los 9 dominios: Core, Partners, Customers, Inventory, Planning, Operations, Notifications, Geocoding, Shared
- **Tables con TanStack Table** para listados paginados y ordenables
- **Forms con React Hook Form + Zod** para validación de schemas del backend
- **Server state con TanStack Query** para cache, optimistic updates y re-fetch automático
- **Client state con Zustand** para auth store y preferencias de UI
- **Dockerfile multi-stage** para dev y prod
- **Proxy Traefik** enrutando tráfico HTTP al frontend (:3000)

## Capabilities

### New Capabilities

- `frontend-auth`: Login LDAP, JWT storage (httpOnly cookies), refresh automático, logout, middleware de protección de rutas
- `frontend-layout`: Layout autenticado con sidebar, topbar, navegación por módulos, breadcrumbs, brutalismo visual
- `frontend-rbac`: Control de acceso por rol (admin/manager/supervisor), ocultamiento de módulos no permitidos, protección de acciones
- `frontend-core`: CRUD de Users y Tech Roles
- `frontend-partners`: CRUD de Partners, Contacts, Agreements, Docs, Forms, SLAs
- `frontend-customers`: CRUD de Customers, Addresses, Address-Partners, Properties, Technicians, Tech Certifications
- `frontend-inventory`: CRUD de Materials, Tools, EPP, Vehicles, Rentals, Cost Rates, Maintenance, Checklists
- `frontend-planning`: CRUD de Daily Plans, Assignments, Loads, Routes
- `frontend-operations`: CRUD de Visits con todo el sub-árbol (assignments, checklists, measurements, checkpoints, photos, reports, vehicle-assignments, rentals, SLA trackings)
- `frontend-notifications`: CRUD de Notification Templates y Notifications
- `frontend-geocoding`: CRUD de Addresses geocodificadas con búsqueda nearby
- `frontend-shared`: CRUD de Report Templates

### Modified Capabilities

- N/A (proyecto nuevo, sin specs existentes)

## Impact

- **Nuevo directorio**: `frontend/` al mismo nivel que `backend/`
- **Docker Compose**: Se agrega servicio `frontend` en `docker-compose.yml` y `docker-compose.prod.yml`
- **Traefik**: Nueva ruta `/` → `frontend:3000` en `traefik/dynamic.yml`
- **Dependencias npm**: next, react, tailwindcss, @tanstack/react-table, @tanstack/react-query, zustand, react-hook-form, zod, @radix-ui/*, shadcn/ui
- **Backend sin cambios**: Todos los endpoints existentes son consumidos tal cual
- **Variables de entorno**: `NEXT_PUBLIC_API_URL` para URL del backend
