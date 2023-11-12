#!/usr/bin/env bash

set -euo pipefail

cert_key="$(sudo kubeadm init phase upload-certs --upload-certs | tail -1)"
echo "$(sudo kubeadm token create --print-join-command --ttl 3h) --control-plane --certificate-key $cert_key"
