CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('order_created', 'order_updated', 'order_completed', 'order_cancelled')),
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk mempercepat query berdasarkan order_id
CREATE INDEX idx_notifications_order_id ON notifications(order_id);

-- Index untuk mempercepat query berdasarkan type
CREATE INDEX idx_notifications_type ON notifications(type);

-- Index untuk mempercepat query berdasarkan is_read
CREATE INDEX idx_notifications_is_read ON notifications(is_read);

-- Index untuk mempercepat query berdasarkan created_at
CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);

-- Composite index untuk query yang sering digunakan (is_read + created_at)
CREATE INDEX idx_notifications_is_read_created_at ON notifications(is_read, created_at DESC);