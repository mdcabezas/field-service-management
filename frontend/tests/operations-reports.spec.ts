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

test.describe('Operations/Reports Module', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows reports', async ({ page }) => {
    await gotoPage(page, MODULES.operationsReports);
    await verifyPageTitle(page, 'Reportes de Visita');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
  });

  test('Create new report flow', async ({ page }) => {
    await gotoPage(page, MODULES.operationsReports);
    await clickCreateButton(page, '+ Crear Reporte');
    await verifyPageTitle(page, 'Crear Reporte de Visita');

    await fillInput(page, 'visit_id', '00000000-0000-0000-0000-000000000001');
    await fillInput(page, 'report_template_id', '00000000-0000-0000-0000-000000000001');
    await fillInput(page, 'source', 'system');

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/operations/reports');
    await waitForToast(page, 'Reporte de visita creado');
    await expectRowCount(page, 2);
  });

  test('View report detail page', async ({ page }) => {
    await gotoPage(page, MODULES.operationsReports);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Detalle de Reporte de Visita');
  });

  test('Edit report flow', async ({ page }) => {
    await gotoPage(page, MODULES.operationsReports);
    await clickFirstRow(page);

    await fillInput(page, 'source', 'manual');
    await page.click('button[type="submit"]:has-text("Guardar")');
    await waitForRedirect(page, '**/operations/reports');
    await waitForToast(page, 'Reporte de visita actualizado');
  });

  test('Delete report flow', async ({ page }) => {
    await gotoPage(page, MODULES.operationsReports);

    await clickCreateButton(page, '+ Crear Reporte');
    await fillInput(page, 'visit_id', '00000000-0000-0000-0000-000000000001');
    await fillInput(page, 'report_template_id', '00000000-0000-0000-0000-000000000001');
    await fillInput(page, 'source', 'delete_test');
    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/operations/reports');

    await clickFirstRow(page);
    await page.click('button:has-text("Eliminar")');
    await waitForRedirect(page, '**/operations/reports');
    await waitForToast(page, 'Reporte de visita eliminado');
  });
});