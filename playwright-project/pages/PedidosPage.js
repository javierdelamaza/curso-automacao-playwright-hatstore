// Page Object de "Meus Pedidos" (/pedidos.html, exige usuário logado).
export class PedidosPage {
  constructor(page) {
    this.page = page;
    this.cards = page.locator('.pedidos-new-card');
    this.vazio = page.locator('#pedidos-new-empty');
  }

  async goto() {
    await this.page.goto('/pedidos.html');
  }

  cardDoPedido(id) {
    return this.cards.filter({ hasText: `Pedido #${id}` });
  }
}
