# Playwright Project — Hat Store

Suíte de testes do curso: tela, E2E, visual, API e contrato.

## Pré-requisitos

- Node.js LTS
- Hat Store rodando em http://localhost:8080 (`cd ..\hatstore-app; podman compose up --build -d`)

## Instalação

```powershell
npm install
npx playwright install chromium
```

## Execução

| Comando | O que faz |
|---|---|
| `npm test` | Suíte completa no Chromium (sem os testes visuais) |
| `npm run test:ui` | UI mode (recomendado para aprender e depurar) |
| `npm run test:headed` | Vendo o navegador |
| `npm run test:front` / `npm run test:api` | Só tela / só API |
| `npm run test:e2e` | Só os fluxos completos (`@e2e`) |
| `npm run test:visual:update` / `npm run test:visual` | Gerar / comparar baseline visual |
| `npm run test:exercicios` / `npm run test:solucoes` | Exercícios / soluções |
| `npm run codegen` | Gravar interações e gerar código |
| `npm run report` | Abrir o relatório HTML |

## Estrutura

```
playwright-project/
├── playwright.config.js
├── api/            HatStoreApi (camada de API), schemas e validador de contrato
├── data/           massa FIXA (catálogo, cupom, endereço)
├── factories/      massa DINÂMICA (usuário com e-mail e CPF únicos)
├── fixtures/       Page Objects, usuário cadastrado, token, página logada, API autenticada
├── pages/          HomePage, LoginPage, CartPage, CheckoutPage, PedidosPage
├── tests/
│   ├── front/      busca, filtros combinados, login, cadastro, carrinho, E2E, visual
│   ├── api/        hats, auth, frete, pedidos, contrato
│   └── exercicios/ exercícios dos módulos 02 a 06 (test.fixme)
└── solucoes/       soluções dos exercícios
```
