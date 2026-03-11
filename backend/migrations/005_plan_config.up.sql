-- Plan config: atur max_products dan report_days per plan (oleh superadmin)
-- Plan: free, premium_1m (1 bulan), premium_3m (3 bulan), premium_6m (6 bulan), premium_1y (1 tahun)
CREATE TABLE plan_config (
    plan_slug TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    duration_months INT NOT NULL DEFAULT 0,  -- 0 = unlimited (free)
    max_products INT NOT NULL DEFAULT 10 CHECK (max_products >= -1),  -- -1 = unlimited
    report_days INT NOT NULL DEFAULT 7 CHECK (report_days >= -1),       -- -1 = unlimited
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Default config (-1 = unlimited)
INSERT INTO plan_config (plan_slug, name, duration_months, max_products, report_days) VALUES
('free', 'Free', 0, 10, 7),
('premium', 'Premium', 1, -1, 30),
('premium_1m', 'Premium 1 Bulan', 1, -1, 30),
('premium_3m', 'Premium 3 Bulan', 3, -1, 90),
('premium_6m', 'Premium 6 Bulan', 6, -1, 180),
('premium_1y', 'Premium 1 Tahun', 12, -1, 365);

-- Perluas tenants.plan dan subscriptions.plan
ALTER TABLE tenants DROP CONSTRAINT IF EXISTS tenants_plan_check;
ALTER TABLE tenants ADD CONSTRAINT tenants_plan_check CHECK (plan IN (
    'free', 'premium', 'premium_1m', 'premium_3m', 'premium_6m', 'premium_1y'
));

ALTER TABLE subscriptions DROP CONSTRAINT IF EXISTS subscriptions_plan_check;
ALTER TABLE subscriptions ADD CONSTRAINT subscriptions_plan_check CHECK (plan IN (
    'free', 'premium', 'premium_1m', 'premium_3m', 'premium_6m', 'premium_1y'
));
