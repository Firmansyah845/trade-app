-- Table stocks
CREATE TABLE stocks (
                        code VARCHAR(10) PRIMARY KEY,
                        name VARCHAR(255) NOT NULL,
                        is_sharia BOOLEAN NOT NULL DEFAULT false,
                        is_active BOOLEAN NOT NULL DEFAULT true,
                        price DECIMAL(20, 2) NOT NULL,
                        created_at TIMESTAMP DEFAULT NOW(),
                        updated_at TIMESTAMP DEFAULT NOW()
);

-- Index untuk query cepat
CREATE INDEX idx_stocks_sharia_active ON stocks(is_sharia, is_active) WHERE is_active = true;

-- Table user_balances
CREATE TABLE user_balances (
                               user_id BIGINT PRIMARY KEY,
                               balance DECIMAL(20, 2) NOT NULL DEFAULT 0,
                               created_at TIMESTAMP DEFAULT NOW(),
                               updated_at TIMESTAMP DEFAULT NOW()
);

-- Index untuk user_id lookup
CREATE INDEX idx_user_balances_user_id ON user_balances(user_id);

-- Table orders
CREATE TABLE orders (
                        id BIGSERIAL PRIMARY KEY,
                        user_id BIGINT NOT NULL,
                        stock_code VARCHAR(10) NOT NULL,
                        quantity BIGINT NOT NULL,
                        price DECIMAL(20, 2) NOT NULL,
                        total_amount DECIMAL(20, 2) NOT NULL,
                        order_type VARCHAR(10) NOT NULL, -- 'buy' or 'sell'
                        status VARCHAR(20) NOT NULL, -- 'pending', 'completed', 'failed'
                        created_at TIMESTAMP DEFAULT NOW(),
                        updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes untuk query performance
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_stock_code ON orders(stock_code);
CREATE INDEX idx_orders_created_at ON orders(created_at DESC);
CREATE INDEX idx_orders_status ON orders(status);

-- Composite index untuk common queries
CREATE INDEX idx_orders_user_status ON orders(user_id, status);