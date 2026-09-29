import { test, expect } from '../../fixtures/index.js';
import { novoUsuario, usuarioMenorDeIdade } from '../../factories/usuarioFactory.js';

test.describe('Cadastro', () => {
  test.beforeEach(async ({ loginPage }) => {
    await loginPage.goto();
    await loginPage.goToRegisterForm();
  });

  test('cadastro com dados válidos é concluído', async ({ loginPage }) => {
    await loginPage.register(novoUsuario());

    await expect(loginPage.registerMessage).toHaveText('Registro bem-sucedido! Você já pode fazer o login.');
  });

  test('cadastro de menor de 18 anos é recusado', async ({ loginPage }) => {
    await loginPage.register(usuarioMenorDeIdade());

    await expect(loginPage.registerMessage).toContainText('18 anos');
  });

  test('cadastro com e-mail já usado é recusado', async ({ loginPage, usuarioCadastrado }) => {
    await loginPage.register(novoUsuario({ email: usuarioCadastrado.email }));

    await expect(loginPage.registerMessage).not.toHaveText('Registro bem-sucedido! Você já pode fazer o login.');
    await expect(loginPage.registerMessage).not.toBeEmpty();
  });
});
