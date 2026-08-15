import { test as base, expect, type APIRequestContext } from '@playwright/test';

export const API_BASE = 'http://localhost:8088';

export const CREDENTIALS = {
  admin: { employee_number: '1001', password: 'admin' },
  operator: { employee_number: '1002', password: 'operador' },
  technician: { employee_number: '1003', password: 'tecnico' },
  invalid: { employee_number: '9999', password: 'wrong' },
};

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
}

export async function loginAPI(
  request: APIRequestContext,
  creds: { employee_number: string; password: string } = CREDENTIALS.admin
): Promise<AuthTokens> {
  const res = await request.post(`${API_BASE}/auth/login`, { data: creds });
  expect(res.status()).toBe(200);
  return res.json();
}

export async function authHeaders(token: string): Promise<Record<string, string>> {
  return { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' };
}

export function uuid(): string {
  return '00000000-0000-0000-0000-' + Date.now().toString(16).padStart(12, '0');
}

export async function createAndReturnId(
  request: APIRequestContext,
  headers: Record<string, string>,
  method: 'POST',
  url: string,
  data: Record<string, unknown>
): Promise<string> {
  const res = await request.post(url, { headers, data });
  expect(res.status()).toBe(201);
  const body = await res.json();
  return body.id;
}

export async function expectNotFound(
  request: APIRequestContext,
  headers: Record<string, string>,
  url: string
): Promise<void> {
  const res = await request.get(url, { headers });
  expect(res.status()).toBe(404);
}

export async function expectForbidden(
  request: APIRequestContext,
  headers: Record<string, string>,
  url: string,
  method: string = 'GET'
): Promise<void> {
  let res;
  switch (method) {
    case 'POST':
      res = await request.post(url, { headers, data: {} });
      break;
    case 'PUT':
      res = await request.put(url, { headers, data: {} });
      break;
    case 'DELETE':
      res = await request.delete(url, { headers });
      break;
    default:
      res = await request.get(url, { headers });
  }
  expect(res.status()).toBe(403);
}

export const test = base.extend<{ authedRequest: APIRequestContext; tokens: AuthTokens }>({
  tokens: async ({}, use) => {
    const context = await base.info.request.newContext();
    const tokens = await loginAPI(context);
    await use(tokens);
    await context.dispose();
  },
  authedRequest: async ({}, use) => {
    const context = await base.info.request.newContext();
    const tokens = await loginAPI(context);
    const authed = await context.newContext({
      extraHTTPHeaders: {
        Authorization: `Bearer ${tokens.access_token}`,
        'Content-Type': 'application/json',
      },
    });
    await use(authed);
    await authed.dispose();
    await context.dispose();
  },
});

export { expect };
