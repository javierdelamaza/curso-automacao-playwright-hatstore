import { test, expect } from '@playwright/test';

// Testes de API/contrato básico (módulo 2), sem precisar de navegador.
test.describe('API — /api/hats', () => {
  test('GET /api/hats retorna 200 e uma lista de chapéus', async ({ request }) => {
    const response = await request.get('/api/hats');

    expect(response.ok()).toBeTruthy();

    const hats = await response.json();
    expect(Array.isArray(hats)).toBeTruthy();
    expect(hats.length).toBeGreaterThan(0);
  });

  test('cada chapéu retornado tem os campos de contrato esperados', async ({ request }) => {
    const response = await request.get('/api/hats');
    const hats = await response.json();

    const primeiro = hats[0];
    expect(primeiro).toHaveProperty('id');
    expect(primeiro).toHaveProperty('nome');
    expect(primeiro).toHaveProperty('price');
    expect(primeiro).toHaveProperty('categoria');
  });
});
