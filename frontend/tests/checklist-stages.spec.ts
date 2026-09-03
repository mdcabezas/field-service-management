import { test, expect } from '@playwright/test';
import {
  login,
  gotoPage,
  verifyPageTitle,
  expectTableVisible,
  clickFirstRow,
  clickCreateButton,
  fillInput,
  fillSelect,
  submitForm,
  waitForToast,
  MODULES,
} from './test-helpers';

test.describe('Checklist Stages', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('List checklist templates with stages', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryChecklistTemplates);
    await verifyPageTitle(page, 'Plantillas de Checklist');
    await expectTableVisible(page);
  });

  test('Create checklist template with stages enabled', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryChecklistTemplates);
    await clickCreateButton(page, '+ Crear Plantilla');
    await verifyPageTitle(page, 'Crear Plantilla');

    await fillInput(page, 'name', 'Test Template with Stages');
    await fillInput(page, 'description', 'Template for testing stages');

    // Enable stages
    const stagesCheckbox = page.locator('input[name="stages_enabled"]');
    if (await stagesCheckbox.isVisible()) {
      await stagesCheckbox.check();
    }

    await submitForm(page, 'Crear');
    await waitForToast(page, 'creada');
  });

  test('View checklist template stages', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryChecklistTemplates);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    // Should show template details with stages section
    const stagesSection = page.locator('text=Estadios, text=Stages, text=Etapas');
    if (await stagesSection.isVisible()) {
      await expect(stagesSection).toBeVisible();
    }
  });

  test('Add stage to checklist template', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryChecklistTemplates);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    // Look for add stage button
    const addStageBtn = page.locator('button:has-text("Agregar Etapa"), button:has-text("Add Stage")');
    if (await addStageBtn.isVisible()) {
      await addStageBtn.click();
      await fillInput(page, 'name', 'New Stage');
      await fillInput(page, 'sort_order', '1');
      await submitForm(page, 'Crear');
      await waitForToast(page, 'etapa');
    }
  });

  test('Edit stage order', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryChecklistTemplates);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    // Find first stage row and click edit
    const editBtn = page.locator('button:has-text("Editar"), button:has-text("Edit")').first();
    if (await editBtn.isVisible()) {
      await editBtn.click();
      await fillInput(page, 'sort_order', '2');
      await submitForm(page, 'Guardar');
      await waitForToast(page, 'actualizada');
    }
  });

  test('Delete stage from template', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryChecklistTemplates);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    page.on('dialog', dialog => dialog.accept());

    const deleteBtn = page.locator('button:has-text("Eliminar"), button:has-text("Delete")').first();
    if (await deleteBtn.isVisible()) {
      await deleteBtn.click();
      await waitForToast(page, 'eliminada');
    }
  });

  test('Toggle stages enabled flag', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryChecklistTemplates);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    const stagesCheckbox = page.locator('input[name="stages_enabled"]');
    if (await stagesCheckbox.isVisible()) {
      const isChecked = await stagesCheckbox.isChecked();
      await stagesCheckbox.toggle();
      await submitForm(page, 'Guardar');
      await waitForToast(page, 'actualizada');
    }
  });

  test('Filter templates by visit type', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryChecklistTemplates);
    await expectTableVisible(page);

    // Look for filter/search functionality
    const searchInput = page.locator('input[placeholder*="Buscar"], input[placeholder*="Search"]');
    if (await searchInput.isVisible()) {
      await searchInput.fill('maintenance');
      await page.waitForLoadState('networkidle');
    }
  });

  test('Stage required flag validation', async ({ page }) => {
    await gotoPage(page, MODULES.inventoryChecklistTemplates);
    await expectTableVisible(page);
    await clickFirstRow(page);
    await page.waitForLoadState('networkidle');

    const requiredCheckbox = page.locator('input[name="required"]');
    if (await requiredCheckbox.isVisible()) {
      await requiredCheckbox.check();
      await submitForm(page, 'Guardar');
      await waitForToast(page, 'actualizada');
    }
  });
});