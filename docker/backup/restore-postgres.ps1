$ErrorActionPreference = 'Stop'

if (-not $env:BACKUP_FILE) {
  throw 'Set BACKUP_FILE before running restore-postgres.ps1'
}

docker compose run --rm -e BACKUP_FILE=$env:BACKUP_FILE -e RESTORE_DATABASE=$env:RESTORE_DATABASE postgres-restore
