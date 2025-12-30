-- Migration: Update payments table for Midtrans payment gateway integration
-- Remove duplicate columns (user_id, total_price) and add gateway-specific fields

-- Step 1: Add new columns for payment gateway
ALTER TABLE payments 
    ADD COLUMN IF NOT EXISTS amount DOUBLE PRECISION,
    ADD COLUMN IF NOT EXISTS currency VARCHAR(3) DEFAULT 'IDR',
    ADD COLUMN IF NOT EXISTS payment_method VARCHAR(50),
    ADD COLUMN IF NOT EXISTS payment_channel VARCHAR(50),
    ADD COLUMN IF NOT EXISTS gateway_name VARCHAR(50) DEFAULT 'midtrans',
    ADD COLUMN IF NOT EXISTS gateway_transaction_id VARCHAR(100),
    ADD COLUMN IF NOT EXISTS gateway_order_id VARCHAR(100),
    ADD COLUMN IF NOT EXISTS gateway_token VARCHAR(255),
    ADD COLUMN IF NOT EXISTS gateway_redirect_url TEXT,
    ADD COLUMN IF NOT EXISTS va_number VARCHAR(50),
    ADD COLUMN IF NOT EXISTS qr_code_url TEXT,
    ADD COLUMN IF NOT EXISTS paid_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS expired_at TIMESTAMP;

-- Step 2: Migrate existing total_price to amount
UPDATE payments SET amount = total_price WHERE amount IS NULL;

-- Step 3: Make amount NOT NULL after migration
ALTER TABLE payments ALTER COLUMN amount SET NOT NULL;

-- Step 4: Drop duplicate columns
ALTER TABLE payments 
    DROP COLUMN IF EXISTS user_id,
    DROP COLUMN IF EXISTS total_price,
    DROP COLUMN IF EXISTS updated_at;

-- Step 5: Drop old indexes
DROP INDEX IF EXISTS idx_payments_user_id;

-- Step 6: Add new indexes
CREATE INDEX IF NOT EXISTS idx_payments_gateway_transaction_id ON payments(gateway_transaction_id);
CREATE INDEX IF NOT EXISTS idx_payments_gateway_order_id ON payments(gateway_order_id);
