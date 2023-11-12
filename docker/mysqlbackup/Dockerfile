FROM restic/restic:0.12.1 AS restic

FROM alpine:3.14

RUN apk add bash curl mysql-client && \
    addgroup -S mybackup && \
    adduser -S mybackup -G mybackup && \
    mkdir /dump && \
    chown mybackup:mybackup /dump
COPY --from=restic /usr/bin/restic /usr/local/bin/restic

USER mybackup
WORKDIR /dump

COPY backup.sh /backup.sh
CMD ["/backup.sh"]
