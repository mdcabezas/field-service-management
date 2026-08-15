import { test, expect, API_BASE, loginAPI, authHeaders, createAndReturnId, expectNotFound, expectForbidden, CREDENTIALS } from './api-helpers';

test.describe('Core - Users API', () => {
  let headers: Record<string, string>;

  test.beforeAll(async ({ request }) => {
    const tokens = await loginAPI(request);
    headers = await authHeaders(tokens.access_token);
  });

  test('GET /api/users returns list', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/users`, { headers });
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(Array.isArray(body) || body.items).toBeTruthy();
  });

  test('GET /api/users/:id returns user', async ({ request }) => {
    const listRes = await request.get(`${API_BASE}/api/users`, { headers });
    const body = await listRes.json();
    const users = Array.isArray(body) ? body : body.items;
    if (users.length === 0) return;

    const res = await request.get(`${API_BASE}/api/users/${users[0].id}`, { headers });
    expect(res.status()).toBe(200);
    const user = await res.json();
    expect(user.id).toBe(users[0].id);
  });

  test('GET /api/users/:id not found returns 404', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/users/00000000-0000-0000-0000-000000000000`, { headers });
    expect(res.status()).toBe(404);
  });

  test('POST /api/users creates user', async ({ request }) => {
    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/users`, {
      email: `test_${Date.now()}@test.com`,
      role: 'technician',
      name: 'Test User',
    });
    expect(id).toBeTruthy();
    await request.delete(`${API_BASE}/api/users/${id}`, { headers });
  });

  test('PUT /api/users/:id updates user', async ({ request }) => {
    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/users`, {
      email: `test_upd_${Date.now()}@test.com`,
      role: 'technician',
      name: 'Test User 2',
    });

    const res = await request.put(`${API_BASE}/api/users/${id}`, {
      headers,
      data: { name: 'Updated User', email: `test_upd2_${Date.now()}@test.com`, role: 'technician' },
    });
    expect(res.status()).toBe(200);
    await request.delete(`${API_BASE}/api/users/${id}`, { headers });
  });

  test('DELETE /api/users/:id deletes user', async ({ request }) => {
    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/users`, {
      email: `test_del_${Date.now()}@test.com`,
      role: 'technician',
      name: 'Test User 3',
    });

    const res = await request.delete(`${API_BASE}/api/users/${id}`, { headers });
    expect(res.status()).toBe(204);
    await expectNotFound(request, headers, `${API_BASE}/api/users/${id}`);
  });

  test('Technician cannot access users', async ({ request }) => {
    const techTokens = await loginAPI(request, CREDENTIALS.technician);
    const techHeaders = await authHeaders(techTokens.access_token);
    await expectForbidden(request, techHeaders, `${API_BASE}/api/users`, 'GET');
  });
});

test.describe('Core - Tech Roles API', () => {
  let headers: Record<string, string>;

  test.beforeAll(async ({ request }) => {
    const tokens = await loginAPI(request);
    headers = await authHeaders(tokens.access_token);
  });

  test('GET /api/tech-roles returns list', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/tech-roles`, { headers });
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(Array.isArray(body) || body.items).toBeTruthy();
  });

  test('GET /api/tech-roles/:id returns role', async ({ request }) => {
    const listRes = await request.get(`${API_BASE}/api/tech-roles`, { headers });
    const body = await listRes.json();
    const roles = Array.isArray(body) ? body : body.items;
    if (roles.length === 0) return;

    const res = await request.get(`${API_BASE}/api/tech-roles/${roles[0].id}`, { headers });
    expect(res.status()).toBe(200);
  });

  test('POST /api/tech-roles creates role', async ({ request }) => {
    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/tech-roles`, {
      name: 'Test Role',
      code: 'test_role_' + Date.now(),
      active: true,
    });
    expect(id).toBeTruthy();
    await request.delete(`${API_BASE}/api/tech-roles/${id}`, { headers });
  });

  test('PUT /api/tech-roles/:id updates role', async ({ request }) => {
    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/tech-roles`, {
      name: 'Test Role 2',
      code: 'test_role2_' + Date.now(),
      active: true,
    });

    const res = await request.put(`${API_BASE}/api/tech-roles/${id}`, {
      headers,
      data: { name: 'Updated Role', code: 'test_role2_upd_' + Date.now(), active: true },
    });
    expect(res.status()).toBe(200);
    await request.delete(`${API_BASE}/api/tech-roles/${id}`, { headers });
  });

  test('DELETE /api/tech-roles/:id deletes role', async ({ request }) => {
    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/tech-roles`, {
      name: 'Test Role 3',
      code: 'test_role3_' + Date.now(),
      active: true,
    });

    const res = await request.delete(`${API_BASE}/api/tech-roles/${id}`, { headers });
    expect(res.status()).toBe(204);
  });
});
