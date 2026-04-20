#!/bin/sh
set -eu

timestamp="$(date +%Y%m%d-%H%M%S)"
target_dir="${BACKUP_ROOT:-/backups}/postgres"
mkdir -p "$target_dir"

pg_dump \
  --host=postgres \
  --username="${POSTGRES_USER}" \
  --dbname="${POSTGRES_DB}" \
  --format=custom \
  --file="${target_dir}/nutritrack-${timestamp}.dump"
