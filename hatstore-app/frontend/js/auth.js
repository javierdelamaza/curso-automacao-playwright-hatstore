
import { fetchApi } from './api.js';

document.addEventListener('DOMContentLoaded', () => {
            // Limpa erro visual e mensagem ao focar no input
            function setupLimparErroAoFocar(prefix, campos) {
                campos.forEach(campo => {
                    const input = document.getElementById(`${prefix}-${campo}`);
                    const erro = document.getElementById(`erro-${prefix}-${campo}`);
                    if (input && erro) {
                        input.addEventListener('focus', () => {
                            erro.textContent = '';
                            input.classList.remove('input-invalid');
                        });
                    }
                });
            }
            setupLimparErroAoFocar('login', ['email', 'password']);
            setupLimparErroAoFocar('register', ['email', 'password', 'nome', 'cpf', 'dob']);
        // Função utilitária para limpar mensagens de erro dos campos
        function limparErrosAuth(prefix) {
            const campos = ['email', 'password', 'nome', 'cpf', 'dob'];
            campos.forEach(campo => {
                const erroEl = document.getElementById(`erro-${prefix}-${campo}`);
                if (erroEl) erroEl.textContent = '';
                const inputEl = document.getElementById(`${prefix}-${campo}`);
                if (inputEl) inputEl.classList.remove('input-invalid');
            });
        }

        // Função utilitária para exibir erro embaixo do campo
        function exibirErroCampoAuth(prefix, campo, mensagem) {
            const erroEl = document.getElementById(`erro-${prefix}-${campo}`);
            const inputEl = document.getElementById(`${prefix}-${campo}`);
            if (erroEl) erroEl.textContent = mensagem;
            if (inputEl) inputEl.classList.add('input-invalid');
        }
    // Seleciona os elementos do DOM
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');
    const showRegisterButton = document.getElementById('show-register-form');
    const showLoginButton = document.getElementById('show-login-form');
    
    const loginMessageEl = document.getElementById('login-message');
    const registerMessageEl = document.getElementById('register-message');

    // Verifica feature flag new_hatstore e exibe campos extras
    const isNewHatstore = window.FEATURE_FLAGS && window.FEATURE_FLAGS.new_hatstore;
    if (isNewHatstore) {
        const nomeGroup = document.getElementById('register-nome-group');
        const nomeInput = document.getElementById('register-nome');
        const cpfGroup = document.getElementById('register-cpf-group');
        const dobGroup = document.getElementById('register-dob-group');
        const cpfInput = document.getElementById('register-cpf');
        const dobInput = document.getElementById('register-dob');
        if (nomeGroup) nomeGroup.style.display = 'block';
        if (nomeInput) nomeInput.required = true;
        if (cpfGroup) cpfGroup.style.display = 'block';
        if (dobGroup) dobGroup.style.display = 'block';
        if (cpfInput) {
            cpfInput.required = true;
            if (window.jQuery && jQuery.fn.mask) {
                jQuery(cpfInput).mask('000.000.000-00');
            }
        }
        if (dobInput) dobInput.required = true;

        // Atualizar requisitos de senha para o novo modelo
        const passwordInput = document.getElementById('register-password');
        const passwordHelp = document.getElementById('password-help');
        if (passwordInput) passwordInput.minLength = 8;
        if (passwordHelp) passwordHelp.textContent = 'A senha deve ter no mínimo 8 caracteres, incluindo maiúscula, minúscula, número e caractere especial.';
    }

    // Função para alternar a visibilidade dos formulários usando display
    const toggleForms = () => {
        if (loginForm.style.display === 'none') {
            loginForm.style.display = 'block';
            registerForm.style.display = 'none';
        } else {
            loginForm.style.display = 'none';
            registerForm.style.display = 'block';
        }
    };
    window.toggleForms = toggleForms;

    // Event listeners para os botões de alternância
    showRegisterButton.addEventListener('click', toggleForms);
    showLoginButton.addEventListener('click', toggleForms);

    // Função para exibir mensagens de feedback de forma acessível
    const showMessage = (element, message, isError = true) => {
        element.textContent = message;
        element.className = 'auth-message'; // Reset class
        element.classList.add(isError ? 'error' : 'success');
        element.style.display = 'block';
    };

    // Event listener para o formulário de REGISTRO
    registerForm.addEventListener('submit', async (event) => {
        event.preventDefault();
        limparErrosAuth('register');
        let isValid = true;

        const email = document.getElementById('register-email').value.trim();
        const password = document.getElementById('register-password').value;
        let nome = '', cpf = '', dob = '';
        if (isNewHatstore) {
            nome = document.getElementById('register-nome').value.trim();
            cpf = document.getElementById('register-cpf').value.trim();
            dob = document.getElementById('register-dob').value.trim();
        }

        // Validação dos campos
        if (!email.match(/^[^@\s]+@[^@\s]+\.[^@\s]+$/)) {
            exibirErroCampoAuth('register', 'email', 'Digite um e-mail válido.');
            isValid = false;
        }
        if (!password || password.length < (isNewHatstore ? 8 : 6)) {
            exibirErroCampoAuth('register', 'password', isNewHatstore ? 'A senha deve ter no mínimo 8 caracteres.' : 'A senha deve ter no mínimo 6 caracteres.');
            isValid = false;
        }
        if (isNewHatstore) {
            if (!nome || nome.length < 3) {
                exibirErroCampoAuth('register', 'nome', 'Nome deve ter pelo menos 3 caracteres.');
                isValid = false;
            }
            if (!cpf.match(/^\d{3}\.\d{3}\.\d{3}-\d{2}$/)) {
                exibirErroCampoAuth('register', 'cpf', 'CPF inválido. Formato: 000.000.000-00');
                isValid = false;
            }
            if (!dob) {
                exibirErroCampoAuth('register', 'dob', 'Preencha a data de nascimento.');
                isValid = false;
            }
        }
        if (!isValid) return;

        const bodyData = { email, password };
        if (isNewHatstore) {
            bodyData.nome = nome;
            bodyData.cpf = cpf;
            bodyData.data_nascimento = dob;
        }

        try {
            await fetchApi('/auth/register', {
                method: 'POST',
                body: JSON.stringify(bodyData),
            });

            showMessage(registerMessageEl, 'Registro bem-sucedido! Você já pode fazer o login.', false);
            setTimeout(() => {
                toggleForms();
                document.getElementById('login-email').value = email;
                document.getElementById('login-password').focus();
            }, 2000);
        } catch (error) {
            showMessage(registerMessageEl, error.message);
        }
    });

    // Event listener para o formulário de LOGIN
    loginForm.addEventListener('submit', async (event) => {
        event.preventDefault();
        limparErrosAuth('login');
        let isValid = true;

        const email = document.getElementById('login-email').value.trim();
        const password = document.getElementById('login-password').value;
        const rememberMe = document.getElementById('remember-me').checked;

        if (!email.match(/^[^@\s]+@[^@\s]+\.[^@\s]+$/)) {
            exibirErroCampoAuth('login', 'email', 'Digite um e-mail válido.');
            isValid = false;
        }
        if (!password) {
            exibirErroCampoAuth('login', 'password', 'Preencha a senha.');
            isValid = false;
        }
        if (!isValid) return;

        try {
            const data = await fetchApi('/auth/login', {
                method: 'POST',
                body: JSON.stringify({ email, password, rememberMe }),
            });

            localStorage.setItem('jwt_token', data.token);
            localStorage.setItem('user_email', email);

            showMessage(loginMessageEl, 'Login bem-sucedido! Redirecionando...', false);
            setTimeout(() => {
                window.location.href = 'index.html';
            }, 1500);

        } catch (error) {
            showMessage(loginMessageEl, `Erro no login: ${error.message || 'Credenciais inválidas.'}`);
        }
    });
});