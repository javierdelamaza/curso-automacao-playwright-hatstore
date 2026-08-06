# Curso de Automação de Testes com Playwright — Hat Store

Curso introdutório de automação de testes end-to-end, usando o site **Hat Store**
(`d:\hatstore\hatstore`) como aplicação alvo. Conteúdo teórico expandido e um módulo extra sobre uso de IA no dia a dia
da automação.

## Como usar este material

1. Suba o Hat Store localmente antes de rodar qualquer teste:
   ```bash
   cd d:\hatstore\hatstore
   podman compose up --build
   ```
   A aplicação fica disponível em `http://localhost:8080`.
2. Siga os módulos na ordem abaixo. Cada pasta tem um `README.md` com teoria,
   exemplos e exercícios.
3. O projeto de automação funcional fica em [playwright-project/](playwright-project/README.md).
   É nele que os exercícios de cada módulo devem ser implementados.

## Módulos

| Módulo | Tema | Conteúdo |
|--------|------|----------|
| [01](01-teoria-automacao-testes/README.md) | Teoria de Automação de Testes | Pirâmide de testes, tipos de teste, quando (não) automatizar, flakiness, métricas |
| [02](02-fundamentos-playwright/README.md) | Fundamentos do Playwright | Instalação, locators, auto-waiting, assertions, Page Object Model, fixtures |
| [03](03-boas-praticas-ia-playwright/README.md) | Boas Práticas de IA aplicada ao Playwright | Prompting, `codegen`, MCP do Playwright, revisão de código gerado, riscos e limites |

## Projeto de exemplo

- [playwright-project/](playwright-project/README.md) — projeto Node.js com Playwright já configurado
  para testar o Hat Store (busca, filtros, login e cadastro), usando Page Object Model.

## Referências oficiais

- Documentação do Playwright: https://playwright.dev/docs/intro
- Repositório de aulas original: https://github.com/galaschi/start-aula-playwright
- Documentação da API do Hat Store: `http://localhost:8080/documentacao.html` e `http://localhost:8080/swagger/index.html`
