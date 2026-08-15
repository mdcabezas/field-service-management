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
  MODULES,
} from './test-helpers';

test.describe('Operations/Visits Module', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows visits', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await verifyPageTitle(page, 'Visitas');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
  });

  test('Create new visit flow', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await clickCreateButton(page, '+ Crear Visita');
    await verifyPageTitle(page, 'Crear Visita');

    const today = new Date().toISOString().split('T')[0];
    await fillInput(page, 'scheduled_at', today + 'T10:00');

    // Wait for visit types to load
    await page.waitForFunction(() => {
      const select = document.querySelector('select[id="type"]');
      return select && select.options.length > 1;
    });

    await fillSelect(page, 'type', 1);
    await fillSelect(page, 'status', 0); // "scheduled" is index 0
    await fillSelect(page, 'priority', 1); // "normal" is index 1 (after "low")
    await fillSelect(page, 'billing_to', 0); // "partner" is index 0
    await fillSelect(page, 'source', 0); // "email" is index 0

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/operations/visits');
    await waitForToast(page, 'Visita creada');
    await expectRowCount(page, 2);
  });

  test('View visit detail page', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Visita -');
  });

  test('Delete visit flow', async ({ page }) => {
    await gotoPage(page, MODULES.operationsVisits);

    await clickCreateButton(page, '+ Crear Visita');

    const today = new Date().toISOString().split('T')[0];
    await fillInput(page, 'scheduled_at', today + 'T10:00');

    // Wait for visit types to load
    await page.waitForFunction(() => {
      const select = document.querySelector('select[id="type"]');
      return select && select.options.length > 1;
    });

    await fillSelect(page, 'type', 1);
    await fillSelect(page, 'status', 0); // "scheduled" is index 0
    await fillSelect(page, 'priority', 1); // "normal" is index 1 (after "low")
    await fillSelect(page, 'billing_to', 0); // "partner" is index 0
    await fillSelect(page, 'source', 0); // "email" is index 0

    // Capture the visit ID from the creation response
    let visitId: string;
    const responsePromise = page.waitForResponse(
      (resp) => resp.url().includes('/api/visits') && resp.request().method() === 'POST'
    );
    await submitForm(page, 'Crear');
    const response = await responsePromise;
    const responseBody = await response.json();
    visitId = responseBody.id;

    await waitForRedirect(page, '**/operations/visits');
    await waitForToast(page, 'Visita creada');

    // Navigate directly to the created visit's detail page
    await page.goto(`http://localhost:3000/operations/visits/${visitId}`);
    await verifyPageTitle(page, 'Visita -');

    // Handle confirm dialog and delete
    page.on('dialog', dialog => dialog.accept());
    await page.click('button:has-text("Eliminar")');
    await waitForRedirect(page, '**/operations/visits');
    await waitForToast(page, 'Visita eliminada');
    await expectRowCount(page, 1);
  });
});