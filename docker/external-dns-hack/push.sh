#!/usr/bin/env bash

set -euo pipefail

TAG=shosti/external-dns:0.9.0-k8s-1.22-v1
docker build . -t "$TAG"
docker push "$TAG"
