#!/usr/bin/env bash

set -euo pipefail

TAG=harbor.eevans.me/library/sambadc:v1
docker build . -t "$TAG" --platform linux/amd64
docker push "$TAG"
