-- Tenant settings for merchant configuration
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS receipt_footer TEXT DEFAULT '';
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS default_payment TEXT DEFAULT 'cash';
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS whatsapp_number TEXT DEFAULT '';

-- Sales: add payment methods (qris, ewallet)
ALTER TABLE sales DROP CONSTRAINT IF EXISTS sales_payment_method_check;
ALTER TABLE sales ADD CONSTRAINT sales_payment_method_check CHECK (payment_method IN ('cash', 'qris', 'transfer', 'ewallet', 'other'));
