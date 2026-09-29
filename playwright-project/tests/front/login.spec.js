import { test, expect } from '../../fixtures/index.js';
import { gerarEmail } from '../../factories/usuarioFactory.js';

test.describe('Login', () => {
  test.beforeEach(async ({ loginPage }) => {
    await loginPage.goto();
  });

  test('login válido redireciona para a home com o usuário logado', async ({ page, loginPage, homePage, usuarioCadastrado }) => {
    await loginPage.login(usuarioCadastrado.email, usuarioCadastrado.password);

    await expect(loginPage.loginMessage).toHaveText('Login bem-sucedido! Redirecionando...');
    await page.waitForURL('**/index.html');
    await expect(homePage.logoutButton).toBeVisible();
  });

  test('login com usuário inexistente mostra erro', async ({ loginPage }) => {
    // E-mail único por teste: a API bloqueia um e-mail após 5 tentativas erradas.
    await loginPage.login(gerarEmail('inexistente'), 'SenhaErrada@123');

    await expect(loginPage.loginMessage).toHaveText('Erro no login: Não autorizado');
  });

  test('login com senha errada mostra erro', async ({ loginPage, usuarioCadastrado }) => {
    await loginPage.login(usuarioCadastrado.email, 'SenhaErrada@123');

    await expect(loginPage.loginMessage).toHaveText('Erro no login: Não autorizado');
  });
});
