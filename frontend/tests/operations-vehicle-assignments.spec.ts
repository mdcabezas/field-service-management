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
    
    // Reset vehicle status to available before each test
    const authStorage = await page.evaluate(() => {
      const stored = localStorage.getItem('auth-storage');
      return stored ? JSON.parse(stored) : null;
    });
    const token = authStorage?.state?.accessToken;
    await page.request.put('http://localhost:8081/api/vehicles/53000000-0000-0000-0000-000000000001', {
      headers: { 
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      data: {
        id: '53000000-0000-0000-0000-000000000001',
        type: '0bf65e18-2316-4870-a9a0-8ceb76bde851',
        license_plate: 'AB-1234',
        name: 'Camioneta 1',
        brand: 'Updated Brand',
        status: 'available',
        created_at: '2026-08-12T06:27:06.82671Z'
      }
    });
  });

  test('List page loads and shows vehicle assignments', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVehicleAssignments);
    await verifyPageTitle(page, 'Asignaciones de Vehículo');
    await expectTableVisible(page);
    await expectRowCount(page, 0);
  });

  test('Create new vehicle assignment flow', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVehicleAssignments);
    await clickCreateButton(page, '+ Crear Asignación');
    await verifyPageTitle(page, 'Crear Asignación de Vehículo');

    // Wait for vehicles to load in the select
    await page.waitForFunction(() => {
      const select = document.querySelector('select[id="vehicle_id"]');
      return select && select.options.length > 1;
    });

    await fillSelect(page, 'vehicle_id', 1);
    await page.dispatchEvent('select[id="vehicle_id"]', 'change');

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/operations/vehicle-assignments');
    await waitForToast(page, 'Asignación de vehículo creada');
    await expectRowCount(page, 1);
  });

  test('View vehicle assignment detail page', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVehicleAssignments);
    await clickCreateButton(page, '+ Crear Asignación');

    // Wait for vehicles to load in the select
    await page.waitForFunction(() => {
      const select = document.querySelector('select[id="vehicle_id"]');
      return select && select.options.length > 1;
    });

    await fillSelect(page, 'vehicle_id', 1);
    await page.dispatchEvent('select[id="vehicle_id"]', 'change');

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/operations/vehicle-assignments');
    await waitForToast(page, 'Asignación de vehículo creada');

    await clickFirstRow(page);
    await verifyPageTitle(page, 'Asignación de Vehículo');
  });

  test('Delete vehicle assignment flow', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVehicleAssignments);

    await clickCreateButton(page, '+ Crear Asignación');

    // Wait for vehicles to load in the select
    await page.waitForFunction(() => {
      const select = document.querySelector('select[id="vehicle_id"]');
      return select && select.options.length > 1;
    });

    await fillSelect(page, 'vehicle_id', 1);
    await page.dispatchEvent('select[id="vehicle_id"]', 'change');

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/operations/vehicle-assignments');
    await waitForToast(page, 'Asignación de vehículo creada');

    await clickFirstRow(page);
    await confirmDialog(page);
    await page.click('button:has-text("Eliminar")');
    await waitForRedirect(page, '**/operations/vehicle-assignments');
    await waitForToast(page, 'Asignación de vehículo eliminada');
    
    await expectRowCount(page, 0);
  });
});