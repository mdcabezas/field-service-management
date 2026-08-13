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

test.describe('Geocoding Module', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows addresses', async ({ page }) => {
    await gotoPage(page, MODULES.geocoding);
    await verifyPageTitle(page, 'Direcciones');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
  });

  test('Create new address flow', async ({ page }) => {
    await gotoPage(page, MODULES.geocoding);
    await clickCreateButton(page, '+ Crear Dirección');
    await verifyPageTitle(page, 'Crear Dirección');

    await fillInput(page, 'street', 'Test Street 123');
    await fillInput(page, 'number', '456');
    await fillInput(page, 'neighborhood', 'Test Neighborhood');
    await fillInput(page, 'city', 'Test City');
    await fillInput(page, 'region', 'Test Region');
    await fillInput(page, 'apartment', 'Apt 101');
    await fillInput(page, 'postal_code', '1234567');

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/geocoding');
    await waitForToast(page, 'Dirección creada');
    await expectRowCount(page, 2);
  });

  test('View address detail page', async ({ page }) => {
    await gotoPage(page, MODULES.geocoding);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Detalle de Dirección');
  });

  test('Edit address flow', async ({ page }) => {
    await gotoPage(page, MODULES.geocoding);
    await clickFirstRow(page);

    await fillInput(page, 'street', 'Updated Street 789');
    await page.click('button[type="submit"]:has-text("Guardar")');
    await waitForRedirect(page, '**/geocoding');
    await waitForToast(page, 'Dirección actualizada');
  });

  test('Delete address flow', async ({ page }) => {
    await gotoPage(page, MODULES.geocoding);

    await clickCreateButton(page, '+ Crear Dirección');
    await fillInput(page, 'street', 'Delete Street');
    await fillInput(page, 'number', '999');
    await fillInput(page, 'neighborhood', 'Delete Neighborhood');
    await fillInput(page, 'city', 'Delete City');
    await fillInput(page, 'region', 'Delete Region');
    await fillInput(page, 'postal_code', '9999999');
    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/geocoding');

    await clickFirstRow(page);
    await page.click('button:has-text("Eliminar")');
    await waitForRedirect(page, '**/geocoding');
    await waitForToast(page, 'Dirección eliminada');
  });
});