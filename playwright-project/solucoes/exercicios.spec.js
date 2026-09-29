import { test, expect } from '../fixtures/index.js';
import { CATEGORIAS, PRODUTOS, ENDERECO_PADRAO, CUPOM_VALIDO, PERCENTUAL_DESCONTO_HATOFF } from '../data/catalogo.js';
import { novoUsuario } from '../factories/usuarioFactory.js';
import { validarContrato } from '../api/validarContrato.js';
import { pedidoSchema } from '../api/schemas.js';

test.describe('Bloco 2 - filtros da home', () => {
  test('filtrar por categoria "Crescer" mostra só o produto dessa categoria', async ({ homePage }) => {
    await homePage.goto();

    await homePage.filterByCategoria(CATEGORIAS.crescer);

    await expect(homePage.hatCards).toHaveCount(1);
    await expect(homePage.hatCardTitles).toHaveText([PRODUTOS.gustavoCarvalho.nome]);
  });

  test('filtrar por faixa de preço 100 a 130 mostra só preços dentro do intervalo', async ({ homePage }) => {
    await homePage.goto();

    await homePage.filterByPreco({ min: 100, max: 130 });

    // Panamá (120), Cowboy (109.90) e Gaúcho (129.90)
    await expect(homePage.hatCards).toHaveCount(3);
    for (const preco of await homePage.precosExibidos()) {
      expect(preco).toBeGreaterThanOrEqual(100);
      expect(preco).toBeLessThanOrEqual(130);
    }
  });
});

test.describe('Bloco 3 - validações do cadastro', () => {
  test.beforeEach(async ({ loginPage }) => {
    await loginPage.goto();
    await loginPage.goToRegisterForm();
  });

  test('senha fraca (sem caractere especial) é recusada', async ({ loginPage }) => {
    await loginPage.register(novoUsuario({ password: 'Senha1234' }));

    await expect(loginPage.registerMessage).not.toBeEmpty();
    await expect(loginPage.registerMessage).not.toHaveText('Registro bem-sucedido! Você já pode fazer o login.');
  });

  test('CPF incompleto mostra erro no campo', async ({ page, loginPage }) => {
    await loginPage.register(novoUsuario({ cpf: '123.456' }));

    await expect(page.locator('#erro-register-cpf')).toHaveText('CPF inválido. Formato: 000.000.000-00');
  });
});

test.describe('Bloco 4 - compra completa com boleto @e2e', () => {
  test('cliente do Nordeste compra com boleto e vê o frete da região', async ({
    paginaLogada,
    homePage,
    cartPage,
    checkoutPage,
    pedidosPage,
  }) => {
    await homePage.goto();
    await homePage.adicionarAoCarrinho(PRODUTOS.gustavoCarvalho.nome);
    await cartPage.finalizarCompra();

    await checkoutPage.preencherEndereco({ ...ENDERECO_PADRAO, cep: '40010-000', cidade: 'Salvador', uf: 'BA' });
    await expect(checkoutPage.freteRegiao).toHaveText('Nordeste');
    await expect(checkoutPage.freteValor).toHaveText('R$ 22.00');

    await checkoutPage.escolherPagamento('boleto');
    await checkoutPage.finalizar();
    await expect(checkoutPage.modalTitulo).toHaveText('Pagamento Boleto');

    await pedidosPage.goto();
    await expect(pedidosPage.cards.first()).toContainText('Boleto');
  });
});

test.describe('Bloco 5 - contrato de pedidos', () => {
  test('GET /api/meus-pedidos respeita o contrato de pedido', async ({ apiAutenticada, usuarioCadastrado }) => {
    const produto = PRODUTOS.gustavoCarvalho;
    const criado = await apiAutenticada.criarPedido({
      nome: usuarioCadastrado.nome,
      cpf: usuarioCadastrado.cpf,
      email: usuarioCadastrado.email,
      ...ENDERECO_PADRAO,
      pagamento: 'pix',
      cupom: '',
      itens: [{ id: produto.id, nome: produto.nome, price: produto.price, quantidade: 1 }],
    });
    expect(criado.status()).toBe(201);

    const pedidos = await (await apiAutenticada.meusPedidos()).json();

    expect(pedidos.length).toBeGreaterThan(0);
    for (const pedido of pedidos) {
      expect(validarContrato(pedidoSchema, pedido)).toEqual([]);
    }
  });
});

test.describe('Bloco 6 - exercício com IA: cupom no carrinho', () => {
  test.beforeEach(async ({ homePage }) => {
    await homePage.goto();
    await homePage.adicionarAoCarrinho(PRODUTOS.panama.nome);
  });

  test('cupom HATOFF aplica 20% de desconto no carrinho', async ({ cartPage }) => {
    await cartPage.aplicarCupom(CUPOM_VALIDO);

    const desconto = PRODUTOS.panama.price * PERCENTUAL_DESCONTO_HATOFF;
    await expect(cartPage.mensagemCupom).toHaveText('Cupom aplicado!');
    await expect(cartPage.desconto).toHaveText(`Cupom HATOFF aplicado: -R$ ${desconto.toFixed(2)}`);
    await expect(cartPage.total).toHaveText(`R$ ${(PRODUTOS.panama.price - desconto).toFixed(2)}`);
  });

  test('cupom inválido mostra "Cupom inválido."', async ({ cartPage }) => {
    await cartPage.aplicarCupom('CUPOM-FALSO');

    await expect(cartPage.mensagemCupom).toHaveText('Cupom inválido.');
    await expect(cartPage.desconto).toHaveCount(0);
  });
});
