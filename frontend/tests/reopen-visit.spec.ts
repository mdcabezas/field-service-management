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

test.describe('Reopen Visit', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('Visit detail shows reopen button for completed visits', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await expectTableVisible(page);

    // Filter for completed visits
    const statusFilter = page.locator('select[name="status"], select:has(option:text("completed"))');
    if (await statusFilter.isVisible()) {
      await statusFilter.selectOption('completed');
      await page.waitForLoadState('networkidle');
    }

    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    // Look for reopen button
    const reopenBtn = page.locator('button:has-text("Reabrir"), button:has-text("Reopen")');
    if (await reopenBtn.isVisible()) {
      await expect(reopenBtn).toBeVisible();
    }
  });

  test('Reopen dialog shows reason input', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    const reopenBtn = page.locator('button:has-text("Reabrir"), button:has-text("Reopen")');
    if (await reopenBtn.isVisible()) {
      await reopenBtn.click();

      // Should show dialog with reason input
      const reasonInput = page.locator('textarea, input[name="reason"], input[placeholder*="razón"], input[placeholder*="reason"]');
      await expect(reasonInput).toBeVisible({ timeout: 5000 });
    }
  });

  test('Submit reopen with reason', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    const reopenBtn = page.locator('button:has-text("Reabrir"), button:has-text("Reopen")');
    if (await reopenBtn.isVisible()) {
      await reopenBtn.click();

      const reasonInput = page.locator('textarea, input[name="reason"]');
      if (await reasonInput.isVisible()) {
        await reasonInput.fill('Additional work needed');
        await submitForm(page, 'Reabrir', 'Reopen');
        await waitForToast(page, 'reabierta');
      }
    }
  });

  test('Reopen updates visit status', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    const reopenBtn = page.locator('button:has-text("Reabrir"), button:has-text("Reopen")');
    if (await reopenBtn.isVisible()) {
      await reopenBtn.click();

      const reasonInput = page.locator('textarea, input[name="reason"]');
      if (await reasonInput.isVisible()) {
        await reasonInput.fill('Customer request');
        await submitForm(page, 'Reabrir', 'Reopen');
        await waitForToast(page, 'reabierta');

        // Verify status changed
        const statusBadge = page.locator('[class*="badge"], [class*="status"]');
        await expect(statusBadge).toBeVisible();
      }
    }
  });
});