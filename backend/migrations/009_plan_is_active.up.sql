-- Soft delete: nonaktifkan plan tanpa menghapus (tenant yang pakai tetap jalan)
ALTER TABLE plan_config ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT true;
