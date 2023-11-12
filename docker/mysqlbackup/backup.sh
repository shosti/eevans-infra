#!/usr/bin/env bash

set -euo pipefail

mysqldump --host "$MYSQL_HOST" \
          --user "$MYSQL_USER" \
          -p"$(cat "$MYSQL_PASS_FILE")" \
          "$MYSQL_DATABASE" > /dump/"$BACKUP_NAME.sql"
restic --verbose backup /dump/
curl "$HEALTH_CHECK_URL"
