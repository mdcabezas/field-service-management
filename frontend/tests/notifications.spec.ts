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
  MODULES,
} from './test-helpers';

test.describe('Notifications/Templates Module', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows notification templates', async ({ page }) => {
    await gotoPage(page, MODULES.notificationsTemplates);
    await verifyPageTitle(page, 'Plantillas de Notificación');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
  });

  test('Create new notification template flow', async ({ page }) => {
    await gotoPage(page, MODULES.notificationsTemplates);
    await clickCreateButton(page, '+ Crear Plantilla');
    await verifyPageTitle(page, 'Crear Plantilla de Notificación');

    await fillSelect(page, 'type', 1);
    await fillSelect(page, 'channel', 1);
    await fillInput(page, 'subject', 'Test Subject');
    await fillInput(page, 'body', 'Test body content');

    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/notifications/templates');
    await waitForToast(page, 'Plantilla de notificación creada');
    await expectRowCount(page, 2);
  });

  test('View notification template detail page', async ({ page }) => {
    await gotoPage(page, MODULES.notificationsTemplates);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Detalle de Plantilla de Notificación');
  });

  test('Edit notification template flow', async ({ page }) => {
    await gotoPage(page, MODULES.notificationsTemplates);
    await clickFirstRow(page);

    await fillSelect(page, 'type', 1);
    await page.click('button[type="submit"]:has-text("Guardar")');
    await waitForRedirect(page, '**/notifications/templates');
    await waitForToast(page, 'Plantilla de notificación actualizada');
  });

  test('Delete notification template flow', async ({ page }) => {
    await gotoPage(page, MODULES.notificationsTemplates);

    await clickCreateButton(page, '+ Crear Plantilla');
    await fillSelect(page, 'type', 1);
    await fillSelect(page, 'channel', 1);
    await fillInput(page, 'subject', 'Delete Subject');
    await fillInput(page, 'body', 'Delete Body');
    await submitForm(page, 'Crear');
    await waitForRedirect(page, '**/notifications/templates');

    await clickFirstRow(page);
    await page.click('button:has-text("Eliminar")');
    await waitForRedirect(page, '**/notifications/templates');
    await waitForToast(page, 'Plantilla de notificación eliminada');
  });
});

test.describe('Notifications/List Module', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List page loads and shows notifications', async ({ page }) => {
    await gotoPage(page, MODULES.notificationsList);
    await verifyPageTitle(page, 'Notificaciones Enviadas');
    await expectTableVisible(page);
    await expectRowCount(page, 1);
  });

  test('View notification detail page', async ({ page }) => {
    await gotoPage(page, MODULES.notificationsList);
    await clickFirstRow(page);
    await verifyPageTitle(page, 'Detalle de Notificación');
  });
});