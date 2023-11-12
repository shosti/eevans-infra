#!/usr/bin/env bash

set -euo pipefail

run_primary() {
    touch /rsyncd.secret
    chmod 0600 /rsyncd.secret
    echo "sysvol-replication:$SYNC_SECRET" >> /rsyncd.secret

    cat <<EOS > /etc/rsyncd.conf
[SysVol]
path = $DATA_DIR
comment = Samba Sysvol Share
uid = root
gid = root
read only = yes
auth users = sysvol-replication
secrets file = /rsyncd.secret
EOS
    exec rsync -v --daemon --no-detach --address "$BIND_ADDRESS"
}

run_replica() {
    touch /rsyncd.secret
    chmod 0600 /rsyncd.secret
    echo "$SYNC_SECRET" > /rsyncd.secret

    cat <<EOS > /etc/crontabs/root
$SCHEDULE /sync.sh >> /dev/stdout 2>> /dev/stderr
EOS
    exec crond -f
}

if [ "$ROLE" = primary ]; then
    run_primary
elif [ "$ROLE" = replica ]; then
    run_replica
fi
