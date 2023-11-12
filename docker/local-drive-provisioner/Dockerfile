FROM ubuntu:jammy

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends ca-certificates curl lvm2 udev parted cryptsetup && \
    curl -L https://github.com/mikefarah/yq/releases/download/v4.35.2/yq_linux_amd64 -o /usr/local/bin/yq && \
    [ "$(sha256sum /usr/local/bin/yq | awk '{print $1}')" = 8afd786b3b8ba8053409c5e7d154403e2d4ed4cf3e93c237462dc9ef75f38c8d ] || exit 1 && \
    chmod +x /usr/local/bin/yq

COPY provision /provision
CMD ["/provision"]
