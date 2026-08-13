## Purpose

Frontend authentication layer that stores JWT tokens in memory and localStorage, attaches Bearer tokens to authenticated requests, handles automatic token refresh, and synchronizes access tokens to cookies for middleware route protection. Covers the complete login flow, token lifecycle, and route protection for the admin dashboard.

## ADDED Requirements

### Requirement: Token storage in auth store

The system SHALL store `access_token` and `refresh_token` in the Zustand auth store alongside user info. Tokens SHALL persist across page reloads via localStorage.

#### Scenario: Tokens stored after login

- **WHEN** user submits valid credentials on login page
- **THEN** `authApi.login()` returns `{ access_token, refresh_token, ... }`
- **AND** tokens are passed to `setTokens(accessToken, refreshToken)` in auth store
- **AND** tokens are available via `useAuthStore.getState().accessToken`

#### Scenario: Tokens cleared on logout

- **WHEN** user clicks logout
- **THEN** `authApi.logout()` is called
- **AND** `clearTokens()` is called in auth store
- **AND** localStorage tokens are removed

### Requirement: Bearer token on authenticated requests

The system SHALL attach `Authorization: Bearer <access_token>` header to every authenticated API request via `apiFetch`.

#### Scenario: Authenticated GET request

- **WHEN** any hook calls `apiFetch('/api/users')`
- **THEN** request includes `Authorization: Bearer <access_token>` header
- **AND** backend returns 200 with data

#### Scenario: No token available

- **WHEN** `apiFetch` is called but auth store has no access token
- **THEN** request proceeds without Authorization header
- **AND** backend returns 401 if auth required

### Requirement: Automatic token refresh

The system SHALL automatically refresh expired access tokens using the refresh token before retrying failed requests.

#### Scenario: Token refresh on 401

- **WHEN** `apiFetch` receives 401 and refresh token exists
- **THEN** `authApi.refresh()` is called with `refresh_token` in JSON body
- **AND** new tokens are stored via `setTokens()`
- **AND** original request is retried with new access token

#### Scenario: Refresh fails

- **WHEN** `authApi.refresh()` returns error or no refresh token
- **THEN** tokens are cleared via `clearTokens()`
- **AND** user is redirected to `/login`

### Requirement: Login captures tokens from response

The login page SHALL capture and store tokens returned by `authApi.login()` before fetching user info.

#### Scenario: Successful login stores tokens

- **WHEN** user submits credentials (1001/admin)
- **THEN** `authApi.login()` returns tokens in JSON body
- **AND** tokens are immediately stored via `setTokens()`
- **AND** `authApi.me()` is called with Authorization header
- **AND** user is redirected to `/`

### Requirement: Middleware cookie sync

The system SHALL synchronize the `access_token` to an `access_token` cookie so middleware can verify authentication without client-side store access.

#### Scenario: Cookie set on login

- **WHEN** `authApi.login()` succeeds
- **THEN** `setCookie('access_token', accessToken, 7)` is called
- **AND** cookie is available to middleware on subsequent requests

#### Scenario: Cookie deleted on logout

- **WHEN** `authApi.logout()` is called
- **THEN** `deleteCookie('access_token')` is called
- **AND** cookie is removed

#### Scenario: Cookie deleted on refresh failure

- **WHEN** token refresh fails
- **THEN** `deleteCookie('access_token')` is called
- **AND** user is redirected to login

### Requirement: Correct dashboard route paths

The frontend SHALL use root paths (`/`) instead of `/dashboard` prefix for all authenticated routes.

#### Scenario: Login redirects to root

- **WHEN** login succeeds
- **THEN** `router.push("/")` is called
- **AND** user lands on dashboard at `/`

#### Scenario: Sidebar links use correct paths

- **WHEN** sidebar renders navigation
- **THEN** all links use paths without `/dashboard` prefix (e.g., `/core/users`)
- **AND** `moduleAccess` in `lib/rbac.ts` matches actual routes

### Requirement: Paginated response handling

All list hooks SHALL correctly extract `items` from backend paginated responses `{ items: T[], total: number }`.

#### Scenario: List hook extracts items

- **WHEN** `useUsers()` calls `apiFetch('/api/users')`
- **THEN** response is `{ items: User[], total: number }`
- **AND** hook returns `items` array to component
- **AND** DataTable renders correct row count