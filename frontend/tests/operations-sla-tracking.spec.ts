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

test.describe('Operations/SLA Tracking Module', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows SLA trackings', async ({ page }) => {
    await gotoPage(page, MODULES.operationsSLATracking);
    await verifyPageTitle(page, 'Seguimiento SLA');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
  });

  test('Create new SLA tracking flow', async ({ page }) => {
    await gotoPage(page, MODULES.operationsSLATracking);
    await clickCreateButton(page, '+ Crear Seguimiento');
    await verifyPageTitle(page, 'Crear Seguimiento SLA');

    await fillInput(page, 'visit_id', '00000000-0000-0000-0000-000000000001');
    await fillInput(page, 'sla_id', '00000000-0000-0000-0000-000000000001');

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/operations/sla-trackings');
    await waitForToast(page, 'Seguimiento SLA creado');
    await expectRowCount(page, 2);
  });

  test('View SLA tracking detail page', async ({ page }) => {
    await gotoPage(page, MODULES.operationsSLATracking);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Detalle de Seguimiento SLA');
  });

  test('Edit SLA tracking flow (update response time)', async ({ page }) => {
    await gotoPage(page, MODULES.operationsSLATracking);
    await clickFirstRow(page);

    await fillInput(page, 'requested_at', new Date().toISOString());
    await page.click('button[type="submit"]:has-text("Guardar")');
    await waitForRedirect(page, '**/operations/sla-trackings');
    await waitForToast(page, 'Seguimiento SLA actualizado');
  });

  test('Delete SLA tracking flow', async ({ page }) => {
    await gotoPage(page, MODULES.operationsSLATracking);

    await clickCreateButton(page, '+ Crear Seguimiento');
    await fillInput(page, 'visit_id', '00000000-0000-0000-0000-000000000001');
    await fillInput(page, 'sla_id', '00000000-0000-0000-0000-000000000001');
    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/operations/sla-trackings');

    await clickFirstRow(page);
    await page.click('button:has-text("Eliminar")');
    await waitForRedirect(page, '**/operations/sla-trackings');
    await waitForToast(page, 'Seguimiento SLA eliminado');
  });
});