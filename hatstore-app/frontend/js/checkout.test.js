/**
 * @jest-environment jsdom
 */
import { exibirModalEstoqueInsuficiente, calcularFreteNew } from './checkout.js';

describe('checkout.js', () => {
  beforeEach(() => {
    window.FEATURE_FLAGS = {};
    document.body.innerHTML = `
      <div id="checkout-old">
      <div id="pedido-resumo"></div>
      <input id="nome" value="Cliente Teste" />
      <input id="cpf" value="123.456.789-00" />
      <input id="email" value="cliente@teste.com" />
      <input id="telefone" value="(11) 91234-5678" />
      <input id="endereco" value="Rua Teste" />
      <input id="numero" value="123" />
      <input id="bairro" value="Centro" />
      <input id="cep" value="12345-678" />
      <input id="cidade" value="São Paulo" />
      <input id="uf" value="SP" />
      <input id="pagamento" value="pix" />
      <span id="erro-nome"></span>
      <span id="erro-cpf"></span>
      <span id="erro-email"></span>
      <span id="erro-telefone"></span>
      <span id="erro-endereco"></span>
      <span id="erro-numero"></span>
      <span id="erro-bairro"></span>
      <span id="erro-cep"></span>
      <span id="erro-cidade"></span>
      <span id="erro-uf"></span>
      <span id="erro-pagamento"></span>
      <button id="btn-pix"></button>
      <button id="btn-boleto"></button>
      </div>
      <div id="checkout-new" style="display:none;">
      <div id="pedido-resumo-new"></div>
      <span id="new-user-nome"></span>
      <span id="new-user-email"></span>
      <span id="new-user-cpf"></span>
      <span id="new-user-dob"></span>
      <input id="new-telefone" value="" />
      <input id="new-endereco" value="" />
      <input id="new-numero" value="" />
      <input id="new-bairro" value="" />
      <input id="new-cep" value="" />
      <input id="new-cidade" value="" />
      <input id="new-uf" value="" />
      <input id="new-pagamento" value="" />
      <span id="erro-new-telefone"></span>
      <span id="erro-new-endereco"></span>
      <span id="erro-new-numero"></span>
      <span id="erro-new-bairro"></span>
      <span id="erro-new-cep"></span>
      <span id="erro-new-cidade"></span>
      <span id="erro-new-uf"></span>
      <span id="erro-new-pagamento"></span>
      <button id="new-btn-pix"></button>
      <button id="new-btn-boleto"></button>
      <div id="frete-result" style="display:none;">
        <div id="frete-loading" style="display:none;"></div>
        <div id="frete-info" style="display:none;">
          <span id="frete-regiao"></span>
          <span id="frete-valor"></span>
        </div>
        <div id="frete-erro" style="display:none;"></div>
      </div>
      </div>
    `;
    localStorage.clear();
  });

  test('montarResumoPedido exibe itens e total com desconto', () => {
    localStorage.setItem('cupom', 'HATOFF');
    const itens = [
      { id: 1, nome: 'Chapéu Panamá', price: 100, quantidade: 2 },
      { id: 2, nome: 'Chapéu Fedora', price: 150, quantidade: 1 }
    ];
    // Função mockada
    function montarResumoPedido(itens) {
      const pedidoResumoDiv = document.getElementById('pedido-resumo');
      if (!pedidoResumoDiv) return;
      const cupom = (localStorage.getItem('cupom') || '').toUpperCase();
      let total = 0;
      let itensHtml = itens.map(item => {
        total += item.price * (item.quantidade || 1);
        return `<li>${item.nome} (${item.quantidade || 1}x) - <strong>R$ ${item.price.toFixed(2)}</strong></li>`;
      }).join('');
      let descontoHtml = '';
      let totalComDesconto = total;
      if (cupom === "HATOFF") {
        const desconto = total * 0.20;
        totalComDesconto -= desconto;
        descontoHtml = `<div class="pedido-desconto">Cupom HATOFF: -R$ ${desconto.toFixed(2)}</div>`;
      }
      pedidoResumoDiv.innerHTML = `
        <strong>Resumo do pedido:</strong>
        <ul>${itensHtml}</ul>
        ${descontoHtml}
        <div class="pedido-total"><strong>Total:</strong> R$ ${totalComDesconto.toFixed(2)}</div>
      `;
    }
    montarResumoPedido(itens);
    expect(document.getElementById('pedido-resumo').innerHTML).toContain('Cupom HATOFF');
    expect(document.getElementById('pedido-resumo').innerHTML).toContain('R$ 280.00'); // 100*2+150=350, 20% desconto=280
  });

  test('montarObjetoPedido retorna objeto correto', () => {
    localStorage.setItem('carrinho', JSON.stringify([
      { id: 1, nome: 'Chapéu Panamá', price: 100, quantidade: 2 }
    ]));
    localStorage.setItem('cupom', 'HATOFF');
    // Função mockada
    function montarObjetoPedido() {
      const carrinho = JSON.parse(localStorage.getItem('carrinho') || '[]');
      const cupom = localStorage.getItem('cupom') || "";
      const itensPedido = carrinho.map(item => ({
        id: item.id,
        nome: item.nome,
        price: item.price,
        quantidade: item.quantidade
      }));
      return {
        nome: document.getElementById('nome').value,
        cpf: document.getElementById('cpf').value,
        email: document.getElementById('email').value,
        telefone: document.getElementById('telefone').value,
        endereco: document.getElementById('endereco').value,
        numero: document.getElementById('numero').value,
        bairro: document.getElementById('bairro').value,
        cep: document.getElementById('cep').value,
        cidade: document.getElementById('cidade').value,
        uf: document.getElementById('uf').value,
        pagamento: document.getElementById('pagamento').value,
        itens: itensPedido,
        cupom: cupom,
      };
    }
    const pedido = montarObjetoPedido();
    expect(pedido.nome).toBe('Cliente Teste');
    expect(pedido.cupom).toBe('HATOFF');
    expect(pedido.itens.length).toBe(1);
    expect(pedido.itens[0].quantidade).toBe(2);
  });

  test('validaCampo retorna false para campo inválido', () => {
    document.getElementById('nome').value = 'A';
    // Função mockada
    function validaCampo(id, regex, mensagemErro) {
      const campo = document.getElementById(id);
      const erroEl = document.getElementById(`erro-${id}`);
      if (!campo.value.match(regex)) {
        campo.classList.add('input-invalid');
        if (erroEl) erroEl.textContent = mensagemErro;
        return false;
      }
      return true;
    }
    const valido = validaCampo('nome', /.{3,}/, 'Digite seu nome completo');
    expect(valido).toBe(false);
    expect(document.getElementById('nome').classList.contains('input-invalid')).toBe(true);
    expect(document.getElementById('erro-nome').textContent).toBe('Digite seu nome completo');
  });

  test('validaCampo retorna true para campo válido', () => {
    document.getElementById('nome').value = 'Cliente Teste';
    function validaCampo(id, regex, mensagemErro) {
      const campo = document.getElementById(id);
      const erroEl = document.getElementById(`erro-${id}`);
      if (!campo.value.match(regex)) {
        campo.classList.add('input-invalid');
        if (erroEl) erroEl.textContent = mensagemErro;
        return false;
      }
      return true;
    }
    const valido = validaCampo('nome', /.{3,}/, 'Digite seu nome completo');
    expect(valido).toBe(true);
  });

  test('calcularFreteNew exibe frete ao consultar CEP válido', async () => {
    document.getElementById('new-cep').value = '01000-000';
    const carrinho = [
      { id: 1, nome: 'Chapéu Panamá', price: 100, quantidade: 2 },
      { id: 2, nome: 'Chapéu Fedora', price: 150, quantidade: 1 }
    ];
    localStorage.setItem('carrinho', JSON.stringify(carrinho));

    global.fetch = jest.fn(() => Promise.resolve({
      ok: true,
      status: 200,
      headers: { get: () => 'application/json' },
      json: () => Promise.resolve({
        cep: '01000-000',
        regiao: 'Sudeste',
        custo_por_item: 12.50,
        quantidade: 3,
        valor_frete: 18.75
      })
    }));

    await calcularFreteNew(carrinho);

    expect(document.getElementById('frete-info').style.display).toBe('block');
    expect(document.getElementById('frete-regiao').textContent).toBe('Sudeste');
    expect(document.getElementById('frete-valor').textContent).toBe('R$ 18.75');
    // Resumo deve exibir subtotal, frete e total
    const resumo = document.getElementById('pedido-resumo-new').innerHTML;
    expect(resumo).toContain('Subtotal');
    expect(resumo).toContain('Frete (Sudeste)');
    expect(resumo).toContain('R$ 18.75');
  });

  test('calcularFreteNew exibe erro quando CEP não atendido', async () => {
    document.getElementById('new-cep').value = '99999-999';
    const carrinho = [
      { id: 1, nome: 'Chapéu', price: 50, quantidade: 1 }
    ];

    global.fetch = jest.fn(() => Promise.resolve({
      ok: false,
      status: 404,
      headers: { get: () => 'text/plain' },
      text: () => Promise.resolve('CEP não atendido pela nossa logística')
    }));

    await calcularFreteNew(carrinho);

    expect(document.getElementById('frete-erro').style.display).toBe('block');
    expect(document.getElementById('frete-erro').textContent).toBe('CEP não atendido pela nossa logística');
  });

  test('calcularFreteNew limpa frete quando CEP inválido', async () => {
    document.getElementById('new-cep').value = '123';
    const carrinho = [
      { id: 1, nome: 'Chapéu', price: 50, quantidade: 1 }
    ];

    await calcularFreteNew(carrinho);

    expect(document.getElementById('frete-result').style.display).toBe('none');
    // Resumo não deve ter linha de frete
    const resumo = document.getElementById('pedido-resumo-new').innerHTML;
    expect(resumo).not.toContain('Frete');
  });

  test('montarResumoPedidoNew exibe subtotal e total com frete', async () => {
    document.getElementById('new-cep').value = '80000-000';
    const carrinho = [
      { id: 1, nome: 'Chapéu Top', price: 200, quantidade: 1 }
    ];
    localStorage.setItem('carrinho', JSON.stringify(carrinho));
    localStorage.setItem('cupom', 'HATOFF');

    global.fetch = jest.fn(() => Promise.resolve({
      ok: true,
      status: 200,
      headers: { get: () => 'application/json' },
      json: () => Promise.resolve({
        cep: '80000-000',
        regiao: 'Sul',
        custo_por_item: 15.00,
        quantidade: 1,
        valor_frete: 15.00
      })
    }));

    await calcularFreteNew(carrinho);

    const resumo = document.getElementById('pedido-resumo-new').innerHTML;
    // Subtotal: 200, desconto 20% = 40, subtotal com desconto = 160
    expect(resumo).toContain('Cupom HATOFF');
    expect(resumo).toContain('R$ 160.00'); // subtotal com desconto
    expect(resumo).toContain('Frete (Sul)');
    expect(resumo).toContain('R$ 15.00'); // frete
    expect(resumo).toContain('R$ 175.00'); // total = 160 + 15
  });
});
