#!/bin/sh
set -eu

if [ -z "${BACKUP_FILE:-}" ]; then
  echo "BACKUP_FILE is required"
  exit 1
fi

rm -rf /data/uploads/*
mkdir -p /data/uploads
tar -xzf "${BACKUP_ROOT:-/backups}/uploads/${BACKUP_FILE}" -C /data
