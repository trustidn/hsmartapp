-- Superadmins: users untuk panel admin (terpisah dari tenant users)
CREATE TABLE superadmins (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_superadmins_email ON superadmins(email);
