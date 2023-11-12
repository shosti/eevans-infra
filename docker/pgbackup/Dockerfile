FROM restic/restic:0.15.1 AS restic

FROM postgres:16.0-bullseye

RUN apt-get update && apt-get install -y curl
COPY --from=restic /usr/bin/restic /usr/local/bin/restic

COPY backup.sh /backup.sh
ENTRYPOINT ["/backup.sh"]
