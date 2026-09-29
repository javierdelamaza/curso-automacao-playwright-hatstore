// Validação do CPF no formulário de pedidos.html (consulta-form-pedidos)
document.addEventListener('DOMContentLoaded', () => {
    const formPedidos = document.getElementById('consulta-form-pedidos');
    if (formPedidos) {
        formPedidos.addEventListener('submit', function(event) {
            const cpfInput = document.getElementById('cpf-consulta-pedidos');
            const erroCpf = document.getElementById('erro-cpf-consulta-pedidos');
            if (erroCpf) erroCpf.textContent = '';
            if (cpfInput) cpfInput.classList.remove('input-invalid');
            const cpf = cpfInput.value.trim();
            if (!cpf.match(/^\d{3}\.\d{3}\.\d{3}-\d{2}$/)) {
                event.preventDefault();
                if (erroCpf) erroCpf.textContent = 'CPF inválido. Formato: 000.000.000-00';
                if (cpfInput) cpfInput.classList.add('input-invalid');
            }
        });
    }
});
import { fetchApi, escapeHtml } from './api.js';

export { showLoader, hideLoader, openModal, closeModal };

let todosChapeus = [];
let carrinho = [];

function applyFeatureFlags() {
    const flags = window.FEATURE_FLAGS || {};
    const filterPanel = document.getElementById('filter-panel');
    if (filterPanel) {
        filterPanel.style.display = flags.filters_enabled === false ? 'none' : '';
    }
}

// Função principal que roda quando o HTML está pronto
document.addEventListener('DOMContentLoaded', async () => {
    applyFeatureFlags();

    document.addEventListener('click', function (e) {
        if (e.target && e.target.id === 'close-cart') {
            const cartSidebar = document.getElementById('cart-sidebar');
            if (cartSidebar) cartSidebar.classList.remove('open');
        }
    });

    const headerCartBtn = document.getElementById('header-cart-btn');
    const cartSidebar = document.getElementById('cart-sidebar');
    if (headerCartBtn && cartSidebar) {
        headerCartBtn.addEventListener('click', () => {
            cartSidebar.classList.add('open');
            atualizarCarrinho();
        });
    }

    if (window.jQuery && window.jQuery.fn.mask) {
        $('#cpf-consulta-sidebar').mask('000.000.000-00');
    }

    function setupAuthUI() {
        const authLinksContainer = document.getElementById('auth-links');
        if (!authLinksContainer) return;

        const token = localStorage.getItem('jwt_token');

        if (token) {
            authLinksContainer.innerHTML = '<button id="logout-button" class="header-link" style="background:none; border:none; cursor:pointer;">Logout</button>';
            document.getElementById('logout-button').addEventListener('click', () => {
                localStorage.removeItem('jwt_token');
                localStorage.removeItem('user_email');
                window.location.reload();
            });
        } else {
            authLinksContainer.innerHTML = '<a href="auth.html" class="header-link">Login</a>';
        }
    }
    setupAuthUI();

    carregarChapeus();
    const searchHatInput = document.getElementById('search-hat');
    if (searchHatInput) {
        searchHatInput.addEventListener('input', filtrarChapeus);
    }

    const carrinhoSalvo = localStorage.getItem('carrinho');
    if (carrinhoSalvo) {
        carrinho = JSON.parse(carrinhoSalvo);
        atualizarCarrinho();
    }


    const cartMenuButton = document.getElementById('cart-menu-button');
    const closeCartButton = document.getElementById('close-cart');

    if (cartMenuButton) {
        cartMenuButton.addEventListener('click', () => {
            cartSidebar.classList.add('open');
            atualizarCarrinho();
        });

        closeCartButton.addEventListener('click', () => cartSidebar.classList.remove('open'));

        window.addEventListener('mousedown', (e) => {
            if (cartSidebar.classList.contains('open') && !cartSidebar.contains(e.target) && e.target !== cartMenuButton) {
                cartSidebar.classList.remove('open');
            }
        });
    }

    const pedidosMenuButton = document.getElementById('pedidos-menu-button');
    const pedidosSidebar = document.getElementById('pedidos-sidebar');
    const closePedidosButton = document.getElementById('close-pedidos');

    if (pedidosMenuButton) {
        pedidosMenuButton.addEventListener('click', () => {
            pedidosSidebar.classList.add('open');
            document.getElementById('pedidos-list').innerHTML = '';
        });

        closePedidosButton.addEventListener('click', () => pedidosSidebar.classList.remove('open'));

        window.addEventListener('mousedown', (e) => {
            if (pedidosSidebar.classList.contains('open') && !pedidosSidebar.contains(e.target) && e.target !== pedidosMenuButton) {
                pedidosSidebar.classList.remove('open');
            }
        });
    }

    const consultaForm = document.getElementById('consulta-form-sidebar');
    if (consultaForm) {
        consultaForm.addEventListener('submit', handleConsultaPedidos);
    }

    const aplicarCupomBtn = document.getElementById('aplicar-cupom');
    if (aplicarCupomBtn) {
        aplicarCupomBtn.addEventListener('click', aplicarCupom);
    }

    const checkoutBtn = document.getElementById('checkout-button');
    if (checkoutBtn) {
        checkoutBtn.addEventListener('click', handleCheckout);
    }

    // Filtro de categoria e preço (configurado no segundo DOMContentLoaded)
});

