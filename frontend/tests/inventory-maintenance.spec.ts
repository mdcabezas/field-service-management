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

test.describe('Inventory/Maintenance Records Module', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows maintenance records', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryMaintenanceRecords);
    await verifyPageTitle(page, 'Registros de Mantención');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
  });

  test('Create new maintenance record flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryMaintenanceRecords);
    await expectTableVisible(page);
    await clickCreateButton(page, ['+ Crear']);
    await verifyPageTitle(page, 'Crear Registro');

    const today = new Date().toISOString().split('T')[0];
    await fillInput(page, 'reference_id', '00000000-0000-0000-0000-000000000001');
    await fillSelect(page, 'type', 1);
    await fillInput(page, 'date', today);
    await fillInput(page, 'description', 'Test maintenance');
    await fillInput(page, 'cost', '150.00');

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/inventory/maintenance-records');
    await waitForToast(page, 'Registro de mantenimiento creado');
    await expectRowCount(page, 2);
  });

  test('View maintenance record detail page', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryMaintenanceRecords);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Editar Registro');
  });

  test('Edit maintenance record flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryMaintenanceRecords);
    await expectTableVisible(page);

    await clickCreateButton(page, ['+ Crear']);
    await page.waitForURL('**/inventory/maintenance-records/new', { timeout: 10000 });

    const today = new Date().toISOString().split('T')[0];
    await fillInput(page, 'reference_id', '00000000-0000-0000-0000-000000000001');
    await fillSelect(page, 'type', 1);
    await fillInput(page, 'date', today);
    await fillInput(page, 'description', 'Record to edit');
    await fillInput(page, 'cost', '200.00');
    await submitForm(page, 'Crear');
    await page.waitForURL((url) => url.pathname === '/inventory/maintenance-records', { timeout: 15000 });
    await expectTableVisible(page);

    await clickFirstRow(page);
    await page.waitForURL((url) => url.pathname.includes('/inventory/maintenance-records/') && !url.pathname.endsWith('/new'), { timeout: 10000 });
    await verifyPageTitle(page, 'Registro');

    await fillInput(page, 'description', 'Updated maintenance description');
    await page.click('button[type="submit"]:has-text("Guardar")');
    await page.waitForURL((url) => url.pathname === '/inventory/maintenance-records', { timeout: 15000 });
    await waitForToast(page, 'Registro actualizado');
  });

  test('Delete maintenance record flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryMaintenanceRecords);
    await expectTableVisible(page);

    await clickCreateButton(page, ['+ Crear']);
    await page.waitForURL('**/inventory/maintenance-records/new', { timeout: 10000 });

    const today = new Date().toISOString().split('T')[0];
    await fillInput(page, 'reference_id', '00000000-0000-0000-0000-000000000002');
    await fillSelect(page, 'type', 1);
    await fillInput(page, 'date', today);
    await fillInput(page, 'description', 'To delete maintenance');
    await fillInput(page, 'cost', '100.00');
    await submitForm(page, 'Crear');
    await page.waitForURL((url) => url.pathname === '/inventory/maintenance-records', { timeout: 15000 });
    await expectTableVisible(page);

    await clickFirstRow(page);
    confirmDialog(page);
    await page.click('button:has-text("Eliminar")');
    await page.waitForURL((url) => url.pathname === '/inventory/maintenance-records', { timeout: 15000 });
    await waitForToast(page, 'Registro de mantenimiento eliminado');
  });
});

test.describe('Inventory/Maintenance Schedules Module', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows maintenance schedules', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryMaintenanceSchedules);
    await verifyPageTitle(page, 'Programación de Mantención');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
  });

  test('Create new maintenance schedule flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryMaintenanceSchedules);
    await expectTableVisible(page);
    await clickCreateButton(page, ['+ Crear']);
    await verifyPageTitle(page, 'Crear Programación');

    await fillInput(page, 'reference_id', '00000000-0000-0000-0000-000000000001');
    await fillSelect(page, 'type', 1);
    await fillInput(page, 'frequency_days', '30');

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/inventory/maintenance-schedules');
    await waitForToast(page, 'Programa de mantenimiento creado');
    await expectRowCount(page, 2);
  });

  test('View maintenance schedule detail page', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryMaintenanceSchedules);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Editar Programación');
  });

  test('Edit maintenance schedule flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryMaintenanceSchedules);
    await expectTableVisible(page);
    await clickFirstRow(page);

    await fillInput(page, 'reference_id', '00000000-0000-0000-0000-000000000001');
    await page.click('button[type="submit"]:has-text("Guardar")');
    await waitForRedirect(page, '**/inventory/maintenance-schedules');
    await waitForToast(page, 'Programa de mantenimiento actualizado');
  });

  test('Delete maintenance schedule flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryMaintenanceSchedules);
    await expectTableVisible(page);

    await clickCreateButton(page, ['+ Crear']);
    await fillInput(page, 'reference_id', '00000000-0000-0000-0000-000000000001');
    await fillSelect(page, 'type', 1);
    await fillInput(page, 'frequency_days', '60');
    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/inventory/maintenance-schedules');

    await expectTableVisible(page);
    await clickFirstRow(page);
    confirmDialog(page);
    await page.click('button:has-text("Eliminar")');
    await waitForRedirect(page, '**/inventory/maintenance-schedules');
    await waitForToast(page, 'Programa de mantenimiento eliminado');
  });
});