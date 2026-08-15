import { test, expect } from '@playwright/test';
import { login, gotoPage, verifyPageTitle } from './test-helpers';

test.describe('RBAC — Non-Admin Roles', () => {
  test.describe('Operator Role (1002)', () => {
    test.beforeEach(async ({ page }) => {
      await login(page, { employee_number: '1002', password: 'operador' });
    });

    test('Operator can access partners module', async ({ page }) => {
      await gotoPage(page, '/partners');
      await verifyPageTitle(page, 'Socios');
    });

    test('Operator can access customers module', async ({ page }) => {
      await gotoPage(page, '/customers');
      await verifyPageTitle(page, 'Clientes');
    });

    test('Operator can access operations module', async ({ page }) => {
      await gotoPage(page, '/operations/visits');
      await verifyPageTitle(page, 'Visitas');
    });

    test('Operator can access planning module', async ({ page }) => {
      await gotoPage(page, '/planning/daily-plans');
      await verifyPageTitle(page, 'Planes Diarios');
    });
  });

  test.describe('Technician Role (1003)', () => {
    test.beforeEach(async ({ page }) => {
      await login(page, { employee_number: '1003', password: 'tecnico' });
    });

    test('Technician can access operations module', async ({ page }) => {
      await gotoPage(page, '/operations/visits');
      await verifyPageTitle(page, 'Visitas');
    });

    test('Technician can access dashboard', async ({ page }) => {
      await gotoPage(page, '/');
      await verifyPageTitle(page, 'Dashboard');
    });
  });
});
