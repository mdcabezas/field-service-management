# Design: Auth Token Flow Fix

## Token Storage Strategy

Store tokens in a Zustand persist store backed by localStorage. This provides:
- Persistence across page reloads
- Easy access from any component via `useAuthStore`
- Automatic cleanup on logout

**Storage shape:**
```typescript
{
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  user: { employeeNumber: string; role: UserRole } | null;
}
```

## Request Flow

```
Login:
  1. POST /auth/login → receive { access_token, refresh_token }
  2. Store both tokens in auth store (persisted to localStorage)
  3. Call GET /auth/me with Authorization header → get user info
  4. Store user info, redirect to /dashboard

Authenticated Request:
  1. apiFetch reads accessToken from auth store
  2. Attaches Authorization: Bearer <token> header
  3. If 401 → attempt refresh
  4. If refresh succeeds → retry original request
  5. If refresh fails → clear tokens, redirect to /login

Refresh:
  1. POST /auth/refresh with { refresh_token } in body
  2. Receive new { access_token, refresh_token }
  3. Store new tokens
  4. Retry queued requests

Logout:
  1. POST /auth/logout (best-effort)
  2. Clear tokens from store
  3. Redirect to /login
```

## Security Considerations

- Tokens in localStorage are vulnerable to XSS. This is acceptable for a dev/internal tool.
- For production, consider httpOnly cookies (requires backend change).
- Refresh token rotation is already implemented backend-side (old token revoked on each refresh).
- Access token TTL is 1 hour; refresh token TTL is 7 days.
