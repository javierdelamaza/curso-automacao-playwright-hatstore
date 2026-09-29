    // Limpa erro visual e mensagem ao focar nos inputs do checkout clássico
    [
        'nome', 'cpf', 'email', 'telefone', 'endereco', 'numero', 'bairro', 'cep', 'cidade', 'uf', 'pagamento'
    ].forEach(id => {
        const input = document.getElementById(id);
        const erro = document.getElementById(`erro-${id}`);
        if (input && erro) {
            input.addEventListener('focus', () => {
                erro.textContent = '';
                input.classList.remove('input-invalid');
            });
        }
    });

    // Limpa erro visual e mensagem ao focar nos inputs do checkout novo
    [
        'new-telefone', 'new-endereco', 'new-numero', 'new-bairro', 'new-cep', 'new-cidade', 'new-uf', 'new-pagamento'
    ].forEach(id => {
        const input = document.getElementById(id);
        const erro = document.getElementById(`erro-${id}`);
        if (input && erro) {
            input.addEventListener('focus', () => {
                erro.textContent = '';
                input.classList.remove('input-invalid');
            });
        }
    });
import { fetchApi, escapeHtml } from './api.js';
import { showLoader, hideLoader, openModal, closeModal } from './app.js';

// Dados do perfil do usuário (preenchido via /api/me)
let userProfile = null;

// Dados do frete calculado
let freteData = null;

document.addEventListener('DOMContentLoaded', async function () {
    // Redireciona para login se não estiver autenticado
    if (!localStorage.getItem('jwt_token')) {
        window.location.href = 'auth.html';
        return;
    }

    const isNewHatstore = window.FEATURE_FLAGS && window.FEATURE_FLAGS.new_hatstore;
    const carrinho = JSON.parse(localStorage.getItem('carrinho') || '[]');

    if (isNewHatstore) {
        // Ativa o checkout novo
        document.getElementById('checkout-old').style.display = 'none';
        document.getElementById('checkout-new').style.display = 'block';

        // Busca perfil do usuário
        try {
            userProfile = await fetchApi('/me');
        } catch (e) {
            userProfile = null;
        }

        // Preenche dados do usuário
        preencherDadosUsuario(userProfile);

        // Monta resumo no checkout novo
        montarResumoPedidoNew(carrinho);

        // Máscaras no checkout novo
        aplicarMascarasNew();

        // Listener para calcular frete ao preencher CEP
        const cepInput = document.getElementById('new-cep');
        if (cepInput) {
            cepInput.addEventListener('blur', () => calcularFreteNew(carrinho));
        }

        // Event listeners
        const formNew = document.getElementById('checkout-form-new');
        if (formNew) {
            formNew.addEventListener('submit', handleFormSubmitNew);
        }
        const btnPix = document.getElementById('new-btn-pix');
        const btnBoleto = document.getElementById('new-btn-boleto');
        if (btnPix && btnBoleto) {
            btnPix.addEventListener('click', () => selectPagamentoNew('pix'));
            btnBoleto.addEventListener('click', () => selectPagamentoNew('boleto'));
        }
    } else {
        // Checkout clássico
        document.getElementById('checkout-old').style.display = 'block';
        document.getElementById('checkout-new').style.display = 'none';

        montarResumoPedido(carrinho);
        aplicarMascaras();

        const form = document.getElementById('checkout-form');
        if (form) {
            form.addEventListener('submit', handleFormSubmit);
        }
        const btnPix = document.getElementById('btn-pix');
        const btnBoleto = document.getElementById('btn-boleto');
        if (btnPix && btnBoleto) {
            btnPix.addEventListener('click', () => selectPagamento('pix'));
            btnBoleto.addEventListener('click', () => selectPagamento('boleto'));
        }
    }
});


// =========================================================================
// CHECKOUT NOVO (new_hatstore = true)
// =========================================================================

function preencherDadosUsuario(profile) {
    const nomeEl = document.getElementById('new-user-nome');
    const emailEl = document.getElementById('new-user-email');
    const cpfEl = document.getElementById('new-user-cpf');
    const dobEl = document.getElementById('new-user-dob');

    if (profile) {
        if (nomeEl) nomeEl.textContent = profile.nome || '—';
        if (emailEl) emailEl.textContent = profile.email || '—';
        if (cpfEl) cpfEl.textContent = profile.cpf || '—';
        if (dobEl) dobEl.textContent = formatarData(profile.data_nascimento) || '—';
    }
}

