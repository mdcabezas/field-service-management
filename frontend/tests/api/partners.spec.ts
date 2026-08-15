import { test, expect, API_BASE, loginAPI, authHeaders, createAndReturnId, expectNotFound, expectForbidden, CREDENTIALS } from './api-helpers';

test.describe('Partners API', () => {
  let headers: Record<string, string>;

  test.beforeAll(async ({ request }) => {
    const tokens = await loginAPI(request);
    headers = await authHeaders(tokens.access_token);
  });

  test('GET /api/partners returns list', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/partners`, { headers });
    expect(res.status()).toBe(200);
  });

  test('GET /api/partners/:id returns partner', async ({ request }) => {
    const listRes = await request.get(`${API_BASE}/api/partners`, { headers });
    const body = await listRes.json();
    const partners = Array.isArray(body) ? body : body.items;
    if (partners.length === 0) return;

    const res = await request.get(`${API_BASE}/api/partners/${partners[0].id}`, { headers });
    expect(res.status()).toBe(200);
  });

  test('POST /api/partners creates partner', async ({ request }) => {
    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/partners`, {
      name: 'Test Partner ' + Date.now(),
      status: 'active',
    });
    expect(id).toBeTruthy();
    await request.delete(`${API_BASE}/api/partners/${id}`, { headers });
  });

  test('PUT /api/partners/:id updates partner', async ({ request }) => {
    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/partners`, {
      name: 'Partner Update Test ' + Date.now(),
      status: 'active',
    });

    const res = await request.put(`${API_BASE}/api/partners/${id}`, {
      headers,
      data: { name: 'Updated Partner', status: 'active' },
    });
    expect(res.status()).toBe(200);
    await request.delete(`${API_BASE}/api/partners/${id}`, { headers });
  });

  test('DELETE /api/partners/:id deletes partner', async ({ request }) => {
    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/partners`, {
      name: 'Partner Delete Test ' + Date.now(),
      status: 'active',
    });

    const res = await request.delete(`${API_BASE}/api/partners/${id}`, { headers });
    expect(res.status()).toBe(204);
    await expectNotFound(request, headers, `${API_BASE}/api/partners/${id}`);
  });

  test('GET /api/partners/:id/contacts returns list', async ({ request }) => {
    const listRes = await request.get(`${API_BASE}/api/partners`, { headers });
    const body = await listRes.json();
    const partners = Array.isArray(body) ? body : body.items;
    if (partners.length === 0) return;

    const res = await request.get(`${API_BASE}/api/partners/${partners[0].id}/contacts`, { headers });
    expect(res.status()).toBe(200);
  });

  test('GET /api/partners/:id/slas returns list', async ({ request }) => {
    const listRes = await request.get(`${API_BASE}/api/partners`, { headers });
    const body = await listRes.json();
    const partners = Array.isArray(body) ? body : body.items;
    if (partners.length === 0) return;

    const res = await request.get(`${API_BASE}/api/partners/${partners[0].id}/slas`, { headers });
    expect(res.status()).toBe(200);
  });

  test('GET /api/partners/:id/agreements returns list', async ({ request }) => {
    const listRes = await request.get(`${API_BASE}/api/partners`, { headers });
    const body = await listRes.json();
    const partners = Array.isArray(body) ? body : body.items;
    if (partners.length === 0) return;

    const res = await request.get(`${API_BASE}/api/partners/${partners[0].id}/agreements`, { headers });
    expect(res.status()).toBe(200);
  });

  test('Technician cannot create partner', async ({ request }) => {
    const techTokens = await loginAPI(request, CREDENTIALS.technician);
    const techHeaders = await authHeaders(techTokens.access_token);
    await expectForbidden(request, techHeaders, `${API_BASE}/api/partners`, 'POST');
  });
});
