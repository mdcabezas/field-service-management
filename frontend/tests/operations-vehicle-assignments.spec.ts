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

test.describe('Operations/Vehicle Assignments Module', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows vehicle assignments', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVehicleAssignments);
    await verifyPageTitle(page, 'Asignaciones de Vehículo');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
  });

  test('Create new vehicle assignment flow', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVehicleAssignments);
    await clickCreateButton(page, '+ Crear Asignación');
    await verifyPageTitle(page, 'Crear Asignación de Vehículo');

    await fillInput(page, 'vehicle_id', '00000000-0000-0000-0000-000000000001');

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/operations/vehicle-assignments');
    await waitForToast(page, 'Asignación de vehículo creada');
    await expectRowCount(page, 2);
  });

  test('View vehicle assignment detail page', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVehicleAssignments);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Detalle de Asignación de Vehículo');
  });

  test('Edit vehicle assignment flow', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVehicleAssignments);
    await clickFirstRow(page);

    await fillInput(page, 'vehicle_id', '00000000-0000-0000-0000-000000000001');
    await page.click('button[type="submit"]:has-text("Guardar")');
    await waitForRedirect(page, '**/operations/vehicle-assignments');
    await waitForToast(page, 'Asignación de vehículo actualizada');
  });

  test('Delete vehicle assignment flow', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVehicleAssignments);

    await clickCreateButton(page, '+ Crear Asignación');
    await fillInput(page, 'vehicle_id', '00000000-0000-0000-0000-000000000001');

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/operations/vehicle-assignments');

    await clickFirstRow(page);
    await page.click('button:has-text("Eliminar")');
    await waitForRedirect(page, '**/operations/vehicle-assignments');
    await waitForToast(page, 'Asignación de vehículo eliminada');
  });
});