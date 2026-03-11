-- Order langganan: tenant buat order, upload bukti, admin verifikasi & approve
CREATE TABLE subscription_orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    plan_slug TEXT NOT NULL,
    amount_rupiah BIGINT NOT NULL DEFAULT 0 CHECK (amount_rupiah >= 0),
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'approved', 'rejected')),
    payment_note TEXT,
    payment_proof_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    approved_at TIMESTAMPTZ,
    approved_by UUID REFERENCES superadmins(id),
    rejection_reason TEXT
);

CREATE INDEX idx_subscription_orders_tenant ON subscription_orders(tenant_id);
CREATE INDEX idx_subscription_orders_status ON subscription_orders(status);
CREATE INDEX idx_subscription_orders_created ON subscription_orders(created_at DESC);
