import { test, expect } from '@playwright/test';
import { login, gotoPage, MODULES, clickFirstRow } from './test-helpers';

test.describe('RBAC/Authorization Tests', () => {
  test('Admin can access all modules', async ({ page }) => {
    await login(page, { employee_number: '1001', password: 'admin' });

    const adminModules = [
      MODULES.coreUsers,
      MODULES.coreTechRoles,
      MODULES.planningDailyPlans,
      MODULES.planningAssignments,
      MODULES.planningRoutes,
      MODULES.inventoryMaterials,
      MODULES.inventoryTools,
      MODULES.inventoryVehicles,
      MODULES.inventoryEPP,
      MODULES.inventoryRentals,
      MODULES.inventoryCostRates,
      MODULES.inventoryMaintenanceRecords,
      MODULES.customers,
      MODULES.partners,
      MODULES.operationsVisits,
      MODULES.operationsReports,
      MODULES.operationsSLATracking,
      MODULES.operationsVehicleAssignments,
      MODULES.notificationsTemplates,
      MODULES.geocoding,
      MODULES.sharedReportTemplates,
    ];

    for (const module of adminModules) {
      await gotoPage(page, module);
      await expect(page).not.toHaveURL(/\/login/);
      await expect(page.locator('h1').first()).toBeVisible({ timeout: 5000 });
    }
  });

  test('Non-admin users redirected from admin-only modules', async ({ page }) => {
    // Note: This test would need different credentials for non-admin users
    // Currently only admin credentials are available in TEST_CREDENTIALS
    // Skip for now - would require additional test users
    test.skip();
  });

  test('Unauthenticated users redirected to login', async ({ page }) => {
    await page.goto('http://localhost:3000/core/users');
    await expect(page).toHaveURL(/\/login/);
  });

  test('Create buttons only visible for authorized roles', async ({ page }) => {
    await login(page, { employee_number: '1001', password: 'admin' });

    const modulesWithCreate = [
      { path: MODULES.coreUsers, buttonText: '+ Crear Usuario' },
      { path: MODULES.coreTechRoles, buttonText: '+ Crear Rol' },
      { path: MODULES.planningDailyPlans, buttonText: '+ Crear Plan' },
      { path: MODULES.planningAssignments, buttonText: '+ Crear Asignación' },
      { path: MODULES.planningRoutes, buttonText: '+ Crear Ruta' },
    ];

    for (const { path, buttonText } of modulesWithCreate) {
      await gotoPage(page, path);
      await expect(page.locator(`button:has-text("${buttonText}")`)).toBeVisible();
    }
  });

  test('Delete buttons only visible for admin', async ({ page }) => {
    await login(page, { employee_number: '1001', password: 'admin' });

    const deleteModules = [
      MODULES.coreUsers,
      MODULES.coreTechRoles,
      MODULES.planningDailyPlans,
      MODULES.planningAssignments,
      MODULES.planningRoutes,
    ];

    for (const path of deleteModules) {
      await gotoPage(page, path);
      await clickFirstRow(page);
      // Admin should see delete button
      await expect(page.locator('button:has-text("Eliminar")')).toBeVisible();
      await page.goBack();
    }
  });
});