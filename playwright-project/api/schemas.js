// Contratos (JSON Schema) montados a partir do Swagger em http://localhost:8080/swagger/index.html
export const chapeuSchema = {
  type: 'object',
  required: ['id', 'nome', 'price', 'quantidade', 'categoria', 'temEstoque'],
  properties: {
    id: { type: 'integer' },
    nome: { type: 'string', minLength: 1 },
    price: { type: 'number', minimum: 0 },
    quantidade: { type: 'integer', minimum: 0 },
    categoria: { enum: ['nacional', 'importado', 'crescer'] },
    temEstoque: { type: 'boolean' },
  },
};

export const listaDeChapeusSchema = { type: 'array', minItems: 1, items: chapeuSchema };

export const loginSchema = {
  type: 'object',
  required: ['message', 'token', 'email'],
  properties: {
    message: { type: 'string' },
    token: { type: 'string', pattern: '^[\\w-]+\\.[\\w-]+\\.[\\w-]+$' },
    email: { type: 'string', format: 'email' },
  },
};

export const freteSchema = {
  type: 'object',
  required: ['cep', 'regiao', 'custo_por_item', 'quantidade', 'valor_frete'],
  properties: {
    cep: { type: 'string' },
    regiao: { enum: ['Sudeste', 'Sul', 'Nordeste', 'Centro-Oeste', 'Norte'] },
    custo_por_item: { type: 'number' },
    quantidade: { type: 'integer' },
    valor_frete: { type: 'number' },
  },
};

const itemPedidoSchema = {
  type: 'object',
  required: ['hat_id', 'nome', 'price', 'quantidade'],
  properties: {
    hat_id: { type: 'integer' },
    nome: { type: 'string' },
    price: { type: 'number' },
    quantidade: { type: 'integer', minimum: 1 },
  },
};

export const pedidoSchema = {
  type: 'object',
  required: ['id', 'nome', 'cpf', 'email', 'pagamento', 'total', 'frete', 'itens', 'created_at'],
  properties: {
    id: { type: 'integer' },
    nome: { type: 'string' },
    cpf: { type: 'string', pattern: '^\\d{3}\\.\\d{3}\\.\\d{3}-\\d{2}$' },
    email: { type: 'string', format: 'email' },
    pagamento: { enum: ['pix', 'boleto'] },
    total: { type: 'number' },
    frete: { type: 'number' },
    cupom: { type: 'string' },
    itens: { type: 'array', minItems: 1, items: itemPedidoSchema },
    created_at: { type: 'string', format: 'date-time' },
  },
};
