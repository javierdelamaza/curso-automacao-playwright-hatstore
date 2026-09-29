# Módulo 01 — Teoria de Automação de Testes

Antes de escrever código, precisamos saber **o que** vale a pena automatizar e **por quê**.

## 1.1 O que é automação de testes

É usar scripts para executar verificações que uma pessoa faria manualmente: clicar,
preencher formulários, conferir textos na tela, validar respostas de API.

O objetivo **não** é substituir o teste humano. A automação serve para:

- repetir a regressão de forma rápida e sempre igual;
- rodar em vários navegadores e em paralelo;
- rodar no CI/CD e barrar uma entrega quebrada;
- liberar o QA para o teste exploratório, que é onde se acham os bugs mais interessantes.

## 1.2 Pirâmide de testes (na Hat Store)

```
        /\        E2E / UI: poucos, lentos, cobrem fluxos críticos
       /  \
      /----\      API / Integração: rápidos, cobrem regras e contratos
     /      \
    /--------\    Unitários: muitos, rapidíssimos (responsabilidade do time de dev)
```

| Camada | Quem costuma fazer | Na Hat Store |
|---|---|---|
| Unitário | Dev (precisa do código-fonte) | Já existe no projeto (Go e Jest). **Não é foco do curso.** |
| API | QA e Dev | `GET /api/hats`, `POST /auth/login`, `POST /api/pedido`... pelo **Swagger** |
| E2E / UI | QA | Buscar, filtrar, cadastrar, logar, comprar com PIX/Boleto |

> O QA consegue atuar nas camadas de **API e E2E sem ler o código-fonte**. Basta o
> ambiente no ar, o Swagger e a tela. É o cenário deste curso.

## 1.3 Tipos de teste que a automação cobre

- **Funcional:** o sistema faz o que deveria. Ex.: o cupom `HATOFF` dá 20% de desconto no carrinho.
- **Regressão:** o que funcionava continua funcionando depois de uma mudança.
- **Contrato (API):** a resposta mantém o formato combinado (campos, tipos).
- **Visual:** a tela não "quebrou" o layout (comparação de screenshots).
- **Não funcionais** (fora do escopo, mas bom conhecer): performance, carga, segurança, acessibilidade.

## 1.4 Quando automatizar (e quando não)

Automatize quando o cenário é:

- **repetitivo:** roda em toda regressão;
- **estável:** a tela ou o fluxo não muda toda semana;
- **de alto risco:** login, checkout, pagamento.

Deixe para depois, ou mantenha manual:

- telas em redesenho constante;
- casos raríssimos ou que exigem julgamento humano ("o botão parece desalinhado");
- testes exploratórios: eles *descobrem* o que vale automatizar.

## 1.5 Flakiness (testes instáveis)

Um teste é *flaky* quando passa e falha sem que nada mude. Causas comuns:

| Causa | Exemplo | Como evitar |
|---|---|---|
| Espera fixa | `waitForTimeout(3000)` | Esperar uma **condição** (`expect(...).toBeVisible()`) |
| Dependência entre testes | Teste B usa usuário criado no teste A | Cada teste cria os próprios dados |
| Seletor frágil | `div > div:nth-child(3)` | Locators por papel, texto, label ou `data-testid` |
| Ambiente | Container ainda subindo | Subir o ambiente antes e conferir que está saudável |

O Playwright resolve boa parte disso com **auto-waiting** (módulo 02). O resto é disciplina
de quem escreve o teste.

## 1.6 Massa de dados e isolamento

- **Cada teste roda sozinho**, em qualquer ordem.
- **Massa fixa:** vem do banco inicial e não muda (catálogo de 15 chapéus, cupom `HATOFF`).
- **Massa dinâmica:** é criada a cada teste (usuário novo com e-mail e CPF únicos).
- Na Hat Store isso importa. Ela **não deixa** repetir e-mail nem CPF, bloqueia o login
  depois de 5 tentativas erradas e o cupom só vale uma vez por CPF.

## 1.7 Métricas de uma suíte saudável

- **Cobertura de fluxos críticos:** login, busca e compra têm teste?
- **Tempo de execução:** dá para rodar a cada Pull Request?
- **Taxa de flakiness:** quantas vezes é preciso "rodar de novo para passar"?
- **Manutenção:** se um botão muda de lugar, quantos arquivos você precisa alterar? (O Page Object, no módulo 03, existe para diminuir esse número.)

## Exercício (sem código)

Com a Hat Store aberta no navegador, escreva:

1. Os **3 fluxos** que você automatizaria primeiro e por quê.
2. **1 fluxo** que você *não* automatizaria agora, com a justificativa.
3. **2 dados** que precisam ser únicos a cada execução de teste.
4. **Bônus:** algo que parece bug. (Dica: compare o que a documentação em `/documentacao.html` diz com o que a tela faz.)

Próximo módulo: [02 — Fundamentos do Playwright](../02-fundamentos-playwright/README.md)
