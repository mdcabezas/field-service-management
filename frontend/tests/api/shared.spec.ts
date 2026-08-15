import { test, expect, API_BASE, loginAPI, authHeaders, createAndReturnId, expectForbidden, CREDENTIALS } from './api-helpers';

test.describe('Shared - Report Templates API', () => {
  let headers: Record<string, string>;

  test.beforeAll(async ({ request }) => {
    const tokens = await loginAPI(request);
    headers = await authHeaders(tokens.access_token);
  });

  test('GET /api/report-templates returns list', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/report-templates`, { headers });
    expect(res.status()).toBe(200);
  });

  test('GET /api/report-templates/:id returns template', async ({ request }) => {
    const listRes = await request.get(`${API_BASE}/api/report-templates`, { headers });
    const body = await listRes.json();
    const templates = Array.isArray(body) ? body : body.items;
    if (templates.length === 0) return;

    const res = await request.get(`${API_BASE}/api/report-templates/${templates[0].id}`, { headers });
    expect(res.status()).toBe(200);
  });

  test('Technician cannot create report template', async ({ request }) => {
    const techTokens = await loginAPI(request, CREDENTIALS.technician);
    const techHeaders = await authHeaders(techTokens.access_token);
    await expectForbidden(request, techHeaders, `${API_BASE}/api/report-templates`, 'POST');
  });
});
