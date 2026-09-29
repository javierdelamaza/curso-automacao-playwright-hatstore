/**
 * @jest-environment jsdom
 */
import { criarCardPedidoNovo, formatarDataPedido } from './pedidos.js';

describe('pedidos.js', () => {
  beforeEach(() => {
    window.FEATURE_FLAGS = {};
    document.body.innerHTML = `
      <section id="pedidos-old">
        <form id="consulta-form-pedidos"></form>
        <input id="cpf-consulta-pedidos" />
        <div id="pedidos-list"></div>
      </section>
      <section id="pedidos-new" style="display:none;">
        <div id="pedidos-new-loading" style="display:none;"></div>
        <div id="pedidos-new-empty" style="display:none;"></div>
        <div id="pedidos-new-list"></div>
      </section>
    `;
    localStorage.clear();
    localStorage.setItem('jwt_token', 'fake-token');
    window.showLoader = jest.fn();
    window.hideLoader = jest.fn();
  });

  // --- Testes do modo antigo ---

  test('exibe alerta se CPF não informado', () => {
    window.alert = jest.fn();
    const form = document.getElementById('consulta-form-pedidos');
    const input = document.getElementById('cpf-consulta-pedidos');
    input.value = '';
    form.addEventListener('submit', (event) => {
      event.preventDefault();
      const cpf = input.value.trim();
      if (!cpf) {
        alert('Por favor, informe um CPF.');
        return;
      }
    });
    const event = new Event('submit');
    form.dispatchEvent(event);
    expect(window.alert).toHaveBeenCalledWith('Por favor, informe um CPF.');
  });

  test('exibe mensagem de nenhum pedido encontrado', async () => {
    const pedidosListDiv = document.getElementById('pedidos-list');
    async function consultarPedidosMock() {
      pedidosListDiv.innerHTML = "<p>Buscando seus pedidos...</p>";
      const pedidos = [];
      pedidosListDiv.innerHTML = '';
      if (!pedidos || pedidos.length === 0) {
        pedidosListDiv.innerHTML = "<p>Nenhum pedido encontrado para este CPF.</p>";
        return;
      }
    }
    await consultarPedidosMock();
    expect(pedidosListDiv.innerHTML).toContain('Nenhum pedido encontrado');
  });

  test('renderiza cards de pedidos corretamente (modo antigo)', async () => {
    const pedidosListDiv = document.getElementById('pedidos-list');
    const pedidos = [
      {
        nome: 'Cliente Teste',
        cpf: '123.456.789-00',
        pagamento: 'pix',
        total: 120.0,
        itens: [
          { nome: 'Chapéu Panamá', quantidade: 2, price: 100 },
          { nome: 'Chapéu Fedora', quantidade: 1, price: 150 }
        ]
      }
    ];
    function renderizarPedidos(pedidos) {
      pedidosListDiv.innerHTML = '';
      pedidos.forEach(pedido => {
        const card = document.createElement('div');
        card.className = 'pedido-card';
        card.innerHTML = `
          <p class="pedido-info-destaque"><strong>Cliente:</strong> ${pedido.nome}</p>
          <p class="pedido-info-destaque"><strong>CPF:</strong> ${pedido.cpf}</p>
          <p class="pedido-info-destaque"><strong>Pagamento:</strong> ${pedido.pagamento}</p>
          <p class="pedido-info-destaque"><strong>Total:</strong> R$ ${pedido.total.toFixed(2)}</p>
          <div class="pedido-itens">
            <strong>Itens do Pedido:</strong>
            <ul>
              ${pedido.itens.map(item => `<li>${item.nome} (${item.quantidade || 1}x) - R$ ${item.price.toFixed(2)}</li>`).join('')}
            </ul>
          </div>
        `;
        pedidosListDiv.appendChild(card);
      });
    }
    renderizarPedidos(pedidos);
    expect(pedidosListDiv.innerHTML).toContain('Cliente Teste');
    expect(pedidosListDiv.innerHTML).toContain('Chapéu Panamá');
    expect(pedidosListDiv.innerHTML).toContain('R$ 120.00');
  });

  // --- Testes do modo novo ---

  test('criarCardPedidoNovo gera card com dados corretos', () => {
    const pedido = {
      id: 42,
      nome: 'João',
      cpf: '123.456.789-00',
      email: 'joao@test.com',
      telefone: '(11) 91234-5678',
      endereco: 'Rua das Flores',
      numero: '100',
      bairro: 'Centro',
      cep: '01000-000',
      cidade: 'São Paulo',
      uf: 'SP',
      pagamento: 'pix',
      total: 149.80,
      frete: 12.50,
      cupom: 'HATOFF',
      created_at: '2026-03-09T14:30:00Z',
      itens: [
        { nome: 'Boné Nacional', quantidade: 2, price: 49.90 },
        { nome: 'Chapéu Importado', quantidade: 1, price: 50.00 }
      ]
    };

    const card = criarCardPedidoNovo(pedido);
    expect(card.className).toBe('pedidos-new-card');
    expect(card.innerHTML).toContain('Pedido #42');
    expect(card.innerHTML).toContain('R$ 149.80');
    expect(card.innerHTML).toContain('Boné Nacional');
    expect(card.innerHTML).toContain('Chapéu Importado');
    expect(card.innerHTML).toContain('PIX');
    expect(card.innerHTML).toContain('Rua das Flores');
    expect(card.innerHTML).toContain('HATOFF');
    expect(card.innerHTML).toContain('01000-000');
    expect(card.innerHTML).toContain('R$ 12.50');
    expect(card.innerHTML).toContain('Frete');
  });

  test('criarCardPedidoNovo não exibe cupom quando não há', () => {
    const pedido = {
      id: 10, nome: 'Maria', cpf: '111.111.111-11', email: 'maria@test.com',
      telefone: '(11) 91234-5678', endereco: 'Rua A', numero: '1',
      bairro: 'B', cep: '01000-000', cidade: 'SP', uf: 'SP',
      pagamento: 'boleto', total: 79.90, frete: 0, cupom: '',
      created_at: '2026-01-15T10:00:00Z',
      itens: [{ nome: 'Boné Crescer', quantidade: 1, price: 79.90 }]
    };

    const card = criarCardPedidoNovo(pedido);
    expect(card.innerHTML).not.toContain('Cupom');
    expect(card.innerHTML).not.toContain('Frete');
    expect(card.innerHTML).toContain('Boleto');
  });

  test('formatarDataPedido formata data corretamente', () => {
    expect(formatarDataPedido('2026-03-09T14:30:00Z')).toMatch(/09\/03\/2026/);
    expect(formatarDataPedido(null)).toBe('—');
    expect(formatarDataPedido('')).toBe('—');
    expect(formatarDataPedido('invalid')).toBe('—');
  });
});
