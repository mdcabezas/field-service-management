import { test, expect, API_BASE, loginAPI, authHeaders, createAndReturnId, expectNotFound, expectForbidden, CREDENTIALS } from './api-helpers';

test.describe('Planning - Daily Plans API', () => {
  let headers: Record<string, string>;

  test.beforeAll(async ({ request }) => {
    const tokens = await loginAPI(request);
    headers = await authHeaders(tokens.access_token);
  });

  test('GET /api/daily-plans returns list', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/daily-plans`, { headers });
    expect(res.status()).toBe(200);
  });

  test('GET /api/daily-plans/:id returns plan', async ({ request }) => {
    const listRes = await request.get(`${API_BASE}/api/daily-plans`, { headers });
    const body = await listRes.json();
    const plans = Array.isArray(body) ? body : body.items;
    if (plans.length === 0) return;

    const res = await request.get(`${API_BASE}/api/daily-plans/${plans[0].id}`, { headers });
    expect(res.status()).toBe(200);
  });

  test('POST /api/daily-plans creates plan', async ({ request }) => {
    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/daily-plans`, {
      name: 'Test Plan ' + Date.now(),
      date: new Date().toISOString(),
    });
    expect(id).toBeTruthy();
    await request.delete(`${API_BASE}/api/daily-plans/${id}`, { headers });
  });

  test('PUT /api/daily-plans/:id updates plan', async ({ request }) => {
    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/daily-plans`, {
      name: 'Plan Update ' + Date.now(),
      date: new Date().toISOString(),
    });

    const res = await request.put(`${API_BASE}/api/daily-plans/${id}`, {
      headers,
      data: { name: 'Updated Plan', date: new Date().toISOString() },
    });
    expect(res.status()).toBe(200);
    await request.delete(`${API_BASE}/api/daily-plans/${id}`, { headers });
  });

  test('DELETE /api/daily-plans/:id deletes plan', async ({ request }) => {
    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/daily-plans`, {
      name: 'Plan Delete ' + Date.now(),
      date: new Date().toISOString(),
    });

    const res = await request.delete(`${API_BASE}/api/daily-plans/${id}`, { headers });
    expect(res.status()).toBe(204);
  });

  test('GET /api/daily-plans/:id/assignments returns list', async ({ request }) => {
    const listRes = await request.get(`${API_BASE}/api/daily-plans`, { headers });
    const body = await listRes.json();
    const plans = Array.isArray(body) ? body : body.items;
    if (plans.length === 0) return;

    const res = await request.get(`${API_BASE}/api/daily-plans/${plans[0].id}/assignments`, { headers });
    expect(res.status()).toBe(200);
  });

  test('GET /api/daily-plans/:id/materials returns list', async ({ request }) => {
    const listRes = await request.get(`${API_BASE}/api/daily-plans`, { headers });
    const body = await listRes.json();
    const plans = Array.isArray(body) ? body : body.items;
    if (plans.length === 0) return;

    const res = await request.get(`${API_BASE}/api/daily-plans/${plans[0].id}/materials`, { headers });
    expect(res.status()).toBe(200);
  });

  test('Technician cannot create daily plan', async ({ request }) => {
    const techTokens = await loginAPI(request, CREDENTIALS.technician);
    const techHeaders = await authHeaders(techTokens.access_token);
    await expectForbidden(request, techHeaders, `${API_BASE}/api/daily-plans`, 'POST');
  });
});

test.describe('Planning - Routes API', () => {
  let headers: Record<string, string>;

  test.beforeAll(async ({ request }) => {
    const tokens = await loginAPI(request);
    headers = await authHeaders(tokens.access_token);
  });

  test('GET /api/route-types returns list', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/route-types`, { headers });
    expect(res.status()).toBe(200);
  });

  test('GET /api/routes returns list', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/routes`, { headers });
    expect(res.status()).toBe(200);
  });

  test('POST /api/routes creates route', async ({ request }) => {
    const typesRes = await request.get(`${API_BASE}/api/route-types`, { headers });
    const typesBody = await typesRes.json();
    const types = Array.isArray(typesBody) ? typesBody : typesBody.items;
    if (types.length === 0) return;

    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/routes`, {
      date: new Date().toISOString(),
      status: 'scheduled',
      type: types[0].id,
    });
    expect(id).toBeTruthy();
    await request.delete(`${API_BASE}/api/routes/${id}`, { headers });
  });
});