//
// FUNÇÕES DE DADOS E API
//
async function carregarChapeus() {
    showLoader();
    try {
        const hats = await fetchApi('/hats'); // Usa a nova função centralizada
        todosChapeus = hats;
        exibirChapeus(hats);
    } catch (error) {
        console.error("Erro ao carregar chapéus:", error);
        openModal("Erro", "<p>Não foi possível carregar os chapéus. Verifique sua conexão ou o código de acesso.</p>");
    } finally {
        hideLoader();
    }
}

async function handleConsultaPedidos(event) {
    event.preventDefault();
    const cpfInput = document.getElementById('cpf-consulta-sidebar');
    const pedidosListDiv = document.getElementById('pedidos-list-sidebar');
    const erroCpf = document.getElementById('erro-cpf-consulta-sidebar');
    if (erroCpf) erroCpf.textContent = '';
    if (cpfInput) cpfInput.classList.remove('input-invalid');

    const cpf = cpfInput.value.trim();
    if (!cpf.match(/^\d{3}\.\d{3}\.\d{3}-\d{2}$/)) {
        if (erroCpf) erroCpf.textContent = 'CPF inválido. Formato: 000.000.000-00';
        if (cpfInput) cpfInput.classList.add('input-invalid');
        return;
    }

    pedidosListDiv.innerHTML = "<p>Buscando pedidos...</p>";

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
            card.innerHTML = `
                <strong>Pedido:</strong> ${escapeHtml(pedido.nome)}<br>
                <strong>CPF:</strong> ${escapeHtml(pedido.cpf)}<br>
                <strong>e-mail:</strong> ${escapeHtml(pedido.email)}<br>
                <strong>Pagamento:</strong> ${escapeHtml(pedido.pagamento)}<br>
                <strong>Total:</strong> R$ ${pedido.total.toFixed(2)}
                <div class="pedido-itens">
                    <strong>Itens:</strong>
                    <ul>
                        ${pedido.itens.map(item => `<li>${escapeHtml(item.nome)} (${item.quantidade || 1}x) - R$ ${item.price.toFixed(2)}</li>`).join('')}
                    </ul>
                </div>
            `;
            pedidosListDiv.appendChild(card);
        });
    } catch (error) {
        console.error("Erro ao consultar pedidos:", error);
        pedidosListDiv.innerHTML = "<p>Erro ao consultar pedidos.</p>";
    }
}


//
// FUNÇÕES DE UI E MANIPULAÇÃO DO DOM
//
function exibirChapeus(hats) {
    const hatsDiv = document.getElementById('hats');
    if (!hatsDiv) return;
    hatsDiv.innerHTML = '';
    if (hats.length === 0) {
        hatsDiv.innerHTML = '<p style="text-align:center;color:#2563eb;font-weight:bold;">Nenhum chapéu encontrado.</p>';
        return;
    }
    // Agrupa os produtos: primeiro com estoque, depois sem estoque
    const hatsComEstoque = hats.filter(hat => hat.temEstoque !== false);
    const hatsSemEstoque = hats.filter(hat => hat.temEstoque === false);
    const hatsOrdenados = [...hatsComEstoque, ...hatsSemEstoque];
    hatsOrdenados.forEach(hat => {
        const hatCard = document.createElement('div');
        hatCard.className = 'hat-card';
        let cardHtml = `
            <img src="imagens/chapeu-${hat.id}.jpg" alt="${escapeHtml(hat.nome)}" onerror="this.src='imagens/logo-hatstore.png'">
            <h3>${escapeHtml(hat.nome)}</h3>
            <p>Preço: R$ ${hat.price.toFixed(2)}</p>
        `;
        if (hat.temEstoque === false) {
            cardHtml += `<p style="color:red;font-weight:bold;">Indisponível no momento</p>`;
        } else {
            cardHtml += `<button class="add-to-cart-btn" data-id="${hat.id}" data-nome="${escapeHtml(hat.nome)}" data-price="${hat.price}">Adicionar ao carrinho</button>`;
        }
        hatCard.innerHTML = cardHtml;
        hatsDiv.appendChild(hatCard);
    });
    document.querySelectorAll('.add-to-cart-btn').forEach(button => {
        button.addEventListener('click', (event) => {
            const { id, nome, price } = event.target.dataset;
            adicionarAoCarrinho(parseInt(id), nome, parseFloat(price));
        });
    });
}

