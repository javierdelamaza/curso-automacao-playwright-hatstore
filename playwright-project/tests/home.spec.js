import { test, expect } from '@playwright/test';
import { HomePage } from '../pages/HomePage.js';

test.describe('Home — busca de produtos', () => {
  test('busca por um produto existente retorna exatamente o resultado esperado', async ({ page }) => {
    const home = new HomePage(page);
    await home.goto();

    await home.search('Panamá');

    await expect(home.hatCards).toHaveCount(1);
    await expect(home.hatCardTitles.first()).toContainText('Panamá');
  });

  test('busca por um termo inexistente não retorna nenhum produto', async ({ page }) => {
    const home = new HomePage(page);
    await home.goto();

    await home.search('produto-que-nao-existe-123');

    await expect(home.hatCards).toHaveCount(0);
  });

  test('catálogo carrega com produtos ao acessar a home', async ({ page }) => {
    const home = new HomePage(page);
    await home.goto();

    await expect(home.hatCards.first()).toBeVisible();
  });
});
