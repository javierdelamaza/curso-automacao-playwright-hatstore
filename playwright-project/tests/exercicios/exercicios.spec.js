// EXERCÍCIOS - remova o `.fixme` de um teste por vez e implemente.
// Soluções em ../../solucoes (só olhe depois de tentar!).
import { test, expect } from '../../fixtures/index.js';

test.describe('Bloco 2 - filtros da home', () => {
  test.fixme('filtrar por categoria "Crescer" mostra só o produto dessa categoria', async ({ homePage }) => {
    await homePage.goto();
    // 1. Aplique o filtro da categoria crescer (veja os métodos do HomePage)
    // 2. Valide que aparece 1 card e que o título é "Chapéu Gustavo Carvalho"
  });

  test.fixme('filtrar por faixa de preço 100 a 130 mostra só preços dentro do intervalo', async ({ homePage }) => {
    await homePage.goto();
    // 1. Filtre com min 100 e max 130
    // 2. Use homePage.precosExibidos() e valide cada preço
    // Dica: o filtro redesenha os cards; espere a quantidade mudar antes de ler os preços
  });
});

test.describe('Bloco 3 - validações do cadastro', () => {
  test.fixme('senha fraca (sem caractere especial) é recusada', async ({ loginPage }) => {
    // Gere um usuário com novoUsuario({ password: 'Senha1234' }) e valide que o cadastro NÃO é concluído
  });

  test.fixme('CPF incompleto mostra erro no campo', async ({ page, loginPage }) => {
    // Descubra pelo DevTools onde a mensagem de erro do CPF aparece
  });
});

test.describe('Bloco 4 - compra completa com boleto @e2e', () => {
  test.fixme('cliente do Nordeste compra com boleto e vê o frete da região', async () => {
    // Use a fixture paginaLogada para já começar logado (sem passar pela tela de login)
    // CEP de Salvador: 40010-000 -> região Nordeste, R$ 22.00 para 1 item
  });
});

test.describe('Bloco 5 - contrato de pedidos', () => {
  test.fixme('GET /api/meus-pedidos respeita o contrato de pedido', async () => {
    // Crie um pedido via apiAutenticada.criarPedido(...) e valide cada item com pedidoSchema
  });
});

test.describe('Bloco 6 - exercício com IA: cupom no carrinho', () => {
  test.fixme('cupom HATOFF aplica 20% de desconto no carrinho', async () => {
    // Faça este teste com a IA seguindo o módulo 06 (prompt com fatos da tela, não do código-fonte)
  });

  test.fixme('cupom inválido mostra "Cupom inválido."', async () => {
    // Idem: use a IA e revise com o checklist do módulo 06
  });
});
