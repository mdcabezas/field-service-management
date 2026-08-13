# Fix Auth Token Flow

## Problem

The frontend login is broken. After valid credentials (1001/admin), the backend authenticates via LDAP and returns JWT tokens in JSON body. But the frontend:

1. Discards the tokens from `authApi.login()` — never stores them
2. Calls `authApi.me()` without `Authorization` header → 401 → "Credenciales inválidas"

The frontend was designed assuming httpOnly cookies, but the backend uses `Authorization: Bearer` headers exclusively.

## Solution

Modify frontend auth layer to:
1. Store `access_token` and `refresh_token` in memory (+ localStorage for persistence across page reloads)
2. Send `Authorization: Bearer <access_token>` on every authenticated request
3. Send `refresh_token` in JSON body for refresh requests
4. Clear tokens on logout

## Scope

- `frontend/src/lib/api.ts` — apiFetch + authApi
- `frontend/src/stores/auth-store.ts` — token storage
- `frontend/src/app/(auth)/login/page.tsx` — capture tokens

## Out of Scope

- Backend changes (backend is correct and complete)
- httpOnly cookie approach (would require backend changes)
