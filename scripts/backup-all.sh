#!/usr/bin/env bash
# Backup database + uploads. Untuk cron berkala.
# Usage: ./scripts/backup-all.sh
# Env: DATABASE_URL, UPLOAD_DIR, BACKUP_DIR (opsional)
#
# Cron contoh (setiap hari jam 02:00):
# 0 2 * * * cd /path/to/hsmartapp && DATABASE_URL=postgres://... ./scripts/backup-all.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$PROJECT_ROOT"

# Default BACKUP_DIR di project root
export BACKUP_DIR="${BACKUP_DIR:-$PROJECT_ROOT/backups}"
export UPLOAD_DIR="${UPLOAD_DIR:-$PROJECT_ROOT/uploads}"

echo "[backup-all] $(date)"

"$SCRIPT_DIR/backup-db.sh"
"$SCRIPT_DIR/backup-uploads.sh"

echo "[backup-all] Selesai."
