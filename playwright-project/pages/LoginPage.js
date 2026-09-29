// Page Object da tela de login/cadastro (/auth.html).
export class LoginPage {
  constructor(page) {
    this.page = page;

    this.loginForm = page.locator('#login-form');
    this.loginEmailInput = page.locator('#login-email');
    this.loginPasswordInput = page.locator('#login-password');
    this.loginSubmitBtn = this.loginForm.getByRole('button', { name: 'Entrar' });
    this.loginMessage = page.locator('#login-message');

    this.showRegisterFormBtn = page.locator('#show-register-form');

    this.registerForm = page.locator('#register-form');
    this.registerEmailInput = page.locator('#register-email');
    this.registerNomeInput = page.locator('#register-nome');
    this.registerCpfInput = page.locator('#register-cpf');
    this.registerDobInput = page.locator('#register-dob');
    this.registerPasswordInput = page.locator('#register-password');
    this.registerSubmitBtn = this.registerForm.getByRole('button', { name: 'Registrar' });
    this.registerMessage = page.locator('#register-message');
  }

  async goto() {
    await this.page.goto('/auth.html');
  }

  async login(email, password) {
    await this.loginEmailInput.fill(email);
    await this.loginPasswordInput.fill(password);
    await this.loginSubmitBtn.click();
  }

  async goToRegisterForm() {
    await this.showRegisterFormBtn.click();
    await this.registerForm.waitFor({ state: 'visible' });
  }

  async register({ email, nome, cpf, dataNascimento, password }) {
    await this.registerEmailInput.fill(email);
    await this.registerNomeInput.fill(nome);
    // O campo tem máscara: digitar tecla a tecla simula o usuário real.
    await this.registerCpfInput.pressSequentially(cpf.replace(/\D/g, ''));
    await this.registerDobInput.fill(dataNascimento);
    await this.registerPasswordInput.fill(password);
    await this.registerSubmitBtn.click();
  }
}
