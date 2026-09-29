import { test, expect } from '../../fixtures/index.js';
import { CATEGORIAS, PRODUTOS } from '../../data/catalogo.js';

test.describe('Home - cenários combinados de filtro e busca', () => {
  test.beforeEach(async ({ homePage }) => {
    await homePage.goto();
  });

  test('categoria + faixa de preço retorna só produtos que atendem aos dois filtros', async ({ homePage }) => {
    await homePage.categoriaCheckbox(CATEGORIAS.nacional).check();
    await homePage.filterByPreco({ min: 100, max: 130 });

    await expect(homePage.hatCardTitles).toHaveText(['Chapéu Gaúcho']);
  });

  test('duas categorias somadas mostram produtos de ambas', async ({ homePage }) => {
    await homePage.filterByCategoria(CATEGORIAS.nacional, CATEGORIAS.crescer);

    await expect(homePage.cardDoProduto(PRODUTOS.gaucho.nome)).toBeVisible();
    await expect(homePage.cardDoProduto(PRODUTOS.gustavoCarvalho.nome)).toBeVisible();
    await expect(homePage.cardDoProduto(PRODUTOS.panama.nome)).toHaveCount(0);
  });

  test('busca funciona sobre o resultado filtrado', async ({ homePage }) => {
    await homePage.filterByCategoria(CATEGORIAS.importado);
    await homePage.search('Fedora');

    await expect(homePage.hatCardTitles).toHaveText([PRODUTOS.fedora.nome]);
  });
});
