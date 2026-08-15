import { test, expect, API_BASE, CREDENTIALS, loginAPI, authHeaders } from './api-helpers';

test.describe('Auth API', () => {
  test('POST /auth/login returns tokens', async ({ request }) => {
    const res = await request.post(`${API_BASE}/auth/login`, { data: CREDENTIALS.admin });
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(body.access_token).toBeTruthy();
    expect(body.refresh_token).toBeTruthy();
    expect(body.token_type).toBe('Bearer');
    expect(body.expires_in).toBeGreaterThan(0);
  });

  test('POST /auth/login with invalid credentials returns 401', async ({ request }) => {
    const res = await request.post(`${API_BASE}/auth/login`, { data: CREDENTIALS.invalid });
    expect(res.status()).toBe(401);
  });

  test('GET /auth/me returns current user', async ({ request }) => {
    const tokens = await loginAPI(request);
    const res = await request.get(`${API_BASE}/auth/me`, {
      headers: await authHeaders(tokens.access_token),
    });
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(body.employee_number).toBeTruthy();
  });

  test('GET /auth/me without token returns 401', async ({ request }) => {
    const res = await request.get(`${API_BASE}/auth/me`);
    expect(res.status()).toBe(401);
  });

  test('POST /auth/refresh returns new tokens', async ({ request }) => {
    const tokens = await loginAPI(request);
    const res = await request.post(`${API_BASE}/auth/refresh`, {
      data: { refresh_token: tokens.refresh_token },
    });
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(body.access_token).toBeTruthy();
    expect(body.refresh_token).toBeTruthy();
  });
});