function filtrarChapeus() {
    const termo = document.getElementById('search-hat').value.toLowerCase();
    const filtrados = todosChapeus.filter(hat => hat.nome.toLowerCase().includes(termo));
    exibirChapeus(filtrados);
}

function aplicarCupom() {
    const cupomInput = document.getElementById('cupom');
    const erroCupom = document.getElementById('erro-cupom');
    if (!cupomInput || !erroCupom) return;

    const cupom = cupomInput.value.trim().toUpperCase();
    if (cupom === "HATOFF") {
        localStorage.setItem('cupom', cupom);
        erroCupom.textContent = "Cupom aplicado!";
        erroCupom.className = 'cupom-sucesso';
        atualizarCarrinho();
    } else {
        erroCupom.textContent = "Cupom inválido.";
        erroCupom.className = 'input-error';
        localStorage.removeItem('cupom');
    }
}


//
// FUNÇÕES DO CARRINHO
//
function adicionarAoCarrinho(id, nome, price) {
    const itemExistente = carrinho.find(item => item.id === id);
    if (itemExistente) {
        itemExistente.quantidade += 1;
    } else {
        carrinho.push({ id, nome, price, quantidade: 1 });
    }
    salvarCarrinho();
    atualizarCarrinho();
    document.getElementById('cart-sidebar').classList.add('open');
}

function alterarQuantidade(id, delta) {
    const item = carrinho.find(item => item.id === id);
    if (!item) return;
    item.quantidade += delta;
    if (item.quantidade <= 0) {
        carrinho = carrinho.filter(i => i.id !== id);
    }
    salvarCarrinho();
    atualizarCarrinho();
}

function removerDoCarrinho(id) {
    carrinho = carrinho.filter(item => item.id !== id);
    salvarCarrinho();
    atualizarCarrinho();
}

function salvarCarrinho() {
    localStorage.setItem('carrinho', JSON.stringify(carrinho));
}

function atualizarCarrinho() {
    const cartDiv = document.getElementById('cart-items');
    const totalDiv = document.getElementById('total-amount');
    const checkoutBtn = document.getElementById('checkout-button');
    if (!cartDiv || !totalDiv || !checkoutBtn) return;

    if (carrinho.length === 0) {
        cartDiv.innerHTML = '<div class="cart-empty-message">Seu carrinho está vazio.</div>';
        totalDiv.innerHTML = '';
        checkoutBtn.style.display = 'none';
        return;
    }

    let total = 0;
    cartDiv.innerHTML = '';
    carrinho.forEach(item => {
        total += item.price * item.quantidade;
        const itemDiv = document.createElement('div');
        itemDiv.className = 'cart-item';
        itemDiv.innerHTML = `
            <div class="cart-item-info">
                <span class="cart-item-nome">${escapeHtml(item.nome)}</span>
                <span class="cart-item-preco">R$ ${item.price.toFixed(2)}</span>
                <div class="cart-item-actions">
                    <button class="qty-btn" data-id="${item.id}" data-delta="-1" title="Diminuir">-</button>
                    <span>${item.quantidade}</span>
                    <button class="qty-btn" data-id="${item.id}" data-delta="1" title="Aumentar">+</button>
                    <button class="remove-btn" data-id="${item.id}" title="Remover" >&#10006;</button>
                </div>
            </div>
        `;
        cartDiv.appendChild(itemDiv);
    });

    // Adiciona event listeners aos botões do carrinho
    cartDiv.querySelectorAll('.qty-btn').forEach(btn =>
        btn.addEventListener('click', e => alterarQuantidade(parseInt(e.target.dataset.id), parseInt(e.target.dataset.delta)))
    );
    cartDiv.querySelectorAll('.remove-btn').forEach(btn =>
        btn.addEventListener('click', e => removerDoCarrinho(parseInt(e.target.dataset.id)))
    );

    const cupom = (localStorage.getItem('cupom') || '').toUpperCase();
    let descontoHtml = '';
    let totalComDesconto = total;

    if (cupom === 'HATOFF') {
        const desconto = total * 0.20;
        totalComDesconto -= desconto;
        descontoHtml = `<div class="cart-cupom-info">Cupom HATOFF aplicado: -R$ ${desconto.toFixed(2)}</div>`;
    }

    totalDiv.innerHTML = `
        <div class="cart-subtotal">Subtotal: <span>R$ ${total.toFixed(2)}</span></div>
        ${descontoHtml}
        <div class="cart-total">Total com desconto: <span>R$ ${totalComDesconto.toFixed(2)}</span></div>
    `;
    checkoutBtn.style.display = 'block';
}

