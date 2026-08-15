import { test, expect } from '@playwright/test';
import {
  login,
  gotoPage,
  verifyPageTitle,
  clickFirstRow,
  MODULES,
} from './test-helpers';

test.describe('Partner SLAs Tab', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('Partner detail page has SLAs tab', async ({ page }) => {
    await gotoPage(page, MODULES.partners);
    await clickFirstRow(page);
    await expect(page.locator('button[role="tab"]:has-text("SLAs")')).toBeVisible();
  });

  test('Click SLAs tab shows SLA table', async ({ page }) => {
    await gotoPage(page, MODULES.partners);
    await clickFirstRow(page);
    await page.click('button[role="tab"]:has-text("SLAs")');
    await expect(page.locator('text=SLAs del Socio')).toBeVisible();
    await expect(page.locator('button:has-text("+ Agregar SLA")')).toBeVisible();
  });

  test('Open add SLA dialog shows form', async ({ page }) => {
    await gotoPage(page, MODULES.partners);
    await clickFirstRow(page);
    await page.click('button[role="tab"]:has-text("SLAs")');
    await page.click('button:has-text("+ Agregar SLA")');

    await page.waitForLoadState('networkidle');

    await expect(page.getByRole('dialog')).toBeVisible();
    await expect(page.getByText('Nombre *')).toBeVisible();
    await expect(page.getByText('Vigencia desde *')).toBeVisible();
  });

  test('Close add SLA dialog on cancel', async ({ page }) => {
    await gotoPage(page, MODULES.partners);
    await clickFirstRow(page);
    await page.click('button[role="tab"]:has-text("SLAs")');
    await page.click('button:has-text("+ Agregar SLA")');

    await page.waitForLoadState('networkidle');

    const dialog = page.getByRole('dialog');
    await expect(dialog).toBeVisible();
    await dialog.locator('button:has-text("Cancelar")').click();
    await expect(dialog).not.toBeVisible();
  });
});
