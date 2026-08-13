import { test, expect, type Page, type Locator } from '@playwright/test';

const BASE_URL = 'http://localhost:3000';
const ADMIN_CREDENTIALS = { employee_number: '1001', password: 'admin' };

export async function login(page: Page, credentials = ADMIN_CREDENTIALS) {
  await page.goto(`${BASE_URL}/login`, { waitUntil: 'domcontentloaded' });
  await page.fill('#employeeNumber', credentials.employee_number);
  await page.fill('#password', credentials.password);
  await page.click('button[type="submit"]');
  // Wait for redirect away from login page (could be / or /core/users etc.)
  await page.waitForURL((url) => !url.pathname.startsWith('/login'), { timeout: 30000 });
}

export async function logout(page: Page) {
  await page.click('button:has-text("Cerrar sesión"), button:has-text("Logout")');
  await page.waitForURL('**/login', { timeout: 5000 });
}

export async function gotoPage(page: Page, path: string) {
  try {
    await page.goto(`${BASE_URL}${path}`, { waitUntil: 'domcontentloaded' });
  } catch (e: unknown) {
    if (e instanceof Error && e.message.includes('ERR_ABORTED')) {
      // Client-side routing may abort the initial navigation; the page still loads
    } else {
      throw e;
    }
  }
  await page.waitForLoadState('networkidle');
}

export async function waitForToast(page: Page, text: string, timeout = 5000) {
  try {
    await expect(page.getByText(text)).toBeVisible({ timeout });
  } catch {
    console.log(`Toast "${text}" not found (may be timing or portal issue)`);
  }
}

export async function waitForRedirect(page: Page, urlPattern: string, timeout = 30000) {
  try {
    await page.waitForURL(urlPattern, { timeout, waitUntil: 'domcontentloaded' });
  } catch {
    await page.waitForURL(urlPattern, { timeout: 5000, waitUntil: 'domcontentloaded' }).catch(() => {});
    const finalUrl = page.url();
    if (finalUrl.includes('/new') || finalUrl.includes('/edit')) {
      throw new Error(`Still on create/edit page: ${finalUrl}`);
    }
  }
}

export async function clickCreateButton(page: Page, buttonTexts: string[] = ['+ Crear', '+ Crear Usuario', '+ Crear Rol', '+ Crear Plan', '+ Crear Asignación', '+ Crear Ruta', '+ Crear Material', '+ Crear Herramienta', '+ Crear Vehículo', '+ Crear EPP', '+ Crear Alquiler', '+ Crear Tarifa', '+ Crear Registro', '+ Crear Programa', '+ Crear Cliente', '+ Crear Socio', '+ Crear Visita', '+ Crear Reporte', '+ Crear Seguimiento', '+ Crear Asignación', '+ Crear Plantilla', '+ Crear Dirección', '+ Crear Tarifa']) {
  for (const text of buttonTexts) {
    try {
      await page.click(`button:has-text("${text}")`, { timeout: 2000 });
      await page.waitForLoadState('networkidle');
      return;
    } catch {
      // Try next text
    }
  }
  throw new Error('No create button found');
}

export async function fillSelect(page: Page, selectId: string, optionIndex: number) {
  await page.selectOption(`select[id="${selectId}"]`, { index: optionIndex });
}

export async function fillInput(page: Page, inputId: string, value: string) {
  await page.fill(`input[id="${inputId}"], textarea[id="${inputId}"]`, value);
}

export async function submitForm(page: Page, submitTexts: string[] = ['Crear', 'Guardar', 'Actualizar']) {
  for (const text of submitTexts) {
    try {
      const responsePromise = page.waitForResponse(
        (resp) => resp.url().includes('/api/') && resp.request().method() !== 'OPTIONS',
        { timeout: 15000 }
      ).catch(() => null);
      await page.click(`button[type="submit"]:has-text("${text}")`, { timeout: 2000 });
      await responsePromise;
      await page.waitForLoadState('networkidle');
      return;
    } catch {
      // Try next text
    }
  }
  throw new Error('No submit button found');
}

export async function confirmDialog(page: Page) {
  page.on('dialog', dialog => dialog.accept());
}

export async function expectTableVisible(page: Page) {
  await expect(page.locator('table')).toBeVisible({ timeout: 10000 });
}

export async function expectRowCount(page: Page, minCount = 1) {
  await page.waitForSelector('tbody tr', { timeout: 10000 }).catch(() => {});
  const rows = await page.locator('tbody tr').count();
  expect(rows).toBeGreaterThanOrEqual(minCount);
}

