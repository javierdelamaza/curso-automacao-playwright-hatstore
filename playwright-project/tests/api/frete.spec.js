import { test, expect } from '../../fixtures/index.js';

test.describe('API - GET /api/frete', () => {
  const casos = [
    { cep: '01310-100', regiao: 'Sudeste', quantidade: 1, valor: 12.5 },
    { cep: '90010-000', regiao: 'Sul', quantidade: 1, valor: 15.0 },
    { cep: '40010-000', regiao: 'Nordeste', quantidade: 1, valor: 22.0 },
  ];

  for (const caso of casos) {
    test(`CEP ${caso.cep} é da região ${caso.regiao}`, async ({ api }) => {
      const response = await api.calcularFrete(caso.cep, caso.quantidade);

      expect(response.status()).toBe(200);
      const frete = await response.json();
      expect(frete.regiao).toBe(caso.regiao);
      expect(frete.valor_frete).toBeCloseTo(caso.valor, 2);
    });
  }

  test('frete é progressivo: o 2º item custa 30% do valor base', async ({ api }) => {
    const frete = await (await api.calcularFrete('01310-100', 2)).json();

    expect(frete.valor_frete).toBeCloseTo(12.5 + 12.5 * 0.3, 2);
  });

  test('CEP em formato inválido devolve 400', async ({ api }) => {
    const response = await api.calcularFrete('123', 1);

    expect(response.status()).toBe(400);
  });
});
