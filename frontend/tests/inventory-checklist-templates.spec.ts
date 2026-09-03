import { test, expect } from '@playwright/test';
import {
  login,
  gotoPage,
  waitForToast,
  verifyPageTitle,
  expectTableVisible,
  expectRowCount,
  clickFirstRow,
  clickCreateButton,
  fillInput,
  fillSelect,
  submitForm,
  confirmDialog,
  createMultipleRecords,
  MODULES,
} from './test-helpers';

// Use existing work_type UUIDs from seed data
const WORK_TYPE_UUIDS = [
  '10000000-0000-0000-0000-000000000001',
  '10000000-0000-0000-0000-000000000002',
  '10000000-0000-0000-0000-000000000003',
];

test.describe('Inventory/Checklist Templates Module', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows checklist templates', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryChecklistTemplates);
    await verifyPageTitle(page, 'Plantillas de Checklist');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
  });

  test('Search filters checklist templates', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryChecklistTemplates);
    await page.fill('input[placeholder*="Buscar"]', 'checklist');
    await expectTableVisible(page);
  });

  test('Create new checklist template flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryChecklistTemplates);
    await clickCreateButton(page, ['+ Crear']);
    await verifyPageTitle(page, 'Crear Plantilla');

    await fillInput(page, 'name', 'Test Checklist Template');
    // Use a valid work_type UUID from seed data
    await fillInput(page, 'work_type', WORK_TYPE_UUIDS[0]);

    await submitForm(page, ['Crear']);
    await waitForToast(page, 'Plantilla creada');
    await page.waitForURL((url) => url.pathname === '/inventory/checklist-templates', { timeout: 15000 });
    await expectRowCount(page, 2);
  });

  test('View checklist template detail page', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryChecklistTemplates);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Editar Plantilla');
  });

  test('Edit checklist template flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryChecklistTemplates);
    await clickFirstRow(page);

    await fillInput(page, 'name', 'Updated Checklist Template Name');
    await submitForm(page, ['Guardar']);
    await page.waitForURL((url) => url.pathname === '/inventory/checklist-templates', { timeout: 30000 });
    await waitForToast(page, 'Plantilla actualizada');
  });

  test('Delete checklist template flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryChecklistTemplates);

    await clickCreateButton(page, ['+ Crear']);
    await fillInput(page, 'name', 'To Delete Template');
    await fillInput(page, 'work_type', WORK_TYPE_UUIDS[1]);
    await submitForm(page, ['Crear']);
    await page.waitForURL((url) => url.pathname === '/inventory/checklist-templates', { timeout: 30000 });
    await waitForToast(page, 'Plantilla creada');

    await clickFirstRow(page);
    page.once('dialog', dialog => dialog.accept());
    await page.click('button:has-text("Eliminar")');
    await page.waitForURL((url) => url.pathname === '/inventory/checklist-templates', { timeout: 30000 });
    await waitForToast(page, 'Plantilla eliminada');
    // After deletion, there should be the original 3 templates (or more)
    await expectRowCount(page, 3);
  });
});
