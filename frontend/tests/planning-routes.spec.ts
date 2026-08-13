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

test.describe('Planning/Routes Module', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows routes', async ({ page }) => {
    await gotoPage(page, MODULES.planningRoutes);
    await verifyPageTitle(page, 'Rutas');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
  });

  test('Search filters routes', async ({ page }) => {
    await gotoPage(page, MODULES.planningRoutes);
    await page.fill('input[placeholder*="Buscar"]', 'test');
    await expectTableVisible(page);
  });

  test('Create new route flow', async ({ page }) => {
    await gotoPage(page, MODULES.planningRoutes);
    await clickCreateButton(page, '+ Crear Ruta');
    await verifyPageTitle(page, 'Crear Ruta');

    await fillSelect(page, 'type', 1);
    const today = new Date().toISOString().split('T')[0];
    await fillInput(page, 'date', today);
    await fillInput(page, 'daily_plan_id', '00000000-0000-0000-0000-000000000001');
    await fillSelect(page, 'status', 1);

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/planning/routes');
    await waitForToast(page, 'Ruta creada');
    await expectRowCount(page, 2);
  });

  test('View route detail page', async ({ page }) => {
    await gotoPage(page, MODULES.planningRoutes);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Detalle de Ruta');
  });

  test('Edit route flow', async ({ page }) => {
    await gotoPage(page, MODULES.planningRoutes);
    await clickFirstRow(page);

    await fillSelect(page, 'type', 1);
    await page.click('button[type="submit"]:has-text("Guardar")');
    await waitForRedirect(page, '**/planning/routes');
    await waitForToast(page, 'Ruta actualizada');
  });

  test('Delete route flow', async ({ page }) => {
    await gotoPage(page, MODULES.planningRoutes);

    await clickCreateButton(page, '+ Crear Ruta');
    await fillSelect(page, 'type', 1);
    const today = new Date().toISOString().split('T')[0];
    await fillInput(page, 'date', today);
    await fillInput(page, 'daily_plan_id', '00000000-0000-0000-0000-000000000001');
    await fillSelect(page, 'status', 1);
    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/planning/routes');

    await clickFirstRow(page);
    await page.click('button:has-text("Eliminar")');
    await waitForRedirect(page, '**/planning/routes');
    await waitForToast(page, 'Ruta eliminada');
  });
});