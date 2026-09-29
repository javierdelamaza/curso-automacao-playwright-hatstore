import { test, expect } from '../../fixtures/index.js';
import { validarContrato } from '../../api/validarContrato.js';
import { listaDeChapeusSchema, loginSchema, freteSchema } from '../../api/schemas.js';

// Teste de contrato: não confere valores, confere o FORMATO da resposta.
// Se um campo sumir, mudar de tipo ou de nome, o teste quebra antes do front quebrar.
test.describe('API - contratos', () => {
  test('GET /api/hats respeita o contrato de lista de chapéus', async ({ api }) => {
    const corpo = await (await api.listarChapeus()).json();

    expect(validarContrato(listaDeChapeusSchema, corpo)).toEqual([]);
  });

  test('POST /auth/login respeita o contrato de login', async ({ api, usuarioCadastrado }) => {
    const corpo = await (await api.login(usuarioCadastrado.email, usuarioCadastrado.password)).json();

    expect(validarContrato(loginSchema, corpo)).toEqual([]);
  });

  test('GET /api/frete respeita o contrato de frete', async ({ api }) => {
    const corpo = await (await api.calcularFrete('01310-100', 1)).json();

    expect(validarContrato(freteSchema, corpo)).toEqual([]);
  });
});
