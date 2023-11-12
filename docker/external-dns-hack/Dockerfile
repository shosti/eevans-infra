FROM golang:1.16.7-buster as build

WORKDIR /build
RUN git clone -b k8s-1.22 https://github.com/rjosephwright/external-dns.git && \
    cd external-dns && \
    make build.amd64

FROM docker.io/bitnami/external-dns:0.9.0-debian-10-r0
COPY --from=build /build/external-dns/build/external-dns /opt/bitnami/external-dns/bin/external-dns
