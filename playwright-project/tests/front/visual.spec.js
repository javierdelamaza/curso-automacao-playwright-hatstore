import { test, expect } from '../../fixtures/index.js';

// Na primeira execução não existe imagem de referência: rode `npm run test:visual:update`.
test.describe('Regressão visual @visual', () => {
  test('home com catálogo completo', async ({ page, homePage }) => {
    await homePage.goto();
    await page.waitForLoadState('networkidle');

    await expect(page).toHaveScreenshot('home.png', { fullPage: true });
  });

  test('tela de login', async ({ page, loginPage }) => {
    await loginPage.goto();

    await expect(page).toHaveScreenshot('login.png');
  });
});
