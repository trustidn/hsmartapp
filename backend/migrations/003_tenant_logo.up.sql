-- Logo URL for tenant (dapat diatur dari superadmin)
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS logo_url TEXT DEFAULT '';
