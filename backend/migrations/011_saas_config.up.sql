-- Pengaturan SaaS (nama, logo, kontak, rekening) - satu baris konfigurasi
CREATE TABLE saas_config (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    app_name TEXT DEFAULT 'HSmart',
    logo_url TEXT,
    admin_contact TEXT,
    bank_name TEXT,
    bank_account_number TEXT,
    bank_account_name TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO saas_config (id) VALUES (uuid_generate_v4());

-- Order: pastikan status paid dan payment_proof_url ada (jika 010 lama tanpa kolom ini)
ALTER TABLE subscription_orders ADD COLUMN IF NOT EXISTS payment_proof_url TEXT;
ALTER TABLE subscription_orders DROP CONSTRAINT IF EXISTS subscription_orders_status_check;
ALTER TABLE subscription_orders ADD CONSTRAINT subscription_orders_status_check 
    CHECK (status IN ('pending', 'paid', 'approved', 'rejected'));
