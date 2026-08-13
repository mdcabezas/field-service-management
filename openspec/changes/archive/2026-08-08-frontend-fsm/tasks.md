## 1. Fundamentos

- [x] 1.1 Scaffold Next.js 15 con TypeScript y App Router
- [x] 1.2 Instalar y configurar Tailwind CSS
- [x] 1.3 Instalar y configurar shadcn/ui con tema brutalista
- [x] 1.4 Crear Dockerfile multi-stage (deps → build → production)
- [x] 1.5 Actualizar docker-compose.yml con servicio frontend
- [x] 1.6 Actualizar docker-compose.prod.yml con servicio frontend
- [x] 1.7 Configurar Traefik dynamic.yml para enrutar / → frontend:3000
- [x] 1.8 Configurar variable de entorno NEXT_PUBLIC_API_URL
- [x] 1.9 Crear tipos TypeScript de la API (types/api.ts)

## 2. Auth

- [x] 2.1 Crear API wrapper con JWT y refresh automático (lib/api.ts)
- [x] 2.2 Crear auth store con Zustand (stores/auth-store.ts)
- [x] 2.3 Crear login page con React Hook Form + Zod
- [x] 2.4 Implementar JWT storage en httpOnly cookies
- [x] 2.5 Implementar refresh automático antes de expiración
- [x] 2.6 Crear middleware de protección de rutas
- [x] 2.7 Implementar logout con revocación de tokens
- [x] 2.8 Crear hook use-auth.ts

## 3. Layout

- [x] 3.1 Crear layout autenticado (app/(dashboard)/layout.tsx)
- [x] 3.2 Crear componente Sidebar con navegación por módulos
- [x] 3.3 Crear componente Topbar con info de usuario y logout
- [x] 3.4 Crear componente Breadcrumbs
- [x] 3.5 Implementar highlight de ruta activa en sidebar
- [x] 3.6 Aplicar estilos brutalistas a todos los componentes UI

## 4. RBAC

- [x] 4.1 Crear lib/rbac.ts con permisos por rol
- [x] 4.2 Implementar ocultamiento de módulos no permitidos en sidebar
- [x] 4.3 Implementar protección de rutas por rol
- [x] 4.4 Implementar ocultamiento de acciones (create/edit/delete) por permiso
- [x] 4.5 Manejar 403 del backend con toast de "No tiene permisos"

## 5. DataTable Genérico

- [x] 5.1 Crear componente DataTable con TanStack Table
- [x] 5.2 Implementar paginación con limit/offset
- [x] 5.3 Implementar búsqueda por columnas
- [x] 5.4 Implementar sorting por columnas
- [x] 5.5 Crear componente DataTableToolbar con filtros

## 6. Core - Users y Tech Roles

- [x] 6.1 Crear hook use-users con TanStack Query
- [x] 6.2 Crear Users list page con DataTable
- [x] 6.3 Crear User form (create/edit) con React Hook Form
- [x] 6.4 Crear User detail page
- [x] 6.5 Crear hook use-tech-roles con TanStack Query
- [x] 6.6 Crear Tech Roles list page con DataTable
- [x] 6.7 Crear Tech Role form (create/edit) con React Hook Form
- [x] 6.8 Crear Tech Role detail page

## 7. Partners

- [x] 7.1 Crear hook use-partners con TanStack Query
- [x] 7.2 Crear Partners list page con DataTable
- [x] 7.3 Crear Partner form (create/edit) con React Hook Form
- [x] 7.4 Crear Partner detail page con tabs
- [x] 7.5 Crear Partner Contacts CRUD (hooks, list, form)
- [x] 7.6 Crear Partner Agreements CRUD (hooks, list, form)
- [x] 7.7 Crear Partner Agreement Docs CRUD (hooks, list, form)
- [x] 7.8 Crear Partner Agreement Forms CRUD (hooks, list, form)
- [x] 7.9 Crear SLAs CRUD (hooks, list, form)

## 8. Customers

