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

test.describe('Customers Module', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows customers', async ({ page }) => {
    await gotoPage(page, MODULES.customers);
    await verifyPageTitle(page, 'Clientes');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
  });

  test('Create new customer flow', async ({ page }) => {
    await gotoPage(page, MODULES.customers);
    await clickCreateButton(page, '+ Crear Cliente');
    await verifyPageTitle(page, 'Crear Cliente');

    await fillInput(page, 'name', 'Test Customer');
    await fillInput(page, 'tax_id', '12345678-9');
    await fillInput(page, 'email', 'test@customer.com');
    await fillInput(page, 'phone', '+56912345678');

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/customers');
    await waitForToast(page, 'Cliente creado');
    await expectRowCount(page, 2);
  });

  test('View customer detail page', async ({ page }) => {
    await gotoPage(page, MODULES.customers);
    await clickFirstRow(page);
    await expect(page.locator('button[role="tab"]:has-text("Direcciones")')).toBeVisible();
    await expect(page.locator('button:has-text("Eliminar")')).toBeVisible();
  });

  test('Edit customer flow', async ({ page }) => {
    await gotoPage(page, MODULES.customers);
    await clickFirstRow(page);

    await fillInput(page, 'name', 'Updated Customer Name');
    await page.click('button[type="submit"]:has-text("Actualizar")');
    await waitForRedirect(page, '**/customers');
    await waitForToast(page, 'Cliente actualizado');
  });

  test('Delete customer flow', async ({ page }) => {
    await gotoPage(page, MODULES.customers);

    await clickCreateButton(page, '+ Crear Cliente');
    await fillInput(page, 'name', 'To Delete Customer');
    await fillInput(page, 'tax_id', '99999999-9');
    await fillInput(page, 'email', 'delete@test.com');
    await fillInput(page, 'phone', '+56911111111');

    // Capture the customer ID from the creation response
    let customerId: string;
    const responsePromise = page.waitForResponse(
      (resp) => resp.url().includes('/api/customers') && resp.request().method() === 'POST'
    );
    await submitForm(page, 'Crear');
    const response = await responsePromise;
    const responseBody = await response.json();
    customerId = responseBody.id;

    await waitForRedirect(page, '**/customers');
    await waitForToast(page, 'Cliente creado');

    // Navigate directly to the created customer's detail page
    await page.goto(`http://localhost:3000/customers/${customerId}`);
    await expect(page.locator('button:has-text("Eliminar")')).toBeVisible();

    await confirmDialog(page);
    await page.click('button:has-text("Eliminar")');
    await waitForRedirect(page, '**/customers');
    await waitForToast(page, 'Cliente eliminado');
  });
});