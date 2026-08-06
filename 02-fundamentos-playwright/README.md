# Módulo 2 — Fundamentos do Playwright

## 1. O que é o Playwright

Framework open-source (Microsoft) para automação de navegadores (Chromium,
Firefox, WebKit) com uma única API, disponível em JS/TS, Python, Java e .NET.
Neste curso usamos **JavaScript**, seguindo o repositório de aulas original.

Diferenciais em relação a ferramentas mais antigas (ex.: Selenium clássico):

- **Auto-waiting**: antes de clicar/preencher, o Playwright espera o elemento
  existir, estar visível, habilitado e estável — sem `sleep` manual.
- **Web-first assertions**: `expect(locator).toBeVisible()` já espera e
  tenta novamente até um timeout, em vez de falhar na primeira checagem.
- **Isolamento de contexto**: cada teste roda em um `BrowserContext` isolado
  (cookies/local storage próprios), sem precisar de um browser novo por teste.
- **Trace viewer / vídeo / screenshot** automáticos em falhas, facilitando debug.

## 2. Instalando o projeto

```bash
mkdir meu-projeto-playwright
cd meu-projeto-playwright
npm init playwright@latest
```

Respostas recomendadas no assistente (mesmo padrão da aula original):

- `TypeScript or JavaScript?` → **JavaScript**
- `Where to put your end-to-end tests?` → **tests**
- `Add a GitHub Actions workflow?` → **false** (adicionamos manualmente depois, se quisermos)
- `Install Playwright browsers?` → **true**

O projeto de exemplo deste curso já vem pronto em
[`playwright-project/`](../playwright-project/README.md) — não precisa reinstalar,
mas é importante saber recriar esse passo em um projeto do zero.

## 3. Locators

Locator é uma referência "preguiçosa" a um elemento: só é resolvido no
navegador no momento da ação/asserção. Prioridade recomendada pela
documentação oficial (do mais para o menos recomendado):

1. `page.getByRole('button', { name: 'Entrar' })` — como um usuário/leitor de tela enxerga o elemento.
2. `page.getByLabel('Email')`, `page.getByPlaceholder(...)`, `page.getByText(...)`.
3. `page.getByTestId('...')` — quando existir `data-testid` no HTML.
4. Seletor CSS/`id` (`page.locator('#search-hat')`) — usar quando não há papel/label semântico claro.

No Hat Store, por exemplo:

```js
page.locator('#search-hat');                 // input de busca (sem label associada)
page.getByRole('button', { name: 'Entrar' }); // botão de login
page.locator('.hat-card');                   // cada card de produto
page.locator('.add-to-cart-btn').first();    // botão "Adicionar ao carrinho"
```

Evite seletores frágeis como `div > div:nth-child(3) > span`, que quebram com
qualquer mudança de layout.

## 4. Ações e assertions básicas

```js
await page.goto('/');
await page.locator('#search-hat').fill('Panamá');
await expect(page.locator('.hat-card')).toHaveCount(1);
await expect(page.locator('.hat-card h3')).toContainText('Panamá');
```

Assertions mais usadas: `toBeVisible`, `toHaveText`, `toContainText`,
`toHaveCount`, `toHaveURL`, `toHaveValue`, `toBeEnabled/Disabled`.

## 5. Page Object Model (POM)

Em vez de repetir seletores em cada teste, encapsulamos a página em uma
classe. Isso evita duplicação e centraliza a manutenção: se um `id` mudar no
front, você ajusta em **um lugar só**.

```js
// pages/HomePage.js
export class HomePage {
  constructor(page) {
    this.page = page;
    this.searchInput = page.locator('#search-hat');
    this.hatCards = page.locator('.hat-card');
  }

  async goto() {
    await this.page.goto('/');
  }

  async search(term) {
    await this.searchInput.fill(term);
  }
}
```

```js
// tests/home.spec.js
import { test, expect } from '@playwright/test';
import { HomePage } from '../pages/HomePage.js';

test('busca por produto existente retorna resultado', async ({ page }) => {
  const home = new HomePage(page);
  await home.goto();
  await home.search('Panamá');
  await expect(home.hatCards).toHaveCount(1);
});
```

## 6. Fixtures

Fixtures do Playwright injetam dependências prontas em cada teste (ex.: um
`page` já logado, ou um Page Object já instanciado), evitando repetir setup
em todo `test(...)`. Exemplo simplificado:

```js
// fixtures.js
import { test as base } from '@playwright/test';
import { HomePage } from './pages/HomePage.js';

export const test = base.extend({
  homePage: async ({ page }, use) => {
    const home = new HomePage(page);
    await home.goto();
    await use(home);
  },
});
```

## 7. Configuração (`playwright.config.js`)

Pontos essenciais:

- `baseURL`: evita repetir a URL completa em cada `page.goto(...)`.
- `use.trace`: `'on-first-retry'` grava trace apenas quando um teste falha e é re-executado — bom equilíbrio custo/benefício.
- `projects`: permite rodar a mesma suíte em Chromium/Firefox/WebKit.
- `webServer` (opcional): sobe a aplicação automaticamente antes dos testes.

## 8. Rodando os testes

```bash
npx playwright test                # roda tudo, headless
npx playwright test --headed       # com navegador visível
npx playwright test --ui           # modo interativo (UI mode)
npx playwright test home.spec.js   # arquivo específico
npx playwright show-report         # abre relatório HTML da última execução
```

## 9. Testes de API com Playwright

O Playwright também tem um cliente HTTP embutido (`request`), útil para
testes de contrato/API sem precisar de navegador:

```js
import { test, expect } from '@playwright/test';

test('GET /api/hats retorna lista de chapéus', async ({ request }) => {
  const response = await request.get('/api/hats');
  expect(response.ok()).toBeTruthy();
  const hats = await response.json();
  expect(Array.isArray(hats)).toBeTruthy();
});
```

## Exercício

No projeto [`playwright-project/`](../playwright-project/README.md):

1. Rode a suíte de exemplo (`npx playwright test`) e confira o relatório.
2. Implemente um teste novo de "dever de casa": filtrar produtos por
   categoria `crescer` e validar que todos os cards exibidos pertencem a essa
   categoria (dica: você pode validar visualmente ou combinar com o retorno de `/api/hats`).
3. Implemente um teste de filtro por faixa de preço.

Próximo módulo: [03 — Boas práticas de IA aplicada ao Playwright](../03-boas-praticas-ia-playwright/README.md)
