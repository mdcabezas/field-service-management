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

    // Wait for visits to load
    await page.waitForFunction(() => {
      const select = document.querySelector('select[id="visit_id"]');
      return select && select.options.length > 1;
    });

    await fillSelect(page, 'visit_id', 1);

    // Wait for partners to load
    await page.waitForFunction(() => {
      const select = document.querySelector('select[id="partner_id"]');
      return select && select.options.length > 1;
    });

    await fillSelect(page, 'partner_id', 1);

    // Wait for SLAs to load after partner selection
    await page.waitForFunction(() => {
      const select = document.querySelector('select[id="sla_id"]');
      return select && select.options.length > 1;
    });

    await fillSelect(page, 'sla_id', 1);

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/operations/sla-trackings');
    await waitForToast(page, 'Seguimiento SLA creado');
    await expectRowCount(page, 2);
  });

  test('View SLA tracking detail page', async ({ page }) => {
    await gotoPage(page, MODULES.operationsSLATracking);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Seguimiento SLA');
  });

  test('Detail page shows SLA tracking info', async ({ page }) => {
    await gotoPage(page, MODULES.operationsSLATracking);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Seguimiento SLA');
    await expect(page.locator('text=Información del Seguimiento')).toBeVisible();
    await expect(page.locator('button:has-text("Eliminar")')).toBeVisible();
  });

  test('Delete SLA tracking flow', async ({ page }) => {
    await gotoPage(page, MODULES.operationsSLATracking);

    await clickCreateButton(page, '+ Crear Seguimiento');

    // Wait for visits to load
    await page.waitForFunction(() => {
      const select = document.querySelector('select[id="visit_id"]');
      return select && select.options.length > 1;
    });

    await fillSelect(page, 'visit_id', 1);

    // Wait for partners to load
    await page.waitForFunction(() => {
      const select = document.querySelector('select[id="partner_id"]');
      return select && select.options.length > 1;
    });

    await fillSelect(page, 'partner_id', 1);

    // Wait for SLAs to load
    await page.waitForFunction(() => {
      const select = document.querySelector('select[id="sla_id"]');
      return select && select.options.length > 1;
    });

    await fillSelect(page, 'sla_id', 1);

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/operations/sla-trackings');

    await clickFirstRow(page);
    await confirmDialog(page);
    await page.click('button:has-text("Eliminar")');
    await waitForRedirect(page, '**/operations/sla-trackings');
    await waitForToast(page, 'Seguimiento SLA eliminado');
  });
});