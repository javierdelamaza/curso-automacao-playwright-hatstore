import { fetchApi, escapeHtml } from './api.js';
import { showLoader, hideLoader } from './app.js';
// frontend/js/pedidos.js

document.addEventListener('DOMContentLoaded', () => {
    // Redireciona para login se não estiver autenticado
    if (!localStorage.getItem('jwt_token')) {
        window.location.href = 'auth.html';
        return;
    }

    const isNewHatstore = window.FEATURE_FLAGS && window.FEATURE_FLAGS.new_hatstore;

    if (isNewHatstore) {
        inicializarNovoPedidos();
    } else {
        inicializarPedidosAntigo();
    }
});

// =====================================================================
// MODO ANTIGO: Busca por CPF
// =====================================================================
function inicializarPedidosAntigo() {
    const pedidosOld = document.getElementById('pedidos-old');
    const pedidosNew = document.getElementById('pedidos-new');
    if (pedidosOld) pedidosOld.style.display = '';
    if (pedidosNew) pedidosNew.style.display = 'none';

    const consultaForm = document.getElementById('consulta-form-pedidos');
    if (!consultaForm) return;

    if (window.jQuery && window.jQuery.fn.mask) {
        jQuery('#cpf-consulta-pedidos').mask('000.000.000-00');
    }

    consultaForm.addEventListener('submit', async (event) => {
        event.preventDefault();

        const cpfInput = document.getElementById('cpf-consulta-pedidos');
        const pedidosListDiv = document.getElementById('pedidos-list');
        const cpf = cpfInput.value.trim();

        if (!cpf) {
            alert('Por favor, informe um CPF.');
            return;
        }

        pedidosListDiv.innerHTML = "<p>Buscando seus pedidos...</p>";
        showLoader();

        try {
            const pedidos = await fetchApi(`/pedidos?cpf=${encodeURIComponent(cpf)}`);

            pedidosListDiv.innerHTML = '';
            if (!pedidos || pedidos.length === 0) {
                pedidosListDiv.innerHTML = "<p>Nenhum pedido encontrado para este CPF.</p>";
                return;
            }

            pedidos.forEach(pedido => {
                const card = document.createElement('div');
                card.className = 'pedido-card';
                const freteHtml = pedido.frete > 0
                    ? `<p class="pedido-info-destaque"><strong>Frete:</strong> R$ ${pedido.frete.toFixed(2)}</p>`
                    : '';
                card.innerHTML = `
                    <p class="pedido-info-destaque"><strong>Cliente:</strong> ${escapeHtml(pedido.nome)}</p>
                    <p class="pedido-info-destaque"><strong>CPF:</strong> ${escapeHtml(pedido.cpf)}</p>
                    <p class="pedido-info-destaque"><strong>Pagamento:</strong> ${escapeHtml(pedido.pagamento)}</p>
                    ${freteHtml}
                    <p class="pedido-info-destaque"><strong>Total:</strong> R$ ${pedido.total.toFixed(2)}</p>
                    <div class="pedido-itens">
                        <strong>Itens do Pedido:</strong>
                        <ul>
                            ${pedido.itens.map(item => `<li>${escapeHtml(item.nome)} (${item.quantidade || 1}x) - R$ ${item.price.toFixed(2)}</li>`).join('')}
                        </ul>
                    </div>
                `;
                pedidosListDiv.appendChild(card);
            });

        } catch (error) {
            console.error('Erro ao consultar pedidos:', error);
            pedidosListDiv.innerHTML = "<p>Ocorreu um erro ao buscar seus pedidos. Tente novamente mais tarde.</p>";
        } finally {
            hideLoader();
        }
    });
}

