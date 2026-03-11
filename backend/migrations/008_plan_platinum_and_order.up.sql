-- Urutan tampilan plan + tambah Platinum
ALTER TABLE plan_config ADD COLUMN IF NOT EXISTS sort_order INT NOT NULL DEFAULT 0;

-- Set urutan: Free, Premium 1 Bulan, Premium 3 Bulan, Premium 6 Bulan, Premium 1 Tahun, Platinum
UPDATE plan_config SET sort_order = 0 WHERE plan_slug = 'free';
UPDATE plan_config SET sort_order = 1 WHERE plan_slug = 'premium_1m';
UPDATE plan_config SET sort_order = 2 WHERE plan_slug = 'premium_3m';
UPDATE plan_config SET sort_order = 3 WHERE plan_slug = 'premium_6m';
UPDATE plan_config SET sort_order = 4 WHERE plan_slug = 'premium_1y';
UPDATE plan_config SET sort_order = 5 WHERE plan_slug = 'platinum';
UPDATE plan_config SET sort_order = 99 WHERE plan_slug = 'premium';  -- premium legacy, tampil terakhir

-- Tambah Platinum jika belum ada
INSERT INTO plan_config (plan_slug, name, duration_months, max_products, report_days, sort_order)
SELECT 'platinum', 'Platinum', 12, -1, 365, 5
WHERE NOT EXISTS (SELECT 1 FROM plan_config WHERE plan_slug = 'platinum');

-- Set price dan duration untuk Platinum (jika kolom ada dari migration 006, 007)
UPDATE plan_config SET price_rupiah = 200000 WHERE plan_slug = 'platinum';
UPDATE plan_config SET duration_days = 365 WHERE plan_slug = 'platinum';

-- Perluas CHECK constraint
ALTER TABLE tenants DROP CONSTRAINT IF EXISTS tenants_plan_check;
ALTER TABLE tenants ADD CONSTRAINT tenants_plan_check CHECK (plan IN (
    'free', 'premium', 'premium_1m', 'premium_3m', 'premium_6m', 'premium_1y', 'platinum'
));

ALTER TABLE subscriptions DROP CONSTRAINT IF EXISTS subscriptions_plan_check;
ALTER TABLE subscriptions ADD CONSTRAINT subscriptions_plan_check CHECK (plan IN (
    'free', 'premium', 'premium_1m', 'premium_3m', 'premium_6m', 'premium_1y', 'platinum'
));
