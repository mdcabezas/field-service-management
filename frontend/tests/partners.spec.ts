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

test.describe('Partners Module', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows partners', async ({ page }) => {
    await gotoPage(page, MODULES.partners);
    await verifyPageTitle(page, 'Socios');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
  });

  test('Create new partner flow', async ({ page }) => {
    await gotoPage(page, MODULES.partners);
    await clickCreateButton(page, '+ Crear Socio');
    await verifyPageTitle(page, 'Crear Socio');

    await fillInput(page, 'name', 'Test Partner');
    await fillInput(page, 'tax_id', '87654321-0');
    await fillSelect(page, 'status', 1);

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/partners');
    await waitForToast(page, 'Socio creado');
    await expectRowCount(page, 2);
  });

  test('View partner detail page', async ({ page }) => {
    await gotoPage(page, MODULES.partners);
    await clickFirstRow(page);
    await expect(page.locator('button[role="tab"]:has-text("SLAs")')).toBeVisible();
    await expect(page.locator('button:has-text("Eliminar")')).toBeVisible();
  });

  test('Edit partner flow', async ({ page }) => {
    await gotoPage(page, MODULES.partners);
    await clickFirstRow(page);

    await fillInput(page, 'name', 'Updated Partner Name');
    await page.click('button[type="submit"]:has-text("Actualizar")');
    await waitForRedirect(page, '**/partners');
    await waitForToast(page, 'Socio actualizado');
  });

  test('Delete partner flow', async ({ page }) => {
    await gotoPage(page, MODULES.partners);

    await clickCreateButton(page, '+ Crear Socio');
    await fillInput(page, 'name', 'To Delete Partner');
    await fillInput(page, 'tax_id', '11111111-1');
    await fillSelect(page, 'status', 1);
    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/partners');

    await clickFirstRow(page);
    await confirmDialog(page);
    await page.click('button:has-text("Eliminar")');
    await waitForRedirect(page, '**/partners');
    await waitForToast(page, 'Socio eliminado');
  });
});