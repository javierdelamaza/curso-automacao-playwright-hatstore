// Camada de API: centraliza endpoints para os testes não repetirem URL e payload.
export class HatStoreApi {
  constructor(request) {
    this.request = request;
  }

  listarChapeus(params) {
    return this.request.get('/api/hats', { params });
  }

  cadastrar(usuario) {
    return this.request.post('/auth/register', {
      data: {
        email: usuario.email,
        password: usuario.password,
        nome: usuario.nome,
        cpf: usuario.cpf,
        data_nascimento: usuario.dataNascimento,
      },
    });
  }

  login(email, password) {
    return this.request.post('/auth/login', { data: { email, password } });
  }

  async obterToken(email, password) {
    const response = await this.login(email, password);
    if (!response.ok()) {
      throw new Error(`Login via API falhou (${response.status()}): ${await response.text()}`);
    }
    return (await response.json()).token;
  }

  calcularFrete(cep, quantidade) {
    return this.request.get('/api/frete', { params: { cep, quantidade } });
  }

  // Os métodos abaixo exigem token: use com o request autenticado da fixture.
  consultarPedidosPorCpf(cpf) {
    return this.request.get('/api/pedidos', { params: { cpf } });
  }

  perfil() {
    return this.request.get('/api/me');
  }

  criarPedido(pedido) {
    return this.request.post('/api/pedido', { data: pedido });
  }

  meusPedidos() {
    return this.request.get('/api/meus-pedidos');
  }
}
