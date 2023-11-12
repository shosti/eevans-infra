FROM alpine:3.16

RUN apk add bash rsync
COPY sync.sh entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]
