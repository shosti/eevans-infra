#!/usr/bin/env bash

set -euo pipefail

systemctl status >/dev/null &&
    systemctl status unbound >/dev/null && \
    [ "$(docker inspect -f "{{.State.Health.Status}}" pihole)" = healthy ]
