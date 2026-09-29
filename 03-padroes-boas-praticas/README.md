# Módulo 03 — Padrões de projeto e boas práticas

## 3.1 Por que Page Object?

Sem padrão, cada teste repete os seletores. Se o placeholder da busca mudar, você altera 20
arquivos. O **Page Object Model (POM)** encapsula uma tela numa classe:

- **locators** ficam num só lugar;
- **ações** ganham nomes de negócio (`search`, `login`, `adicionarAoCarrinho`);
- o teste fica legível: descreve *o que* o usuário faz, não *como* a tela é montada.

## 3.2 Refatorando a busca para `HomePage`

[pages/HomePage.js](../playwright-project/pages/HomePage.js) (trecho):

```js
export class HomePage {
  constructor(page) {
    this.page = page;
    this.searchInput = page.getByPlaceholder('Buscar chapéu...');
    this.hatCards = page.locator('.hat-card');
    this.hatCardTitles = page.locator('.hat-card h3');
  }

  async goto() {
    await this.page.goto('/');
    await this.hatCards.first().waitFor();
  }

  async search(term) {
    await this.searchInput.fill(term);
  }
}
```

Compare o [primeiro teste](../playwright-project/tests/front/01-primeiro-teste.spec.js) com a
versão refatorada em [busca.spec.js](../playwright-project/tests/front/busca.spec.js).

**Regras práticas:**
- O Page Object **não faz asserções** (`expect`). Quem verifica é o teste.
- Métodos devolvem o resultado ou nada. Não escondem regras de negócio.
- Um Page Object por tela ou componente relevante: `HomePage`, `LoginPage`, `CartPage`,
  `CheckoutPage` e `PedidosPage`.

## 3.3 `LoginPage`: cadastro e login

[pages/LoginPage.js](../playwright-project/pages/LoginPage.js) cobre os dois formulários da
tela `/auth.html`. Um detalhe que só aparece testando: o campo CPF tem **máscara**.
Por isso usamos `pressSequentially` (digita tecla a tecla, como o usuário) em vez de `fill`.

```js
await this.registerCpfInput.pressSequentially(cpf.replace(/\D/g, ''));
```

## 3.4 Cenários positivos e negativos

Um bom conjunto de testes cobre o **caminho feliz** e as **regras de rejeição**.
Na Hat Store, pela tela e pela documentação em `/documentacao.html`, descobrimos:

| Regra | Cenário negativo |
|---|---|
| Idade mínima 18 anos | Cadastro de menor é recusado |
| E-mail único | Cadastro com e-mail já usado é recusado |
| Senha forte (8+, maiúscula, minúscula, número, especial) | Senha fraca é recusada (exercício) |
| CPF no formato 000.000.000-00 | CPF incompleto mostra erro (exercício) |
| Login | Senha errada ou usuário inexistente → "Erro no login: Não autorizado" |

Veja [login.spec.js](../playwright-project/tests/front/login.spec.js) e
[cadastro.spec.js](../playwright-project/tests/front/cadastro.spec.js).

> Repare que a mensagem de login inválido é **genérica**: não diz se o erro foi no e-mail ou
> na senha. Isso é uma boa prática de segurança, e o teste confirma que continua assim.

## 3.5 Fixtures: injetando Page Objects nos testes

Em vez de `const home = new HomePage(page)` em todo teste, criamos **fixtures**
([fixtures/index.js](../playwright-project/fixtures/index.js)):

```js
import { test as base } from '@playwright/test';
import { HomePage } from '../pages/HomePage.js';

export const test = base.extend({
  homePage: async ({ page }, use) => use(new HomePage(page)),
});
```

E no teste:

```js
import { test, expect } from '../../fixtures/index.js';

test('busca por produto existente retorna só ele', async ({ homePage }) => {
  await homePage.goto();
  await homePage.search('Panamá');
  await expect(homePage.hatCards).toHaveCount(1);
});
```

O Playwright só cria a fixture se o teste **pedir** por ela no parâmetro.

## 3.6 Exercício

Em [tests/exercicios/exercicios.spec.js](../playwright-project/tests/exercicios/exercicios.spec.js), bloco **"Bloco 3"**:

1. **Senha fraca** (`Senha1234`, sem caractere especial) deve ser recusada.
2. **CPF incompleto** deve mostrar erro no campo. Descubra no DevTools onde a mensagem aparece.

## 3.7 Dica de IA: refatorar para Page Object

Cole no Copilot Chat (modo **Ask**) um teste "cru" e o HTML do elemento copiado do DevTools
(botão direito → *Copy → Copy outerHTML*):

> "Tenho este teste Playwright em JavaScript [cole o teste]. O HTML real da tela é [cole o HTML].
> Crie um Page Object `CartPage` com os locators e métodos usados, priorizando
> `getByRole`/`getByLabel`. Não invente seletores que não estejam no HTML que colei."

Depois, **revise**: os locators existem? Os nomes dos métodos fazem sentido para o negócio?

Próximo módulo: [04 — Frontend avançado](../04-frontend-avancado/README.md)
