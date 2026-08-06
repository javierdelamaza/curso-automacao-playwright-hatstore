import { test, expect } from '@playwright/test';
import { HomePage } from '../pages/HomePage.js';

// Exercícios do módulo 2 ("dever de casa"). Remova o `test.fixme` e
// implemente a asserção proposta no comentário de cada teste.

test.describe('Home — filtros', () => {
  test.fixme(
    'filtrar por categoria "crescer" exibe apenas produtos dessa categoria',
    async ({ page }) => {
      const home = new HomePage(page);
      await home.goto();

      await home.filterByCategoria('crescer');

      // Dica: a massa de dados seed (backend/migrations/001_init.sql) só tem
      // 1 produto na categoria "crescer" ("Chapéu Gustavo Carvalho").
      await expect(home.hatCards).toHaveCount(1);
      await expect(home.hatCardTitles.first()).toContainText('Gustavo Carvalho');
    }
  );

  test.fixme(
    'filtrar por faixa de preço retorna apenas produtos dentro do intervalo',
    async ({ page }) => {
      const home = new HomePage(page);
      await home.goto();

      await home.filterByPreco({ min: 100, max: 130 });

      // Dica: valide que cada card exibido tem preço entre 100 e 130.
      // Pode usar home.page.locator('.hat-card p') para ler o texto "Preço: R$ ...".
    }
  );
});
