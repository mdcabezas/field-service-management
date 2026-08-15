import { test, expect } from '@playwright/test';
import {
  login,
  gotoPage,
  clickFirstRow,
  MODULES,
} from './test-helpers';

test.describe('Visit Detail Tabs', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('Visit detail has all tabs', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await clickFirstRow(page);
    
    const tabNames = ['Información', 'Asignaciones', 'Checklist', 'Mediciones', 'Checkpoints', 'Fotos', 'Reportes', 'Vehículos', 'SLAs'];
    for (const name of tabNames) {
      await expect(page.locator(`button[role="tab"]:has-text("${name}")`)).toBeVisible();
    }
  });

  test('Click through all tabs without error', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await clickFirstRow(page);
    
    const tabs = ['Asignaciones', 'Checklist', 'Mediciones', 'SLAs', 'Información'];
    for (const tab of tabs) {
      await page.click(`button[role="tab"]:has-text("${tab}")`);
      await page.waitForLoadState('networkidle');
    }
  });

  test('Informacion tab shows visit fields', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await clickFirstRow(page);
    await expect(page.locator('text=Fecha')).toBeVisible();
    await expect(page.locator('text=Prioridad')).toBeVisible();
    await expect(page.locator('text=Estado')).toBeVisible();
  });
});
