#!/usr/bin/env bash

set -euo pipefail

cat <<EOS > /etc/krb5.conf
[libdefaults]
default_realm = ${REALM}

[realms]
${REALM} = {
           kdc = ${KDC_ADDRESS}
           admin_server = ${KDC_ADDRESS}
}
EOS

if ! [ -f /secrets/heimdal.mkey ]; then
    echo "Please mount a master key at /secrets/heimdal.mkey"
    exit 1
fi

if ! [ -f /data/heimdal.db ]; then
    kadmin --config-file=/config/kdc.conf -l init \
           --realm-max-ticket-life=unlimited  \
           --realm-max-renewable-life=unlimited \
           "$REALM"
fi

exec /usr/lib/heimdal-servers/kdc --config-file=/config/kdc.conf "$@"
