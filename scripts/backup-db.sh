#!/usr/bin/env bash
# Backup database PostgreSQL (HSmart).
# Usage: ./scripts/backup-db.sh
# Env: DATABASE_URL (wajib), BACKUP_DIR (default: ./backups/db)

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
: "${DATABASE_URL:?Set DATABASE_URL}"
: "${BACKUP_DIR:=$PROJECT_ROOT/backups/db}"
[[ "$BACKUP_DIR" != /* ]] && BACKUP_DIR="$PROJECT_ROOT/$BACKUP_DIR"

mkdir -p "$BACKUP_DIR"
TIMESTAMP=$(date +%Y%m%d-%H%M)
OUTPUT="$BACKUP_DIR/hsmart-${TIMESTAMP}.sql.gz"

echo "[backup-db] Starting backup to $OUTPUT"
pg_dump "$DATABASE_URL" | gzip > "$OUTPUT"
echo "[backup-db] Done: $OUTPUT ($(du -h "$OUTPUT" | cut -f1))"
