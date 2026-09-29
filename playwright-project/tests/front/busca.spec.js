import { test, expect } from '../../fixtures/index.js';
import { PRODUTOS, TOTAL_PRODUTOS_CATALOGO } from '../../data/catalogo.js';

test.describe('Home - busca de produtos', () => {
  test.beforeEach(async ({ homePage }) => {
    await homePage.goto();
  });

  test('catálogo exibe todos os produtos ao abrir a home', async ({ homePage }) => {
    await expect(homePage.hatCards).toHaveCount(TOTAL_PRODUTOS_CATALOGO);
  });

  test('busca por produto existente retorna só ele', async ({ homePage }) => {
    await homePage.search('Panamá');

    await expect(homePage.hatCards).toHaveCount(1);
    await expect(homePage.hatCardTitles).toHaveText([PRODUTOS.panama.nome]);
  });

  test('busca não diferencia maiúsculas de minúsculas', async ({ homePage }) => {
    await homePage.search('FEDORA');

    await expect(homePage.hatCardTitles).toHaveText([PRODUTOS.fedora.nome]);
  });

  test('busca por termo inexistente mostra mensagem de vazio', async ({ homePage }) => {
    await homePage.search('produto-que-nao-existe');

    await expect(homePage.hatCards).toHaveCount(0);
    await expect(homePage.nenhumResultado).toBeVisible();
  });

  test('produto sem estoque aparece como indisponível e sem botão de compra', async ({ homePage }) => {
    const card = homePage.cardDoProduto(PRODUTOS.porkPieSemEstoque.nome);

    await expect(card.getByText('Indisponível no momento')).toBeVisible();
    await expect(card.getByRole('button', { name: 'Adicionar ao carrinho' })).toHaveCount(0);
  });
});
