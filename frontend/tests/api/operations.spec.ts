import { test, expect, API_BASE, loginAPI, authHeaders, createAndReturnId, expectNotFound, expectForbidden, CREDENTIALS } from './api-helpers';

test.describe('Operations - Visits API', () => {
  let headers: Record<string, string>;

  test.beforeAll(async ({ request }) => {
    const tokens = await loginAPI(request);
    headers = await authHeaders(tokens.access_token);
  });

  test('GET /api/visits returns list', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/visits`, { headers });
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(Array.isArray(body) || body.items).toBeTruthy();
  });

  test('GET /api/visits/:id returns visit', async ({ request }) => {
    const listRes = await request.get(`${API_BASE}/api/visits`, { headers });
    const body = await listRes.json();
    const visits = Array.isArray(body) ? body : body.items;
    if (visits.length === 0) return;

    const res = await request.get(`${API_BASE}/api/visits/${visits[0].id}`, { headers });
    expect(res.status()).toBe(200);
  });

  test('GET /api/visits/:id/assignments returns list', async ({ request }) => {
    const listRes = await request.get(`${API_BASE}/api/visits`, { headers });
    const body = await listRes.json();
    const visits = Array.isArray(body) ? body : body.items;
    if (visits.length === 0) return;

    const res = await request.get(`${API_BASE}/api/visits/${visits[0].id}/assignments`, { headers });
    expect(res.status()).toBe(200);
  });

  test('GET /api/visits/:id/checklist-materials returns list', async ({ request }) => {
    const listRes = await request.get(`${API_BASE}/api/visits`, { headers });
    const body = await listRes.json();
    const visits = Array.isArray(body) ? body : body.items;
    if (visits.length === 0) return;

    const res = await request.get(`${API_BASE}/api/visits/${visits[0].id}/checklist-materials`, { headers });
    expect(res.status()).toBe(200);
  });

  test('GET /api/visits/:id/measurements returns list', async ({ request }) => {
    const listRes = await request.get(`${API_BASE}/api/visits`, { headers });
    const body = await listRes.json();
    const visits = Array.isArray(body) ? body : body.items;
    if (visits.length === 0) return;

    const res = await request.get(`${API_BASE}/api/visits/${visits[0].id}/measurements`, { headers });
    expect(res.status()).toBe(200);
  });

  test('GET /api/visits/:id/sla-trackings returns list', async ({ request }) => {
    const listRes = await request.get(`${API_BASE}/api/visits`, { headers });
    const body = await listRes.json();
    const visits = Array.isArray(body) ? body : body.items;
    if (visits.length === 0) return;

    const res = await request.get(`${API_BASE}/api/visits/${visits[0].id}/sla-trackings`, { headers });
    expect(res.status()).toBe(200);
  });

  test('GET /api/visits/:id/reports returns list', async ({ request }) => {
    const listRes = await request.get(`${API_BASE}/api/visits`, { headers });
    const body = await listRes.json();
    const visits = Array.isArray(body) ? body : body.items;
    if (visits.length === 0) return;

    const res = await request.get(`${API_BASE}/api/visits/${visits[0].id}/reports`, { headers });
    expect(res.status()).toBe(200);
  });

  test('GET /api/visits/:id/photos returns list', async ({ request }) => {
    const listRes = await request.get(`${API_BASE}/api/visits`, { headers });
    const body = await listRes.json();
    const visits = Array.isArray(body) ? body : body.items;
    if (visits.length === 0) return;

    const res = await request.get(`${API_BASE}/api/visits/${visits[0].id}/photos`, { headers });
    expect(res.status()).toBe(200);
  });

  test('GET /api/visits/:id/checkpoints returns list', async ({ request }) => {
    const listRes = await request.get(`${API_BASE}/api/visits`, { headers });
    const body = await listRes.json();
    const visits = Array.isArray(body) ? body : body.items;
    if (visits.length === 0) return;

    const res = await request.get(`${API_BASE}/api/visits/${visits[0].id}/checkpoints`, { headers });
    expect(res.status()).toBe(200);
  });
});

test.describe('Operations - SLA Trackings API', () => {
  let headers: Record<string, string>;

  test.beforeAll(async ({ request }) => {
    const tokens = await loginAPI(request);
    headers = await authHeaders(tokens.access_token);
  });

  test('GET /api/visit-sla-trackings returns list', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/visit-sla-trackings`, { headers });
    expect(res.status()).toBe(200);
  });
});

test.describe('Operations - Vehicle Assignments API', () => {
  let headers: Record<string, string>;

  test.beforeAll(async ({ request }) => {
    const tokens = await loginAPI(request);
    headers = await authHeaders(tokens.access_token);
  });

  test('GET /api/vehicle-assignments returns list', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/vehicle-assignments`, { headers });
    expect(res.status()).toBe(200);
  });
});

test.describe('Operations - Visit Reports API', () => {
  let headers: Record<string, string>;

  test.beforeAll(async ({ request }) => {
    const tokens = await loginAPI(request);
    headers = await authHeaders(tokens.access_token);
  });

  test('GET /api/visit-reports returns list with resolved names', async ({ request }) => {
    const res = await request.get(`${API_BASE}/api/visit-reports`, { headers });
    expect(res.status()).toBe(200);
    const body = await res.json();
    const reports = Array.isArray(body) ? body : body.items;
    if (reports.length > 0) {
      expect(reports[0]).toHaveProperty('report_template_name');
      expect(reports[0]).toHaveProperty('visit_label');
    }
  });

  test('POST /api/visit-reports creates report with source and recorded_at', async ({ request }) => {
    const visitsRes = await request.get(`${API_BASE}/api/visits`, { headers });
    const visitsBody = await visitsRes.json();
    const visits = Array.isArray(visitsBody) ? visitsBody : visitsBody.items;
    if (visits.length === 0) return;

    const templatesRes = await request.get(`${API_BASE}/api/report-templates`, { headers });
    const templatesBody = await templatesRes.json();
    const templates = Array.isArray(templatesBody) ? templatesBody : templatesBody.items;
    if (templates.length === 0) return;

    const createRes = await request.post(`${API_BASE}/api/visit-reports`, {
      headers,
      data: {
        visit_id: visits[0].id,
        report_template_id: templates[0].id,
        source: "paper",
        recorded_at: "2026-01-15T10:30:00Z",
      },
    });
    expect(createRes.status()).toBe(201);
    const created = await createRes.json();
    expect(created.source).toBe("paper");
    expect(created.visit_id).toBe(visits[0].id);
  });
});
