import { test, expect } from '../../fixtures/index.js';
import { novoUsuario, usuarioMenorDeIdade, gerarEmail } from '../../factories/usuarioFactory.js';

test.describe('API - autenticação', () => {
  test('POST /auth/register cria usuário novo', async ({ api }) => {
    const response = await api.cadastrar(novoUsuario());

    expect(response.status()).toBe(201);
    expect(await response.json()).toEqual({ message: 'Usuário criado com sucesso!' });
  });

  test('POST /auth/register recusa e-mail duplicado com 409', async ({ api, usuarioCadastrado }) => {
    const response = await api.cadastrar(novoUsuario({ email: usuarioCadastrado.email }));

    expect(response.status()).toBe(409);
  });

  test('POST /auth/register recusa menor de idade com 400', async ({ api }) => {
    const response = await api.cadastrar(usuarioMenorDeIdade());

    expect(response.status()).toBe(400);
  });

  test('POST /auth/login devolve token para credenciais válidas', async ({ api, usuarioCadastrado }) => {
    const response = await api.login(usuarioCadastrado.email, usuarioCadastrado.password);

    expect(response.status()).toBe(200);
    const corpo = await response.json();
    expect(corpo.token).toBeTruthy();
    expect(corpo.email).toBe(usuarioCadastrado.email);
  });

  test('POST /auth/login devolve 401 para credenciais inválidas', async ({ api }) => {
    const response = await api.login(gerarEmail('inexistente'), 'SenhaErrada@123');

    expect(response.status()).toBe(401);
  });

  test('GET /api/me sem token devolve 401', async ({ api }) => {
    const response = await api.perfil();

    expect(response.status()).toBe(401);
  });

  test('GET /api/me com token devolve o perfil do usuário', async ({ apiAutenticada, usuarioCadastrado }) => {
    const response = await apiAutenticada.perfil();

    expect(response.status()).toBe(200);
    const perfil = await response.json();
    expect(perfil).toMatchObject({
      email: usuarioCadastrado.email,
      nome: usuarioCadastrado.nome,
      cpf: usuarioCadastrado.cpf,
    });
    expect(perfil).not.toHaveProperty('password');
  });
});
