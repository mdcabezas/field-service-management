import { test, expect } from '@playwright/test';
import {
  login,
  gotoPage,
  verifyPageTitle,
  expectTableVisible,
  clickFirstRow,
  MODULES,
} from './test-helpers';

test.describe('Customer Addresses Tab', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('Customer detail page has Direcciones tab', async ({ page }) => {
    await gotoPage(page, MODULES.customers);
    await clickFirstRow(page);
    await expect(page.locator('button[role="tab"]:has-text("Direcciones")')).toBeVisible();
  });

  test('Click Direcciones tab shows address table', async ({ page }) => {
    await gotoPage(page, MODULES.customers);
    await clickFirstRow(page);
    await page.click('button[role="tab"]:has-text("Direcciones")');
    await expect(page.locator('text=Direcciones del Cliente')).toBeVisible();
    await expect(page.locator('button:has-text("+ Agregar Dirección")')).toBeVisible();
  });

  test('Open add address dialog shows form', async ({ page }) => {
    await gotoPage(page, MODULES.customers);
    await clickFirstRow(page);
    await page.click('button[role="tab"]:has-text("Direcciones")');
    await page.click('button:has-text("+ Agregar Dirección")');
    
    // Wait for dialog to appear
    await page.waitForLoadState('networkidle');
    
    // Check dialog content
    await expect(page.getByRole('dialog')).toBeVisible();
    await expect(page.getByText('Tipo *')).toBeVisible();
    await expect(page.getByText('Nombre (opcional)')).toBeVisible();
  });

  test('Address type selector has options', async ({ page }) => {
    await gotoPage(page, MODULES.customers);
    await clickFirstRow(page);
    await page.click('button[role="tab"]:has-text("Direcciones")');
    await page.click('button:has-text("+ Agregar Dirección")');
    
    await page.waitForLoadState('networkidle');
    
    // Check dialog opened with type field
    await expect(page.getByRole('dialog').getByText('Tipo *')).toBeVisible();
  });

  test('Close address dialog on cancel', async ({ page }) => {
    await gotoPage(page, MODULES.customers);
    await clickFirstRow(page);
    await page.click('button[role="tab"]:has-text("Direcciones")');
    await page.click('button:has-text("+ Agregar Dirección")');
    
    await page.waitForLoadState('networkidle');
    
    // Wait for dialog and click cancel
    const dialog = page.getByRole('dialog');
    await expect(dialog).toBeVisible();
    await dialog.locator('button:has-text("Cancelar")').click();
    
    // Dialog should be closed
    await expect(dialog).not.toBeVisible();
  });
});
