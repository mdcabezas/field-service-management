import { test, expect } from '@playwright/test';
import { login, gotoPage, verifyPageTitle } from './test-helpers';

test.describe('Section Index Pages', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('Dashboard loads with metrics', async ({ page }) => {
    await gotoPage(page, '/');
    await verifyPageTitle(page, 'Dashboard');
    await expect(page.locator('h1:has-text("Dashboard")')).toBeVisible();
  });

  test('Core section index redirects to users', async ({ page }) => {
    await gotoPage(page, '/core');
    await page.waitForURL('**/core/users', { timeout: 5000 });
    await expect(page.locator('h1:has-text("Usuarios")')).toBeVisible();
  });

  test('Planning section index shows module cards', async ({ page }) => {
    await gotoPage(page, '/planning');
    await verifyPageTitle(page, 'Planificación');
    await expect(page.locator('h1:has-text("Planificación")')).toBeVisible();
  });

  test('Inventory section index shows module cards', async ({ page }) => {
    await gotoPage(page, '/inventory');
    await verifyPageTitle(page, 'Inventario');
    await expect(page.locator('h1:has-text("Inventario")')).toBeVisible();
  });

  test('Operations section index shows module cards', async ({ page }) => {
    await gotoPage(page, '/operations');
    await verifyPageTitle(page, 'Operaciones');
    await expect(page.locator('h1:has-text("Operaciones")')).toBeVisible();
  });

  test('Notifications section index shows module cards', async ({ page }) => {
    await gotoPage(page, '/notifications');
    await verifyPageTitle(page, 'Notificaciones');
    await expect(page.locator('h1:has-text("Notificaciones")')).toBeVisible();
  });

  test('Shared section index shows module cards', async ({ page }) => {
    await gotoPage(page, '/shared');
    await verifyPageTitle(page, 'Compartido');
    await expect(page.locator('h1:has-text("Compartido")')).toBeVisible();
  });
});
