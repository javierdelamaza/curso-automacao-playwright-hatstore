import { test as base, expect } from '@playwright/test';
import { HomePage } from '../pages/HomePage.js';
import { LoginPage } from '../pages/LoginPage.js';
import { CartPage } from '../pages/CartPage.js';
import { CheckoutPage } from '../pages/CheckoutPage.js';
import { PedidosPage } from '../pages/PedidosPage.js';
import { HatStoreApi } from '../api/HatStoreApi.js';
import { novoUsuario } from '../factories/usuarioFactory.js';

export const test = base.extend({
  // Page Objects prontos para uso
  homePage: async ({ page }, use) => use(new HomePage(page)),
  loginPage: async ({ page }, use) => use(new LoginPage(page)),
  cartPage: async ({ page }, use) => use(new CartPage(page)),
  checkoutPage: async ({ page }, use) => use(new CheckoutPage(page)),
  pedidosPage: async ({ page }, use) => use(new PedidosPage(page)),

  api: async ({ request }, use) => use(new HatStoreApi(request)),

  // Usuário novo e já cadastrado (via API, que é mais rápido que pela tela)
  usuarioCadastrado: async ({ api }, use) => {
    const usuario = novoUsuario();
    const response = await api.cadastrar(usuario);
    expect(response.status(), await response.text()).toBe(201);
    await use(usuario);
  },

  token: async ({ api, usuarioCadastrado }, use) => {
    await use(await api.obterToken(usuarioCadastrado.email, usuarioCadastrado.password));
  },

  // Cliente HTTP que já envia "Authorization: Bearer <token>" em toda chamada
  apiAutenticada: async ({ playwright, baseURL, token }, use) => {
    const contexto = await playwright.request.newContext({
      baseURL,
      extraHTTPHeaders: { Authorization: `Bearer ${token}` },
    });
    await use(new HatStoreApi(contexto));
    await contexto.dispose();
  },

  // Página já logada: grava o token no localStorage antes de carregar o site, sem passar pela tela de login
  paginaLogada: async ({ page, usuarioCadastrado, token }, use) => {
    await page.addInitScript(
      ({ jwt, email }) => {
        window.localStorage.setItem('jwt_token', jwt);
        window.localStorage.setItem('user_email', email);
      },
      { jwt: token, email: usuarioCadastrado.email },
    );
    await use(page);
  },
});

export { expect };
