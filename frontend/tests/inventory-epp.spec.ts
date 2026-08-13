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

test.describe('Inventory/EPP Module', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows EPP items', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryEPP);
    await verifyPageTitle(page, 'Elementos de Protección Personal');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
  });

  test('Create new EPP item flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryEPP);
    await expectTableVisible(page);
    await clickCreateButton(page, ['+ Crear']);
    await verifyPageTitle(page, 'Crear EPP');

    await fillInput(page, 'name', 'Test EPP');
    await fillSelect(page, 'type', 1);
    await fillSelect(page, 'lifecycle', 1);

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/inventory/epp');
    await waitForToast(page, 'EPP creado');
    await expectRowCount(page, 2);
  });

  test('View EPP detail page', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryEPP);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Editar EPP');
  });

  test('Edit EPP flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryEPP);
    await expectTableVisible(page);
    await clickFirstRow(page);

    await fillInput(page, 'name', 'Updated EPP Name');
    await page.click('button[type="submit"]:has-text("Guardar")');
    await waitForRedirect(page, '**/inventory/epp');
    await waitForToast(page, 'EPP actualizado');
  });

  test('Delete EPP flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryEPP);
    await expectTableVisible(page);

    await clickCreateButton(page, ['+ Crear']);
    await fillInput(page, 'name', 'To Delete EPP');
    await fillSelect(page, 'type', 1);
    await fillSelect(page, 'lifecycle', 1);
    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/inventory/epp');

    await page.locator('tbody tr', { hasText: 'To Delete EPP' }).last().click();
    confirmDialog(page);
    await page.click('button:has-text("Eliminar")');
    await waitForRedirect(page, '**/inventory/epp');
    await waitForToast(page, 'EPP eliminado');
  });
});