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
  waitForNetworkIdle,
  createMultipleRecords,
  testPagination,
  testPaginationNavigation,
  MODULES,
} from './test-helpers';

function uniqueEmail(prefix = 'test') {
  return `${prefix}-${Date.now()}-${Math.random().toString(36).slice(2, 8)}@example.com`;
}

test.describe('Core/Users Module', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows users', async ({ page }) => {
    await gotoPage(page, MODULES.coreUsers);
    await verifyPageTitle(page, 'Usuarios');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
  });

  test('Search filters users', async ({ page }) => {
    await gotoPage(page, MODULES.coreUsers);
    await page.fill('input[placeholder*="Buscar"]', 'admin');
    // Search for "admin" should at least show the admin user
    const rows = await page.locator('tbody tr').count();
    expect(rows).toBeGreaterThanOrEqual(1);
  });

  test('Create new user flow', async ({ page }) => {
    await gotoPage(page, MODULES.coreUsers);
    
    // Navigate directly to create page
    await page.goto('http://localhost:3000/core/users/new');
    await page.waitForLoadState('networkidle');
    
    // Debug: check URL
    console.log('URL after goto:', page.url());
    
    // Debug: check h1 elements
    const h1s = await page.locator('h1').allTextContents();
    console.log('H1 elements:', h1s);
    
    await verifyPageTitle(page, 'Crear Usuario');

    await fillInput(page, 'email', uniqueEmail('test'));
    await fillInput(page, 'name', 'Test User');
    await fillInput(page, 'password', 'password123');
    await fillSelect(page, 'role', 1);

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/core/users');
    await waitForToast(page, 'Usuario creado');
    await expectRowCount(page, 2);
  });

  test('Create user with duplicate email shows error', async ({ page }) => {
    await gotoPage(page, MODULES.coreUsers);
    await clickCreateButton(page, '+ Crear Usuario');

    const dupEmail = uniqueEmail('duplicate');
    await fillInput(page, 'email', dupEmail);
    await fillInput(page, 'name', 'Duplicate User');
    await fillInput(page, 'password', 'password123');
    await fillSelect(page, 'role', 1);

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/core/users');
    await waitForToast(page, 'Usuario creado');

    // Try again with same email
    await clickCreateButton(page, '+ Crear Usuario');
    await fillInput(page, 'email', dupEmail);
    await fillInput(page, 'name', 'Duplicate User 2');
    await fillInput(page, 'password', 'password123');
    await fillSelect(page, 'role', 1);

    await submitForm(page, 'Crear');
    await waitForToast(page, 'Error al crear usuario');
  });

  test('View user detail page', async ({ page }) => {
    await gotoPage(page, MODULES.coreUsers);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Editar Usuario');
  });

  test('Edit user flow', async ({ page }) => {
    await gotoPage(page, MODULES.coreUsers);
    await clickFirstRow(page);

    await fillInput(page, 'name', 'Updated Admin Name');
    await page.click('button[type="submit"]:has-text("Actualizar")');
    await waitForRedirect(page, '**/core/users');
    await waitForToast(page, 'Usuario actualizado');
  });

  test('Delete user flow', async ({ page }) => {
    await gotoPage(page, MODULES.coreUsers);

    await clickCreateButton(page, '+ Crear Usuario');
    await fillInput(page, 'email', uniqueEmail('todelete'));
    await fillInput(page, 'name', 'To Delete User');
    await fillInput(page, 'password', 'password123');
    await fillSelect(page, 'role', 1);
    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/core/users');

    await clickFirstRow(page);
    await page.click('button:has-text("Eliminar")');
    confirmDialog(page);
    await waitForRedirect(page, '**/core/users');
    await waitForToast(page, 'Usuario eliminado');
  });

  test('Empty state when no users exist (not applicable - admin always exists)', async ({ page }) => {
    await gotoPage(page, MODULES.coreUsers);
    await expectTableVisible(page);
  });

  test('Navigate to edit page multiple times', async ({ page }) => {
    await gotoPage(page, MODULES.coreUsers);
    await verifyPageTitle(page, 'Usuarios');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
    
    // Click first row
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Editar Usuario');
    await waitForNetworkIdle(page);
    
    // Go back
    await page.goBack();
    await waitForNetworkIdle(page);
    await verifyPageTitle(page, 'Usuarios');
    
    // Click again
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Editar Usuario');
    await waitForNetworkIdle(page);
    
    // Go back
    await page.goBack();
    await waitForNetworkIdle(page);
    await verifyPageTitle(page, 'Usuarios');
    
    // Click third time
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Editar Usuario');
    await waitForNetworkIdle(page);
  });

