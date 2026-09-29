# Módulo 07 — Fechamento e próximos passos

## 7.1 Resolução do exercício de IA

Compare o seu teste do cupom com [solucoes/exercicios.spec.js](../playwright-project/solucoes/exercicios.spec.js)
(bloco "Bloco 6") e responda:

- Quantas mensagens você gastou com a IA? O que mudaria no primeiro prompt para acertar de primeira?
- A IA inventou algum locator? Colocou algum `waitForTimeout`?
- A asserção dela verificava o **valor** (R$ 24.00 de desconto, total R$ 96.00) ou só "algo apareceu"?
- Você fez o teste falhar de propósito para garantir que ele pega o erro?

## 7.2 O que você construiu

- Uma suíte com **testes de tela, E2E, visuais, de API e de contrato**, organizada com Page Objects, fixtures, massa fixa e factory.
- **Um bug real encontrado e documentado** (`test.fail` do cupom HATOFF) e uma divergência de documentação (25% x 20%).
- Um método para usar IA **gratuita** sem depender do código-fonte: fatos da tela → prompt → revisão com checklist → MCP.

## Próximos passos (para continuar estudando)

| Tema | Por onde começar |
|---|---|
| **CI/CD** | Rodar a suíte no GitHub Actions a cada push: https://playwright.dev/docs/ci-intro |
| **TypeScript** | Migrar o projeto para TS (tipos ajudam a IA a errar menos): https://playwright.dev/docs/test-typescript |
| **Autenticação reaproveitada** | `storageState` para logar uma vez por execução: https://playwright.dev/docs/auth |
| **Mock de rede** | `page.route` para simular erros da API: https://playwright.dev/docs/mock |
| **Acessibilidade** | `@axe-core/playwright`: https://playwright.dev/docs/accessibility-testing |
| **Visual no CI** | Gerar baselines dentro de um container para não depender do SO |
| **Test Agents** | Aprofundar planner, generator e healer: https://playwright.dev/docs/test-agents |

Documentação oficial: https://playwright.dev/docs/intro

Voltar ao [início do curso](../README.md)
