-- 1. Kembalikan type kolom payment_status ke VARCHAR
ALTER TABLE orders
ALTER COLUMN payment_status TYPE VARCHAR(30);

-- 2. Hapus default value
ALTER TABLE orders
ALTER COLUMN payment_status DROP DEFAULT;

-- 3. Hapus kolom-kolom payment
ALTER TABLE orders
DROP COLUMN payment_method,
DROP COLUMN payment_status,
DROP COLUMN payment_id;

-- 4. Hapus enum type (SETELAH kolom tidak pakai enum lagi)
DROP TYPE IF EXISTS payment_status_enum;
