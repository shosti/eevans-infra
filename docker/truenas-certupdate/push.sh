#!/usr/bin/env bash

set -euo pipefail

TAG=harbor.eevans.me/library/truenas-certupdate:v3
docker build . -t "$TAG"
docker push "$TAG"
