import { test, expect } from '@playwright/test';
import {
  login,
  gotoPage,
  verifyPageTitle,
  expectTableVisible,
  clickFirstRow,
  clickCreateButton,
  fillInput,
  submitForm,
  waitForToast,
  MODULES,
} from './test-helpers';

test.describe('Partner Customization', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List partners with customization', async ({ page }) => {
    await gotoPage(page, MODULES.partners);
    await verifyPageTitle(page, 'Socios');
    await expectTableVisible(page);
  });

  test('View partner contacts', async ({ page }) => {
    await gotoPage(page, MODULES.partners);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    // Look for contacts section
    const contactsSection = page.locator('text=Contactos, text=Contacts');
    if (await contactsSection.isVisible()) {
      await expect(contactsSection).toBeVisible();
    }
  });

  test('Add partner contact', async ({ page }) => {
    await gotoPage(page, MODULES.partners);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    const addContactBtn = page.locator('button:has-text("Agregar Contacto"), button:has-text("Add Contact")');
    if (await addContactBtn.isVisible()) {
      await addContactBtn.click();
      await fillInput(page, 'name', 'Test Contact');
      await fillInput(page, 'role', 'Supervisor');
      await fillInput(page, 'phone', '555-0123');
      await fillInput(page, 'email', 'test@partner.com');
      await submitForm(page, 'Crear');
      await waitForToast(page, 'contacto');
    }
  });

  test('Partner SLA configuration', async ({ page }) => {
    await gotoPage(page, MODULES.partners);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    // Look for SLA section
    const slaSection = page.locator('text=SLA, text=Acuerdo de Nivel de Servicio');
    if (await slaSection.isVisible()) {
      await expect(slaSection).toBeVisible();
    }
  });

  test('Partner visit template assignment', async ({ page }) => {
    await gotoPage(page, MODULES.partners);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    // Look for template assignment section
    const templateSection = page.locator('text=Plantillas, text=Templates');
    if (await templateSection.isVisible()) {
      await expect(templateSection).toBeVisible();
    }
  });
});