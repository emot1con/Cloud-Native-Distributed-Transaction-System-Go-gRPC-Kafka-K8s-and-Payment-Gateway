-- Rollback: Restore original payments table structure

-- Step 1: Add back old columns
ALTER TABLE payments 
    ADD COLUMN IF NOT EXISTS user_id INTEGER,
    ADD COLUMN IF NOT EXISTS total_price DOUBLE PRECISION,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Step 2: Migrate amount back to total_price
UPDATE payments SET total_price = amount WHERE total_price IS NULL;

-- Step 3: Get user_id from orders table
UPDATE payments p 
SET user_id = o.user_id 
FROM orders o 
WHERE p.order_id = o.id AND p.user_id IS NULL;

-- Step 4: Make columns NOT NULL
ALTER TABLE payments ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE payments ALTER COLUMN total_price SET NOT NULL;

-- Step 5: Drop new gateway columns
ALTER TABLE payments 
    DROP COLUMN IF EXISTS amount,
    DROP COLUMN IF EXISTS currency,
    DROP COLUMN IF EXISTS payment_method,
    DROP COLUMN IF EXISTS payment_channel,
    DROP COLUMN IF EXISTS gateway_name,
    DROP COLUMN IF EXISTS gateway_transaction_id,
    DROP COLUMN IF EXISTS gateway_order_id,
    DROP COLUMN IF EXISTS gateway_token,
    DROP COLUMN IF EXISTS gateway_redirect_url,
    DROP COLUMN IF EXISTS va_number,
    DROP COLUMN IF EXISTS qr_code_url,
    DROP COLUMN IF EXISTS paid_at,
    DROP COLUMN IF EXISTS expired_at;

-- Step 6: Drop new indexes
DROP INDEX IF EXISTS idx_payments_gateway_transaction_id;
DROP INDEX IF EXISTS idx_payments_gateway_order_id;

-- Step 7: Recreate old index
CREATE INDEX IF NOT EXISTS idx_payments_user_id ON payments(user_id);