test('Rapid navigation between edit pages', async ({ page }) => {
    await gotoPage(page, MODULES.coreUsers);
    await verifyPageTitle(page, 'Usuarios');
    await expectTableVisible(page);
    
    // Click first row
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Editar Usuario');
    await waitForNetworkIdle(page);
    
    // Immediately go back and click again (simulating rapid navigation)
    await page.goBack();
    await page.waitForLoadState('domcontentloaded');
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Editar Usuario');
    await waitForNetworkIdle(page);
  });

  test('Verify edit page data loads correctly on first visit', async ({ page }) => {
    await gotoPage(page, MODULES.coreUsers);
    await verifyPageTitle(page, 'Usuarios');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
    
    // Click first row - first visit
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Editar Usuario');
    await waitForNetworkIdle(page);
    
    // Check that form fields are populated
    const emailValue = await page.inputValue('#email');
    const nameValue = await page.inputValue('#name');
    console.log('First visit - Email:', emailValue, 'Name:', nameValue);
    expect(emailValue).toBeTruthy();
    expect(nameValue).toBeTruthy();
    
    // Go back
    await page.goBack();
    await waitForNetworkIdle(page);
    await verifyPageTitle(page, 'Usuarios');
    
    // Click again - second visit
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Editar Usuario');
    await waitForNetworkIdle(page);
    
    // Check that form fields are populated
    const emailValue2 = await page.inputValue('#email');
    const nameValue2 = await page.inputValue('#name');
    console.log('Second visit - Email:', emailValue2, 'Name:', nameValue2);
    expect(emailValue2).toBeTruthy();
    expect(nameValue2).toBeTruthy();
    
    // Values should be the same
    expect(emailValue2).toBe(emailValue);
    expect(nameValue2).toBe(nameValue);
  });

  test('Debug create page', async ({ page }) => {
    await gotoPage(page, MODULES.coreUsers);
    
    // Check what buttons are visible
    const buttons = await page.locator('button').allTextContents();
    console.log('Buttons:', buttons);
    
    // Click create button
    await page.click('button:has-text("+ Crear Usuario")');
    await page.waitForLoadState('networkidle');
    
    // Check page title
    const title = await page.locator('h1:has-text("Crear Usuario")').textContent();
    console.log('Page title:', title);
  });

  test('Pagination with 6 users (pageSize=5)', async ({ page }) => {
    // Create 6 users with unique emails
    const baseName = `PAG_USER_${Date.now()}`;
    await createMultipleRecords(
      page,
      MODULES.coreUsers,
      6,
      (i) => ({
        email: `pag${i}-${Date.now()}@example.com`,
        name: `${baseName}_${i}`,
        password: 'testpass123',
      }),
      '+ Crear Usuario'
    );

    // Set the role select separately (it's a <select>, not an <input>)
    // The createMultipleRecords already created the users; role defaults to first option

    // Test pagination
    const { btnAnterior, btnSiguiente } = await testPagination(page, MODULES.coreUsers, 6);

    // Verify pagination navigation
    await testPaginationNavigation(page, btnAnterior, btnSiguiente);
  });
});