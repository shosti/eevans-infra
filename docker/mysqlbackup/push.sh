#!/usr/bin/env bash

set -euo pipefail

TAG=harbor.eevans.me/library/mysqlbackup:v2
docker build . -t "$TAG"
docker push "$TAG"
