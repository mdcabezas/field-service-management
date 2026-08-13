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
  fillInput,
  fillSelect,
  submitForm,
  confirmDialog,
  createMultipleRecords,
  testPagination,
  testPaginationNavigation,
  MODULES,
} from './test-helpers';

test.describe('Inventory/Materials Module', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows materials', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryMaterials);
    await verifyPageTitle(page, 'Materiales');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
  });

  test('Search filters materials', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryMaterials);
    await page.fill('input[placeholder*="Buscar"]', 'material');
    await expectTableVisible(page);
  });

  test('Create new material flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryMaterials);
    await expectTableVisible(page);
    await clickCreateButton(page, '+ Crear');
    await verifyPageTitle(page, 'Crear Material');

    await fillInput(page, 'sku', `MAT_${Date.now()}`);
    await fillInput(page, 'name', 'Test Material');
    await fillInput(page, 'unit', 'UN');
    await fillInput(page, 'unit_cost', '25.50');

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/inventory/materials');
    await waitForToast(page, 'Material creado');
    await expectRowCount(page, 2);
  });

  test('View material detail page', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryMaterials);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Editar Material');
  });

  test('Edit material flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryMaterials);
    await expectTableVisible(page);
    await clickFirstRow(page);

    await fillInput(page, 'name', 'Updated Material Name');
    await page.click('button[type="submit"]:has-text("Guardar")');
    await waitForRedirect(page, '**/inventory/materials');
    await waitForToast(page, 'Material actualizado');
  });

  test('Delete material flow', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryMaterials);
    await expectTableVisible(page);

    await page.locator('button:has-text("+ Crear")').waitFor({ timeout: 10000 });
    await clickCreateButton(page, '+ Crear');
    await page.waitForURL('**/inventory/materials/new', { timeout: 10000 });

    const sku = `DEL_${Date.now()}`;
    await fillInput(page, 'sku', sku);
    await fillInput(page, 'name', 'AAAAA-Delete Material');
    await fillInput(page, 'unit', 'UN');
    await fillInput(page, 'unit_cost', '5.00');
    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/inventory/materials');
    await page.waitForLoadState('networkidle');
    await gotoPage(page, MODULES.inventoryMaterials);

    await page.fill('input[placeholder*="Buscar"]', 'AAAAA-Delete Material');
    await page.waitForLoadState('networkidle');

    const row = page.locator('tbody tr', { hasText: 'AAAAA-Delete Material' });
    await row.waitFor({ state: 'visible', timeout: 15000 });
    await row.click();
    confirmDialog(page);
    await page.click('button:has-text("Eliminar")');
    await waitForRedirect(page, '**/inventory/materials');
    await waitForToast(page, 'Material eliminado');
  });

  test('Pagination with 6 records (pageSize=5)', async ({ page }) => {
    // Create 6 materials with unique names
    const baseName = `PAG_MAT_${Date.now()}`;
    await createMultipleRecords(
      page,
      MODULES.inventoryMaterials,
      6,
      (i) => ({
        sku: `SKU${i.toString().padStart(3, '0')}`,
        name: `${baseName}_${i}`,
        unit: 'UN',
        unit_cost: '10.50',
      }),
      '+ Crear'
    );

    // Test pagination
    await gotoPage(page, MODULES.inventoryMaterials);
    await expectTableVisible(page);
    const { btnAnterior, btnSiguiente } = await testPagination(page, MODULES.inventoryMaterials, 5);

    // Verify pagination navigation
    await testPaginationNavigation(page, btnAnterior, btnSiguiente);
  });
});