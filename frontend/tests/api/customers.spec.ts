import { test, expect, API_BASE, loginAPI, authHeaders, createAndReturnId, expectNotFound, expectForbidden, CREDENTIALS } from './api-helpers';

test.describe('Customers API', () => {
  let headers: Record<string, string>;

  test.beforeAll(async ({ request }) => {
    const tokens = await loginAPI(request);
    headers = await authHeaders(tokens.access_token);
  });

  test('GET /api/customers returns list', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/customers`, { headers });
    expect(res.status()).toBe(200);
  });

  test('GET /api/customers/:id returns customer', async ({ request }) => {
    const listRes = await request.get(`${API_BASE}/api/customers`, { headers });
    const body = await listRes.json();
    const customers = Array.isArray(body) ? body : body.items;
    if (customers.length === 0) return;

    const res = await request.get(`${API_BASE}/api/customers/${customers[0].id}`, { headers });
    expect(res.status()).toBe(200);
  });

  test('POST /api/customers creates customer', async ({ request }) => {
    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/customers`, {
      name: 'Test Customer ' + Date.now(),
      status: 'active',
    });
    expect(id).toBeTruthy();
    await request.delete(`${API_BASE}/api/customers/${id}`, { headers });
  });

  test('PUT /api/customers/:id updates customer', async ({ request }) => {
    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/customers`, {
      name: 'Customer Update Test ' + Date.now(),
      status: 'active',
    });

    const res = await request.put(`${API_BASE}/api/customers/${id}`, {
      headers,
      data: { name: 'Updated Customer', status: 'active' },
    });
    expect(res.status()).toBe(200);
    await request.delete(`${API_BASE}/api/customers/${id}`, { headers });
  });

  test('DELETE /api/customers/:id deletes customer', async ({ request }) => {
    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/customers`, {
      name: 'Customer Delete Test ' + Date.now(),
      status: 'active',
    });

    const res = await request.delete(`${API_BASE}/api/customers/${id}`, { headers });
    expect(res.status()).toBe(204);
    await expectNotFound(request, headers, `${API_BASE}/api/customers/${id}`);
  });

  test('GET /api/customers/:id/addresses returns list', async ({ request }) => {
    const listRes = await request.get(`${API_BASE}/api/customers`, { headers });
    const body = await listRes.json();
    const customers = Array.isArray(body) ? body : body.items;
    if (customers.length === 0) return;

    const res = await request.get(`${API_BASE}/api/customers/${customers[0].id}/addresses`, { headers });
    expect(res.status()).toBe(200);
  });

  test('POST /api/customer-addresses creates address', async ({ request }) => {
    const custList = await request.get(`${API_BASE}/api/customers`, { headers });
    const custBody = await custList.json();
    const customers = Array.isArray(custBody) ? custBody : custBody.items;
    if (customers.length === 0) return;

    const addrList = await request.get(`${API_BASE}/api/addresses`, { headers });
    const addrBody = await addrList.json();
    const addresses = Array.isArray(addrBody) ? addrBody : addrBody.items;
    if (addresses.length === 0) return;

    const id = await createAndReturnId(request, headers, 'POST', `${API_BASE}/api/customer-addresses`, {
      customer_id: customers[0].id,
      address_id: addresses[0].id,
      type: 'commercial',
    });
    expect(id).toBeTruthy();
    await request.delete(`${API_BASE}/api/customer-addresses/${id}`, { headers });
  });

  test('GET /api/technicians returns list', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/technicians`, { headers });
    expect(res.status()).toBe(200);
  });

  test('Technician cannot create customer', async ({ request }) => {
    const techTokens = await loginAPI(request, CREDENTIALS.technician);
    const techHeaders = await authHeaders(techTokens.access_token);
    await expectForbidden(request, techHeaders, `${API_BASE}/api/customers`, 'POST');
  });
});
