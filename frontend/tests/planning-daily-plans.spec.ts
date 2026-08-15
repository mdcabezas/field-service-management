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

test.describe('Planning/Daily Plans Module', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows daily plans', async ({ page }) => {
    await gotoPage(page, MODULES.planningDailyPlans);
    await verifyPageTitle(page, 'Planes Diarios');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
  });

  test('Search filters daily plans', async ({ page }) => {
    await gotoPage(page, MODULES.planningDailyPlans);
    await page.fill('input[placeholder*="Buscar"]', '2026');
    await expect(page.locator('tbody tr')).toHaveCount(1);
  });

  test('Create new daily plan flow', async ({ page }) => {
    await gotoPage(page, MODULES.planningDailyPlans);
    await clickCreateButton(page, '+ Crear Plan');
    await verifyPageTitle(page, 'Crear Plan Diario');

    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    const dateStr = tomorrow.toISOString().split('T')[0];

    await fillInput(page, 'date', dateStr);
    await fillInput(page, 'notes', 'Test plan notes');
    await page.click('button[type="submit"]:has-text("Crear")');
    await waitForRedirect(page, '**/planning/daily-plans');
    await waitForToast(page, 'Plan diario creado');
    await expectRowCount(page, 2);
  });

  test('View daily plan detail page', async ({ page }) => {
    await gotoPage(page, MODULES.planningDailyPlans);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Editar Plan Diario');
  });

  test('Edit daily plan flow', async ({ page }) => {
    await gotoPage(page, MODULES.planningDailyPlans);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Editar Plan Diario');
    await page.waitForLoadState('networkidle');

    await fillInput(page, 'notes', 'Updated test plan notes');
    await page.click('button[type="submit"]:has-text("Guardar")');
    await waitForToast(page, 'Plan diario actualizado');
    await waitForRedirect(page, '**/planning/daily-plans');

    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');
    await expect(page.locator('#notes')).toHaveValue('Updated test plan notes');
  });

  test('Delete daily plan flow', async ({ page }) => {
    await gotoPage(page, MODULES.planningDailyPlans);

    await clickCreateButton(page, '+ Crear Plan');
    const nextWeek = new Date();
    nextWeek.setDate(nextWeek.getDate() + 7);
    const dateStr = nextWeek.toISOString().split('T')[0];
    await fillInput(page, 'date', dateStr);
    await fillInput(page, 'notes', 'To delete plan');
    await page.click('button[type="submit"]:has-text("Crear")');
    await waitForRedirect(page, '**/planning/daily-plans');

    await clickFirstRow(page);
    await page.click('button:has-text("Eliminar")');
    await waitForRedirect(page, '**/planning/daily-plans');
    await waitForToast(page, 'Plan diario eliminado');
  });
});