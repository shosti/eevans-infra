#!/usr/bin/env bash

set -euo pipefail

TAG=harbor.eevans.me/library/dnsupdate:v1
docker build . -t "$TAG"
docker push "$TAG"
