import { test, expect } from '../../fixtures/index.js';
import { PRODUTOS, ENDERECO_PADRAO, FRETE_SUDESTE_1_ITEM } from '../../data/catalogo.js';

// Fluxo completo do cliente: login -> busca -> carrinho -> checkout -> consulta do pedido.
test.describe('Compra completa @e2e', () => {
  test.beforeEach(async ({ page, loginPage, usuarioCadastrado }) => {
    await loginPage.goto();
    await loginPage.login(usuarioCadastrado.email, usuarioCadastrado.password);
    await page.waitForURL('**/index.html');
  });

  test('cliente compra um chapéu com PIX e vê o pedido em Meus Pedidos', async ({
    homePage,
    cartPage,
    checkoutPage,
    pedidosPage,
    usuarioCadastrado,
  }) => {
    const produto = PRODUTOS.gustavoCarvalho;

    await test.step('buscar o produto e adicionar ao carrinho', async () => {
      await homePage.search('Gustavo');
      await homePage.adicionarAoCarrinho(produto.nome);
      await expect(cartPage.itemNomes).toHaveText([produto.nome]);
    });

    await test.step('ir para o checkout com os dados do usuário logado', async () => {
      await cartPage.finalizarCompra();
      await expect(checkoutPage.usuarioNome).toHaveText(usuarioCadastrado.nome);
      await expect(checkoutPage.usuarioCpf).toHaveText(usuarioCadastrado.cpf);
    });

    await test.step('preencher endereço e ver o frete calculado', async () => {
      await checkoutPage.preencherEndereco(ENDERECO_PADRAO);
      await expect(checkoutPage.freteRegiao).toHaveText(FRETE_SUDESTE_1_ITEM.regiao);
      await expect(checkoutPage.freteValor).toHaveText(`R$ ${FRETE_SUDESTE_1_ITEM.valor.toFixed(2)}`);
    });

    await test.step('pagar com PIX', async () => {
      await checkoutPage.escolherPagamento('pix');
      await checkoutPage.finalizar();
      await expect(checkoutPage.modalTitulo).toHaveText('Pagamento PIX');
      await expect(checkoutPage.codigoPagamento).not.toBeEmpty();
    });

    await test.step('consultar o pedido em Meus Pedidos', async () => {
      await pedidosPage.goto();
      await expect(pedidosPage.cards).toHaveCount(1);
      await expect(pedidosPage.cards.first()).toContainText(produto.nome);
      await expect(pedidosPage.cards.first()).toContainText('PIX');
    });
  });
});
