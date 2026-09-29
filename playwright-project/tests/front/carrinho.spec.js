import { test, expect } from '../../fixtures/index.js';
import { PRODUTOS } from '../../data/catalogo.js';

test.describe('Carrinho', () => {
  test.beforeEach(async ({ homePage }) => {
    await homePage.goto();
  });

  test('adicionar produto abre o carrinho com o item e o total', async ({ homePage, cartPage }) => {
    await homePage.adicionarAoCarrinho(PRODUTOS.panama.nome);

    await expect(cartPage.itemNomes).toHaveText([PRODUTOS.panama.nome]);
    await expect(cartPage.total).toHaveText('R$ 120.00');
  });

  test('aumentar a quantidade recalcula o total', async ({ homePage, cartPage }) => {
    await homePage.adicionarAoCarrinho(PRODUTOS.panama.nome);
    await cartPage.aumentarQuantidade(PRODUTOS.panama.nome);

    await expect(cartPage.total).toHaveText('R$ 240.00');
  });

  test('remover o único item deixa o carrinho vazio', async ({ homePage, cartPage }) => {
    await homePage.adicionarAoCarrinho(PRODUTOS.panama.nome);
    await cartPage.remover(PRODUTOS.panama.nome);

    await expect(cartPage.vazio).toBeVisible();
  });
});
