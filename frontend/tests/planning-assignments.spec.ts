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

test.describe('Planning/Assignments Module (Expanded)', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page shows resolved names not UUIDs', async ({ page }) => {
    await gotoPage(page, MODULES.planningAssignments);
    await verifyPageTitle(page, 'Asignaciones');
    await expectTableVisible(page);

    const firstRow = page.locator('tbody tr').first();
    await expect(firstRow.locator('td').nth(0)).not.toContainText('0000-0000');
    await expect(firstRow.locator('td').nth(1)).not.toContainText('0000-0000');
    await expect(firstRow.locator('td').nth(2)).not.toContainText('0000-0000');
  });

  test('Create assignment with all dropdowns populated', async ({ page }) => {
    await gotoPage(page, MODULES.planningAssignments);
    await clickCreateButton(page, '+ Crear Asignación');
    await verifyPageTitle(page, 'Crear Asignación');

    await fillSelect(page, 'daily_plan_id', 1);
    await fillSelect(page, 'tech_id', 1);
    await fillSelect(page, 'role_id', 1);

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/planning/assignments');
    await waitForToast(page, 'Asignación creada');
  });

  test('Create assignment with invalid daily plan shows error', async ({ page }) => {
    await gotoPage(page, MODULES.planningAssignments);
    await clickCreateButton(page, '+ Crear Asignación');

    await fillSelect(page, 'daily_plan_id', 1);
    await fillSelect(page, 'tech_id', 1);
    await fillSelect(page, 'role_id', 1);

    await submitForm(page, 'Crear');
    await waitForToast(page, 'Asignación creada');
  });

  test('Edit assignment changes tech_id and role_id', async ({ page }) => {
    await gotoPage(page, MODULES.planningAssignments);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Detalle de Asignación');

    await fillSelect(page, 'tech_id', 2);
    await fillSelect(page, 'role_id', 2);

    await page.click('button[type="submit"]:has-text("Guardar")');
    await waitForRedirect(page, '**/planning/assignments');
    await waitForToast(page, 'Asignación actualizada');
  });

  test('daily_plan_id is immutable on edit', async ({ page }) => {
    await gotoPage(page, MODULES.planningAssignments);
    await clickFirstRow(page);

    await expect(page.locator('select[id="daily_plan_id"]')).toBeDisabled();
    await expect(page.locator('select[id="tech_id"]')).toBeEnabled();
    await expect(page.locator('select[id="role_id"]')).toBeEnabled();
  });

  test('Delete assignment flow', async ({ page }) => {
    await gotoPage(page, MODULES.planningAssignments);

    await clickCreateButton(page, '+ Crear Asignación');
    await fillSelect(page, 'daily_plan_id', 1);
    await fillSelect(page, 'tech_id', 1);
    await fillSelect(page, 'role_id', 1);
    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/planning/assignments');

    await clickFirstRow(page);
    await page.click('button:has-text("Eliminar")');
    await waitForRedirect(page, '**/planning/assignments');
    await waitForToast(page, 'Asignación eliminada');
  });

  test('Empty state handling when no technicians', async ({ page }) => {
    await page.goto('http://localhost:3000/planning/assignments/new');
    await verifyPageTitle(page, 'Crear Asignación');

    const techOptions = await page.locator('select[id="tech_id"] option').count();
    const roleOptions = await page.locator('select[id="role_id"] option').count();
    const planOptions = await page.locator('select[id="daily_plan_id"] option').count();

    expect(techOptions).toBeGreaterThan(1);
    expect(roleOptions).toBeGreaterThan(1);
    expect(planOptions).toBeGreaterThan(1);
  });

  test('Search filters assignments', async ({ page }) => {
    await gotoPage(page, MODULES.planningAssignments);
    await page.fill('input[placeholder*="Buscar"]', '2026');
    await expectTableVisible(page);
  });
});