// Page Object do carrinho lateral (abre na própria home).
export class CartPage {
  constructor(page) {
    this.page = page;
    this.sidebar = page.locator('#cart-sidebar');
    this.itens = this.sidebar.locator('.cart-item');
    this.itemNomes = this.sidebar.locator('.cart-item-nome');
    this.vazio = this.sidebar.getByText('Seu carrinho está vazio.');
    this.cupomInput = this.sidebar.getByLabel('Cupom de desconto:');
    this.aplicarCupomBtn = this.sidebar.locator('#aplicar-cupom');
    this.mensagemCupom = this.sidebar.locator('#erro-cupom');
    this.subtotal = this.sidebar.locator('.cart-subtotal span');
    this.desconto = this.sidebar.locator('.cart-cupom-info');
    this.total = this.sidebar.locator('.cart-total span');
    this.finalizarCompraBtn = this.sidebar.getByRole('button', { name: 'Finalizar compra' });
  }

  item(nome) {
    return this.itens.filter({ hasText: nome });
  }

  async aumentarQuantidade(nome) {
    await this.item(nome).getByTitle('Aumentar').click();
  }

  async remover(nome) {
    await this.item(nome).getByTitle('Remover').click();
  }

  async aplicarCupom(codigo) {
    await this.cupomInput.fill(codigo);
    await this.aplicarCupomBtn.click();
  }

  async finalizarCompra() {
    await this.finalizarCompraBtn.click();
    await this.page.waitForURL('**/checkout.html');
  }
}