- [x] 8.1 Crear hook use-customers con TanStack Query
- [x] 8.2 Crear Customers list page con DataTable
- [x] 8.3 Crear Customer form (create/edit) con React Hook Form
- [x] 8.4 Crear Customer detail page con tabs
- [x] 8.5 Crear Customer Addresses CRUD (hooks, list, form)
- [x] 8.6 Crear Address Partners CRUD (hooks, list, form)
- [x] 8.7 Crear Properties CRUD (hooks, list, form)
- [x] 8.8 Crear Technicians CRUD (hooks, list, form)
- [x] 8.9 Crear Tech Certifications CRUD (hooks, list, form)

## 9. Inventory

- [x] 9.1 Crear Materials CRUD (hooks, list, form, detail)
- [x] 9.2 Crear Tools CRUD (hooks, list, form, detail)
- [x] 9.3 Crear EPP Items CRUD (hooks, list, form, detail)
- [x] 9.4 Crear Vehicles CRUD (hooks, list, form, detail)
- [x] 9.5 Crear Rentals CRUD (hooks, list, form, detail)
- [x] 9.6 Crear Cost Rates CRUD (hooks, list, form, detail)
- [x] 9.7 Crear Maintenance Records CRUD (hooks, list, form)
- [x] 9.8 Crear Maintenance Schedules CRUD (hooks, list, form)
- [x] 9.9 Crear Checklist Templates CRUD (hooks, list, form, detail)
- [x] 9.10 Crear Checklist Template Materials/Tools/EPPs CRUD

## 10. Planning

- [x] 10.1 Crear Daily Plans CRUD (hooks, list, form, detail)
- [x] 10.2 Crear Daily Plan Assignments CRUD (hooks, list, form)
- [x] 10.3 Crear Daily Load Materials CRUD (hooks, list, form)
- [x] 10.4 Crear Daily Load Tools CRUD (hooks, list, form)
- [x] 10.5 Crear Daily Load EPPs CRUD (hooks, list, form)
- [x] 10.6 Crear Routes CRUD (hooks, list, form, detail)

## 11. Operations

- [x] 11.1 Crear Visits CRUD (hooks, list, form, detail)
- [x] 11.2 Crear Visit Assignments CRUD (hooks, list, form)
- [x] 11.3 Crear Visit Checklist Materials/Tools/EPPs CRUD
- [x] 11.4 Crear Visit Measurements CRUD (hooks, list, form)
- [x] 11.5 Crear Visit Checkpoints CRUD (hooks, list, form)
- [x] 11.6 Crear Visit Photos CRUD (hooks, list, form)
- [x] 11.7 Crear Visit Reports CRUD (hooks, list, form, detail)
- [x] 11.8 Crear Report Entries CRUD (hooks, list, form)
- [x] 11.9 Crear Report Images CRUD (hooks, list, form)
- [x] 11.10 Crear Vehicle Assignments CRUD (hooks, list, form)
- [x] 11.11 Crear Visit Rentals CRUD (hooks, list, form)
- [x] 11.12 Crear SLA Trackings CRUD (hooks, list, form)

## 12. Notifications

- [x] 12.1 Crear Notification Templates CRUD (hooks, list, form, detail)
- [x] 12.2 Crear Notifications CRUD (hooks, list, form, detail)

## 13. Geocoding

- [x] 13.1 Crear Geocoding Addresses CRUD (hooks, list, form, detail)
- [x] 13.2 Implementar búsqueda nearby con mapa

## 14. Shared

- [x] 14.1 Crear Report Templates CRUD (hooks, list, form, detail)

## 15. Pulido

- [x] 15.1 Agregar loading states (skeletons) en todas las páginas
- [x] 15.2 Agregar error boundaries
- [x] 15.3 Implementar responsive layout (tablet/mobile)
- [x] 15.4 Agregar toast notifications para éxito/error en CRUD
- [x] 15.5 Crear dashboard principal con métricas resumen
- [x] 15.6 Documentar setup y desarrollo en README.md
