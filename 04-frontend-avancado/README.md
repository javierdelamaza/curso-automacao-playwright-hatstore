# Módulo 04 — Automação frontend avançada

## 4.1 Organizando testes por contexto de negócio

```
tests/front/
├── 01-primeiro-teste.spec.js   (histórico do módulo 02)
├── busca.spec.js               catálogo e busca
├── combinados.spec.js          filtros + busca juntos
├── login.spec.js
├── cadastro.spec.js
├── carrinho.spec.js
├── compra-e2e.spec.js          fluxo completo (@e2e)
└── visual.spec.js              regressão visual (@visual)
```

Um arquivo por funcionalidade, nomes de teste em linguagem de negócio e `test.describe` para agrupar.
**Tags** no nome (`@e2e`, `@visual`) permitem rodar subconjuntos:

```powershell
npm run test:e2e        # só @e2e
npm run test:front      # só a pasta front
```

## 4.2 Massa de dados: fixa x dinâmica

**Fixa:** [data/catalogo.js](../playwright-project/data/catalogo.js) reúne o que vem do banco
inicial e não muda: produtos, cupom, endereço padrão, frete do Sudeste.
Foi levantada **pela tela e pelo Swagger**, sem abrir o código.

```js
export const PRODUTOS = {
  panama: { id: 1, nome: 'Chapéu Panamá', price: 120.0, categoria: 'importado' },
  // ...
};
```

**Dinâmica:** [factories/usuarioFactory.js](../playwright-project/factories/usuarioFactory.js)
gera um usuário novo por teste: e-mail único, **CPF válido** (com dígitos verificadores) e
data de nascimento de 25 anos.

```js
const usuario = novoUsuario();                         // padrão válido
const menor = usuarioMenorDeIdade();                   // variação pronta
const senhaFraca = novoUsuario({ password: '123' });   // sobrescreve só o que interessa
```

Por que dinâmica? A Hat Store **não aceita** e-mail nem CPF repetidos, e o cupom só vale
uma vez por CPF. Com dados fixos, o 2º run já falharia.

## 4.3 Fixtures de massa dinâmica

Em [fixtures/index.js](../playwright-project/fixtures/index.js):

| Fixture | O que entrega | Como |
|---|---|---|
| `usuarioCadastrado` | Um usuário novo, já cadastrado | `POST /auth/register` (API é mais rápida que a tela) |
| `token` | JWT desse usuário | `POST /auth/login` |
| `paginaLogada` | `page` já logada | Grava o token no `localStorage` antes de abrir o site |

Duas formas de começar logado:

```js
// 1) Pela tela (testa o login, mas é mais lento)
test.beforeEach(async ({ page, loginPage, usuarioCadastrado }) => {
  await loginPage.goto();
  await loginPage.login(usuarioCadastrado.email, usuarioCadastrado.password);
  await page.waitForURL('**/index.html');
});

// 2) Injetando o token (rápido, para testes que NÃO são sobre login)
test('...', async ({ paginaLogada, homePage }) => { /* já está logado */ });
```

Como descobrimos a chave `jwt_token`? No DevTools: **Application → Local Storage**, depois de
fazer login pela tela. Sem código-fonte.

## 4.4 Cenários combinados

[combinados.spec.js](../playwright-project/tests/front/combinados.spec.js) testa a **interação
entre componentes**: categoria + faixa de preço, duas categorias juntas, filtro + busca.
É nas combinações que aparecem os bugs que os testes isolados não pegam.

## 4.5 Regressão visual

```js
test('home com catálogo completo', async ({ page, homePage }) => {
  await homePage.goto();
  await page.waitForLoadState('networkidle');
  await expect(page).toHaveScreenshot('home.png', { fullPage: true });
});
```

1. **Primeira execução:** gere a imagem de referência (*baseline*):
   ```powershell
   npm run test:visual:update
   ```
2. **Próximas execuções:** compara pixel a pixel com a baseline:
   ```powershell
   npm run test:visual
   ```
3. Se falhar, o relatório mostra **Expected / Actual / Diff** lado a lado.

Cuidados: a baseline depende do sistema operacional e das fontes (por isso cada aluno gera a
sua e ela não vai para o Git neste curso). Dados que mudam, como datas e estoque, precisam de
`mask`. A tolerância está configurada em `playwright.config.js` (`maxDiffPixelRatio`).

## 4.6 e 4.7 Fluxo E2E completo

[compra-e2e.spec.js](../playwright-project/tests/front/compra-e2e.spec.js) faz o caminho do
cliente de ponta a ponta, dividido em `test.step` (cada passo aparece no relatório):

1. login no `beforeEach` (pela tela);
2. buscar "Gustavo" e adicionar ao carrinho;
3. checkout mostra nome e CPF do usuário logado;
4. preencher endereço → frete calculado ao sair do campo CEP (Sudeste, R$ 12.50);
5. pagar com PIX → modal "Pagamento PIX" com o código;
6. **Meus Pedidos** mostra o pedido com o produto e "PIX".

> **E2E é caro:** tenha poucos, cobrindo os fluxos que dão dinheiro ou dão problema.
> As variações de regra ficam nos testes de API (módulo 05).

## 4.8 Scripts do `package.json`

| Script | O que faz |
|---|---|
| `npm test` | Toda a suíte no Chromium (sem visual) |
| `npm run test:front` / `test:api` | Só tela / só API |
| `npm run test:e2e` | Só os fluxos `@e2e` |
| `npm run test:visual` / `test:visual:update` | Comparar / gerar baseline visual |
| `npm run test:all-browsers` | Chromium + Firefox (instale com `npx playwright install firefox`) |
| `npm run test:exercicios` / `test:solucoes` | Exercícios / soluções |

## Exercício

Bloco **"Bloco 4"** em [exercicios.spec.js](../playwright-project/tests/exercicios/exercicios.spec.js):
compra de um cliente do **Nordeste** (CEP `40010-000`) pagando com **Boleto**, começando já
logado com a fixture `paginaLogada`. Valide região, valor do frete, o modal de pagamento e
o pedido em Meus Pedidos.

## 4.9 Dica de IA: gerar massa de dados

Peça à IA:

> "Crie uma função JavaScript que gere um CPF **válido** (com dígitos verificadores) no formato
> 000.000.000-00, e explique o algoritmo em 5 linhas."

Depois **teste a resposta**: gere 5 CPFs e valide num validador online. A IA erra algoritmos com
frequência. Compare com [factories/usuarioFactory.js](../playwright-project/factories/usuarioFactory.js).

Próximo módulo: [05 — Testes de API](../05-testes-api/README.md)