function formatarData(dateStr) {
    if (!dateStr) return null;
    const parts = dateStr.split('-');
    if (parts.length === 3) return `${parts[2]}/${parts[1]}/${parts[0]}`;
    return dateStr;
}

function montarResumoPedidoNew(itens) {
    const div = document.getElementById('pedido-resumo-new');
    if (!div) return;

    const cupom = (localStorage.getItem('cupom') || '').toUpperCase();
    let subtotal = 0;

    let itensHtml = itens.map(item => {
        subtotal += item.price * (item.quantidade || 1);
        return `<div class="checkout-new-item">
            <span>${escapeHtml(item.nome)} <small>(${item.quantidade || 1}x)</small></span>
            <strong>R$ ${item.price.toFixed(2)}</strong>
        </div>`;
    }).join('');

    let descontoHtml = '';
    let subtotalComDesconto = subtotal;

    if (cupom === "HATOFF") {
        const desconto = subtotal * 0.20;
        subtotalComDesconto -= desconto;
        descontoHtml = `<div class="checkout-new-discount">Cupom HATOFF: <strong>-R$ ${desconto.toFixed(2)}</strong></div>`;
    }

    let freteHtml = '';
    let totalFinal = subtotalComDesconto;

    if (freteData && freteData.valor_frete != null) {
        freteHtml = `<div class="checkout-new-frete-line"><span>Frete (${escapeHtml(freteData.regiao)})</span><strong>R$ ${freteData.valor_frete.toFixed(2)}</strong></div>`;
        totalFinal += freteData.valor_frete;
    }

    div.innerHTML = `
        ${itensHtml}
        ${descontoHtml}
        <div class="checkout-new-subtotal"><span>Subtotal</span><strong>R$ ${subtotalComDesconto.toFixed(2)}</strong></div>
        ${freteHtml}
        <div class="checkout-new-total"><span>Total</span><strong>R$ ${totalFinal.toFixed(2)}</strong></div>
    `;
}

async function calcularFreteNew(carrinho) {
    const cepInput = document.getElementById('new-cep');
    const freteResult = document.getElementById('frete-result');
    const freteLoading = document.getElementById('frete-loading');
    const freteInfo = document.getElementById('frete-info');
    const freteErro = document.getElementById('frete-erro');

    if (!cepInput || !freteResult) return;

    const cep = cepInput.value.trim();
    if (!/^\d{5}-\d{3}$/.test(cep)) {
        freteData = null;
        freteResult.style.display = 'none';
        montarResumoPedidoNew(carrinho);
        return;
    }

    const quantidade = carrinho.reduce((sum, item) => sum + (item.quantidade || 1), 0);

    freteResult.style.display = 'block';
    freteLoading.style.display = 'block';
    freteInfo.style.display = 'none';
    freteErro.style.display = 'none';

    try {
        const data = await fetchApi(`/frete?cep=${encodeURIComponent(cep)}&quantidade=${quantidade}`);
        freteData = data;
        freteLoading.style.display = 'none';
        freteInfo.style.display = 'block';
        document.getElementById('frete-regiao').textContent = data.regiao;
        document.getElementById('frete-valor').textContent = `R$ ${data.valor_frete.toFixed(2)}`;
        montarResumoPedidoNew(carrinho);
    } catch (error) {
        freteData = null;
        freteLoading.style.display = 'none';
        freteErro.style.display = 'block';
        freteErro.textContent = error.message || 'Erro ao calcular frete';
        montarResumoPedidoNew(carrinho);
    }
}

function aplicarMascarasNew() {
    if (window.jQuery && window.jQuery.fn.mask) {
        $('#new-telefone').mask('(00) 00000-0000');
        $('#new-cep').mask('00000-000');
        $('#new-numero').mask('000000');
    }
}

function selectPagamentoNew(tipo) {
    document.getElementById('new-btn-pix').classList.remove('selected');
    document.getElementById('new-btn-boleto').classList.remove('selected');

    if (tipo === 'pix') document.getElementById('new-btn-pix').classList.add('selected');
    if (tipo === 'boleto') document.getElementById('new-btn-boleto').classList.add('selected');

    document.getElementById('new-pagamento').value = tipo;
    const erroEl = document.getElementById('erro-new-pagamento');
    if (erroEl) erroEl.textContent = '';
}

