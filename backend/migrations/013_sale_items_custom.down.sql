-- Remove custom items support
DELETE FROM sale_items WHERE product_id IS NULL;
ALTER TABLE sale_items DROP COLUMN IF EXISTS product_name;
ALTER TABLE sale_items ALTER COLUMN product_id SET NOT NULL;
