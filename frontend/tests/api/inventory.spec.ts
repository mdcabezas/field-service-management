import { test, expect, API_BASE, loginAPI, authHeaders, createAndReturnId, expectNotFound, expectForbidden, CREDENTIALS } from './api-helpers';

test.describe('Inventory - Materials API', () => {
  let headers: Record<string, string>;

  test.beforeAll(async ({ request }) => {
    const tokens = await loginAPI(request);
    headers = await authHeaders(tokens.access_token);
  });

  test('GET /api/materials returns list', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/materials`, { headers });
    expect(res.status()).toBe(200);
  });

  test('GET /api/materials/:id returns material', async ({ request }) => {
    const listRes = await request.get(`${API_BASE}/api/materials`, { headers });
    const body = await listRes.json();
    const items = Array.isArray(body) ? body : body.items;
    if (items.length === 0) return;

    const res = await request.get(`${API_BASE}/api/materials/${items[0].id}`, { headers });
    expect(res.status()).toBe(200);
  });

  test('POST /api/materials creates material', async ({ request }) => {
    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/materials`, {
      name: 'Test Material ' + Date.now(),
      unit: 'unit',
    });
    expect(id).toBeTruthy();
    await request.delete(`${API_BASE}/api/materials/${id}`, { headers });
  });

  test('PUT /api/materials/:id updates material', async ({ request }) => {
    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/materials`, {
      name: 'Material Update ' + Date.now(),
      unit: 'unit',
    });

    const res = await request.put(`${API_BASE}/api/materials/${id}`, {
      headers,
      data: { name: 'Updated Material', unit: 'meter' },
    });
    expect(res.status()).toBe(200);
    await request.delete(`${API_BASE}/api/materials/${id}`, { headers });
  });

  test('DELETE /api/materials/:id deletes material', async ({ request }) => {
    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/materials`, {
      name: 'Material Delete ' + Date.now(),
      unit: 'unit',
    });

    const res = await request.delete(`${API_BASE}/api/materials/${id}`, { headers });
    expect(res.status()).toBe(204);
  });
});

test.describe('Inventory - Tools API', () => {
  let headers: Record<string, string>;

  test.beforeAll(async ({ request }) => {
    const tokens = await loginAPI(request);
    headers = await authHeaders(tokens.access_token);
  });

  test('GET /api/tools returns list', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/tools`, { headers });
    expect(res.status()).toBe(200);
  });

  test('POST /api/tools creates tool', async ({ request }) => {
    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/tools`, {
      name: 'Test Tool ' + Date.now(),
      status: 'available',
    });
    expect(id).toBeTruthy();
    await request.delete(`${API_BASE}/api/tools/${id}`, { headers });
  });
});

test.describe('Inventory - EPP API', () => {
  let headers: Record<string, string>;

  test.beforeAll(async ({ request }) => {
    const tokens = await loginAPI(request);
    headers = await authHeaders(tokens.access_token);
  });

  test('GET /api/epp-items returns list', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/epp-items`, { headers });
    expect(res.status()).toBe(200);
  });

  test('POST /api/epp-items creates epp', async ({ request }) => {
    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/epp-items`, {
      name: 'Test EPP ' + Date.now(),
      type: 'eyes',
      lifecycle: 'reusable',
    });
    expect(id).toBeTruthy();
    await request.delete(`${API_BASE}/api/epp-items/${id}`, { headers });
  });
});

test.describe('Inventory - Vehicles API', () => {
  let headers: Record<string, string>;

  test.beforeAll(async ({ request }) => {
    const tokens = await loginAPI(request);
    headers = await authHeaders(tokens.access_token);
  });

  test('GET /api/vehicles returns list', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/vehicles`, { headers });
    expect(res.status()).toBe(200);
  });

  test('POST /api/vehicles creates vehicle', async ({ request }) => {
    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/vehicles`, {
      name: 'Test Vehicle ' + Date.now(),
      status: 'available',
    });
    expect(id).toBeTruthy();
    await request.delete(`${API_BASE}/api/vehicles/${id}`, { headers });
  });
});

test.describe('Inventory - Rentals API', () => {
  let headers: Record<string, string>;

  test.beforeAll(async ({ request }) => {
    const tokens = await loginAPI(request);
    headers = await authHeaders(tokens.access_token);
  });

  test('GET /api/rentals returns list', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/rentals`, { headers });
    expect(res.status()).toBe(200);
  });
});

test.describe('Inventory - Cost Rates API', () => {
  let headers: Record<string, string>;

  test.beforeAll(async ({ request }) => {
    const tokens = await loginAPI(request);
    headers = await authHeaders(tokens.access_token);
  });

  test('GET /api/cost-rates returns list', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/cost-rates`, { headers });
    expect(res.status()).toBe(200);
  });
});

test.describe('Inventory - Maintenance Records API', () => {
  let headers: Record<string, string>;

  test.beforeAll(async ({ request }) => {
    const tokens = await loginAPI(request);
    headers = await authHeaders(tokens.access_token);
  });

  test('GET /api/maintenance-records returns list', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/maintenance-records`, { headers });
    expect(res.status()).toBe(200);
  });
});

test.describe('Inventory - Maintenance Schedules API', () => {
  let headers: Record<string, string>;

  test.beforeAll(async ({ request }) => {
    const tokens = await loginAPI(request);
    headers = await authHeaders(tokens.access_token);
  });

  test('GET /api/maintenance-schedules returns list', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/maintenance-schedules`, { headers });
    expect(res.status()).toBe(200);
  });
});

test.describe('Inventory - Checklist Templates API', () => {
  let headers: Record<string, string>;

  test.beforeAll(async ({ request }) => {
    const tokens = await loginAPI(request);
    headers = await authHeaders(tokens.access_token);
  });

  test('GET /api/checklist-templates returns list', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/checklist-templates`, { headers });
    expect(res.status()).toBe(200);
  });
});
