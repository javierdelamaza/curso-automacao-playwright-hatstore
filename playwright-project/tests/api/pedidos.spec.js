import { test, expect } from '../../fixtures/index.js';
import { PRODUTOS, ENDERECO_PADRAO, CUPOM_VALIDO } from '../../data/catalogo.js';

function montarPedido(usuario, sobrescrever = {}) {
  const produto = PRODUTOS.gustavoCarvalho;
  return {
    nome: usuario.nome,
    cpf: usuario.cpf,
    email: usuario.email,
    ...ENDERECO_PADRAO,
    pagamento: 'pix',
    cupom: '',
    itens: [{ id: produto.id, nome: produto.nome, price: produto.price, quantidade: 1 }],
    ...sobrescrever,
  };
}

test.describe('API - pedidos', () => {
  test('POST /api/pedido sem token devolve 401', async ({ api, usuarioCadastrado }) => {
    const response = await api.criarPedido(montarPedido(usuarioCadastrado));

    expect(response.status()).toBe(401);
  });

  test('POST /api/pedido cria o pedido e ele aparece em /api/meus-pedidos', async ({ apiAutenticada, usuarioCadastrado }) => {
    const criado = await apiAutenticada.criarPedido(montarPedido(usuarioCadastrado));
    expect(criado.status(), await criado.text()).toBe(201);
    const pedido = await criado.json();

    const meus = await (await apiAutenticada.meusPedidos()).json();
    expect(meus.map((p) => p.id)).toContain(pedido.id);
  });

  test('GET /api/pedidos?cpf= encontra o pedido pelo CPF', async ({ apiAutenticada, usuarioCadastrado }) => {
    await apiAutenticada.criarPedido(montarPedido(usuarioCadastrado));

    const response = await apiAutenticada.consultarPedidosPorCpf(usuarioCadastrado.cpf);

    expect(response.status()).toBe(200);
    const pedidos = await response.json();
    expect(pedidos).toHaveLength(1);
    expect(pedidos[0].cpf).toBe(usuarioCadastrado.cpf);
  });

  test('cupom HATOFF só pode ser usado uma vez por CPF', async ({ apiAutenticada, usuarioCadastrado }) => {
    // Bug encontrado pelo próprio teste: qualquer pedido com HATOFF devolve 500 "Erro ao buscar chapéu".
    // test.fail() documenta o bug; quando corrigirem, o teste passa a "falhar" e avisa para remover a marcação.
    test.fail(true, 'BUG: POST /api/pedido com cupom HATOFF devolve 500');

    const pedidoComCupom = montarPedido(usuarioCadastrado, { cupom: CUPOM_VALIDO });

    const primeiro = await apiAutenticada.criarPedido(pedidoComCupom);
    expect(primeiro.status(), await primeiro.text()).toBe(201);

    const segundo = await apiAutenticada.criarPedido(pedidoComCupom);
    expect(segundo.ok()).toBeFalsy();
  });

  test('pedido sem itens devolve 400', async ({ apiAutenticada, usuarioCadastrado }) => {
    const response = await apiAutenticada.criarPedido(montarPedido(usuarioCadastrado, { itens: [] }));

    expect(response.status()).toBe(400);
  });
});
