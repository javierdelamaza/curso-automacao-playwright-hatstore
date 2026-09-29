# Módulo 02 — Fundamentos do Playwright

## 2.1 O que é o Playwright e como criar um projeto do zero

O Playwright é um framework open source da Microsoft para automatizar navegadores
(Chromium, Firefox e WebKit) com uma única API. Tem versões em JavaScript/TypeScript, Python,
Java e .NET. Neste curso usamos **JavaScript**.

O que ele tem de diferente:

- **Auto-waiting:** antes de clicar ou preencher, espera o elemento existir, ficar visível, habilitado e estável.
- **Web-first assertions:** `expect(locator).toBeVisible()` tenta de novo até o timeout, sem `sleep`.
- **Isolamento:** cada teste roda num contexto de navegador novo (cookies e localStorage limpos).
- **Ferramentas de debug:** UI mode, trace viewer, vídeo e screenshot na falha, `codegen`.

**Criando um projeto do zero** (para aprender; o projeto do curso já vem pronto):

```powershell
mkdir meu-projeto-playwright
cd meu-projeto-playwright
npm init playwright@latest
```

Respostas no assistente:

| Pergunta | Resposta |
|---|---|
| TypeScript or JavaScript? | **JavaScript** |
| Where to put your end-to-end tests? | **tests** |
| Add a GitHub Actions workflow? | **N** |
| Install Playwright browsers? | **Y** |

## 2.2 Tour pelo projeto do curso

```
playwright-project/
├── playwright.config.js   configuração (baseURL, navegadores, trace, relatório)
├── package.json           dependências e scripts (npm test, npm run test:ui...)
├── pages/                 Page Objects (módulo 03)
├── fixtures/              fixtures customizadas (módulos 03 e 04)
├── data/                  massa de dados fixa (módulo 04)
├── factories/             massa de dados dinâmica (módulo 04)
├── api/                   cliente da API e contratos (módulo 05)
├── tests/
│   ├── front/             testes de tela
│   ├── api/               testes de API
│   └── exercicios/        exercícios (test.fixme)
└── solucoes/              soluções dos exercícios
```

Pontos importantes do `playwright.config.js`:

- `baseURL: 'http://localhost:8080'` → no teste basta escrever `page.goto('/')`.
- `trace: 'on-first-retry'`, `screenshot: 'only-on-failure'`, `video: 'retain-on-failure'` → evidências quando falha.
- `projects` → `chromium` (padrão), `firefox`, `visual` (módulo 04) e `solucoes`.

## 2.3 Descobrindo locators **sem o código-fonte**: DevTools

Você não precisa do código da aplicação para saber como localizar um elemento. O navegador
mostra tudo.

1. Abra http://localhost:8080 no Chrome ou no Edge.
2. Clique com o botão direito no campo de busca → **Inspecionar** (ou `F12`).
3. Você vê o HTML real:
   ```html
   <input type="text" id="search-hat" placeholder="Buscar chapéu...">
   ```
4. Aba **Accessibility**, no painel do DevTools: mostra o **papel** (role) e o **nome acessível** do elemento.
   É isso que o `getByRole` usa.

**Ordem de preferência de locators** (recomendação oficial do Playwright):

| Prioridade | Locator | Exemplo na Hat Store |
|---|---|---|
| 1 | `getByRole` | `page.getByRole('button', { name: 'Adicionar ao carrinho' })` |
| 2 | `getByLabel` | `page.getByLabel('Cupom de desconto:')` |
| 3 | `getByPlaceholder` / `getByText` | `page.getByPlaceholder('Buscar chapéu...')` |
| 4 | `getByTestId` | `page.getByTestId('...')`, quando o dev adiciona `data-testid` |
| 5 | CSS / id | `page.locator('#login-email')`, `page.locator('.hat-card')` |

**Evite:** `div > div:nth-child(3) > span` e XPath longos. Quebram com qualquer mudança de layout.

