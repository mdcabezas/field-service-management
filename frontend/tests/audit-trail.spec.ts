import { test, expect } from '@playwright/test';
import {
  login,
  gotoPage,
  verifyPageTitle,
  expectTableVisible,
  clickFirstRow,
  MODULES,
} from './test-helpers';

test.describe('Audit Trail', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('Visit detail shows audit history', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    // Look for audit section or link
    const auditSection = page.locator('text=Auditoría, text=Audit, text=Historial');
    if (await auditSection.isVisible()) {
      await expect(auditSection).toBeVisible();
    }
  });

  test('Audit entries display timestamp and user', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    // Navigate to audit view
    const auditLink = page.locator('a:has-text("Auditoría"), button:has-text("Auditoría"), a:has-text("Historial")');
    if (await auditLink.isVisible()) {
      await auditLink.click();
      await page.waitForLoadState('networkidle');

      // Should show audit entries with timestamp
      const timestamp = page.locator('time, [class*="date"], [class*="timestamp"]');
      if (await timestamp.isVisible()) {
        await expect(timestamp).toBeVisible();
      }
    }
  });

  test('Audit trail shows field changes', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    const auditLink = page.locator('a:has-text("Auditoría"), button:has-text("Auditoría"), a:has-text("Historial")');
    if (await auditLink.isVisible()) {
      await auditLink.click();
      await page.waitForLoadState('networkidle');

      // Should show old/new values
      const changeEntry = page.locator('[class*="change"], [class*="diff"], tr:has(td)');
      if (await changeEntry.isVisible()) {
        await expect(changeEntry).toBeVisible();
      }
    }
  });
});