// Bloco 2: primeiro teste "cru", sem Page Object. No Bloco 3 ele é refatorado em busca.spec.js.
import { test, expect } from '@playwright/test';

test('busca por "Panamá" retorna exatamente um chapéu', async ({ page }) => {
  await page.goto('/');

  await page.getByPlaceholder('Buscar chapéu...').fill('Panamá');

  await expect(page.locator('.hat-card')).toHaveCount(1);
  await expect(page.locator('.hat-card h3')).toHaveText('Chapéu Panamá');
});