// =====================================================================
// MODO NOVO: Pedidos automáticos pelo login
// =====================================================================
async function inicializarNovoPedidos() {
    const pedidosOld = document.getElementById('pedidos-old');
    const pedidosNew = document.getElementById('pedidos-new');
    if (pedidosOld) pedidosOld.style.display = 'none';
    if (pedidosNew) pedidosNew.style.display = '';

    const loadingEl = document.getElementById('pedidos-new-loading');
    const emptyEl = document.getElementById('pedidos-new-empty');
    const listEl = document.getElementById('pedidos-new-list');

    loadingEl.style.display = 'flex';
    emptyEl.style.display = 'none';
    listEl.innerHTML = '';

    try {
        const pedidos = await fetchApi('/meus-pedidos');

        loadingEl.style.display = 'none';

        if (!pedidos || pedidos.length === 0) {
            emptyEl.style.display = 'flex';
            return;
        }

        pedidos.forEach(pedido => {
            listEl.appendChild(criarCardPedidoNovo(pedido));
        });

    } catch (error) {
        console.error('Erro ao buscar pedidos:', error);
        loadingEl.style.display = 'none';
        listEl.innerHTML = '<div class="pedidos-new-error">Ocorreu um erro ao carregar seus pedidos. Tente novamente mais tarde.</div>';
    }
}

function criarCardPedidoNovo(pedido) {
    const card = document.createElement('div');
    card.className = 'pedidos-new-card';

    const dataFormatada = formatarDataPedido(pedido.created_at);
    const statusPagamento = pedido.pagamento === 'pix' ? 'PIX' : 'Boleto';
    const temCupom = pedido.cupom && pedido.cupom.length > 0;
    const temFrete = pedido.frete > 0;

    let itensHtml = pedido.itens.map(item => `
        <div class="pedidos-new-item">
            <span class="pedidos-new-item-nome">${escapeHtml(item.nome)}</span>
            <span class="pedidos-new-item-qty">${item.quantidade}x</span>
            <span class="pedidos-new-item-price">R$ ${(item.price * item.quantidade).toFixed(2)}</span>
        </div>
    `).join('');

    card.innerHTML = `
        <div class="pedidos-new-card-header">
            <div class="pedidos-new-card-info">
                <span class="pedidos-new-card-id">Pedido #${pedido.id}</span>
                <span class="pedidos-new-card-date">${escapeHtml(dataFormatada)}</span>
            </div>
            <div class="pedidos-new-card-total">R$ ${pedido.total.toFixed(2)}</div>
        </div>
        <div class="pedidos-new-card-body">
            <div class="pedidos-new-itens">
                <h4>Itens</h4>
                ${itensHtml}
            </div>
            <div class="pedidos-new-card-details">
                <div class="pedidos-new-detail-row">
                    <span class="pedidos-new-detail-label">Pagamento</span>
                    <span class="pedidos-new-detail-value pedidos-new-badge">${escapeHtml(statusPagamento)}</span>
                </div>
                <div class="pedidos-new-detail-row">
                    <span class="pedidos-new-detail-label">Entrega</span>
                    <span class="pedidos-new-detail-value">${escapeHtml(pedido.endereco)}, ${escapeHtml(pedido.numero)} — ${escapeHtml(pedido.bairro)}, ${escapeHtml(pedido.cidade)}/${escapeHtml(pedido.uf)}</span>
                </div>
                <div class="pedidos-new-detail-row">
                    <span class="pedidos-new-detail-label">CEP</span>
                    <span class="pedidos-new-detail-value">${escapeHtml(pedido.cep)}</span>
                </div>
                ${temCupom ? `
                <div class="pedidos-new-detail-row">
                    <span class="pedidos-new-detail-label">Cupom</span>
                    <span class="pedidos-new-detail-value pedidos-new-cupom">${escapeHtml(pedido.cupom)}</span>
                </div>
                ` : ''}
                ${temFrete ? `
                <div class="pedidos-new-detail-row">
                    <span class="pedidos-new-detail-label">Frete</span>
                    <span class="pedidos-new-detail-value">R$ ${pedido.frete.toFixed(2)}</span>
                </div>
                ` : ''}
            </div>
        </div>
    `;

    return card;
}

function formatarDataPedido(dateStr) {
    if (!dateStr) return '—';
    const d = new Date(dateStr);
    if (isNaN(d.getTime())) return '—';
    const dia = String(d.getDate()).padStart(2, '0');
    const mes = String(d.getMonth() + 1).padStart(2, '0');
    const ano = d.getFullYear();
    const hora = String(d.getHours()).padStart(2, '0');
    const min = String(d.getMinutes()).padStart(2, '0');
    return `${dia}/${mes}/${ano} às ${hora}:${min}`;
}

export { criarCardPedidoNovo, formatarDataPedido };