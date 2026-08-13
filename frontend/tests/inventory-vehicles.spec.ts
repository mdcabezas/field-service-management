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

test.describe('Inventory/Vehicles Module', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows vehicles', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryVehicles);
    await verifyPageTitle(page, 'Vehículos');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
  });

  test('Create new vehicle flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryVehicles);
    await expectTableVisible(page);
    await clickCreateButton(page, ['+ Crear']);
    await verifyPageTitle(page, 'Crear Vehículo');

    await fillInput(page, 'license_plate', 'TEST123');
    await fillInput(page, 'name', 'Test Vehicle');
    await fillInput(page, 'brand', 'TestBrand');
    await fillSelect(page, 'status', 1);

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/inventory/vehicles');
    await waitForToast(page, 'Vehículo creado');
    await expectRowCount(page, 2);
  });

  test('View vehicle detail page', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryVehicles);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Editar Vehículo');
  });

  test('Edit vehicle flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryVehicles);
    await expectTableVisible(page);
    await clickFirstRow(page);

    await fillInput(page, 'brand', 'Updated Brand');
    await page.click('button[type="submit"]:has-text("Guardar")');
    await waitForRedirect(page, '**/inventory/vehicles');
    await waitForToast(page, 'Vehículo actualizado');
  });

  test('Delete vehicle flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryVehicles);
    await expectTableVisible(page);

    await clickCreateButton(page, ['+ Crear']);
    await fillInput(page, 'license_plate', 'DELETE456');
    await fillInput(page, 'name', 'Delete Vehicle');
    await fillInput(page, 'brand', 'DeleteBrand');
    await fillSelect(page, 'status', 1);
    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/inventory/vehicles');

    await page.locator('tbody tr', { hasText: 'DELETE456' }).last().click();
    confirmDialog(page);
    await page.click('button:has-text("Eliminar")');
    await waitForRedirect(page, '**/inventory/vehicles');
    await waitForToast(page, 'Vehículo eliminado');
  });
});