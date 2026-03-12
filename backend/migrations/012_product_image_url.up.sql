-- Add image_url to products (foto produk, path relatif /uploads/products/xxx)
ALTER TABLE products ADD COLUMN IF NOT EXISTS image_url TEXT;
