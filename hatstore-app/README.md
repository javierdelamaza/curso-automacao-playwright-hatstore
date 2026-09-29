# Hat Store (aplicação que vamos testar)

Loja virtual fictícia de chapéus, criada para treinamento de QA (programa CWI Crescer).
É o **sistema sob teste** do curso.

> [!IMPORTANT]
> Você **não precisa ler nem mexer** no código desta pasta. Ela está aqui só para você
> subir a aplicação na sua máquina. No curso, testamos a Hat Store como um usuário e
> como um consumidor da API: pela tela, pelo DevTools e pelo Swagger. É o mesmo cenário
> de um QA que recebe só a URL de um ambiente de homologação.

## Subir a aplicação

Pré-requisito: Podman Desktop instalado e a máquina Podman iniciada (veja o [módulo 00](../00-setup/README.md)).

```powershell
cd hatstore-app
podman compose up --build -d
```

Na primeira vez demora alguns minutos (baixa as imagens e compila o backend).
Quando terminar, abra **http://localhost:8080**.

| O quê | Endereço |
|---|---|
| Loja | http://localhost:8080 |
| Documentação funcional | http://localhost:8080/documentacao.html |
| Swagger (contrato da API) | http://localhost:8080/swagger/index.html |

## Parar, reiniciar e zerar os dados

```powershell
podman compose stop          # para os containers (mantém os dados)
podman compose start         # sobe de novo
podman compose down -v       # remove containers E banco (volta ao estado inicial)
podman compose up --build -d # sobe do zero
```

Use o `down -v` quando quiser um banco limpo (por exemplo, antes de gravar uma aula
ou quando o estoque de algum chapéu acabar de tanto rodar testes).

## O que a loja faz (visão de usuário)

- Catálogo com 15 chapéus em 3 categorias (Nacional, Internacional, Crescer).
- Busca por nome e filtros por categoria e faixa de preço.
- Cadastro (nome, e-mail, CPF, data de nascimento, senha forte, maior de 18 anos) e login.
- Carrinho com cupom de desconto `HATOFF`.
- Checkout com frete por região (CEP) e pagamento PIX ou Boleto.
- "Meus Pedidos" para o usuário logado.
