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

test.describe('Operations/Visits Module', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows visits', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await verifyPageTitle(page, 'Visitas');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
  });

  test('Create new visit flow', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await clickCreateButton(page, '+ Crear Visita');
    await verifyPageTitle(page, 'Crear Visita');

    const today = new Date().toISOString().split('T')[0];
    await fillInput(page, 'scheduled_at', today + 'T10:00');
    await fillInput(page, 'type', 'maintenance');
    await fillInput(page, 'status', 'pending');
    await fillInput(page, 'priority', 'high');
    await fillInput(page, 'billing_to', 'Test Customer');

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/operations/visits');
    await waitForToast(page, 'Visita creada');
    await expectRowCount(page, 2);
  });

  test('View visit detail page', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Detalle de Visita');
  });

  test('Edit visit flow', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await clickFirstRow(page);

    await fillInput(page, 'notes', 'Updated visit notes');
    await page.click('button[type="submit"]:has-text("Guardar")');
    await waitForRedirect(page, '**/operations/visits');
    await waitForToast(page, 'Visita actualizada');
  });

  test('Delete visit flow', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);

    await clickCreateButton(page, '+ Crear Visita');
    const today = new Date().toISOString().split('T')[0];
    await fillInput(page, 'scheduled_at', today + 'T10:00');
    await fillInput(page, 'type', 'maintenance');
    await fillInput(page, 'status', 'pending');
    await fillInput(page, 'priority', 'high');
    await fillInput(page, 'billing_to', 'Delete Customer');
    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/operations/visits');

    await clickFirstRow(page);
    await page.click('button:has-text("Eliminar")');
    await waitForRedirect(page, '**/operations/visits');
    await waitForToast(page, 'Visita eliminada');
  });
});