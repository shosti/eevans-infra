FROM alpine:3.18
RUN apk add keepalived bash && \
    adduser -S keepalived_script
COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
