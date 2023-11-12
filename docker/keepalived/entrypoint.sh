#!/usr/bin/env bash

set -euo pipefail

exec keepalived --log-console --dont-fork "$@"