function validarFormularioNew() {
    let isValid = true;
    document.querySelectorAll('#checkout-form-new .input-error').forEach(span => span.textContent = '');
    document.querySelectorAll('#checkout-form-new .input-invalid').forEach(input => input.classList.remove('input-invalid'));

    isValid = validaCampo('new-telefone', /^\(\d{2}\) \d{5}-\d{4}$/, 'Telefone inválido. Formato: (99) 99999-9999') && isValid;
    isValid = validaCampo('new-endereco', /.+/, 'Preencha a rua') && isValid;
    isValid = validaCampo('new-numero', /^\d+$/, 'Digite apenas números') && isValid;
    isValid = validaCampo('new-bairro', /.+/, 'Preencha o bairro') && isValid;
    isValid = validaCampo('new-cep', /^\d{5}-\d{3}$/, 'CEP inválido. Formato: 00000-000') && isValid;
    isValid = validaCampo('new-cidade', /.+/, 'Preencha a cidade') && isValid;
    isValid = validaCampo('new-uf', /^[A-Za-z]{2}$/, 'Digite a sigla do estado (ex: SP)') && isValid;
    isValid = validaCampo('new-pagamento', /.+/, 'Selecione a forma de pagamento') && isValid;

    return isValid;
}

function montarObjetoPedidoNew() {
    const carrinho = JSON.parse(localStorage.getItem('carrinho') || '[]');
    const cupom = (localStorage.getItem('cupom') || "").toUpperCase();

    if (!carrinho || carrinho.length === 0) {
        throw new Error('Seu carrinho está vazio. Adicione produtos antes de finalizar o pedido.');
    }

    const itensPedido = carrinho.map(item => ({
        id: item.id,
        nome: item.nome,
        price: item.price,
        quantidade: item.quantidade
    }));

    let total = 0;
    carrinho.forEach(item => {
        total += item.price * (item.quantidade || 1);
    });

    if (cupom === "HATOFF") {
        total = total * 0.80;
    }

    // Nome, CPF e email vêm do perfil do usuário
    const nome = (userProfile && userProfile.nome) || '';
    const cpf = (userProfile && userProfile.cpf) || '';
    const email = (userProfile && userProfile.email) || '';

    const telefone = document.getElementById('new-telefone').value.trim();
    const endereco = document.getElementById('new-endereco').value.trim();
    const numero = document.getElementById('new-numero').value.trim();
    const bairro = document.getElementById('new-bairro').value.trim();
    const cep = document.getElementById('new-cep').value.trim();
    const cidade = document.getElementById('new-cidade').value.trim();
    const uf = document.getElementById('new-uf').value.trim();
    const pagamento = document.getElementById('new-pagamento').value.trim();

    if (!nome || !cpf || !email || !telefone || !endereco || !numero || !bairro || !cep || !cidade || !uf || !pagamento) {
        throw new Error('Por favor, preencha todos os campos antes de finalizar o pedido.');
    }

    const pedido = { nome, cpf, email, telefone, endereco, numero, bairro, cep, cidade, uf, pagamento, itens: itensPedido };
    pedido.cupom = cupom === "HATOFF" ? cupom : "";

    return pedido;
}

async function handleFormSubmitNew(event) {
    event.preventDefault();

    if (!validarFormularioNew()) {
        alert('Por favor, corrija os campos inválidos.');
        return;
    }

    let pedido;
    try {
        pedido = montarObjetoPedidoNew();
    } catch (error) {
        alert(error.message || 'Erro ao preparar pedido');
        return;
    }

    showLoader();
    try {
        await fetchApi('/pedido', {
            method: 'POST',
            body: JSON.stringify(pedido)
        });

        localStorage.removeItem('carrinho');
        localStorage.removeItem('cupom');
        exibirCodigoPagamento(pedido.pagamento);
    } catch (error) {
        if (error.message.includes('Cupom HATOFF já utilizado') || error.message.includes('HATOFF')) {
            exibirModalCompraNegada();
        } else if (error.message.toLowerCase().includes('estoque')) {
            exibirModalEstoqueInsuficiente(error.message);
        } else {
            console.error('Erro ao finalizar pedido:', error);
            alert(`Erro ao finalizar pedido: ${error.message}`);
        }
    } finally {
        hideLoader();
    }
}


// =========================================================================
// CHECKOUT CLÁSSICO (new_hatstore = false)
// =========================================================================

