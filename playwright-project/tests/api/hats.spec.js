import { test, expect } from '../../fixtures/index.js';
import { CATEGORIAS, TOTAL_PRODUTOS_CATALOGO } from '../../data/catalogo.js';

test.describe('API - GET /api/hats', () => {
  test('lista o catálogo completo', async ({ api }) => {
    const response = await api.listarChapeus();

    expect(response.status()).toBe(200);
    const chapeus = await response.json();
    expect(chapeus).toHaveLength(TOTAL_PRODUTOS_CATALOGO);
  });

  test('filtra por categoria', async ({ api }) => {
    const response = await api.listarChapeus({ categoria: CATEGORIAS.crescer });

    const chapeus = await response.json();
    expect(chapeus.length).toBeGreaterThan(0);
    for (const chapeu of chapeus) {
      expect(chapeu.categoria).toBe(CATEGORIAS.crescer);
    }
  });

  test('filtra por faixa de preço', async ({ api }) => {
    const response = await api.listarChapeus({ min: 100, max: 130 });

    const chapeus = await response.json();
    expect(chapeus.length).toBeGreaterThan(0);
    for (const chapeu of chapeus) {
      expect(chapeu.price).toBeGreaterThanOrEqual(100);
      expect(chapeu.price).toBeLessThanOrEqual(130);
    }
  });

  test('produto sem estoque vem com temEstoque = false', async ({ api }) => {
    const chapeus = await (await api.listarChapeus()).json();

    const semEstoque = chapeus.filter((c) => c.quantidade === 0);
    expect(semEstoque.length).toBeGreaterThan(0);
    for (const chapeu of semEstoque) {
      expect(chapeu.temEstoque).toBe(false);
    }
  });
});
