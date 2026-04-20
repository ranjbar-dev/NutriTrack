$ErrorActionPreference = 'Stop'

if (-not $env:BACKUP_FILE) {
  throw 'Set BACKUP_FILE before running restore-uploads.ps1'
}

docker compose run --rm -e BACKUP_FILE=$env:BACKUP_FILE uploads-restore