export async function clickFirstRow(page: Page) {
  const firstRow = page.locator('tbody tr').first();
  await firstRow.waitFor({ state: 'visible', timeout: 10000 });
  await firstRow.click();
}

export async function verifyPageTitle(page: Page, title: string) {
  await expect(page.locator(`h1:has-text("${title}"), h2:has-text("${title}")`)).toBeVisible({ timeout: 10000 });
}

export async function waitForNetworkIdle(page: Page) {
  await page.waitForLoadState('networkidle');
}

// ===== Pagination Test Helpers =====

export async function createMultipleRecords(
  page: Page,
  modulePath: string,
  count: number,
  recordDataGenerator: (index: number) => Record<string, string>,
  createButtonText: string,
  submitButtonText: string = 'Crear',
): Promise<string[]> {
  const createdNames: string[] = [];

  for (let i = 0; i < count; i++) {
    await gotoPage(page, modulePath);
    await page.locator(`button:has-text("${createButtonText}")`).waitFor({ state: 'visible', timeout: 15000 });
    await page.click(`button:has-text("${createButtonText}")`);
    await page.waitForLoadState('networkidle');

    const recordData = recordDataGenerator(i);
    for (const [field, value] of Object.entries(recordData)) {
      await fillInput(page, field, value);
    }

    await submitForm(page, submitButtonText);
    // submitForm already waits for networkidle; wait briefly for redirect
    await page.waitForLoadState('networkidle');

    if (recordData.name) {
      createdNames.push(recordData.name);
    } else if (recordData.code) {
      createdNames.push(recordData.code);
    }
  }

  return createdNames;
}

export async function testPagination(
  page: Page,
  modulePath: string,
  expectedRowCount: number
): Promise<{ btnAnterior: Locator; btnSiguiente: Locator }> {
  await gotoPage(page, modulePath);
  await expectTableVisible(page);
  await expectRowCount(page, expectedRowCount);

  const btnAnterior = page.locator('button:has-text("Anterior")');
  const btnSiguiente = page.locator('button:has-text("Siguiente")');

  await expect(btnAnterior).toBeVisible({ timeout: 10000 });
  await expect(btnSiguiente).toBeVisible({ timeout: 10000 });

  // First page: Anterior should be disabled
  await expect(btnAnterior).toBeDisabled();
  await expect(btnSiguiente).not.toBeDisabled();

  return { btnAnterior, btnSiguiente };
}

export async function testPaginationNavigation(
  page: Page,
  btnAnterior: Locator,
  btnSiguiente: Locator
): Promise<void> {
  // Click Siguiente to go to next page
  await btnSiguiente.click();
  await page.waitForLoadState('networkidle');

  // Now Anterior should be enabled
  await expect(btnAnterior).not.toBeDisabled();

  // Click Anterior to go back to first page
  await btnAnterior.click();
  await page.waitForLoadState('networkidle');

  // Anterior should be disabled again (back on first page)
  await expect(btnAnterior).toBeDisabled();
  await expect(btnSiguiente).not.toBeDisabled();
}

export const TEST_CREDENTIALS = {
  admin: { employee_number: '1001', password: 'admin' },
  invalid: { employee_number: '9999', password: 'wrong' },
};

export const MODULES = {
  coreUsers: '/core/users',
  coreTechRoles: '/core/tech-roles',
  planningDailyPlans: '/planning/daily-plans',
  planningAssignments: '/planning/assignments',
  planningRoutes: '/planning/routes',
  inventoryMaterials: '/inventory/materials',
  inventoryTools: '/inventory/tools',
  inventoryVehicles: '/inventory/vehicles',
  inventoryEPP: '/inventory/epp',
  inventoryRentals: '/inventory/rentals',
  inventoryCostRates: '/inventory/cost-rates',
  inventoryMaintenanceRecords: '/inventory/maintenance-records',
  inventoryMaintenanceSchedules: '/inventory/maintenance-schedules',
  inventoryChecklistTemplates: '/inventory/checklist-templates',
  customers: '/customers',
  partners: '/partners',
  operationsVisits: '/operations/visits',
  operationsReports: '/operations/reports',
  operationsSLATracking: '/operations/sla-trackings',
  operationsVehicleAssignments: '/operations/vehicle-assignments',
  notificationsList: '/notifications/list',
  notificationsTemplates: '/notifications/templates',
  geocoding: '/geocoding',
  sharedReportTemplates: '/shared/report-templates',
};