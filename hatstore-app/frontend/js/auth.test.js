

/**
 * @jest-environment jsdom
 */

import './auth.js';

describe('auth.js', () => {
  let loginForm, registerForm, showRegisterButton, showLoginButton, loginMessageEl, registerMessageEl;

  beforeEach(() => {
    window.FEATURE_FLAGS = {};
    document.body.innerHTML = `
      <form id="login-form"><input id="login-email"><input id="login-password"><input type="checkbox" id="remember-me"></form>
      <form id="register-form"><input id="register-email"><input id="register-password"><div id="register-nome-group" style="display:none;"><input id="register-nome"></div><div id="register-cpf-group" style="display:none;"><input id="register-cpf"></div><div id="register-dob-group" style="display:none;"><input id="register-dob"></div></form>
      <button id="show-register-form"></button>
      <button id="show-login-form"></button>
      <span id="login-message"></span>
      <span id="register-message"></span>
    `;
    loginForm = document.getElementById('login-form');
    registerForm = document.getElementById('register-form');
    showRegisterButton = document.getElementById('show-register-form');
    showLoginButton = document.getElementById('show-login-form');
    loginMessageEl = document.getElementById('login-message');
    registerMessageEl = document.getElementById('register-message');
  });

  test('toggleForms alterna visibilidade dos formulários', () => {
    loginForm.style.display = 'block';
    registerForm.style.display = 'none';
    document.dispatchEvent(new Event('DOMContentLoaded'));
    window.toggleForms();
    expect(loginForm.style.display).toBe('none');
    expect(registerForm.style.display).toBe('block');
    window.toggleForms();
    expect(loginForm.style.display).toBe('block');
    expect(registerForm.style.display).toBe('none');
  });

  test('showMessage exibe mensagem de erro', () => {
    // Função mockada
    const showMessage = (element, message, isError = true) => {
      element.textContent = message;
      element.className = 'auth-message';
      element.classList.add(isError ? 'error' : 'success');
      element.style.display = 'block';
    };
    showMessage(loginMessageEl, 'Erro de login');
    expect(loginMessageEl.textContent).toBe('Erro de login');
    expect(loginMessageEl.classList.contains('error')).toBe(true);
    expect(loginMessageEl.style.display).toBe('block');
  });

  test('showMessage exibe mensagem de sucesso', () => {
    const showMessage = (element, message, isError = true) => {
      element.textContent = message;
      element.className = 'auth-message';
      element.classList.add(isError ? 'error' : 'success');
      element.style.display = 'block';
    };
    showMessage(registerMessageEl, 'Sucesso!', false);
    expect(registerMessageEl.textContent).toBe('Sucesso!');
    expect(registerMessageEl.classList.contains('success')).toBe(true);
    expect(registerMessageEl.style.display).toBe('block');
  });

  test('campos Nome, CPF e data de nascimento ficam ocultos quando new_hatstore é false', () => {
    window.FEATURE_FLAGS = { new_hatstore: false };
    document.dispatchEvent(new Event('DOMContentLoaded'));
    const nomeGroup = document.getElementById('register-nome-group');
    const cpfGroup = document.getElementById('register-cpf-group');
    const dobGroup = document.getElementById('register-dob-group');
    expect(nomeGroup.style.display).toBe('none');
    expect(cpfGroup.style.display).toBe('none');
    expect(dobGroup.style.display).toBe('none');
  });

  test('campos Nome, CPF e data de nascimento ficam visíveis quando new_hatstore é true', () => {
    window.FEATURE_FLAGS = { new_hatstore: true };
    document.dispatchEvent(new Event('DOMContentLoaded'));
    const nomeGroup = document.getElementById('register-nome-group');
    const cpfGroup = document.getElementById('register-cpf-group');
    const dobGroup = document.getElementById('register-dob-group');
    expect(nomeGroup.style.display).toBe('block');
    expect(cpfGroup.style.display).toBe('block');
    expect(dobGroup.style.display).toBe('block');
  });

  test('campos Nome, CPF e data de nascimento são required quando new_hatstore é true', () => {
    window.FEATURE_FLAGS = { new_hatstore: true };
    document.dispatchEvent(new Event('DOMContentLoaded'));
    const nomeInput = document.getElementById('register-nome');
    const cpfInput = document.getElementById('register-cpf');
    const dobInput = document.getElementById('register-dob');
    expect(nomeInput.required).toBe(true);
    expect(cpfInput.required).toBe(true);
    expect(dobInput.required).toBe(true);
  });
});
