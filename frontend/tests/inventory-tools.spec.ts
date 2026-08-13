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

test.describe('Inventory/Tools Module', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows tools', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryTools);
    await verifyPageTitle(page, 'Herramientas');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
  });

  test('Create new tool flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryTools);
    await expectTableVisible(page);
    await page.locator('button:has-text("+ Crear")').waitFor({ timeout: 10000 });
    await clickCreateButton(page, ['+ Crear']);
    await verifyPageTitle(page, 'Crear Herramienta');

    await fillInput(page, 'code', 'TOOL_TEST');
    await fillInput(page, 'name', 'Test Tool');
    await fillSelect(page, 'status', 1);

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/inventory/tools');
    await waitForToast(page, 'Herramienta creada');
    await expectRowCount(page, 2);
  });

  test('View tool detail page', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryTools);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForURL('**/inventory/tools/**', { timeout: 10000 });
    await verifyPageTitle(page, 'Herramienta');
  });

  test('Edit tool flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryTools);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForURL('**/inventory/tools/**', { timeout: 10000 });
    await fillInput(page, 'name', 'Updated Tool Name');
    await page.click('button[type="submit"]:has-text("Guardar")');
    await waitForRedirect(page, '**/inventory/tools');
    await waitForToast(page, 'Herramienta actualizada');
  });

  test('Delete tool flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryTools);
    await expectTableVisible(page);

    await page.locator('button:has-text("+ Crear")').waitFor({ timeout: 10000 });
    await clickCreateButton(page, ['+ Crear']);
    await fillInput(page, 'code', 'TOOL_DELETE');
    await fillInput(page, 'name', 'To Delete Tool');
    await fillSelect(page, 'status', 1);
    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/inventory/tools');

    const row = page.locator('tbody tr', { hasText: 'To Delete Tool' }).last();
    await row.waitFor({ timeout: 10000 });
    await row.click();
    confirmDialog(page);
    await page.click('button:has-text("Eliminar")');
    await waitForRedirect(page, '**/inventory/tools');
    await waitForToast(page, 'Herramienta eliminada');
  });
});