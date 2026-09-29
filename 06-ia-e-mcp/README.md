# Módulo 06 — IA aplicada a testes e MCP do Playwright

A IA (Copilot, ChatGPT, Gemini...) acelera muito a escrita de testes. Mas **quem responde pela
qualidade da suíte é você**, que revisa e roda o teste. Aqui a IA é copiloto, não piloto automático.

Tudo neste módulo funciona com **contas gratuitas** e **sem acesso ao código-fonte** da aplicação.

---

## 6.1 Onde a IA ajuda de verdade

- Gerar o esqueleto de um teste ou Page Object a partir de uma descrição **com fatos reais**.
- Sugerir locators melhores quando você mostra o HTML.
- Explicar uma mensagem de erro, um stack trace ou um timeout.
- Refatorar testes repetidos em fixtures e helpers.
- Gerar massa de dados (CPF válido, e-mails únicos) e JSON Schema.
- Transformar um cenário em linguagem natural ou Gherkin em código.

## 6.2 Onde a IA erra (e por quê)

| Erro típico | Exemplo | Por que acontece |
|---|---|---|
| **Locator inventado** | `page.locator('#btn-login')` que não existe | A IA **não vê a página**, a menos que você mostre (HTML, `codegen`) ou dê um navegador a ela (MCP) |
| **Espera fixa** | `await page.waitForTimeout(3000)` | Padrão antigo, muito comum nos dados de treino |
| **Teste genérico** | `#username`, `#password` | Ela "completa" com o que é comum, não com o seu app |
| **Asserção fraca** | `expect(msg).not.toBeEmpty()` | Passa até com mensagem de erro |
| **"Fazer passar"** | Mockar a API, `test.skip`, afrouxar a asserção | Ela otimiza para "ficar verde", não para achar bug |

## 6.3 Prompt com **fatos da tela**, não do código

Um bom prompt dá à IA **fatos verificáveis** tirados da aplicação rodando. Não é preciso ter o código-fonte.

**Prompt fraco:**
> Cria um teste de login pro meu site.

**Prompt forte (modelo):**

```
Contexto: projeto Playwright em JavaScript, testes em tests/front, usando as fixtures de
fixtures/index.js (homePage, cartPage...) e Page Objects em pages/.

Fatos da tela (copiados do DevTools/codegen):
- Campo de cupom: <input id="cupom" placeholder="Digite o cupom"> com <label>Cupom de desconto:</label>
- Botão: <button id="aplicar-cupom">Aplicar</button>
- Após aplicar "HATOFF", aparece "Cupom aplicado!" e a linha
  "Cupom HATOFF aplicado: -R$ 24.00" para um item de R$ 120.00.

Tarefa: gere um teste que adiciona o "Chapéu Panamá", aplica HATOFF e valida a mensagem,
o desconto e o total.

Restrições: use os métodos existentes do CartPage; sem waitForTimeout; asserções com
toHaveText (valor exato); não invente seletores fora dos que listei.
```

Os quatro blocos são sempre os mesmos: **contexto → fatos → tarefa → restrições**.

**De onde tirar os fatos (sem código-fonte):**
- DevTools → botão direito no elemento → *Copy → Copy outerHTML*
- saída do `npm run codegen`
- `error-context.md` que o Playwright gera quando um teste falha (traz o *snapshot de acessibilidade* da página)
- Swagger e aba Network, para API

## 6.4 Mão na massa com o Copilot Free (VS Code)

1. Abra o chat (`Ctrl+Alt+I`) no modo **Ask**.
2. Anexe contexto: digite `#` e escolha arquivos (ex.: `#file:CartPage.js`, `#file:fixtures/index.js`).
   A IA passa a usar o **seu** padrão.
3. Cole o prompt forte do item 6.3 (troque pela tarefa que quiser).
4. Copie o código para um arquivo em `tests/front/`, rode (`npm test` ou ▶ na aba Testing) e **revise com o checklist (6.10)**.

> Economize a cota do plano Free: um prompt bem feito resolve em 1 mensagem o que 5 prompts vagos não resolvem.

## 6.5 O mesmo fluxo no chat web (ChatGPT/Gemini free)

Sem Copilot? O fluxo é o mesmo, só que por copiar e colar:

1. Abra o arquivo do Page Object e das fixtures e **cole o conteúdo** no chat, junto com o prompt.
2. Copie a resposta para o VS Code.
3. Rode, revise e, se falhar, cole o **erro** de volta no chat (item 6.6).

A diferença é só a comodidade. As boas práticas são as mesmas.

## 6.6 Debug com IA

Quando um teste falha, **não descreva o erro de memória**. Mande o erro real:

- a mensagem do terminal (Expected x Received, linha com `>`);
- o arquivo `test-results/<teste>/error-context.md`, que traz o estado da página no momento da falha;
- no relatório HTML e no UI mode existe o botão **Copy prompt**, que monta esse pacote pronto para colar na IA.

Prompt:
> "Este teste Playwright falhou. Erro e contexto abaixo. Explique a causa provável e proponha a
> correção **sem** usar waitForTimeout e **sem** afrouxar a asserção. [cole]"

## 6.7 O que é MCP e como ligar o MCP do Playwright

