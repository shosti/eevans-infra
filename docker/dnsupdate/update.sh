#!/usr/bin/env bash

set -euo pipefail

API_ROOT=https://api.cloudflare.com/client/v4

cloudflare_api_call() {
    curl --fail \
         -H "Authorization: Bearer $API_TOKEN" \
         "$@"
}

create_dns_record() {
    zone_id="$1"
    ipv4="$2"

    cloudflare_api_call \
        -H "Content-Type: Application/json" \
        -X POST \
        "$API_ROOT/zones/$zone_id/dns_records" \
        --data "{\"type\": \"A\", \"name\": \"$HOST\", \"content\": \"$ipv4\", \"ttl\": 1}"
}

update_dns_record() {
    zone_id="$1"
    record_id="$2"
    ipv4="$3"

    cloudflare_api_call \
        -H "Content-Type: Application/json" \
        -X PUT \
        "$API_ROOT/zones/$zone_id/dns_records/$record_id" \
        --data "{\"type\": \"A\", \"name\": \"$HOST\", \"content\": \"$ipv4\", \"ttl\": 1}"
}

ipv4="$(curl --fail ipv4.icanhazip.com)"
zone_id="$(cloudflare_api_call "$API_ROOT/zones?name=$ZONE_NAME" | jq -r .result[0].id)"
record="$(cloudflare_api_call "$API_ROOT/zones/$zone_id/dns_records?name=$HOST&type=A" | jq -r .result[0] || true)"
if ! [ "$record" ] || [ "$record" = null ]; then
    create_dns_record "$zone_id" "$ipv4"
    exit 0
fi

if [ "$(echo "$record" | jq -r .content)" = "$ipv4" ]; then
    echo "Record is up to date"
    exit 0
fi

update_dns_record "$zone_id" "$(echo "$record" | jq -r .id)" "$ipv4"
