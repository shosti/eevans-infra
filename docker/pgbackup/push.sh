#!/usr/bin/env bash

set -euo pipefail

TAG=gitea.eevans.me/shosti/pgbackup:v7-pg16
docker buildx build . -t "$TAG" --platform linux/amd64 --progress=plain
docker push "$TAG"
