#!/usr/bin/env bash

set -euo pipefail

TAG=harbor.eevans.me/library/sambasync:v1
docker build . -t "$TAG" --platform linux/amd64
docker push "$TAG"
