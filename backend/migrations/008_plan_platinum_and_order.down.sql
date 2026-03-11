-- Hapus Platinum
DELETE FROM plan_config WHERE plan_slug = 'platinum';

ALTER TABLE plan_config DROP COLUMN IF EXISTS sort_order;

ALTER TABLE tenants DROP CONSTRAINT IF EXISTS tenants_plan_check;
ALTER TABLE tenants ADD CONSTRAINT tenants_plan_check CHECK (plan IN (
    'free', 'premium', 'premium_1m', 'premium_3m', 'premium_6m', 'premium_1y'
));

ALTER TABLE subscriptions DROP CONSTRAINT IF EXISTS subscriptions_plan_check;
ALTER TABLE subscriptions ADD CONSTRAINT subscriptions_plan_check CHECK (plan IN (
    'free', 'premium', 'premium_1m', 'premium_3m', 'premium_6m', 'premium_1y'
));
