import { test, expect } from '@playwright/test';
import {
  login,
  logout,
  gotoPage,
  waitForToast,
  verifyPageTitle,
  waitForRedirect,
  expectTableVisible,
  expectRowCount,
  clickFirstRow,
  TEST_CREDENTIALS,
  MODULES,
} from './test-helpers';

test.describe('Auth Module', () => {
  test('Valid login redirects to dashboard', async ({ page }) => {
    await login(page);
    await verifyPageTitle(page, 'LOCALIS FSM');
  });

  test('Invalid credentials shows error', async ({ page }) => {
    await page.goto('http://localhost:3000/login');
    await page.fill('#employeeNumber', TEST_CREDENTIALS.invalid.employee_number);
    await page.fill('#password', TEST_CREDENTIALS.invalid.password);
    await page.click('button[type="submit"]');
    await waitForToast(page, 'Credenciales inválidas');
  });

  test('Login with empty fields shows validation errors', async ({ page }) => {
    await page.goto('http://localhost:3000/login');
    await page.click('button[type="submit"]');
    await expect(page.locator('text=Número de empleado requerido')).toBeVisible();
    await expect(page.locator('text=Contraseña requerida')).toBeVisible();
  });

  test('Logout works correctly', async ({ page }) => {
    await login(page);
    await logout(page);
    await expect(page).toHaveURL(/\/login/);
  });

  test('Protected routes redirect to login when unauthenticated', async ({ page }) => {
    await page.goto('http://localhost:3000/core/users');
    await expect(page).toHaveURL(/\/login/);
  });
});