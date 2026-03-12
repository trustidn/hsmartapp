-- Support custom items: product_id nullable, product_name for display
ALTER TABLE sale_items ALTER COLUMN product_id DROP NOT NULL;

ALTER TABLE sale_items ADD COLUMN IF NOT EXISTS product_name TEXT;

-- Backfill product_name for existing rows (from products)
UPDATE sale_items
SET product_name = p.name
FROM products p
WHERE sale_items.product_id = p.id AND sale_items.product_name IS NULL;
