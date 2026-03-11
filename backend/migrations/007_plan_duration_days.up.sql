-- Durasi langganan per plan (hari). 0 atau -1 = tanpa batas (free).
ALTER TABLE plan_config ADD COLUMN IF NOT EXISTS duration_days INT NOT NULL DEFAULT 0 CHECK (duration_days >= -1);

-- Default durasi
UPDATE plan_config SET duration_days = 0 WHERE plan_slug = 'free';
UPDATE plan_config SET duration_days = 30 WHERE plan_slug = 'premium';
UPDATE plan_config SET duration_days = 30 WHERE plan_slug = 'premium_1m';
UPDATE plan_config SET duration_days = 90 WHERE plan_slug = 'premium_3m';
UPDATE plan_config SET duration_days = 180 WHERE plan_slug = 'premium_6m';
UPDATE plan_config SET duration_days = 365 WHERE plan_slug = 'premium_1y';
