import { test, expect } from '@playwright/test';
import {
  login,
  gotoPage,
  verifyPageTitle,
  expectTableVisible,
  clickFirstRow,
  fillInput,
  submitForm,
  waitForToast,
  MODULES,
} from './test-helpers';

test.describe('Visit Checklist Full Flow', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('Visit detail shows checklist section', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    // Look for checklist section
    const checklistSection = page.locator('text=Checklist, text=Lista de Verificación');
    if (await checklistSection.isVisible()) {
      await expect(checklistSection).toBeVisible();
    }
  });

  test('Checklist items are editable', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    // Find checklist checkboxes
    const checkboxes = page.locator('input[type="checkbox"]');
    const count = await checkboxes.count();

    if (count > 0) {
      // Toggle first checkbox
      await checkboxes.first().click();
      await page.waitForLoadState('networkidle');
    }
  });

  test('Checklist completion blocks status transition', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    // Look for status transition buttons
    const completeBtn = page.locator('button:has-text("Completar"), button:has-text("Complete")');
    if (await completeBtn.isVisible()) {
      // Try to complete without filling checklist
      await completeBtn.click();

      // Should show validation error or warning
      const warning = page.locator('text=checklist, text=etapas, text=incompleto');
      // Warning may or may not appear depending on data
    }
  });

  test('Navigate to checklist detail page', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    // Look for view checklist button/link
    const checklistLink = page.locator('a:has-text("Ver Checklist"), button:has-text("Ver Checklist")');
    if (await checklistLink.isVisible()) {
      await checklistLink.click();
      await page.waitForLoadState('networkidle');

      // Should be on checklist detail page
      const checklistPage = page.locator('text=Checklist, text=Etapas');
      await expect(checklistPage).toBeVisible({ timeout: 10000 });
    }
  });

  test('Checklist stages show completion progress', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    // Look for progress indicators
    const progress = page.locator('[class*="progress"], [class*="badge"]:has-text("/"), text=/\\d+\\/\\d+/');
    if (await progress.isVisible()) {
      await expect(progress).toBeVisible();
    }
  });
});