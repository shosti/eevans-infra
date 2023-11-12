#!/usr/bin/env bash

set -euo pipefail

DNS_REALM="$(echo "$REALM" | tr '[:upper:]' '[:lower:]')"

provision_dc() {
    rm -f /etc/samba/smb.conf
    rm -f /etc/krb5.conf

    samba-tool domain provision \
               --server-role dc \
               --use-rfc2307 \
               --realm "$REALM"\
               --domain "$DOMAIN" \
               --dns-backend SAMBA_INTERNAL \
               --adminpass "$ADMIN_PASS" \
               --option="dns forwarder = 1.1.1.1" \
               --option="interfaces = $SELF_IP" \
               --option="vfs objects = dfs_samba4 acl_xattr streams_xattr xattr_tdb" # see https://github.com/lxc/lxc/issues/2708#issuecomment-473466062

    cp /etc/samba/smb.conf /var/lib/samba/smb.conf
    touch /var/lib/samba/PROVISIONED
}

join_dc() {
    cat <<EOS > /etc/krb5.conf
[libdefaults]
    dns_lookup_realm = false
    dns_lookup_kdc = true
    default_realm = $REALM
EOS
    rm -f /etc/samba/smb.conf

    samba-tool domain join "$DNS_REALM" DC \
               -U"$DOMAIN\administrator" \
               --password="$ADMIN_PASS" \
               --option='idmap_ldb:use rfc2307 = yes' \
               --option="dns forwarder = 1.1.1.1" \
               --option="interfaces = $SELF_IP"
    cp /etc/samba/smb.conf /var/lib/samba/smb.conf
    touch /var/lib/samba/PROVISIONED
}

if [ "$ROLE" = primary ] && ! [ -e /var/lib/samba/PROVISIONED ]; then
    provision_dc
fi

if [ "$ROLE" = secondary ] && ! [ -e /var/lib/samba/PROVISIONED ]; then
    join_dc
fi

cp /var/lib/samba/private/krb5.conf /etc/krb5.conf
cp /var/lib/samba/smb.conf /etc/samba/smb.conf

exec /usr/sbin/samba --foreground --debuglevel=1 --debug-stdout
