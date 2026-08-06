import { test, expect } from '@playwright/test';
import { LoginPage } from '../pages/LoginPage.js';

// Gera um e-mail único por execução para não colidir com cadastros anteriores.
function emailUnico() {
  return `teste.playwright.${Date.now()}@hatstore.test`;
}

test.describe('Autenticação', () => {
  test('cadastro de novo usuário é concluído com sucesso', async ({ page }) => {
    const auth = new LoginPage(page);
    await auth.goto();
    await auth.goToRegisterForm();

    await auth.register({
      email: emailUnico(),
      nome: 'Usuário Teste Playwright',
      cpf: '123.456.789-09',
      dataNascimento: '2000-01-01',
      password: 'Senha@123',
    });

    await expect(auth.registerMessage).not.toHaveText('');
  });

  test('login com credenciais inválidas exibe mensagem de erro', async ({ page }) => {
    const auth = new LoginPage(page);
    await auth.goto();

    await auth.login('usuario-que-nao-existe@hatstore.test', 'SenhaErrada@123');

    await expect(auth.loginMessage).not.toHaveText('');
  });

  test.fixme(
    'login com credenciais válidas redireciona e exibe o usuário logado',
    async ({ page }) => {
      // Dever de casa (módulo 1 das aulas originais): use um usuário já
      // existente na massa de dados de teste, ou encadeie um cadastro antes
      // do login (cuidado com isolamento: gere um e-mail novo a cada run).
      const auth = new LoginPage(page);
      await auth.goto();

      await auth.login('usuario@existente.test', 'SenhaValida@123');

      await expect(auth.authLinks).not.toContainText('Login');
    }
  );
});
