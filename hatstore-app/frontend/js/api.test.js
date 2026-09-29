/**
 * @jest-environment jsdom
 */

// Mock do fetch global
beforeEach(() => {
  global.fetch = jest.fn();
  localStorage.clear();
  window.location.href = '';
});

afterEach(() => {
  jest.resetAllMocks();
});

// Importa a função para teste
const { fetchApi } = require('./api.js');

test('adiciona token e email nos headers quando presentes', async () => {
  localStorage.setItem('jwt_token', 'token123');
  localStorage.setItem('user_email', 'user@teste.com');
  global.fetch.mockResolvedValue({
    ok: true,
    status: 200,
    headers: { get: () => 'application/json' },
    json: async () => ({ sucesso: true })
  });
  await fetchApi('/teste');
  const headers = global.fetch.mock.calls[0][1].headers;
  expect(headers['Authorization']).toBe('Bearer token123');
  expect(headers['X-User-Email']).toBe('user@teste.com');
});

test('redireciona para login e lança erro se status 401', async () => {
  global.fetch.mockResolvedValue({
    ok: false,
    status: 401,
    headers: { get: () => 'application/json' },
    text: async () => 'Não autorizado'
  });
  localStorage.setItem('jwt_token', 'token123');
  await expect(fetchApi('/teste')).rejects.toThrow('Não autorizado');
  expect(localStorage.getItem('jwt_token')).toBe(null);
});

test('lança erro se resposta não ok', async () => {
  global.fetch.mockResolvedValue({
    ok: false,
    status: 500,
    headers: { get: () => 'application/json' },
    text: async () => 'Erro na API'
  });
  await expect(fetchApi('/teste')).rejects.toThrow('Erro na API');
});

test('retorna objeto vazio se resposta sem json', async () => {
  global.fetch.mockResolvedValue({
    ok: true,
    status: 204,
    headers: { get: () => '' },
    json: async () => ({})
  });
  const result = await fetchApi('/teste');
  expect(result).toEqual({});
});

test('sempre envia header X-Source: frontend', async () => {
  global.fetch.mockResolvedValue({
    ok: true,
    status: 200,
    headers: { get: () => 'application/json' },
    json: async () => ({})
  });
  await fetchApi('/teste');
  const headers = global.fetch.mock.calls[0][1].headers;
  expect(headers['X-Source']).toBe('frontend');
});

test('envia header X-Source: frontend em rotas /auth', async () => {
  global.fetch.mockResolvedValue({
    ok: true,
    status: 200,
    headers: { get: () => 'application/json' },
    json: async () => ({})
  });
  await fetchApi('/auth/login');
  const [url, options] = global.fetch.mock.calls[0];
  expect(url).toBe('/auth/login');
  expect(options.headers['X-Source']).toBe('frontend');
});

test('envia header X-Source: frontend em /auth/register', async () => {
  global.fetch.mockResolvedValue({
    ok: true,
    status: 200,
    headers: { get: () => 'application/json' },
    json: async () => ({})
  });
  await fetchApi('/auth/register');
  const [url, options] = global.fetch.mock.calls[0];
  expect(url).toBe('/auth/register');
  expect(options.headers['X-Source']).toBe('frontend');
});

