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

  test('List page shows names instead of UUIDs in columns', async ({ page }) => {
    await gotoPage(page, MODULES.operationsReports);
    await expectTableVisible(page);
    // Verify column headers use name fields, not UUID fields
    await expect(page.locator('th:has-text("Plantilla")')).toBeVisible();
    await expect(page.locator('th:has-text("Visita")')).toBeVisible();
  });

  test('Create new report flow', async ({ page }) => {
    await gotoPage(page, MODULES.operationsReports);
    await clickCreateButton(page, '+ Crear Reporte');
    await verifyPageTitle(page, 'Crear Reporte');

    await fillSelect(page, 'visit_id', 1);
    await fillSelect(page, 'report_template_id', 1);
    await fillSelect(page, 'source', 1);

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/operations/reports');
    await waitForToast(page, 'Reporte creado');
    await expectRowCount(page, 2);
  });

  test('View report detail page', async ({ page }) => {
    await gotoPage(page, MODULES.operationsReports);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Reporte');
  });

  test('Detail page shows names instead of UUIDs', async ({ page }) => {
    await gotoPage(page, MODULES.operationsReports);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Reporte');
    // Verify labels are visible
    await expect(page.locator('label:has-text("Visita")')).toBeVisible();
    await expect(page.locator('label:has-text("Plantilla")')).toBeVisible();
  });

  test('Delete report flow', async ({ page }) => {
    await gotoPage(page, MODULES.operationsReports);

    await clickCreateButton(page, '+ Crear Reporte');

    await fillSelect(page, 'visit_id', 1);
    await fillSelect(page, 'report_template_id', 1);
    await fillSelect(page, 'source', 1);

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/operations/reports');

    await clickFirstRow(page);
    await confirmDialog(page);
    await page.click('button:has-text("Eliminar")');
    await waitForRedirect(page, '**/operations/reports');
    await waitForToast(page, 'Reporte eliminado');
  });
});