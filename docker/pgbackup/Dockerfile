FROM restic/restic:0.16.2 AS restic

FROM postgres:16.1-bullseye

RUN apt-get update && apt-get install -y curl
COPY --from=restic /usr/bin/restic /usr/local/bin/restic

COPY backup.sh /backup.sh
ENTRYPOINT ["/backup.sh"]
