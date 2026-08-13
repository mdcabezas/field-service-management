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

test.describe('Inventory/Rentals Module', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows rentals', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryRentals);
    await verifyPageTitle(page, 'Arriendos');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
  });

  test('Create new rental flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryRentals);
    await expectTableVisible(page);
    await clickCreateButton(page, ['+ Crear']);
    await verifyPageTitle(page, 'Crear Arriendo');

    await fillSelect(page, 'type', 1);
    await fillInput(page, 'supplier', 'Test Supplier');
    await fillInput(page, 'item_description', 'Test rental description');
    await fillInput(page, 'daily_cost', '100');

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/inventory/rentals');
    await waitForToast(page, 'Alquiler creado');
    await expectRowCount(page, 2);
  });

  test('View rental detail page', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryRentals);
    await expectTableVisible(page);
    const rows = await page.locator('tbody tr').count();
    if (rows === 0) {
      await clickCreateButton(page, ['+ Crear']);
      await fillSelect(page, 'type', 1);
      await fillInput(page, 'supplier', 'View Rental');
      await fillInput(page, 'item_description', 'View Rental Item');
      await fillInput(page, 'daily_cost', '75');
      await submitForm(page, 'Crear');
      await waitForRedirect(page, '**/inventory/rentals');
      await gotoPage(page, MODULES.inventoryRentals);
      await expectTableVisible(page);
    }
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Editar Arriendo');
  });

  test('Edit rental flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryRentals);
    await expectTableVisible(page);
    await clickFirstRow(page);

    await fillInput(page, 'supplier', 'Updated Supplier');
    await page.click('button[type="submit"]:has-text("Guardar")');
    await waitForRedirect(page, '**/inventory/rentals');
    await waitForToast(page, 'Alquiler actualizado');
  });

  test('Delete rental flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryRentals);
    await expectTableVisible(page);

    await clickCreateButton(page, ['+ Crear']);
    await fillSelect(page, 'type', 1);
    await fillInput(page, 'supplier', 'AAAAA-Delete Supplier');
    await fillInput(page, 'item_description', 'AAAAA-Delete Rental');
    await fillInput(page, 'daily_cost', '50');
    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/inventory/rentals');
    await expectTableVisible(page);

    await clickFirstRow(page);
    confirmDialog(page);
    await page.click('button:has-text("Eliminar")');
    await waitForRedirect(page, '**/inventory/rentals');
    await waitForToast(page, 'Alquiler eliminado');
  });
});