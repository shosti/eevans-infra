#!/usr/bin/env bash
#
# See https://wiki.samba.org/index.php/Rsync_based_SysVol_replication_workaround

set -euo pipefail

rsync -XAavz \
      --delete-after \
      --password-file=/rsyncd.secret \
      rsync://sysvol-replication@"$RSYNC_MASTER"/SysVol/ "$DATA_DIR/"
