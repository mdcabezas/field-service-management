import { test, expect } from '@playwright/test';
import {
  login,
  gotoPage,
  waitForToast,
  verifyPageTitle,
  waitForRedirect,
  expectTableVisible,
  expectRowCount,
  clickFirstRow,
  clickCreateButton,
  fillInput,
  fillSelect,
  submitForm,
  confirmDialog,
  MODULES,
} from './test-helpers';

test.describe('Core/Tech Roles Module', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows tech roles', async ({ page }) => {
    await gotoPage(page, MODULES.coreTechRoles);
    await verifyPageTitle(page, 'Roles Técnicos');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
  });

  test('Search filters tech roles', async ({ page }) => {
    await gotoPage(page, MODULES.coreTechRoles);
    await page.fill('input[placeholder*="Buscar"]', 'admin');
    await expect(page.locator('tbody tr')).toHaveCount(1);
  });

  test('Create new tech role flow', async ({ page }) => {
    await page.goto('http://localhost:3000/core/tech-roles/new');
    await page.waitForLoadState('networkidle');
    await verifyPageTitle(page, 'Crear Rol Técnico');

    await fillInput(page, 'code', 'TEST_ROLE');
    await fillInput(page, 'name', 'Test Role');
    await page.check('#active');
    
    // Debug: check form state
    console.log('URL before submit:', page.url());
    const form = page.locator('form');
    await form.evaluate(f => console.log('Form action:', f.action, 'method:', f.method));
    
    // Check for validation errors before submit
    const errorTexts = await page.locator('.text-red-600').allTextContents();
    console.log('Errors before submit:', errorTexts);
    
    await page.click('button[type="submit"]:has-text("Crear")');
    
    // Wait a bit and check URL
    await page.waitForTimeout(3000);
    console.log('URL after submit:', page.url());
    
    // Check for validation errors after submit
    const errorTextsAfter = await page.locator('.text-red-600').allTextContents();
    console.log('Errors after submit:', errorTextsAfter);
    
    // Check for toasts
    const toasts = await page.locator('[role="status"], [role="alert"]').allTextContents();
    console.log('Toasts:', toasts);
    
    // Wait for redirect
    await page.waitForURL('**/core/tech-roles', { timeout: 30000 });
    await waitForToast(page, 'Rol técnico creado');
    await expectRowCount(page, 2);
  });

  test('Create tech role with duplicate code shows error', async ({ page }) => {
    await gotoPage(page, MODULES.coreTechRoles);
    await clickCreateButton(page, '+ Crear Rol');

    await fillInput(page, 'code', 'ELECTRICAL');
    await fillInput(page, 'name', 'Duplicate Electrical');
    await page.click('button[type="submit"]:has-text("Crear")');
    await waitForToast(page, 'Error al crear rol técnico');
  });

  test('View tech role detail page', async ({ page }) => {
    await gotoPage(page, MODULES.coreTechRoles);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Detalle de Rol Técnico');
  });

  test('Edit tech role flow', async ({ page }) => {
    await gotoPage(page, MODULES.coreTechRoles);
    await clickFirstRow(page);

    await fillInput(page, 'name', 'Updated Electrical');
    await page.click('button[type="submit"]:has-text("Guardar")');
    await waitForRedirect(page, '**/core/tech-roles');
    await waitForToast(page, 'Rol técnico actualizado');
  });

  test('Delete tech role flow', async ({ page }) => {
    await gotoPage(page, MODULES.coreTechRoles);

    await clickCreateButton(page, '+ Crear Rol');
    await fillInput(page, 'code', 'TO_DELETE');
    await fillInput(page, 'name', 'To Delete Role');
    await page.click('button[type="submit"]:has-text("Crear")');
    await waitForRedirect(page, '**/core/tech-roles');

    await clickFirstRow(page);
    await page.click('button:has-text("Eliminar")');
    await waitForRedirect(page, '**/core/tech-roles');
    await waitForToast(page, 'Rol técnico eliminado');
  });

  test('Direct navigation to create page', async ({ page }) => {
    await page.goto('http://localhost:3000/core/tech-roles/new');
    await page.waitForLoadState('networkidle');
    
    const h1s = await page.locator('h1').allTextContents();
    console.log('H1 elements:', h1s);
    
    await verifyPageTitle(page, 'Crear Rol Técnico');
  });

  test('Debug form submission with network', async ({ page }) => {
    await page.goto('http://localhost:3000/core/tech-roles/new');
    await page.waitForLoadState('networkidle');
    
    // Fill form
    await page.fill('#code', 'TEST_ROLE');
    await page.fill('#name', 'Test Role');
    await page.check('#active');
    
    // Listen for network requests
    page.on('request', request => {
      if (request.url().includes('/api/tech-roles')) {
        console.log('API Request:', request.method(), request.url(), request.postData());
      }
    });
    
    page.on('response', response => {
      if (response.url().includes('/api/tech-roles')) {
        console.log('API Response:', response.status(), response.url());
        response.text().then(t => console.log('Response body:', t));
      }
    });
    
    await page.click('button[type="submit"]:has-text("Crear")');
    await page.waitForTimeout(5000);
    
    console.log('Final URL:', page.url());
  });
});