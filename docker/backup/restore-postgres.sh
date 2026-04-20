#!/bin/sh
set -eu

if [ -z "${BACKUP_FILE:-}" ]; then
  echo "BACKUP_FILE is required"
  exit 1
fi

target_db="${RESTORE_DATABASE:-$POSTGRES_DB}"

pg_restore \
  --clean \
  --if-exists \
  --no-owner \
  --host=postgres \
  --username="${POSTGRES_USER}" \
  --dbname="${target_db}" \
  "${BACKUP_ROOT:-/backups}/postgres/${BACKUP_FILE}"
