-- Harga per plan (Rupiah). Free = 0.
ALTER TABLE plan_config ADD COLUMN IF NOT EXISTS price_rupiah BIGINT NOT NULL DEFAULT 0 CHECK (price_rupiah >= 0);

-- Default harga
UPDATE plan_config SET price_rupiah = 0 WHERE plan_slug = 'free';
UPDATE plan_config SET price_rupiah = 10000 WHERE plan_slug = 'premium';
UPDATE plan_config SET price_rupiah = 10000 WHERE plan_slug = 'premium_1m';
UPDATE plan_config SET price_rupiah = 27000 WHERE plan_slug = 'premium_3m';
UPDATE plan_config SET price_rupiah = 50000 WHERE plan_slug = 'premium_6m';
UPDATE plan_config SET price_rupiah = 100000 WHERE plan_slug = 'premium_1y';
