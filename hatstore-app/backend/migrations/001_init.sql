-- Migration única: schema completo da Hat Store
-- Criada em 2026-03-16 (consolidação das migrations 001-007)

-- ============================================================================
-- TABELAS
-- ============================================================================

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    nome VARCHAR(255),
    cpf VARCHAR(14) UNIQUE,
    data_nascimento VARCHAR(10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS hats (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    quantidade INTEGER NOT NULL DEFAULT 0,
    categoria VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pedidos (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    cpf VARCHAR(14) NOT NULL,
    email VARCHAR(255) NOT NULL,
    telefone VARCHAR(15) NOT NULL,
    endereco VARCHAR(255) NOT NULL,
    numero VARCHAR(10) NOT NULL,
    bairro VARCHAR(100) NOT NULL,
    cep VARCHAR(9) NOT NULL,
    cidade VARCHAR(100) NOT NULL,
    uf VARCHAR(2) NOT NULL,
    pagamento VARCHAR(20) NOT NULL,
    total DECIMAL(10, 2) NOT NULL,
    frete REAL NOT NULL DEFAULT 0,
    cupom VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pedido_items (
    id SERIAL PRIMARY KEY,
    pedido_id INTEGER NOT NULL REFERENCES pedidos(id) ON DELETE CASCADE,
    hat_id INTEGER NOT NULL REFERENCES hats(id),
    nome VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    quantidade INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS cupom_hatoff_usado (
    cpf VARCHAR(14) PRIMARY KEY,
    usado BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id SERIAL PRIMARY KEY,
    action VARCHAR(50) NOT NULL,
    user_email VARCHAR(255),
    details JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS login_attempts (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    count INTEGER NOT NULL DEFAULT 0,
    last_attempt TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS settings (
    key VARCHAR(255) PRIMARY KEY,
    value VARCHAR(255) NOT NULL
);

-- ============================================================================
-- ÍNDICES
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_hats_categoria ON hats(categoria);

CREATE INDEX IF NOT EXISTS idx_pedidos_cpf ON pedidos(cpf);
CREATE INDEX IF NOT EXISTS idx_pedidos_created_at ON pedidos(created_at);
CREATE INDEX IF NOT EXISTS idx_pedidos_email ON pedidos(email);
CREATE INDEX IF NOT EXISTS idx_pedidos_email_created_at ON pedidos(email, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_pedido_items_pedido_id ON pedido_items(pedido_id);
CREATE INDEX IF NOT EXISTS idx_pedido_items_hat_id ON pedido_items(hat_id);

CREATE INDEX IF NOT EXISTS idx_audit_logs_user_email ON audit_logs(user_email);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);

CREATE INDEX IF NOT EXISTS idx_login_attempts_email ON login_attempts(email);

-- ============================================================================
-- SEED: Chapéus
-- ============================================================================

INSERT INTO hats (nome, price, quantidade, categoria) VALUES
    ('Chapéu Panamá', 120.00, 100, 'importado'),
    ('Chapéu Fedora', 150.00, 80, 'importado'),
    ('Chapéu Bucket', 49.90, 150, 'importado'),
    ('Chapéu Cowboy', 109.90, 5, 'importado'),
    ('Chapéu Floppy', 79.90, 120, 'importado'),
    ('Chapéu Bowler', 69.90, 7, 'importado'),
    ('Chapéu Beanie', 39.90, 20, 'importado'),
    ('Chapéu Pork Pie', 59.90, 0, 'importado'),
    ('Chapéu Trilby', 84.90, 9, 'importado'),
    ('Chapéu Snapback', 44.90, 0, 'nacional'),
    ('Chapéu Sertanejo', 99.90, 110, 'nacional'),
    ('Chapéu Gaúcho', 129.90, 130, 'nacional'),
    ('Chapéu Cangaceiro', 139.90, 40, 'nacional'),
    ('Chapéu de Pescador', 29.90, 16, 'nacional'),
    ('Chapéu Gustavo Carvalho', 60.00, 1000, 'crescer')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- SEED: Settings e Feature Flags
-- ============================================================================

INSERT INTO settings (key, value) VALUES
    ('audit_frontend_enabled', 'true'),
    ('filters_enabled', 'true'),
    ('platform_enabled', 'true'),
    ('fix_password_audit_log', 'true'),
    ('fix_hatoff_discount', 'true'),
    ('new_hatstore', 'true')
ON CONFLICT (key) DO NOTHING;
