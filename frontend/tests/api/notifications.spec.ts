import { test, expect, API_BASE, loginAPI, authHeaders, createAndReturnId, expectForbidden, CREDENTIALS } from './api-helpers';

test.describe('Notifications - Templates API', () => {
  let headers: Record<string, string>;

  test.beforeAll(async ({ request }) => {
    const tokens = await loginAPI(request);
    headers = await authHeaders(tokens.access_token);
  });

  test('GET /api/notification-templates returns list', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/notification-templates`, { headers });
    expect(res.status()).toBe(200);
  });

  test('GET /api/notification-templates/:id returns template', async ({ request }) => {
    const listRes = await request.get(`${API_BASE}/api/notification-templates`, { headers });
    const body = await listRes.json();
    const templates = Array.isArray(body) ? body : body.items;
    if (templates.length === 0) return;

    const res = await request.get(`${API_BASE}/api/notification-templates/${templates[0].id}`, { headers });
    expect(res.status()).toBe(200);
  });
});

test.describe('Notifications - List API', () => {
  let headers: Record<string, string>;

  test.beforeAll(async ({ request }) => {
    const tokens = await loginAPI(request);
    headers = await authHeaders(tokens.access_token);
  });

  test('GET /api/notifications returns list', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/notifications`, { headers });
    expect(res.status()).toBe(200);
  });

  test('Technician cannot create notification template', async ({ request }) => {
    const techTokens = await loginAPI(request, CREDENTIALS.technician);
    const techHeaders = await authHeaders(techTokens.access_token);
    await expectForbidden(request, techHeaders, `${API_BASE}/api/notification-templates`, 'POST');
  });
});
