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

test.describe('Shared/Report Templates Module', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows report templates', async ({ page }) => {
    await gotoPage(page, MODULES.sharedReportTemplates);
    await verifyPageTitle(page, 'Plantillas de Reporte');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
  });

  test('Create new report template flow', async ({ page }) => {
    await gotoPage(page, MODULES.sharedReportTemplates);
    await clickCreateButton(page, '+ Crear Plantilla');
    await verifyPageTitle(page, 'Crear Plantilla de Reporte');

    await fillInput(page, 'name', 'Test Report Template');
    await fillInput(page, 'fields_json', '{"field1": "value1", "field2": "value2"}');

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/shared/report-templates');
    await waitForToast(page, 'Plantilla de reporte creada');
    await expectRowCount(page, 2);
  });

  test('View report template detail page', async ({ page }) => {
    await gotoPage(page, MODULES.sharedReportTemplates);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Detalle de Plantilla de Reporte');
  });

  test('Edit report template flow', async ({ page }) => {
    await gotoPage(page, MODULES.sharedReportTemplates);
    await clickFirstRow(page);

    await fillInput(page, 'name', 'Updated Template Name');
    await page.click('button[type="submit"]:has-text("Guardar")');
    await waitForRedirect(page, '**/shared/report-templates');
    await waitForToast(page, 'Plantilla de reporte actualizada');
  });

  test('Delete report template flow', async ({ page }) => {
    await gotoPage(page, MODULES.sharedReportTemplates);

    await clickCreateButton(page, '+ Crear Plantilla');
    await fillInput(page, 'name', 'To Delete Template');
    await fillInput(page, 'fields_json', '{"delete": "true"}');
    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/shared/report-templates');

    await clickFirstRow(page);
    await page.click('button:has-text("Eliminar")');
    await waitForRedirect(page, '**/shared/report-templates');
    await waitForToast(page, 'Plantilla de reporte eliminada');
  });
});