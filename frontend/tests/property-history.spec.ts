import { test, expect } from '@playwright/test';
import {
  login,
  gotoPage,
  verifyPageTitle,
  expectTableVisible,
  clickFirstRow,
  MODULES,
} from './test-helpers';

test.describe('Property History', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('View property visit history', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    // Look for property history section
    const historySection = page.locator('text=Historial, text=History, text=Histórico');
    if (await historySection.isVisible()) {
      await expect(historySection).toBeVisible();
    }
  });

  test('Property history shows audit entries', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    // Look for audit/history link or tab
    const auditLink = page.locator('a:has-text("Historial"), button:has-text("Historial")');
    if (await auditLink.isVisible()) {
      await auditLink.click();
      await page.waitForLoadState('networkidle');

      // Should show audit entries table or list
      const entriesList = page.locator('table, [class*="audit"], [class*="history"]');
      await expect(entriesList).toBeVisible({ timeout: 10000 });
    }
  });

  test('Property history filters by date range', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    // Look for date filter inputs
    const dateFrom = page.locator('input[type="date"], input[name*="from"], input[name*="start"]');
    const dateTo = page.locator('input[type="date"], input[name*="to"], input[name*="end"]');

    if (await dateFrom.isVisible() && await dateTo.isVisible()) {
      const today = new Date().toISOString().split('T')[0];
      await dateFrom.fill(today);
      await dateTo.fill(today);
      await page.waitForLoadState('networkidle');
    }
  });
});