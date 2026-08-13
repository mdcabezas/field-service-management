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

test.describe('Inventory/Cost Rates Module', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows cost rates', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryCostRates);
    await verifyPageTitle(page, 'Tarifas de Costo');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
  });

  test('Create new cost rate flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryCostRates);
    await expectTableVisible(page);
    await clickCreateButton(page, '+ Crear Tarifa');
    await verifyPageTitle(page, 'Crear Tarifa');

    const today = new Date().toISOString().split('T')[0];
    const nextMonth = new Date(Date.now() + 30*24*60*60*1000).toISOString().split('T')[0];

    await fillSelect(page, 'type', 1);
    await fillInput(page, 'value', '50.50');
    await fillSelect(page, 'unit', 1);
    await fillSelect(page, 'reference_type', 1);
    await fillInput(page, 'valid_from', today);
    await fillInput(page, 'valid_until', nextMonth);

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/inventory/cost-rates');
    await waitForToast(page, 'Tarifa de costo creada');
    await expectRowCount(page, 2);
  });

  test('View cost rate detail page', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryCostRates);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Editar Tarifa');
  });

  test('Edit cost rate flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryCostRates);
    await expectTableVisible(page);
    await clickFirstRow(page);

    await fillInput(page, 'value', '75.00');
    await page.click('button[type="submit"]:has-text("Guardar")');
    await waitForRedirect(page, '**/inventory/cost-rates');
    await waitForToast(page, 'Tarifa de costo actualizada');
  });

  test('Delete cost rate flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryCostRates);
    await expectTableVisible(page);

    await clickCreateButton(page, '+ Crear Tarifa');
    const today = new Date().toISOString().split('T')[0];
    const nextMonth = new Date(Date.now() + 30*24*60*60*1000).toISOString().split('T')[0];
    await fillSelect(page, 'type', 1);
    await fillInput(page, 'value', '25.00');
    await fillSelect(page, 'unit', 1);
    await fillSelect(page, 'reference_type', 1);
    await fillInput(page, 'valid_from', today);
    await fillInput(page, 'valid_until', nextMonth);
    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/inventory/cost-rates');

    await expectTableVisible(page);
    await clickFirstRow(page);
    confirmDialog(page);
    await page.click('button:has-text("Eliminar")');
    await waitForRedirect(page, '**/inventory/cost-rates');
    await waitForToast(page, 'Tarifa de costo eliminada');
  });
});