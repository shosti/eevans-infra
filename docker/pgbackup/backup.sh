#!/usr/bin/env bash

set -euo pipefail

pg_dump > /dump/"$BACKUP_NAME.sql"
restic --verbose backup /dump/
curl "$HEALTH_CHECK_URL"