### Quando não existe um locator estável: peça um `data-testid` ao dev

Os **cards de produto** da Hat Store não têm role, label nem id. Só a classe CSS `.hat-card`.
Se o time de front mudar essa classe por causa de estilo, todos os testes quebram.

Você descobre isso **só olhando o DevTools**, sem ler código. A ação de QA é abrir um pedido
para o time de desenvolvimento:

> **Título:** Adicionar `data-testid` nos cards de produto da home
> **Onde:** Home → lista de chapéus (cada card)
> **Pedido:** `data-testid="hat-card"` no card e `data-testid="hat-card-price"` no preço
> **Motivo:** hoje o único seletor disponível é a classe CSS `.hat-card`, que pode mudar por estilo e quebrar a automação.

Isso é **testabilidade**, e é trabalho do QA também.

## 2.4 `codegen` e "Pick locator"

O Playwright grava o que você faz e gera o código com locators **reais**:

```powershell
npm run codegen      # = npx playwright codegen http://localhost:8080
```

- Abre um navegador e a janela do **Inspector**.
- Clique e digite na loja → o código aparece no Inspector.
- Botão **Pick locator**: passe o mouse sobre um elemento e veja o melhor locator.

Na extensão do VS Code (aba Testing): **Record new** e **Pick locator** fazem o mesmo, direto no editor.

> O `codegen` é um **ponto de partida**. Revise: ele não cria asserções boas e às vezes
> escolhe um locator pior do que você escolheria.

## 2.5 Primeiro teste: ações e asserções

Arquivo [tests/front/01-primeiro-teste.spec.js](../playwright-project/tests/front/01-primeiro-teste.spec.js):

```js
import { test, expect } from '@playwright/test';

test('busca por "Panamá" retorna exatamente um chapéu', async ({ page }) => {
  await page.goto('/');

  await page.getByPlaceholder('Buscar chapéu...').fill('Panamá');

  await expect(page.locator('.hat-card')).toHaveCount(1);
  await expect(page.locator('.hat-card h3')).toHaveText('Chapéu Panamá');
});
```

Estrutura de todo teste: **preparar → agir → verificar** (Arrange, Act, Assert).

| Ações mais usadas | Asserções mais usadas |
|---|---|
| `goto`, `click`, `fill`, `check`, `selectOption`, `press`, `pressSequentially` | `toBeVisible`, `toHaveText`, `toContainText`, `toHaveCount`, `toHaveURL`, `toHaveValue`, `toBeEnabled` |

## 2.6 Rodando e depurando

```powershell
npm test                  # todos os testes, sem abrir navegador (headless)
npm run test:headed       # vendo o navegador
npm run test:ui           # UI mode: linha do tempo, DOM de cada passo (recomendado para aprender)
npx playwright test tests/front/busca.spec.js      # um arquivo
npx playwright test -g "Panamá"                    # testes cujo nome contém "Panamá"
npm run report            # relatório HTML da última execução
```

**Quando um teste falha:**

1. Leia a mensagem: *Expected* x *Received* e a linha marcada com `>`.
2. Abra o relatório → screenshot, vídeo e **trace** (linha do tempo com o DOM de cada passo).
3. Rode com `--debug` para executar passo a passo.

## Exercício (dever de casa)

Em [tests/exercicios/exercicios.spec.js](../playwright-project/tests/exercicios/exercicios.spec.js), bloco **"Bloco 2 - filtros da home"**:

1. Remova o `.fixme` do teste de categoria **Crescer** e implemente.
2. Remova o `.fixme` do teste de **faixa de preço 100 a 130** e implemente.
3. Rode com `npm run test:exercicios`.

Descubra os locators pelo DevTools ou pelo `codegen`. A solução está em
`solucoes/exercicios.spec.js`; só olhe depois de tentar.

Próximo módulo: [03 — Padrões e boas práticas](../03-padroes-boas-praticas/README.md)
