#!/usr/bin/env bash
# Backup folder uploads (logos, payment-proofs).
# Usage: ./scripts/backup-uploads.sh
# Env: UPLOAD_DIR (default: ./uploads), BACKUP_DIR (default: ./backups/uploads)

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
: "${UPLOAD_DIR:=$PROJECT_ROOT/uploads}"
: "${BACKUP_DIR:=$PROJECT_ROOT/backups/uploads}"

# Resolve relative UPLOAD_DIR from project root
[[ "$UPLOAD_DIR" != /* ]] && UPLOAD_DIR="$PROJECT_ROOT/$UPLOAD_DIR"
[[ "$BACKUP_DIR" != /* ]] && BACKUP_DIR="$PROJECT_ROOT/$BACKUP_DIR"

mkdir -p "$BACKUP_DIR"
TIMESTAMP=$(date +%Y%m%d-%H%M)
OUTPUT="$BACKUP_DIR/uploads-${TIMESTAMP}.tar.gz"

if [[ ! -d "$UPLOAD_DIR" ]]; then
  echo "[backup-uploads] WARN: $UPLOAD_DIR tidak ada, skip."
  exit 0
fi

echo "[backup-uploads] Starting backup to $OUTPUT"
PARENT="$(cd "$(dirname "$UPLOAD_DIR")" && pwd)"
NAME="$(basename "$UPLOAD_DIR")"
tar -czf "$OUTPUT" -C "$PARENT" "$NAME"
echo "[backup-uploads] Done: $OUTPUT ($(du -h "$OUTPUT" | cut -f1))"
