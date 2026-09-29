// Massa FIXA: dados que vêm do seed do banco e não mudam entre execuções.
// Foi levantada olhando a tela e o GET /api/hats (Swagger), não o código-fonte.
export const CATEGORIAS = {
  nacional: 'nacional',
  importado: 'importado',
  crescer: 'crescer',
};

export const PRODUTOS = {
  panama: { id: 1, nome: 'Chapéu Panamá', price: 120.0, categoria: CATEGORIAS.importado },
  fedora: { id: 2, nome: 'Chapéu Fedora', price: 150.0, categoria: CATEGORIAS.importado },
  gaucho: { id: 12, nome: 'Chapéu Gaúcho', price: 129.9, categoria: CATEGORIAS.nacional },
  gustavoCarvalho: { id: 15, nome: 'Chapéu Gustavo Carvalho', price: 60.0, categoria: CATEGORIAS.crescer },
  porkPieSemEstoque: { id: 8, nome: 'Chapéu Pork Pie', price: 59.9, categoria: CATEGORIAS.importado },
};

export const TOTAL_PRODUTOS_CATALOGO = 15;

export const CUPOM_VALIDO = 'HATOFF';
export const PERCENTUAL_DESCONTO_HATOFF = 0.2;

export const ENDERECO_PADRAO = {
  telefone: '(11) 98765-4321',
  endereco: 'Rua das Flores',
  numero: '123',
  bairro: 'Centro',
  cep: '01310-100',
  cidade: 'São Paulo',
  uf: 'SP',
};

export const FRETE_SUDESTE_1_ITEM = { regiao: 'Sudeste', valor: 12.5 };
