-- 1. Tambahkan kolom dulu
ALTER TABLE orders
ADD COLUMN payment_method VARCHAR(50),
ADD COLUMN payment_status VARCHAR(30),
ADD COLUMN payment_id VARCHAR(100);

-- 2. Isi data lama
UPDATE orders
SET payment_status = 'PENDING'
WHERE payment_status IS NULL;

-- 3. Buat ENUM dulu
CREATE TYPE payment_status_enum AS ENUM (
  'PENDING',
  'SUCCESS',
  'FAILED',
  'EXPIRED'
);

-- 4. HAPUS default SEBELUM ganti type (INI KUNCI)
ALTER TABLE orders
ALTER COLUMN payment_status DROP DEFAULT;

-- 5. Ubah type ke ENUM
ALTER TABLE orders
ALTER COLUMN payment_status
TYPE payment_status_enum
USING payment_status::payment_status_enum;

-- 6. Baru set default ENUM
ALTER TABLE orders
ALTER COLUMN payment_status SET DEFAULT 'PENDING';
