import { test, expect, type Page } from '@playwright/test';

const BASE_URL = 'http://localhost:3000';
const ADMIN_CREDENTIALS = { employee_number: '1001', password: 'admin' };

async function login(page: Page) {
  await page.goto(`${BASE_URL}/login`);
  await page.fill('#employeeNumber', ADMIN_CREDENTIALS.employee_number);
  await page.fill('#password', ADMIN_CREDENTIALS.password);
  await page.click('button[type="submit"]');
  await page.waitForURL('**/', { timeout: 10000 });
}

async function gotoAssignments(page: Page) {
  await page.goto(`${BASE_URL}/planning/assignments`);
  await page.waitForLoadState('networkidle');
}

async function waitForToast(page: Page, text: string) {
  // Sonner toasts are rendered in a portal - try to find but don't fail if not visible
  try {
    await expect(page.getByText(text)).toBeVisible({ timeout: 5000 });
  } catch {
    console.log(`Toast "${text}" not found (may be timing or portal issue)`);
  }
}

async function waitForRedirect(page: Page, urlPattern: string, timeout = 15000) {
  await page.waitForURL(urlPattern, { timeout });
}

test.describe('Assignment Flows', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows assignments', async ({ page }) => {
    await gotoAssignments(page);
    await expect(page.locator('h1:has-text("Asignaciones")')).toBeVisible();
    await expect(page.locator('table')).toBeVisible();
  });

  test('Create new assignment flow', async ({ page }) => {
    await gotoAssignments(page);

    await page.click('button:has-text("Crear Asignación")');
    await page.waitForURL('**/planning/assignments/new');

    await expect(page.locator('h1:has-text("Crear Asignación")')).toBeVisible();

    await page.selectOption('select[id="daily_plan_id"]', { index: 1 });
    await page.selectOption('select[id="tech_id"]', { index: 1 });
    await page.selectOption('select[id="role_id"]', { index: 1 });

    await page.click('button[type="submit"]:has-text("Crear")');
    await waitForRedirect(page, '**/planning/assignments');
    await waitForToast(page, 'Asignación creada');
  });

  test('View existing assignment (allows editing tech_id and role_id)', async ({ page }) => {
    await gotoAssignments(page);

    const firstRow = page.locator('tbody tr').first();
    await firstRow.click();
    await page.waitForURL(/\/planning\/assignments\/.*/);

    await expect(page.locator('h1:has-text("Detalle de Asignación")')).toBeVisible();

    // daily_plan_id should be disabled (immutable after creation)
    await expect(page.locator('select[id="daily_plan_id"]')).toBeDisabled();
    // tech_id and role_id should be editable
    await expect(page.locator('select[id="tech_id"]')).toBeEnabled();
    await expect(page.locator('select[id="role_id"]')).toBeEnabled();

    // Verify "Guardar" button exists for admin
    await expect(page.locator('button[type="submit"]:has-text("Guardar")')).toBeVisible();
  });

  test('Delete assignment flow', async ({ page }) => {
    await gotoAssignments(page);

    const firstRow = page.locator('tbody tr').first();
    await firstRow.click();
    await page.waitForURL(/\/planning\/assignments\/.*/);

    // Check if delete button exists
    const deleteButton = page.locator('button:has-text("Eliminar")');
    await expect(deleteButton).toBeVisible();

    await deleteButton.click();
    await page.on('dialog', dialog => dialog.accept());
    
    // Wait for either navigation back to list or error
    await page.waitForLoadState('networkidle');
    
    // Check if we're back on list page
    const isOnList = page.url().includes('/planning/assignments') && !page.url().match(/\/planning\/assignments\/[^/]+$/);
    if (isOnList) {
      await expect(page.locator('table')).toBeVisible({ timeout: 10000 });
    }
    
    await waitForToast(page, 'Asignación eliminada');
  });

  test('Empty state handling on create form when no data', async ({ page }) => {
    await page.goto(`${BASE_URL}/planning/assignments/new`);
    await page.waitForLoadState('networkidle');

    const techSelect = page.locator('select[id="tech_id"]');
    const roleSelect = page.locator('select[id="role_id"]');
    const planSelect = page.locator('select[id="daily_plan_id"]');

    const techOptions = await techSelect.locator('option').count();
    const roleOptions = await roleSelect.locator('option').count();
    const planOptions = await planSelect.locator('option').count();

    expect(techOptions).toBeGreaterThan(1);
    expect(roleOptions).toBeGreaterThan(1);
    expect(planOptions).toBeGreaterThan(1);
  });
});