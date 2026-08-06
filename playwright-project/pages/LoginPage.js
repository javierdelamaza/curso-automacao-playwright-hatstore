// Page Object da tela de login/cadastro (auth.html)
export class LoginPage {
  constructor(page) {
    this.page = page;

    // Login
    this.loginEmailInput = page.locator('#login-email');
    this.loginPasswordInput = page.locator('#login-password');
    this.loginForm = page.locator('#login-form');
    this.loginSubmitBtn = this.loginForm.getByRole('button', { name: 'Entrar' });
    this.loginMessage = page.locator('#login-message');

    // Alternar para o formulário de cadastro
    this.showRegisterFormBtn = page.locator('#show-register-form');

    // Cadastro
    this.registerForm = page.locator('#register-form');
    this.registerEmailInput = page.locator('#register-email');
    this.registerNomeInput = page.locator('#register-nome');
    this.registerCpfInput = page.locator('#register-cpf');
    this.registerDobInput = page.locator('#register-dob');
    this.registerPasswordInput = page.locator('#register-password');
    this.registerSubmitBtn = this.registerForm.getByRole('button', { name: 'Registrar' });
    this.registerMessage = page.locator('#register-message');

    this.authLinks = page.locator('#auth-links');
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

    // Campos de nome/CPF/nascimento só aparecem quando a feature flag
    // "new_hatstore" está ativa — preenche apenas se estiverem visíveis.
    if (await this.registerNomeInput.isVisible()) {
      await this.registerNomeInput.fill(nome);
      await this.registerCpfInput.fill(cpf);
      await this.registerDobInput.fill(dataNascimento);
    }

    await this.registerPasswordInput.fill(password);
    await this.registerSubmitBtn.click();
  }
}
