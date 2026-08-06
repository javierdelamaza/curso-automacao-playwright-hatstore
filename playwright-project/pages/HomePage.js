// Page Object da home (catálogo, busca e filtros)
export class HomePage {
  constructor(page) {
    this.page = page;
    this.searchInput = page.locator('#search-hat');
    this.hatCards = page.locator('.hat-card');
    this.hatCardTitles = page.locator('.hat-card h3');
    this.categoriaCheckbox = (categoria) =>
      page.locator(`.filter-cat-input[value="${categoria}"]`);
    this.precoMinInput = page.locator('#preco-min');
    this.precoMaxInput = page.locator('#preco-max');
    this.aplicarFiltroBtn = page.locator('.filter-apply-btn');
    this.cartButton = page.locator('#header-cart-btn');
  }

  async goto() {
    await this.page.goto('/');
    // espera o catálogo inicial carregar antes de qualquer interação
    await this.hatCards.first().waitFor();
  }

  async search(term) {
    await this.searchInput.fill(term);
  }

  async filterByCategoria(categoria) {
    await this.categoriaCheckbox(categoria).check();
    await this.aplicarFiltroBtn.click();
  }

  async filterByPreco({ min, max }) {
    if (min !== undefined) await this.precoMinInput.fill(String(min));
    if (max !== undefined) await this.precoMaxInput.fill(String(max));
    await this.aplicarFiltroBtn.click();
  }
}
