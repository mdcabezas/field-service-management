## Context

El backend Go ya expone 163 endpoints REST con JWT auth (LDAP + GLAuth). No existe frontend. Se necesita una UI web para que operadores, supervisores y administradores gestionen el sistema.

Stack decidido: Next.js 15, TypeScript, Tailwind CSS, shadcn/ui, TanStack Table, React Hook Form + Zod, TanStack Query, Zustand.

Diseño visual: Brutalista (borders gruesos, sin rounded, alto contraste, monospace).

## Goals / Non-Goals

**Goals:**
- Login LDAP con JWT (access + refresh tokens en httpOnly cookies)
- Layout autenticado con sidebar, topbar y navegación por módulos
- RBAC client-side (admin/manager/supervisor) con ocultamiento de módulos
- CRUD completo para los 9 dominios del backend
- Docker multi-stage para dev y prod
- Consumo directo de todos los endpoints existentes sin modificar el backend

**Non-Goals:**
- App móvil nativa
- Multi-tenant (el backend es single-tenant)
- WebSocket/real-time (uso futuro, no inicial)
- Internacionalización (solo español)
- Modo offline/PWA
- Testing E2E inicial (se agrega después)

## Decisions

### 1. Next.js App Router (no Pages Router)

**Decisión:** Usar App Router con Server Components por defecto.

**Razón:** Server Components reducen bundle size, mejoran SEO inicial (aunque es app autenticada), y facilitan streaming. Client Components solo donde hay interactividad (forms, tables).

**Alternativa considerada:** Pages Router más maduro pero sin beneficios de Server Components.

### 2. Estructura de directorios por dominio

**Decisión:** Organizar por feature, no por tipo de archivo.

```
frontend/
├── app/
│   ├── (auth)/login/page.tsx
│   └── (dashboard)/
│       ├── layout.tsx           # sidebar + topbar
│       ├── partners/
│       │   ├── page.tsx         # list
│       │   └── [id]/page.tsx    # detail
│       └── ...
├── components/
│   ├── ui/                      # shadcn primitives
│   ├── layout/                  # sidebar, topbar, breadcrumbs
│   └── data-table/              # TanStack Table wrapper
├── lib/
│   ├── api.ts                   # fetch wrapper con JWT
│   ├── auth.ts                  # cookie helpers
│   └── rbac.ts                  # permisos por rol
├── hooks/
│   ├── use-auth.ts              # auth state
│   └── use-entity.ts            # TanStack Query hooks genéricos
├── stores/
│   └── auth-store.ts            # Zustand
└── types/
    └── api.ts                   # TypeScript types del backend
```

**Razón:** Co-locación de feature-related files. Fácil de encontrar y mantener.

**Alternativa considerada:** Por tipo (components/, hooks/, lib/) - más difícil de mantener a escala.

### 3. API wrapper centralizado

**Decisión:** Un solo `lib/api.ts` que maneja JWT, refresh automático y error handling.

```typescript
// Pseudocódigo del wrapper
async function apiFetch(path, options) {
  const token = getAccessToken()
  const res = await fetch(NEXT_PUBLIC_API_URL + path, {
    ...options,
    headers: { Authorization: `Bearer ${token}` }
  })
  if (res.status === 401) {
    await refreshToken()
    return apiFetch(path, options) // retry
  }
  return res
}
```

**Razón:** DRY. Un solo punto para manejar refresh, rate limiting, error logging.

**Alternativa considerada:** Axios interceptor - más pesado, innecesario con fetch nativo.

### 4. TanStack Query para server state

**Decisión:** Cada endpoint del backend tiene un hook en `hooks/use-entity.ts` que usa TanStack Query.

**Razón:** Cache automático, optimistic updates, re-fetch en focus/stale, loading/error states gratis.

**Alternativa considerada:** SWR - más simple pero menos features (no hay mutation optimistic tan fácil).

### 5. Zustand para client state

**Decisión:** Zustand store para auth (tokens, user info, role) y UI preferences (sidebar collapsed).

**Razón:** Simple, sin boilerplate, funciona bien con Server Components. Redux overkill para este caso.

**Alternativa considerada:** React Context - re-renders innecesarios, no persiste entre navegaciones.

### 6. React Hook Form + Zod para forms

**Decisión:** Cada form usa React Hook Form con Zod schema para validación.

**Razón:** Validación type-safe, performance (no re-renderiza en cada keystroke), integración natural con TypeScript.

**Alternativa considerada:** Formik - más pesado, menos type-safe.

### 7. TanStack Table para listados

**Decisión:** Wrapper genérico `<DataTable columns={} data={} />` que todos los módulos reutilizan.

**Razón:** Pagination, sorting, filtering, column visibility - todo incluido. Headless = control total sobre UI brutalista.

**Alternativa considerada:** Table de shadcn (usa TanStack Table internamente) - menos control.

### 8. Brutalismo como design system

**Decisión:** Estilos consistentes via Tailwind classes:
- `border-2 border-black` en todos los componentes
- `rounded-none` global (override de shadcn defaults)
- `font-mono` para labels, headings, buttons
- Colores: black/white/primary (sin grays suaves)
- Sin sombras, sin gradients, sin animaciones suaves

**Razón:** Consistencia visual, distinguishable, rápido de implementar.

### 9. Docker multi-stage

**Decisión:** Dockerfile con stages: deps → build → production.

**Razón:** Image final pequeña (~150MB), solo archivos estáticos + Node runtime para SSR.

### 10. Proxy Traefik

**Decisión:** Traefik enruta `/` → `frontend:3000` y `/api/*` → `backend:8080`.

**Razón:** Single origin para el usuario, CORS simplificado, SSL termination en Traefik.

## Risks / Trade-offs

- **[JWT en cookies]** → httpOnly cookies son seguras pero difíciles de debuggear. Mitigación: logging server-side del refresh flow.

- **[RBAC client-side]** → El backend ya valida permisos (403). RBAC client-side es UX (ocultar botones), no seguridad. Mitigación: nunca confiar solo en client-side RBAC.

- **[163 endpoints]** → Muchos hooks a crear. Mitigación: hooks genéricos (`useCrudEntity`) que aceptan path y columnas.

- **[Brutalismo]** → Puede parecer "feo" a usuarios no técnicos. Mitigación: es una decisión de diseño válida, documentar en README.

- **[Server Components]** → Forms y tables necesitan Client Components. Mitigación: solo los componentes interactivos llevan "use client".

- **[Sin testing E2E]** → Riesgo de regressions. Mitigación: agregar Playwright después del MVP.
