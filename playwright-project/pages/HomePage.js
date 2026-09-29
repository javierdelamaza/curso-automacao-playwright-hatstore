// Page Object da home: catálogo, busca, filtros e cabeçalho.
export class HomePage {
  constructor(page) {
    this.page = page;

    this.searchInput = page.getByPlaceholder('Buscar chapéu...');
    // Cards não têm role/label: bom candidato a pedir `data-testid` ao dev.
    this.hatCards = page.locator('.hat-card');
    this.hatCardTitles = page.locator('.hat-card h3');
    this.hatCardPrices = page.locator('.hat-card p', { hasText: 'Preço: R$' });
    this.nenhumResultado = page.getByText('Nenhum chapéu encontrado.');

    this.filterForm = page.locator('#filter-form');
    this.precoMinInput = this.filterForm.getByLabel('Mínimo');
    this.precoMaxInput = this.filterForm.getByLabel('Máximo');
    this.aplicarFiltroBtn = this.filterForm.getByRole('button', { name: 'Aplicar' });

    this.cartButton = page.getByRole('button', { name: 'Carrinho' });
    this.meusPedidosLink = page.getByRole('link', { name: 'Meus Pedidos' });
    this.loginLink = page.locator('#auth-links').getByRole('link', { name: 'Login' });
    this.logoutButton = page.locator('#auth-links').getByRole('button', { name: 'Logout' });
  }

  // O texto do checkbox é diferente do valor: "importado" aparece como "Internacional".
  categoriaCheckbox(categoria) {
    const rotulos = { nacional: 'Nacional', importado: 'Internacional', crescer: 'Crescer' };
    return this.filterForm.getByLabel(rotulos[categoria], { exact: true });
  }

  cardDoProduto(nome) {
    return this.hatCards.filter({ has: this.page.getByRole('heading', { name: nome, exact: true }) });
  }

  async goto() {
    await this.page.goto('/');
    await this.hatCards.first().waitFor();
  }

  async search(term) {
    await this.searchInput.fill(term);
  }

  async filterByCategoria(...categorias) {
    for (const categoria of categorias) {
      await this.categoriaCheckbox(categoria).check();
    }
    await this.aplicarFiltroBtn.click();
  }

  async filterByPreco({ min, max }) {
    if (min !== undefined) await this.precoMinInput.fill(String(min));
    if (max !== undefined) await this.precoMaxInput.fill(String(max));
    await this.aplicarFiltroBtn.click();
  }

  async adicionarAoCarrinho(nome) {
    await this.cardDoProduto(nome).getByRole('button', { name: 'Adicionar ao carrinho' }).click();
  }

  async precosExibidos() {
    const textos = await this.hatCardPrices.allTextContents();
    return textos.map((t) => Number(t.replace('Preço: R$', '').trim()));
  }
}
