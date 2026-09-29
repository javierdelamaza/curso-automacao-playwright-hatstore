// Page Object do checkout (layout "novo", ativo por padrão no Hat Store).
export class CheckoutPage {
  constructor(page) {
    this.page = page;
    this.form = page.locator('#checkout-form-new');

    this.usuarioNome = page.locator('#new-user-nome');
    this.usuarioCpf = page.locator('#new-user-cpf');
    this.resumo = page.locator('#pedido-resumo-new');

    this.telefoneInput = this.form.getByLabel('Telefone');
    this.ruaInput = this.form.getByLabel('Rua');
    this.numeroInput = this.form.getByLabel('Número');
    this.bairroInput = this.form.getByLabel('Bairro');
    this.cepInput = this.form.getByLabel('CEP');
    this.cidadeInput = this.form.getByLabel('Cidade');
    this.ufInput = this.form.getByLabel('UF');

    this.freteInfo = page.locator('#frete-info');
    this.freteRegiao = page.locator('#frete-regiao');
    this.freteValor = page.locator('#frete-valor');

    this.pixBtn = page.locator('#new-btn-pix');
    this.boletoBtn = page.locator('#new-btn-boleto');
    this.finalizarBtn = this.form.getByRole('button', { name: /Finalizar Compra/ });

    this.modalTitulo = page.locator('#modal-title');
    this.codigoPagamento = page.locator('#codigo-pagamento');
  }

  async preencherEndereco({ telefone, endereco, numero, bairro, cep, cidade, uf }) {
    await this.telefoneInput.pressSequentially(telefone.replace(/\D/g, ''));
    await this.ruaInput.fill(endereco);
    await this.numeroInput.fill(numero);
    await this.bairroInput.fill(bairro);
    await this.cepInput.pressSequentially(cep.replace(/\D/g, ''));
    // O frete é calculado quando o campo CEP perde o foco.
    await this.cidadeInput.fill(cidade);
    await this.ufInput.fill(uf);
  }

  async escolherPagamento(tipo) {
    await (tipo === 'pix' ? this.pixBtn : this.boletoBtn).click();
  }

  async finalizar() {
    await this.finalizarBtn.click();
  }
}
