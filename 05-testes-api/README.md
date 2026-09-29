# Módulo 05 — Testes de API com Playwright

## 5.1 Swagger: a fonte da verdade (sem código-fonte)

Abra http://localhost:8080/swagger/index.html. Ali estão todos os endpoints, os parâmetros,
os corpos de requisição e de resposta e os códigos de status. Use **Try it out** para chamar a API pelo navegador.

| Endpoint | Auth? | Para quê |
|---|---|---|
| `GET /api/hats` | não | Catálogo (filtros `categoria`, `min`, `max`) |
| `POST /auth/register` | não | Cadastro |
| `POST /auth/login` | não | Login → devolve `token` (JWT) |
| `GET /api/frete?cep=&quantidade=` | não | Cálculo de frete por região |
| `GET /api/me` | **sim** | Perfil do usuário logado |
| `POST /api/pedido` | **sim** | Criar pedido |
| `GET /api/meus-pedidos` | **sim** | Pedidos do usuário logado |
| `GET /api/pedidos?cpf=` | **sim** | Pedidos por CPF |

> Outra fonte sem código: DevTools → aba **Network**. Faça login pela tela e veja a
> requisição `login`, com URL, corpo enviado e resposta. É assim que se descobre uma API que não tem Swagger.

## 5.2 A fixture `request` e o primeiro GET

O Playwright tem um cliente HTTP embutido. Não precisa de navegador, então roda muito rápido.

```js
import { test, expect } from '@playwright/test';

test('GET /api/hats retorna o catálogo', async ({ request }) => {
  const response = await request.get('/api/hats');
  expect(response.status()).toBe(200);
  const chapeus = await response.json();
  expect(chapeus).toHaveLength(15);
});
```

## 5.3 Organização: `tests/front` x `tests/api` e a camada de API

Assim como o Page Object encapsula telas, a classe
[api/HatStoreApi.js](../playwright-project/api/HatStoreApi.js) encapsula endpoints:

```js
export class HatStoreApi {
  constructor(request) { this.request = request; }
  listarChapeus(params) { return this.request.get('/api/hats', { params }); }
  login(email, password) { return this.request.post('/auth/login', { data: { email, password } }); }
  // ...
}
```

O teste fica: `await api.listarChapeus({ categoria: 'crescer' })`.

## 5.4 Autenticação: fixture de API com JWT

Endpoints protegidos exigem o cabeçalho `Authorization: Bearer <token>`. A fixture
`apiAutenticada` ([fixtures/index.js](../playwright-project/fixtures/index.js)) cria um
cliente HTTP que envia esse cabeçalho em toda chamada:

```js
apiAutenticada: async ({ playwright, baseURL, token }, use) => {
  const contexto = await playwright.request.newContext({
    baseURL,
    extraHTTPHeaders: { Authorization: `Bearer ${token}` },
  });
  await use(new HatStoreApi(contexto));
  await contexto.dispose();
},
```

A cadeia de fixtures resolve tudo sozinha: `apiAutenticada` → `token` → `usuarioCadastrado`.

## 5.5 Cobrindo os endpoints principais

| Arquivo | O que cobre |
|---|---|
| [hats.spec.js](../playwright-project/tests/api/hats.spec.js) | Catálogo, filtros, produto sem estoque |
| [auth.spec.js](../playwright-project/tests/api/auth.spec.js) | Cadastro (201, 409, 400), login (200, 401), `/api/me` com e sem token |
| [frete.spec.js](../playwright-project/tests/api/frete.spec.js) | Região por CEP (teste parametrizado), frete progressivo, CEP inválido |
| [pedidos.spec.js](../playwright-project/tests/api/pedidos.spec.js) | Criar pedido, sem token, meus pedidos, por CPF, cupom, pedido vazio |

**Teste parametrizado** (um teste por linha da tabela):

```js
const casos = [
  { cep: '01310-100', regiao: 'Sudeste', valor: 12.5 },
  { cep: '90010-000', regiao: 'Sul', valor: 15.0 },
];
for (const caso of casos) {
  test(`CEP ${caso.cep} é da região ${caso.regiao}`, async ({ api }) => { /* ... */ });
}
```

### Bug encontrado pelos testes

O teste *"cupom HATOFF só pode ser usado uma vez por CPF"* revelou que **qualquer pedido com o
cupom HATOFF devolve 500** ("Erro ao buscar chapéu"). O teste fica marcado assim:

```js
test.fail(true, 'BUG: POST /api/pedido com cupom HATOFF devolve 500');
```

`test.fail` diz "espero que este teste falhe". A suíte continua verde **e o bug fica
documentado**. Quando o dev corrigir, o teste passa a "falhar" e avisa para remover a marcação.
Nunca apague um teste só porque ele encontrou um bug.

Outra divergência que a exploração revela, sem abrir código: a documentação
(`/documentacao.html`) diz **25%** de desconto num lugar e **20%** em outro. A tela aplica 20%.
Pergunta de QA: qual é a regra certa? Isso vira uma dúvida para o PO, não um "chute" no teste.

## 5.6 Testes de contrato

O teste de contrato não confere **valores**, confere o **formato**: campos obrigatórios, tipos e
formatos. Se o backend renomear `price` para `preco`, o teste quebra **antes** do front.

Usamos **JSON Schema** ([api/schemas.js](../playwright-project/api/schemas.js)) e a biblioteca
gratuita **Ajv** ([api/validarContrato.js](../playwright-project/api/validarContrato.js)):

```js
export const chapeuSchema = {
  type: 'object',
  required: ['id', 'nome', 'price', 'quantidade', 'categoria', 'temEstoque'],
  properties: {
    id: { type: 'integer' },
    price: { type: 'number', minimum: 0 },
    categoria: { enum: ['nacional', 'importado', 'crescer'] },
    // ...
  },
};

expect(validarContrato(listaDeChapeusSchema, corpo)).toEqual([]);
```

Quando falha, a mensagem lista exatamente o que quebrou (ex.: `/0/price must be number`).

## 5.7 Exercício

Bloco **"Bloco 5"** em [exercicios.spec.js](../playwright-project/tests/exercicios/exercicios.spec.js):
crie um pedido via API e valide que **cada pedido** de `GET /api/meus-pedidos` respeita o
`pedidoSchema`.

## 5.8 Dica de IA: gerar o schema a partir da resposta

No Swagger, execute `GET /api/frete` e copie o JSON de resposta. Peça à IA:

> "Gere um JSON Schema (draft-07) para esta resposta [cole o JSON]. Marque todos os campos como
> obrigatórios, use `enum` para `regiao` com os valores Sudeste, Sul, Nordeste, Centro-Oeste e Norte."

Revise: tipos corretos? `integer` x `number`? Compare com o `freteSchema` do projeto.

Próximo módulo: [06 — IA e MCP](../06-ia-e-mcp/README.md)
