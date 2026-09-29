# Curso de Automação de Testes com Playwright + IA — Hat Store

Curso introdutório de automação de testes **E2E (tela) e de API** com Playwright, com um
módulo de **IA aplicada a testes**, incluindo MCP. Foi pensado para quem tem pouca ou
nenhuma experiência com automação ou com IA.

Baseado no repositório de aulas [start-aula-playwright](https://github.com/galaschi/start-aula-playwright).

## Três premissas do curso

1. **Sem acesso ao código-fonte.** Você testa a Hat Store do mesmo jeito que um QA que só
   recebeu a URL do ambiente: pela tela, pelo DevTools, pelo `codegen`, pelo Swagger e com a IA
   olhando a página no navegador (MCP). O código da loja fica em `hatstore-app/` só para você
   conseguir subir a aplicação.
2. **Só ferramentas gratuitas.** Node, VS Code, Podman, Playwright e **GitHub Copilot Free**.
   Como plano B, ChatGPT ou Gemini na versão gratuita pelo navegador.
3. **Aulas curtas.** Cada vídeo tem até 10 minutos e trata de um assunto só.

## Estrutura do repositório

```
├── 00-setup/                    instalação passo a passo (comece aqui)
├── 01-teoria-automacao-testes/
├── 02-fundamentos-playwright/
├── 03-padroes-boas-praticas/
├── 04-frontend-avancado/
├── 05-testes-api/
├── 06-ia-e-mcp/
├── 07-fechamento/
├── hatstore-app/                aplicação sob teste (sobe com Podman)
├── playwright-project/          projeto de automação (onde você programa)
└── docs/                        Guia do Aluno em PDF
```

## Módulos

| # | Módulo | O que você aprende |
|---|---|---|
| 00 | [Setup](00-setup/README.md) | Instalar Git, Node, VS Code, Podman; subir a Hat Store; rodar o 1º teste; ativar Copilot Free |
| 01 | [Teoria](01-teoria-automacao-testes/README.md) | Pirâmide, tipos de teste, quando automatizar, flakiness, massa de dados |
| 02 | [Fundamentos do Playwright](02-fundamentos-playwright/README.md) | Projeto do zero, locators sem código-fonte, `codegen`, assertions, UI mode, trace |
| 03 | [Padrões e boas práticas](03-padroes-boas-praticas/README.md) | Page Object, login e cadastro, cenários negativos, fixtures |
| 04 | [Frontend avançado](04-frontend-avancado/README.md) | Massa de dados, factory, cenários combinados, regressão visual, E2E completo |
| 05 | [Testes de API](05-testes-api/README.md) | Swagger, `request`, autenticação JWT, testes de contrato |
| 06 | [IA e MCP](06-ia-e-mcp/README.md) | Prompts com fatos da tela, Copilot Free, chat web, MCP do Playwright, Test Agents |
| 07 | [Fechamento](07-fechamento/README.md) | Resolução do exercício de IA e próximos passos |

## Início rápido (depois do setup)

```powershell
cd hatstore-app
podman compose up --build -d
cd ..\playwright-project
npm install
npx playwright install chromium
npm test
```