function handleCheckout() {
    if (carrinho.length === 0) {
        alert('Seu carrinho está vazio!');
        return;
    }
    window.location.href = "checkout.html";
}


function showLoader() {
    const loader = document.getElementById('loader');
    if (loader) loader.style.display = 'flex';
}
function hideLoader() {
    const loader = document.getElementById('loader');
    if (loader) loader.style.display = 'none';
}

function openModal(title, bodyHtml) {
    const modal = document.getElementById('modal');
    if (modal) {
        document.getElementById('modal-title').innerText = title;
        document.getElementById('modal-body').innerHTML = bodyHtml;
        modal.style.display = 'flex';
    }
}
function closeModal() {
    const modal = document.getElementById('modal');
    if (modal) modal.style.display = 'none';
}
// Expondo closeModal globalmente para o onclick no HTML
window.closeModal = closeModal;

// FILTRO DE CHAPÉUS (categoria e preço)
document.addEventListener('DOMContentLoaded', () => {
    const filterForm = document.getElementById('filter-form');
    if (filterForm) {
        // Função de filtro
        function aplicarFiltro() {
            const categorias = Array.from(filterForm.querySelectorAll('input[name="categoria"]:checked')).map(cb => cb.value);
            let precoMin = document.getElementById('preco-min').value.trim();
            let precoMax = document.getElementById('preco-max').value.trim();
            precoMin = precoMin === '' ? null : Math.max(0, parseInt(precoMin, 10));
            precoMax = precoMax === '' ? null : Math.max(0, parseInt(precoMax, 10));
            // Envia para o backend como centavos zerados (ex: 120 -> 120.00)
            const precoMinStr = precoMin !== null ? (precoMin + '.00') : null;
            const precoMaxStr = precoMax !== null ? (precoMax + '.00') : null;

            // Validação: mínimo não pode ser maior que máximo
            if (precoMin !== null && precoMax !== null && precoMin > precoMax) {
                openModal('Filtro inválido', '<p>O valor mínimo não pode ser maior que o valor máximo.</p>');
                return;
            }

            // Termo de busca
            const termo = (document.getElementById('search-hat')?.value || '').toLowerCase();

            // Monta query string
            let params = [];
            if (categorias.length > 0) params.push('categoria=' + encodeURIComponent(categorias.join(',')));
            if (precoMinStr !== null) params.push('min=' + precoMinStr);
            if (precoMaxStr !== null) params.push('max=' + precoMaxStr);
            let url = '/hats';
            if (params.length > 0) url += '?' + params.join('&');

            showLoader();
            fetchApi(url).then(hats => {
                // Filtro de pesquisa (nome) ainda é feito no frontend
                let filtrados = hats.filter(hat => termo === '' || hat.nome.toLowerCase().includes(termo));
                exibirChapeus(filtrados);
            }).catch(() => {
                openModal('Erro', '<p>Não foi possível carregar os chapéus filtrados.</p>');
            }).finally(() => {
                hideLoader();
            });
        }

        // Aplica filtro ao marcar/desmarcar qualquer categoria
        filterForm.querySelectorAll('input[name="categoria"]').forEach(cb => {
            cb.addEventListener('change', aplicarFiltro);
        });

        // Aplica filtro de faixa de preço apenas ao clicar em "Aplicar"
        filterForm.addEventListener('submit', function (e) {
            e.preventDefault();
            aplicarFiltro();
        });
    }
});