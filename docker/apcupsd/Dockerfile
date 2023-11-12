FROM debian:bullseye
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y --no-install-recommends apcupsd
COPY ./apcupsd.conf ./apccontrol /etc/apcupsd/

CMD ["apcupsd", "-b"]
