import { test, expect, API_BASE, loginAPI, authHeaders, createAndReturnId, expectNotFound, expectForbidden, CREDENTIALS } from './api-helpers';

test.describe('Geocoding API', () => {
  let headers: Record<string, string>;

  test.beforeAll(async ({ request }) => {
    const tokens = await loginAPI(request);
    headers = await authHeaders(tokens.access_token);
  });

  test('GET /api/addresses returns list', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/addresses`, { headers });
    expect(res.status()).toBe(200);
  });

  test('GET /api/addresses/:id returns address', async ({ request }) => {
    const listRes = await request.get(`${API_BASE}/api/addresses`, { headers });
    const body = await listRes.json();
    const addresses = Array.isArray(body) ? body : body.items;
    if (addresses.length === 0) return;

    const res = await request.get(`${API_BASE}/api/addresses/${addresses[0].id}`, { headers });
    expect(res.status()).toBe(200);
  });

  test('POST /api/addresses creates address', async ({ request }) => {
    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/addresses`, {
      street: 'Test Street ' + Date.now(),
    });
    expect(id).toBeTruthy();
    await request.delete(`${API_BASE}/api/addresses/${id}`, { headers });
  });

  test('PUT /api/addresses/:id updates address', async ({ request }) => {
    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/addresses`, {
      street: 'Address Update ' + Date.now(),
    });

    const res = await request.put(`${API_BASE}/api/addresses/${id}`, {
      headers,
      data: { street: 'Updated Street' },
    });
    expect(res.status()).toBe(200);
    await request.delete(`${API_BASE}/api/addresses/${id}`, { headers });
  });

  test('DELETE /api/addresses/:id deletes address', async ({ request }) => {
    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/addresses`, {
      street: 'Address Delete ' + Date.now(),
    });

    const res = await request.delete(`${API_BASE}/api/addresses/${id}`, { headers });
    expect(res.status()).toBe(204);
  });

  test('Technician cannot create address', async ({ request }) => {
    const techTokens = await loginAPI(request, CREDENTIALS.technician);
    const techHeaders = await authHeaders(techTokens.access_token);
    await expectForbidden(request, techHeaders, `${API_BASE}/api/addresses`, 'POST');
  });
});
