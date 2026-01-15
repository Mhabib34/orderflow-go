CREATE TYPE payment_status_enum AS ENUM (
  'PENDING',
  'SUCCESS',
  'FAILED',
  'EXPIRED'
);

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    order_id UUID NOT NULL,
    
    amount BIGINT NOT NULL,

    method VARCHAR(50) NOT NULL,          -- qris, va, ewallet
    provider_ref_id VARCHAR(100) UNIQUE,  -- order_id dari provider

    status payment_status_enum DEFAULT 'PENDING',

    expired_at TIMESTAMP,
    paid_at TIMESTAMP,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
