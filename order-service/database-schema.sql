CREATE TYPE payment_status_enum AS ENUM (
  'PENDING',
  'SUCCESS',
  'FAILED',
  'EXPIRED'
);

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    status DEFAULT 'pending' VARCHAR(50) NOT NULL,
    total_amount NUMERIC(12,2) NOT NULL,
    email VARCHAR(255) NOT NULL;

    payment_method VARCHAR(50),
    payment_status payment_status_enum DEFAULT 'PENDING',
    payment_id VARCHAR(100),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
