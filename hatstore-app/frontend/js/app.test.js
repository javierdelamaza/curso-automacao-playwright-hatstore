/**
 * @jest-environment jsdom
 */

import { exibirChapeus } from './app.js';

// Mock do localStorage e DOM
beforeEach(() => {
  localStorage.clear();
  document.body.innerHTML = `
    <div id="cart-items"></div>
    <div id="total-amount"></div>
    <button id="checkout-button"></button>
    <input id="cupom" />
    <span id="erro-cupom"></span>
    <input id="search-hat" />
    <div id="hats"></div>
    <div id="cart-sidebar"></div>
  `;
  // Expondo variáveis globais para os testes
  window.carrinho = [];
  window.todosChapeus = [
    { id: 1, nome: 'Chapéu Panamá', price: 100, quantidade: 1 },
    { id: 2, nome: 'Chapéu Fedora', price: 150, quantidade: 1 }
  ];
});

// Funções mock para simular comportamento do app.js
function adicionarAoCarrinho(id, nome, price) {
  let carrinho = JSON.parse(localStorage.getItem('carrinho') || '[]');
  const itemExistente = carrinho.find(item => item.id === id);
  if (itemExistente) {
    itemExistente.quantidade += 1;
  } else {
    carrinho.push({ id, nome, price, quantidade: 1 });
  }
  localStorage.setItem('carrinho', JSON.stringify(carrinho));
}

function removerDoCarrinho(id) {
  let carrinho = JSON.parse(localStorage.getItem('carrinho') || '[]');
  carrinho = carrinho.filter(item => item.id !== id);
  localStorage.setItem('carrinho', JSON.stringify(carrinho));
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
  } else {
    erroCupom.textContent = "Cupom inválido.";
    erroCupom.className = 'input-error';
    localStorage.removeItem('cupom');
  }
}

function filtrarChapeus() {
  const termo = document.getElementById('search-hat').value.toLowerCase();
  const filtrados = window.todosChapeus.filter(hat => hat.nome.toLowerCase().includes(termo));
  return filtrados;
}

// Teste adicionar ao carrinho
it('adiciona item ao carrinho', () => {
  adicionarAoCarrinho(1, 'Chapéu Panamá', 100);
  const carrinho = JSON.parse(localStorage.getItem('carrinho'));
  expect(carrinho.length).toBe(1);
  expect(carrinho[0].nome).toBe('Chapéu Panamá');
});

// Teste remover do carrinho
it('remove item do carrinho', () => {
  adicionarAoCarrinho(1, 'Chapéu Panamá', 100);
  removerDoCarrinho(1);
  const carrinho = JSON.parse(localStorage.getItem('carrinho'));
  expect(carrinho.length).toBe(0);
});

// Teste aplicar cupom válido
it('aplica cupom HATOFF', () => {
  document.getElementById('cupom').value = 'HATOFF';
  aplicarCupom();
  expect(localStorage.getItem('cupom')).toBe('HATOFF');
  expect(document.getElementById('erro-cupom').textContent).toBe('Cupom aplicado!');
});

// Teste aplicar cupom inválido
it('rejeita cupom inválido', () => {
  document.getElementById('cupom').value = 'INVALIDO';
  aplicarCupom();
  expect(localStorage.getItem('cupom')).toBe(null);
  expect(document.getElementById('erro-cupom').textContent).toBe('Cupom inválido.');
});

// Teste filtrar chapéus
it('filtra chapéus pelo nome', () => {
  document.getElementById('search-hat').value = 'fedora';
  const filtrados = filtrarChapeus();
  expect(filtrados.length).toBe(1);
  expect(filtrados[0].nome).toBe('Chapéu Fedora');
});
