FROM alpine:3.14
RUN apk add keepalived bash && \
    adduser -S keepalived_script
COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
