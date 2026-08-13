# Tasks: Fix Auth Token Flow

## 1. Token Storage

- [x] 1.1 Update `auth-store.ts` to store `accessToken` and `refreshToken` alongside user info
- [x] 1.2 Add `setTokens(accessToken, refreshToken)` and `clearTokens()` actions to the store
- [x] 1.3 Persist tokens to localStorage via zustand persist middleware

## 2. API Layer

- [x] 2.1 Update `apiFetch` to read `accessToken` from auth store and attach `Authorization: Bearer` header
- [x] 2.2 Update `authApi.login()` to return tokens and store them via `setTokens()`
- [x] 2.3 Update `authApi.refresh()` to send `refresh_token` in JSON body and store new tokens
- [x] 2.4 Update `authApi.me()` to use authenticated `apiFetch` (no `skipAuth`)
- [x] 2.5 Update `authApi.logout()` to clear tokens after request

## 3. Login Page

- [x] 3.1 Update `login/page.tsx` to capture tokens from login response before calling `me()`

## 4. Verify

- [x] 4.1 Run `npm run build` to verify no type errors
- [x] 4.2 Test login flow end-to-end with curl-verified credentials (1001/admin)

## 5. Middleware Cookie Sync

- [x] 5.1 Add `setCookie(name, value, days)` and `deleteCookie(name)` helpers to `lib/api.ts`
- [x] 5.2 Set `access_token` cookie in `authApi.login()` after storing tokens
- [x] 5.3 Delete `access_token` cookie in `authApi.logout()` and refresh failure paths
- [x] 5.4 Run `npm run build` to verify

## 6. Fix Dashboard Route

- [x] 6.1 Update `login/page.tsx` — change `router.push("/dashboard")` to `router.push("/")`
- [x] 6.2 Update `middleware.ts` — change `protectedRoutes` from `/dashboard` to match actual routes
- [x] 6.3 Update `lib/api.ts` — change all `window.location.href = "/login"` to use correct path
- [x] 6.4 Run `npm run build` to verify

## 7. Fix Sidebar Routes

- [x] 7.1 Update `sidebar.tsx` — remove `/dashboard` prefix from all paths
- [x] 7.2 Update `lib/rbac.ts` — remove `/dashboard` prefix from `moduleAccess` keys
- [x] 7.3 Update `(dashboard)/page.tsx` — fix quick link paths
- [x] 7.4 Run `npm run build` to verify

## 8. Fix Paginated Response Handling

- [x] 8.1 Update all list hooks to extract `items` from `{ items, total }` response
- [x] 8.2 Create `core/page.tsx` redirect to `/core/users`
- [x] 8.3 Create `core/users/new/page.tsx` for user creation
- [x] 8.4 Create `core/tech-roles/new/page.tsx` for tech-role creation
- [x] 8.5 Run `npm run build` to verify
