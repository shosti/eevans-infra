FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y \
      acl \
      attr \
      dnsutils \
      samba \
      samba-dsdb-modules \
      samba-vfs-modules \
      winbind \
      krb5-config \
      krb5-user

COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
