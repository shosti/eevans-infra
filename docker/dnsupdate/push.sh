#!/usr/bin/env bash

set -euo pipefail

TAG=gitea.eevans.me/shosti/dnsupdate:v1
docker build . -t "$TAG"
docker push "$TAG"
