-- database/migrations/001_create_indexes.sql

-- Index on user_id in cart_items (frequent lookup)
CREATE INDEX IF NOT EXISTS idx_cart_user_id ON cart_items(user_id);

-- Index on product_id in cart_items (frequent lookup)
CREATE INDEX IF NOT EXISTS idx_cart_product_id ON cart_items(product_id);

-- Index on email in users (for login / lookup)
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Index on user_id in addresses
CREATE INDEX IF NOT EXISTS idx_address_user_id ON addresses(user_id);

-- Index on user_id in orders
CREATE INDEX IF NOT EXISTS idx_order_user_id ON orders(user_id);

-- Index on order status (e.g. for filter: 'placed', 'cancelled')
CREATE INDEX IF NOT EXISTS idx_order_status ON orders(status);
