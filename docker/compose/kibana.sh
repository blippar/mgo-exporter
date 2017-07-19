#!/bin/sh
set -e

echo "> Waiting for elastic to become available"
until curl -sS --fail "http://elastic:9200/_cat/health" >/dev/null 2>/dev/null; do
    sleep 1
done

echo "> Check if kibana index exists"

if curl -sS --fail "http://elastic:9200/_cat/health" >/dev/null 2>/dev/null; then
    echo "> Create kibana index"
    curl -sS -XPUT "http://elastic:9200/.kibana/?pretty" -d @/opt/provision/kibana_index.json

    echo "> Starting index-patterns and visualizations imports"
    cd /opt/kibana
    /opt/kibana/import.sh -imvd -e "http://elastic:9200" MgoTest

else
    echo "> Kibana index already exists, skipping provisionning"
fi
