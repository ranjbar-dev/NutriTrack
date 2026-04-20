#!/bin/sh
set -eu

timestamp="$(date +%Y%m%d-%H%M%S)"
target_dir="${BACKUP_ROOT:-/backups}/uploads"
mkdir -p "$target_dir"

tar -czf "${target_dir}/uploads-${timestamp}.tar.gz" -C /data uploads