function montarResumoPedido(itens) {
    const pedidoResumoDiv = document.getElementById('pedido-resumo');
    if (!pedidoResumoDiv) return;

    const cupom = (localStorage.getItem('cupom') || '').toUpperCase();
    let total = 0;

    let itensHtml = itens.map(item => {
        total += item.price * (item.quantidade || 1);
        return `<li>${escapeHtml(item.nome)} (${item.quantidade || 1}x) - <strong>R$ ${item.price.toFixed(2)}</strong></li>`;
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

async function handleFormSubmit(event) {
    event.preventDefault();

    if (!validarFormulario()) {
        alert('Por favor, corrija os campos inválidos.');
        return;
    }

    let pedido;
    try {
        pedido = montarObjetoPedido();
    } catch (error) {
        alert(error.message || 'Erro ao preparar pedido');
        return;
    }

    showLoader();
    try {
        await fetchApi('/pedido', {
            method: 'POST',
            body: JSON.stringify(pedido)
        });

        localStorage.removeItem('carrinho');
        localStorage.removeItem('cupom');
        exibirCodigoPagamento(pedido.pagamento);
    } catch (error) {
        if (error.message.includes('Cupom HATOFF já utilizado') || error.message.includes('HATOFF')) {
            exibirModalCompraNegada();
        } else if (error.message.toLowerCase().includes('estoque')) {
            exibirModalEstoqueInsuficiente(error.message);
        } else {
            console.error('Erro ao finalizar pedido:', error);
            alert(`Erro ao finalizar pedido: ${error.message}`);
        }
    } finally {
        hideLoader();
    }
}

function montarObjetoPedido() {
    const carrinho = JSON.parse(localStorage.getItem('carrinho') || '[]');
    const cupom = (localStorage.getItem('cupom') || "").toUpperCase();

    if (!carrinho || carrinho.length === 0) {
        throw new Error('Seu carrinho está vazio. Adicione produtos antes de finalizar o pedido.');
    }

    const itensPedido = carrinho.map(item => ({
        id: item.id,
        nome: item.nome,
        price: item.price,
        quantidade: item.quantidade
    }));

    if (!itensPedido || itensPedido.length === 0) {
        throw new Error('Erro ao processar itens do carrinho. Tente novamente.');
    }

    let total = 0;
    carrinho.forEach(item => {
        total += item.price * (item.quantidade || 1);
    });

    if (cupom === "HATOFF") {
        total = total * 0.80;
    }

    const nome = document.getElementById('nome').value.trim();
    const cpf = document.getElementById('cpf').value.trim();
    const email = document.getElementById('email').value.trim();
    const telefone = document.getElementById('telefone').value.trim();
    const endereco = document.getElementById('endereco').value.trim();
    const numero = document.getElementById('numero').value.trim();
    const bairro = document.getElementById('bairro').value.trim();
    const cep = document.getElementById('cep').value.trim();
    const cidade = document.getElementById('cidade').value.trim();
    const uf = document.getElementById('uf').value.trim();
    const pagamento = document.getElementById('pagamento').value.trim();

    if (!nome || !cpf || !email || !telefone || !endereco || !numero || !bairro || !cep || !cidade || !uf || !pagamento) {
        throw new Error('Por favor, preencha todos os campos antes de finalizar o pedido.');
    }

    const pedido = { nome, cpf, email, telefone, endereco, numero, bairro, cep, cidade, uf, pagamento, itens: itensPedido };
    pedido.cupom = cupom === "HATOFF" ? cupom : "";

    return pedido;
}


// =========================================================================
// FUNÇÕES COMPARTILHADAS
// =========================================================================

function validarFormulario() {
    let isValid = true;
    document.querySelectorAll('#checkout-form .input-error').forEach(span => span.textContent = '');
    document.querySelectorAll('#checkout-form .input-invalid').forEach(input => input.classList.remove('input-invalid'));

    isValid = validaCampo('nome', /.{3,}/, 'Digite seu nome completo') && isValid;
    isValid = validaCampo('cpf', /^\d{3}\.\d{3}\.\d{3}-\d{2}$/, 'CPF inválido. Formato: 000.000.000-00') && isValid;
    isValid = validaCampo('email', /^[^@\s]+@[^@\s]+\.[^@\s]+$/, 'Digite um e-mail válido') && isValid;
    isValid = validaCampo('telefone', /^\(\d{2}\) \d{5}-\d{4}$/, 'Telefone inválido. Formato: (99) 99999-9999') && isValid;
    isValid = validaCampo('endereco', /.+/, 'Preencha a rua') && isValid;
    isValid = validaCampo('numero', /^\d+$/, 'Digite apenas números') && isValid;
    isValid = validaCampo('bairro', /.+/, 'Preencha o bairro') && isValid;
    isValid = validaCampo('cep', /^\d{5}-\d{3}$/, 'CEP inválido. Formato: 00000-000') && isValid;
    isValid = validaCampo('cidade', /.+/, 'Preencha a cidade') && isValid;
    isValid = validaCampo('uf', /^[A-Za-z]{2}$/, 'Digite a sigla do estado (ex: SP)') && isValid;
    isValid = validaCampo('pagamento', /.+/, 'Selecione a forma de pagamento') && isValid;

    return isValid;
}

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

function selectPagamento(tipo) {
    document.getElementById('btn-pix').classList.remove('selected');
    document.getElementById('btn-boleto').classList.remove('selected');

    if (tipo === 'pix') document.getElementById('btn-pix').classList.add('selected');
    if (tipo === 'boleto') document.getElementById('btn-boleto').classList.add('selected');

    document.getElementById('pagamento').value = tipo;
    document.getElementById('erro-pagamento').textContent = '';
}

function exibirCodigoPagamento(metodo) {
    const titulo = `Pagamento ${metodo === "pix" ? "PIX" : "Boleto"}`;
    const codigo = metodo === "pix"
        ? "00020126580014BR.GOV.BCB.PIX0136b1e1f2e3d4c5b6a7f8e9d0c1b2a3f4g5h6i7j8k9l5204000053039865405120.005802BR5920Had Store6009SAO PAULO61080540900062070503***6304ABCD"
        : "34191.79001 01043.510047 91020.150008 7 92180011000";

    const bodyHtml = `
        <p>Utilize o código abaixo para realizar o pagamento:</p>
        <pre id="codigo-pagamento" class="confirmacao-codigo">${codigo}</pre>
        <button id="copiar-codigo" class="auth-button" style="margin-bottom: 1rem;">Copiar código</button>
        <button id="fechar-modal-pagamento" class="auth-button" style="background: #444;">Fechar</button>
    `;
    openModal(titulo, bodyHtml);

    document.getElementById('copiar-codigo').addEventListener('click', function () {
        navigator.clipboard.writeText(codigo).then(() => {
            this.innerText = "Copiado!";
            setTimeout(() => { this.innerText = "Copiar código"; }, 2000);
        });
    });
    document.getElementById('fechar-modal-pagamento').addEventListener('click', () => {
        closeModal();
        window.location.href = 'index.html';
    });
}

function exibirModalCompraNegada() {
    const bodyHtml = `
    <p style="color:#e11d48;font-weight:bold;">O cupom HATOFF só pode ser usado uma vez por CPF.</p>
        <p>Por favor, remova o cupom no carrinho e tente finalizar sua compra novamente.</p>
        <button id="fechar-modal-negada" class="auth-button">Entendi</button>
    `;
    openModal("Compra negada", bodyHtml);
    document.getElementById('fechar-modal-negada').addEventListener('click', closeModal);
}

function aplicarMascaras() {
    if (window.jQuery && window.jQuery.fn.mask) {
        $('#cpf').mask('000.000.000-00');
        $('#telefone').mask('(00) 00000-0000');
        $('#cep').mask('00000-000');
        $('#numero').mask('000000');
    }
}

function exibirModalEstoqueInsuficiente(msg) {
    const bodyHtml = `
        <p style="color:#e11d48;font-weight:bold;">${escapeHtml(msg || 'Não há estoque suficiente para um ou mais produtos do pedido.')}</p>
        <p>Por favor, revise seu carrinho ou escolha outros produtos disponíveis.</p>
        <button id="fechar-modal-estoque" class="auth-button" style="background: #444; margin-top: 1rem;">Fechar e voltar à loja</button>
    `;
    openModal('Estoque insuficiente', bodyHtml);
    document.getElementById('fechar-modal-estoque').addEventListener('click', function () {
        let produtoId = null;
        if (msg) {
            const match = msg.match(/id\s*:?\s*(\d+)/i);
            if (match) produtoId = parseInt(match[1]);
        }
        let carrinho = JSON.parse(localStorage.getItem('carrinho') || '[]');
        if (produtoId) {
            carrinho = carrinho.filter(item => item.id !== produtoId);
            localStorage.setItem('carrinho', JSON.stringify(carrinho));
        } else {
            localStorage.removeItem('carrinho');
        }
        closeModal();
        window.location.href = 'index.html';
    });
}

export { exibirModalEstoqueInsuficiente, calcularFreteNew };
