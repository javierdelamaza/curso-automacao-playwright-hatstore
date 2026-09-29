# Módulo 00 — Setup do ambiente (Windows)

Ao final deste módulo você terá:

- a **Hat Store rodando** em http://localhost:8080;
- o **projeto de testes** rodando e o relatório aberto;
- o **VS Code com Copilot Free** configurado (e um plano B de IA no navegador).

> Tempo estimado: 40–60 min (a maior parte é download).
> Você vai precisar de **permissão de administrador** para instalar o WSL/Podman.
> Se não tiver, peça ao suporte antes de começar.

---

## 0.1 Contas necessárias (gratuitas)

| Conta | Para quê | Link |
|---|---|---|
| GitHub | Clonar o curso e usar o Copilot Free | https://github.com/signup |

> [!NOTE]
> O repositório do curso é **privado**. Envie seu usuário do GitHub para os instrutores
> e aceite o convite que chega por e-mail antes de seguir para o passo 0.5.

---

## 0.2 Git e VS Code

**Opção A — pelo terminal (PowerShell):**

```powershell
winget install --id Git.Git -e
winget install --id Microsoft.VisualStudioCode -e
```

**Opção B — pelo site:** baixe e instale com as opções padrão:
- Git: https://git-scm.com/download/win
- VS Code: https://code.visualstudio.com/

Feche e abra o PowerShell de novo e confira:

```powershell
git --version
code --version
```

**Extensões do VS Code** (aba *Extensions*, `Ctrl+Shift+X`):

| Extensão | ID |
|---|---|
| Playwright Test for VSCode | `ms-playwright.playwright` |
| GitHub Copilot Chat | `github.copilot-chat` |

> Ao abrir a pasta do curso, o VS Code sugere essas extensões automaticamente
> (arquivo `.vscode/extensions.json`).

---

## 0.3 Node.js (LTS)

```powershell
winget install --id OpenJS.NodeJS.LTS -e
```

Ou baixe o instalador **LTS** em https://nodejs.org/. Feche e abra o PowerShell e confira:

```powershell
node --version   # v20 ou superior
npm --version
```

> Se o PowerShell reclamar de *execution policy* ao rodar `npm`/`npx`, execute uma vez:
> `Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned`

---

## 0.4 Podman Desktop (para subir a Hat Store)

O Podman roda containers, como o Docker, e é gratuito. A Hat Store usa dois containers:
banco PostgreSQL + aplicação.

1. Instale:
   ```powershell
   winget install --id RedHat.Podman-Desktop -e
   ```
   ou baixe em https://podman-desktop.io/downloads.
2. Abra o **Podman Desktop**. Na primeira execução ele mostra um assistente:
   - Se pedir para instalar o **WSL2**, aceite. Pode ser necessário **reiniciar o Windows**.
   - Clique em **Create new Podman machine** e aceite os valores padrão.
   - Em *Settings → Resources*, confirme que a máquina está **Running**.
3. Instale o **Compose**: no Podman Desktop, vá em *Settings → Resources → Compose → Setup*
   (ou no banner "Compose not installed"). Siga o assistente.
4. Confira no PowerShell:
   ```powershell
   podman --version
   podman machine list        # a máquina deve aparecer como "Currently running"
   podman compose version
   ```

> Se você reiniciar o computador, abra o Podman Desktop (ou rode `podman machine start`)
> antes de subir a Hat Store.

---

## 0.5 Clonar o curso e subir a Hat Store

```powershell
cd $HOME\Documents
git clone https://github.com/javierdelamaza/curso-automacao-playwright-hatstore.git
cd curso-automacao-playwright-hatstore
code .
```

> Na primeira vez, o Git abre uma janela para você **entrar na sua conta do GitHub**.

No terminal do VS Code (`` Ctrl+` ``):

```powershell
cd hatstore-app
podman compose up --build -d
```

Espere terminar (na primeira vez leva alguns minutos). Depois abra **http://localhost:8080**.

**Explore a loja como usuário** (é o seu primeiro "teste exploratório"):
- busque "Panamá", filtre por categoria e por preço;
- crie uma conta, faça login, adicione um chapéu ao carrinho, aplique o cupom `HATOFF`;
- finalize uma compra com PIX e veja em **Meus Pedidos**;
- abra o Swagger: http://localhost:8080/swagger/index.html.

> [!TIP]
> Anote o que você achou estranho. Testar não é só automatizar, e o curso vai mostrar
> que a Hat Store tem bugs de verdade.

---

## 0.6 Rodar o primeiro teste

Ainda no terminal do VS Code:

```powershell
cd ..\playwright-project
npm install
npx playwright install chromium
npm test
```

O resultado esperado é algo como `43 passed` (alguns aparecem como *skipped*: são os exercícios).
Abra o relatório:

```powershell
npm run report
```

Também dá para rodar pela extensão: ícone de **tubo de ensaio** (Testing) na barra lateral
→ ▶ ao lado de um teste. Marque **Show browser** para ver o navegador.

---

## 0.7 IA gratuita: GitHub Copilot Free + plano B

### Copilot Free (principal)

1. No VS Code, clique no ícone do **Copilot** na barra de status (canto inferior direito)
   → **Sign in** / **Use Copilot for free** → entre com sua conta do GitHub.
2. Abra o chat: `Ctrl+Alt+I`. No seletor de modo, você vai usar:
   - **Ask**: perguntas e explicações;
   - **Agent**: a IA edita arquivos e usa ferramentas, como o MCP do Playwright (módulo 06).
3. O plano Free tem **cota mensal limitada** de completions e de mensagens no chat.
   Confira os limites atuais em https://github.com/features/copilot/plans.
   Dica: use o modo **Ask** para dúvidas e guarde o **Agent** para as práticas do módulo 06.

> [!WARNING]
> Verifique a política da sua empresa sobre o uso de contas pessoais de IA. No curso só
> usamos a Hat Store, que é fictícia. **Nunca** cole código, dados ou credenciais de
> projetos reais em uma IA de conta pessoal ou gratuita.

### Plano B — chat no navegador (sem instalar nada)

Se a cota do Copilot acabar ou ele não estiver disponível: ChatGPT (https://chatgpt.com) ou
Gemini (https://gemini.google.com) no plano gratuito. O fluxo é *copiar e colar*: você leva
para o chat o HTML ou o erro e traz a resposta para o VS Code. O módulo 06 ensina esse fluxo também.

---

## Problemas comuns

| Sintoma | Causa provável | Solução |
|---|---|---|
| `Cannot connect to Podman` | Máquina Podman parada | `podman machine start` ou abra o Podman Desktop |
| `podman compose` não encontrado | Compose não instalado | Podman Desktop → Settings → Resources → Compose → Setup |
| `port 8080 already in use` | Outra aplicação na porta | Feche a aplicação ou pare containers antigos (`podman ps`, `podman stop <id>`) |
| Todos os testes falham com `ECONNREFUSED` | Hat Store não está no ar | Suba com `podman compose up -d` na pasta `hatstore-app` |
| `npx` bloqueado no PowerShell | Execution policy | `Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned` |
| Testes falham por estoque/dados estranhos | Banco "sujo" de muitas execuções | `podman compose down -v` e `podman compose up --build -d` |
| Firefox: `Host system is missing dependencies` / `Failed to launch` | Antivírus ou política corporativa bloqueando o Firefox do Playwright | O curso usa só o Chromium (`npm test`). O Firefox é opcional; ignore ou peça liberação ao suporte |

Próximo módulo: [01 — Teoria de Automação de Testes](../01-teoria-automacao-testes/README.md)
