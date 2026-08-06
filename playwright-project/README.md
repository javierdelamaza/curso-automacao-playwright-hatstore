# Playwright Project — Hat Store

Suíte de testes E2E em Playwright (JavaScript) para o site Hat Store,
usada como exercício prático dos módulos do curso.

## Pré-requisitos

- Node.js LTS instalado.
- Hat Store rodando (veja o README principal do repositório `hatstore`):
  ```bash
  cd d:\hatstore\hatstore
  podman compose up --build
  ```
  Aplicação disponível em `http://localhost:8080`.

## Instalação

```bash
npm install
npx playwright install
```

## Execução

```bash
npm test              # roda toda a suíte, headless
npm run test:headed   # com navegador visível
npm run test:ui       # modo interativo (UI mode) — recomendado para aprender
npm run test:home     # apenas os testes de busca/filtro da home
npm run test:login    # apenas os testes de login/cadastro
npm run test:api      # apenas os testes de API
npm run report        # abre o relatório HTML da última execução
```

Para rodar contra outro ambiente (ex.: produção), defina `BASE_URL`:

```bash
BASE_URL=https://hatstore-prd.fly.dev npm test
```

## Estrutura

```
playwright-project/
├── playwright.config.js
├── package.json
├── pages/              # Page Object Model
│   ├── HomePage.js
│   └── LoginPage.js
└── tests/
    ├── home.spec.js              # busca por produto
    ├── search-and-filter.spec.js # exercícios de filtro (dever de casa, módulo 2)
    ├── login.spec.js             # login e cadastro
    └── api/
        └── hats.spec.js          # teste de API/contrato básico
```

## Exercícios pendentes (ver módulo 2 e 3 do curso)

Os testes marcados com `test.fixme(...)` em `search-and-filter.spec.js` são o
"dever de casa": implemente-os seguindo o padrão de Page Object já usado em
`home.spec.js` e `login.spec.js`.