**MCP (Model Context Protocol)** é um padrão aberto que dá **ferramentas** a um agente de IA.
O **Playwright MCP** dá à IA um **navegador de verdade**: ela abre a página, clica, digita e lê o
*snapshot de acessibilidade* (o mesmo que o `getByRole` usa).

**Resultado:** acabam os locators inventados. A IA vê a página como ela é, **sem código-fonte**.

O repositório já traz a configuração em [.vscode/mcp.json](../.vscode/mcp.json):

```json
{
  "servers": {
    "playwright": { "type": "stdio", "command": "npx", "args": ["@playwright/mcp@latest"] }
  }
}
```

Para ativar:
1. Abra `.vscode/mcp.json` no VS Code e clique em **Start**, acima de `"playwright"`.
   Ou use `Ctrl+Shift+P` → **MCP: List Servers** → playwright → Start.
2. No chat, mude para o modo **Agent**.
3. Clique no ícone de **ferramentas** do chat e confira se as ferramentas `browser_*` do Playwright estão marcadas.

## 6.8 MCP na prática: a IA explora a Hat Store e gera o teste

Com a Hat Store no ar e o chat em modo **Agent**:

```
Use o Playwright MCP para abrir http://localhost:8080.
1. Busque por "Fedora", adicione ao carrinho e aumente a quantidade para 2.
2. Anote os locators que você usou (prefira role/label/placeholder).
3. Depois gere um teste em tests/front/ia-carrinho.spec.js que reproduz o fluxo e valida
   o total de R$ 300.00, usando as fixtures de #file:fixtures/index.js e os métodos de
   #file:pages/CartPage.js e #file:pages/HomePage.js. Não use waitForTimeout.
```

Acompanhe: o navegador abre e a IA "narra" o que viu. **Aprove cada ação** quando o VS Code pedir.
Depois, rode o teste gerado e revise com o checklist.

> Cada passo do agente consome requisições da sua cota Free. Peça tarefas pequenas e bem definidas.

## 6.9 Playwright Test Agents (planner, generator, healer)

O Playwright traz três agentes prontos, que usam o MCP por baixo:

| Agente | Faz o quê |
|---|---|
| **planner** | Explora a aplicação e escreve um **plano de testes** em Markdown |
| **generator** | Transforma o plano em arquivos `.spec.js` e verifica os locators ao vivo |
| **healer** | Roda um teste que falhou, investiga a página e propõe a correção |

Para instalar (dentro de `playwright-project`):

```powershell
npx playwright init-agents --loop=vscode
```

O comando cria as definições dos agentes e um *seed test*. No chat do VS Code, selecione o
agente (ex.: **planner**) e peça: *"Crie um plano de testes para o carrinho da Hat Store"*.

> Os agentes geram **rascunhos**. O plano precisa da sua revisão de QA (faltou cenário negativo?)
> e o código passa pelo mesmo checklist.

## 6.10 Checklist de revisão de código gerado por IA

Antes de aceitar um teste gerado ou editado pela IA:

- [ ] Os locators **existem** na página (confira no DevTools, no `codegen` ou no UI mode)?
- [ ] Não há `waitForTimeout` nem `sleep`?
- [ ] O teste roda **sozinho** e em qualquer ordem (massa própria)?
- [ ] As asserções verificam **o valor esperado** (`toHaveText('Cupom aplicado!')`), e não só "algo apareceu"?
- [ ] O teste **falha quando deveria**? Mude o valor esperado de propósito e veja ficar vermelho.
- [ ] Nenhum `test.skip`, mock ou asserção afrouxada foi colocado "para passar"?
- [ ] Usa os padrões do projeto (fixtures, Page Objects, `data/` e `factories/`)?
- [ ] Nome do teste em linguagem de negócio?

## Segurança e IA

- **Nunca** cole segredos (senhas reais, tokens, `.env`) num prompt.
- **Nunca** cole código, dados de clientes ou documentos de projetos reais numa IA de **conta pessoal ou gratuita**. Siga a política da empresa.
- Dados de teste são sempre **fictícios** (use as factories).
- Se a IA sugerir desligar validações ("ignora o HTTPS", "pula a verificação do token"), trate como **alerta**. O teste existe para validar o comportamento real.

## 6.11 Se você tiver acesso ao código-fonte

> [!TIP]
> **Bônus, não pré-requisito.** Tudo acima funciona sem o código. Se no seu projeto você tiver acesso:
> - dá para **confirmar regras** (ex.: o percentual real do cupom) direto na fonte, em vez de depender só da documentação;
> - dá para anexar o componente da tela como **contexto extra** para a IA;
> - mas **não** cole código proprietário numa IA de conta pessoal ou gratuita. Com a IA corporativa aprovada, siga as regras da empresa.

## Exercício com IA

Bloco **"Bloco 6"** em [exercicios.spec.js](../playwright-project/tests/exercicios/exercicios.spec.js).
Faça os dois testes do cupom **com a IA**:

1. Levante os fatos da tela (DevTools ou MCP).
2. Escreva o prompt no formato **contexto → fatos → tarefa → restrições**.
3. Gere, rode e passe pelo checklist.
4. **Anote** quantas mensagens você gastou e o que precisou corrigir. Vamos discutir no módulo 07.

Próximo módulo: [07 — Fechamento](../07-fechamento/README.md)
