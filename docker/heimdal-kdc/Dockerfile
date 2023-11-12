FROM ubuntu:focal
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && \
    apt-get install -y --no-install-recommends heimdal-kdc bash && \
    mkdir /data

COPY genmaster /usr/local/bin/genmaster
COPY kdc.conf /config/kdc.conf
COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
